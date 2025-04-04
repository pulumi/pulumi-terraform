using System.IO;
using System.Collections.Generic;
using System.Threading.Tasks;

using Pulumi;
using Pulumi.Terraform.State;

class Program
{
    static Task<int> Main()
    {
        return Deployment.RunAsync(() => {

            var config = new Config();
            var bucketName = config.Require("bucketName");
            var prefix = config.Require("prefix");
            var region = config.Require("region");
            var remoteState = new RemoteStateReference("ossstate", new OssRemoteStateReferenceArgs
            {
                Bucket = bucketName,
                Prefix = prefix,
                Key = "terraform.tfstate",
                Region = region,
            });

            return new Dictionary<string, object?>
            {
                { "vpcId", remoteState.GetOutput("vpc_id") },
                { "publicSubnetIds", remoteState.GetOutput("public_subnet_ids") },
                { "bucketArn", remoteState.GetOutput("bucket_arn") },
            };
        });
    }
}
