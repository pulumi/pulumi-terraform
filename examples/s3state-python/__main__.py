import pulumi
from pulumi_terraform import state

config = pulumi.Config()
bucket_name = config.require("bucketName")
key = config.require("key")
region = config.require("region")

s = state.RemoteStateReference("s3state", "s3", state.S3BackendArgs(
    bucket = bucket_name,
    key = key + "/terraform.tfstate",
    region = region,
))

pulumi.export('vpcId', s.get_output("vpc_id"))
pulumi.export('publicSubnetIds', s.get_output("public_subnet_ids"))
pulumi.export('bucketArn', s.get_output("bucket_arn"))
