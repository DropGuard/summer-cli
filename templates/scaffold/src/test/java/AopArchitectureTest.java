package ${package};

import com.github.dropguard.summer.tx.Transactional;
import com.tngtech.archunit.core.domain.JavaCall;
import com.tngtech.archunit.core.domain.JavaClass;
import com.tngtech.archunit.core.domain.JavaClasses;
import com.tngtech.archunit.core.importer.ClassFileImporter;
import com.tngtech.archunit.lang.ArchCondition;
import com.tngtech.archunit.lang.ConditionEvents;
import com.tngtech.archunit.lang.SimpleConditionEvent;
import org.junit.jupiter.api.Test;

import static com.tngtech.archunit.lang.syntax.ArchRuleDefinition.noClasses;

class AopArchitectureTest {

    private final JavaClasses classes = new ClassFileImporter().importPackages("${package}");

    @Test
    void preventAopSelfInvocation() {
        noClasses()
            .should(callAopAnnotatedMethodOnSelf())
            .because("JDK Dynamic Proxies only intercept external calls. " +
                     "Self-invocation bypasses AOP annotations like @Transactional. " +
                     "Please extract the method to another @Component.")
            .allowEmptyShould(true)
            .check(classes);
    }

    @Test
    void preventAnonymousClasses() {
        noClasses()
            .should(beAnonymousClass())
            .because("Anonymous classes are not allowed in Summer applications. " +
                     "Use lambdas or explicit top-level / static member types instead.")
            .allowEmptyShould(true)
            .check(classes);
    }

    @Test
    void preventLocalClasses() {
        noClasses()
            .should(beLocalClass())
            .because("Local classes and local records inside methods are not allowed in Summer applications. " +
                     "Declare explicit static member records or top-level classes instead.")
            .allowEmptyShould(true)
            .check(classes);
    }

    private static ArchCondition<JavaClass> beAnonymousClass() {
        return new ArchCondition<>("be an anonymous class") {
            @Override
            public void check(JavaClass clazz, ConditionEvents events) {
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
}
