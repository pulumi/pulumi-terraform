---
title: Terraform
meta_desc: Consume Terraform state in Pulumi
layout: package
---

The Terraform provider allows you to reference Terraform state from Pulumi programs.


## Example

{{< chooser language "typescript,python,go,csharp,java,yaml" >}}

{{% choosable language typescript %}}

```typescript
import { state as tf_state } from "@pulumi/terraform";

let outputs = tf_state.getLocalReferenceOutput({
  path: "./terraform.0-12-24.tfstate",
});

export const state = outputs.outputs;
```

{{% /choosable %}}

{{% choosable language python %}}

```python
import pulumi
import pulumi_terraform as terraform

outputs = terraform.state.get_local_reference(path="./terraform.0-12-24.tfstate")

pulumi.export("state", outputs.outputs)
```

{{% /choosable %}}

{{% choosable language go %}}

```go
package main

import (
	"github.com/pulumi/pulumi-terraform/sdk/v6/go/terraform/state"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		state := state.GetLocalReferenceOutput(ctx, state.GetLocalReferenceOutputArgs{
			Path: pulumi.String("./terraform.0-12-24.tfstate"),
		})
		ctx.Export("state", state.Outputs())
		return nil
	})
}
```

{{% /choosable %}}

{{% choosable language csharp %}}

```csharp
using System.Collections.Generic;
using Pulumi;
using Pulumi.Terraform.State;

return await Deployment.RunAsync(() =>
{
    var outputs = GetLocalReference.Invoke(new GetLocalReferenceInvokeArgs
    {
        Path = "./terraform.0-12-24.tfstate",
    });

    return new Dictionary<string, object?>
    {
        ["state"] = outputs.Apply(x => x.Outputs),
    };
});
```

{{% /choosable %}}

{{% choosable language java %}}

```java
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
```

{{% /choosable %}}

{{% choosable language yaml %}}

```yaml
name: terraform-local-state-with-yaml
runtime: yaml
outputs:
  state:
    fn::invoke:
      function: terraform:state:getLocalReference
      arguments:
        path: ./terraform.0-12-24.tfstate
      return: outputs
```

{{% /choosable %}}

{{< /chooser >}}
