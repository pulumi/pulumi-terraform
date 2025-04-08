package myproject;

import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.terraform.TerraformFunctions;
import com.pulumi.terraform.inputs.LocalStateReferenceArgs;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var output = TerraformFunctions.localStateReference(LocalStateReferenceArgs.builder().path("./terraform.0-12-24.tfstate").build());
            ctx.export("state", output.applyValue(x -> x.outputs()));
        });
    }
}
