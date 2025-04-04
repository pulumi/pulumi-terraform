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
            var tfeToken = config.RequireSecret("tfeToken");
            var organization = config.Require("organization");
            var workspaceName = config.Require("workspaceName");
            var remoteState = new RemoteStateReference("remote-backend-state", new RemoteBackendRemoteStateReferenceArgs()
            {
                Token = tfeToken,
                Organization = organization,
                Workspaces = new RemoteBackendWorkspaceConfig()
                {
                    Name = workspaceName,
                }
            });

            return new Dictionary<string, object?>
            {
                { "password", remoteState.GetOutput("password") },
            };
        });
    }
}
