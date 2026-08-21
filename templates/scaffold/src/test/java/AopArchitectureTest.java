package ${package};

import com.github.dropguard.summer.tx.Transactional;
import com.tngtech.archunit.core.domain.JavaCall;
import com.tngtech.archunit.core.domain.JavaClass;
import com.tngtech.archunit.lang.ArchCondition;
import com.tngtech.archunit.lang.ConditionEvents;
import com.tngtech.archunit.lang.SimpleConditionEvent;
import org.junit.jupiter.api.Test;

import static com.tngtech.archunit.lang.syntax.ArchRuleDefinition.noClasses;

class AopArchitectureTest {

    @Test
    void preventAopSelfInvocation() {
        noClasses()
            .should(callAopAnnotatedMethodOnSelf())
            .because("JDK Dynamic Proxies only intercept external calls. " +
                     "Self-invocation bypasses AOP annotations like @Transactional. " +
                     "Please extract the method to another @Component.");
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
