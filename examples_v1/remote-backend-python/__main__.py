import pulumi
from pulumi_terraform import state

config = pulumi.Config()
token = config.require_secret("tfeToken")
workspace_name = config.require("workspaceName")
organization = config.require("organization")

s = state.RemoteStateReference("remote-backend", "remote", state.RemoteBackendArgs(
    token=token,
    organization=organization,
    workspace_name=workspace_name))

pulumi.export('password', s.get_output("password"))
