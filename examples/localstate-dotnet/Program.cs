﻿using System.IO;
using System.Collections.Generic;
using System.Threading.Tasks;

using Pulumi;
using Pulumi.Terraform.State;

class Program
{
    static Task<int> Main()
    {
        return Deployment.RunAsync(() => {

            var remoteState = new RemoteStateReference("localstate", new LocalBackendRemoteStateReferenceArgs
            {
                Path = Path.GetFullPath("terraform.tfstate"),
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
