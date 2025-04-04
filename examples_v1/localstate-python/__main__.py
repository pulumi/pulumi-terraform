import os

import pulumi
from pulumi_terraform import state

config = pulumi.Config()
statefile = config.require("statefile")

script_dir = os.path.dirname(os.path.realpath(__file__))

s = state.RemoteStateReference("localstate", "local", state.LocalBackendArgs(path = os.path.join(script_dir, statefile)))

pulumi.export('vpcId', s.get_output("vpc_id"))
pulumi.export('publicSubnetIds', s.get_output("public_subnet_ids"))
pulumi.export('bucketArn', s.get_output("bucket_arn"))
