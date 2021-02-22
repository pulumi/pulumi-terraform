import pulumi
from pulumi_terraform import state

config = pulumi.Config()
bucket_name = config.require("bucketName")
region = config.require("region")
prefix = config.require("prefix")

s = state.RemoteStateReference("ossstate", "oss", state.OssBackendArgs(
    bucket = bucket_name,
    prefix = prefix,
    key = "terraform.tfstate",
    region = region,
))

pulumi.export('vpcId', s.get_output("vpc_id"))
pulumi.export('publicSubnetIds', s.get_output("public_subnet_ids"))
pulumi.export('bucketArn', s.get_output("bucket_arn"))
