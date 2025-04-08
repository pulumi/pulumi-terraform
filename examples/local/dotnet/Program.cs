using System.Collections.Generic;
using Pulumi;
using Pulumi.Terraform;

return await Deployment.RunAsync(() =>
{
    var outputs = LocalStateReference.Invoke(new LocalStateReferenceInvokeArgs
    {
        Path = "./terraform.0-12-24.tfstate",
    });

    // Export outputs here
    return new Dictionary<string, object?>
    {
        ["state"] = outputs.Apply(x => x.Outputs),
    };
});
