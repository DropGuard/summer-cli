package ${package};

import com.github.dropguard.summer.tx.Transactional;
import com.tngtech.archunit.core.domain.JavaCall;
import com.tngtech.archunit.core.domain.JavaClass;
import com.tngtech.archunit.core.domain.JavaClasses;
import com.tngtech.archunit.core.domain.JavaMethod;
import com.tngtech.archunit.core.importer.ClassFileImporter;
import com.tngtech.archunit.lang.ArchCondition;
import com.tngtech.archunit.lang.ConditionEvents;
import com.tngtech.archunit.lang.SimpleConditionEvent;
import org.junit.jupiter.api.Test;

import static com.tngtech.archunit.lang.syntax.ArchRuleDefinition.classes;

class AopArchitectureTest {

    private final JavaClasses classes = new ClassFileImporter().importPackages("${package}");

    @Test
    void preventAopSelfInvocation() {
        classes()
            .should(callAopAnnotatedMethodOnSelf())
            .because("JDK Dynamic Proxies only intercept external calls. " +
                     "Self-invocation bypasses AOP annotations like @Transactional. " +
                     "Please extract the method to another @Component.")
            .allowEmptyShould(true)
            .check(classes);
    }

    @Test
    void preventAnonymousClasses() {
        classes()
            .should(beAnonymousClass())
            .because("Anonymous classes are not allowed in Summer applications. " +
                     "Use lambdas or explicit top-level / static member types instead.")
            .allowEmptyShould(true)
            .check(classes);
    }

    @Test
    void preventLocalClasses() {
        classes()
            .should(beLocalClass())
            .because("Local classes and local records inside methods are not allowed in Summer applications. " +
                     "Declare explicit static member records or top-level classes instead.")
            .allowEmptyShould(true)
            .check(classes);
    }

    @Test
    void preventJdbcTemplateInThreadHoppingMethods() {
        classes()
            .should(callJdbcTemplateFromThreadSpawningMethod())
            .because("The transaction context is bound via ScopedValue to the calling virtual " +
                     "thread and deliberately does NOT cross threads. Work dispatched to another " +
                     "thread runs OUTSIDE the transaction: it gets its own auto-commit connection, " +
                     "and its writes survive a rollback of the outer boundary. Keep transactional " +
                     "data access on the one thread, or move the concurrent work outside the " +
                     "@Transactional boundary.")
            .allowEmptyShould(true)
            .check(classes);
    }

    private static ArchCondition<JavaClass> beAnonymousClass() {
        return new ArchCondition<>("be an anonymous class") {
            @Override
            public void check(JavaClass clazz, ConditionEvents events) {
                if (isRuleInfrastructure(clazz)) {
                    return;
                }
                if (clazz.isAnonymousClass()) {
                    String message = String.format("Class %s is an anonymous class", clazz.getName());
                    events.add(SimpleConditionEvent.violated(clazz, message));
                }
            }
        };
    }

    private static ArchCondition<JavaClass> beLocalClass() {
        return new ArchCondition<>("be a local class") {
            @Override
            public void check(JavaClass clazz, ConditionEvents events) {
                if (isRuleInfrastructure(clazz)) {
                    return;
                }
                if (clazz.isLocalClass()) {
                    String message = String.format("Class %s is a local class or local record", clazz.getName());
                    events.add(SimpleConditionEvent.violated(clazz, message));
                }
            }
        };
    }

    private static ArchCondition<JavaClass> callAopAnnotatedMethodOnSelf() {
        return new ArchCondition<>("call AOP-annotated methods on themselves") {
            @Override
            public void check(JavaClass clazz, ConditionEvents events) {
                if (isRuleInfrastructure(clazz)) {
                    return;
                }
                for (JavaCall<?> call : clazz.getMethodCallsFromSelf()) {
                    if (call.getTargetOwner().equals(clazz) &&
                        call.getTarget().isAnnotatedWith(Transactional.class)) {

                        String message = String.format(
                            "Method %s calls AOP-annotated method %s in the same class",
                            call.getOrigin().getName(), call.getTarget().getName());
                        events.add(SimpleConditionEvent.violated(call, message));
                    }
                }
            }
        };
    }

    private static ArchCondition<JavaClass> callJdbcTemplateFromThreadSpawningMethod() {
        return new ArchCondition<>(
                "call JdbcTemplate from a method that also dispatches to another thread") {
            @Override
            public void check(JavaClass clazz, ConditionEvents events) {
                if (isRuleInfrastructure(clazz)) {
                    return;
                }
                boolean classLevelTx = clazz.isAnnotatedWith(Transactional.class);
                for (JavaMethod method : clazz.getMethods()) {
                    // The write escapes the boundary only when the DISPATCHING method itself is
                    // inside one: a dispatch from a non-transactional method is just normal
                    // independent work.
                    boolean methodLevelTx = method.isAnnotatedWith(Transactional.class);
                    if (!classLevelTx && !methodLevelTx) {
                        continue;
                    }
                    if (!dispatchesToAnotherThread(method)) {
                        continue;
                    }
                    if (touchesJdbc(effectiveCalls(clazz, method))) {
                        events.add(SimpleConditionEvent.violated(
                            method,
                            String.format(
                                "Method %s.%s() dispatches work to another thread whose body"
                                    + " calls JdbcTemplate. Dispatched writes run OUTSIDE the"
                                    + " transaction (ScopedValue confinement is by design): they"
                                    + " get their own auto-commit connection and survive a rollback"
                                    + " of the outer boundary. Keep the transactional data access on"
                                    + " the one thread, or move the concurrent work outside the"
                                    + " @Transactional boundary.",
                                clazz.getName(), method.getName())));
                    }
                }
            }
        };
    }

    /**
     * javac moves a lambda body into a SYNTHETIC method ({@code lambda$<enclosing>$<n>}), so the
     * calls inside a dispatched lambda are NOT in the dispatching method's own call set — a
     * plain per-method check misses exactly the primary violating shape. The effective call set
     * of a method is therefore its own calls plus those of its derived lambdas. Known gap:
     * lambdas nested inside lambdas are named {@code lambda$null$…} and cannot be attributed by
     * name; anonymous classes are a separate class (and banned by another rule anyway).
     */
    private static boolean touchesJdbc(java.util.Set<JavaCall<?>> calls) {
        for (JavaCall<?> call : calls) {
            if (targetsJdbcTemplate(call)) {
                return true;
            }
        }
        return false;
    }

    private static java.util.Set<JavaCall<?>> effectiveCalls(JavaClass clazz, JavaMethod method) {
        java.util.Set<JavaCall<?>> calls = new java.util.HashSet<>(method.getMethodCallsFromSelf());
        String lambdaPrefix = "lambda$" + method.getName() + "$";
        for (JavaMethod candidate : clazz.getMethods()) {
            if (candidate.getName().startsWith(lambdaPrefix)) {
                calls.addAll(candidate.getMethodCallsFromSelf());
            }
        }
        return calls;
    }

    /** The rules must never judge their own test class (or its anonymous conditions). */
    private static boolean isRuleInfrastructure(JavaClass clazz) {
        return clazz.getName().equals("${package}.AopArchitectureTest")
                || clazz.getName().startsWith("${package}.AopArchitectureTest$");
    }

    private static boolean dispatchesToAnotherThread(JavaMethod method) {
        for (JavaCall<?> call : method.getMethodCallsFromSelf()) {
            if (dispatchesToAnotherThread(call)) {
                return true;
            }
        }
        return false;
    }

    /** Matched by NAME so the rule compiles in projects that do not use the data-jdbc starter. */
    private static boolean targetsJdbcTemplate(JavaCall<?> call) {
        return call.getTargetOwner().getName().startsWith("com.github.dropguard.summer.data.jdbc.");
    }

    private static boolean dispatchesToAnotherThread(JavaCall<?> call) {
        String owner = call.getTargetOwner().getName();
        String target = call.getTarget().getName();
        if (owner.equals("java.lang.Thread")
                && (target.equals("startVirtualThread") || target.equals("start"))) {
            return true;
        }
        if ((owner.equals("java.util.concurrent.ExecutorService")
                    || owner.equals("java.util.concurrent.Executor"))
                && (target.equals("submit")
                    || target.equals("execute")
                    || target.equals("invokeAll")
                    || target.equals("invokeAny"))) {
            return true;
        }
        return owner.equals("java.util.concurrent.CompletableFuture")
                && (target.equals("supplyAsync") || target.equals("runAsync"));
    }
}
