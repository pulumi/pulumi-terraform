package myproject;

import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.terraform.state.StateFunctions;
import com.pulumi.terraform.state.inputs.GetLocalReferenceArgs;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var output = StateFunctions.getLocalReference(GetLocalReferenceArgs.builder().path("./terraform.0-12-24.tfstate").build());
            ctx.export("state", output.applyValue(x -> x.outputs()));
        });
    }
}
