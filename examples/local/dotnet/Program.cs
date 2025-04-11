using System.Collections.Generic;
using Pulumi;
using Pulumi.Terraform.State;

return await Deployment.RunAsync(() =>
{
    var outputs = GetLocalReference.Invoke(new GetLocalReferenceInvokeArgs
    {
        Path = "./terraform.0-12-24.tfstate",
    });

    // Export outputs here
    return new Dictionary<string, object?>
    {
        ["state"] = outputs.Apply(x => x.Outputs),
    };
});
