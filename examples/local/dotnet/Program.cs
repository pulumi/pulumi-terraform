using System.Collections.Generic;
using Pulumi;
using Pulumi.Terraform.State;

return await Deployment.RunAsync(() =>
{
    var outputs = GetLocalReference.Invoke(new GetLocalReferenceInvokeArgs
    {
        Path = "./terraform.0-12-24.tfstate",
    }).Apply(x => x.Outputs);

    // Export outputs here
    return new Dictionary<string, object?>
    {
        ["state"] = outputs,
        ["bucketArn"] = outputs.Apply(x => x["bucket_arn"]),
        ["firstSubnetId"] = outputs.Apply(x => (x["public_subnet_ids"] as IList<object>)[0]),
    };
});
