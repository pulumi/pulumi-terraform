import pulumi
import pulumi_terraform as terraform

outputs = terraform.state.get_local_reference(path="./terraform.0-12-24.tfstate").outputs

pulumi.export("state", outputs)
pulumi.export("bucketArn", outputs["bucket_arn"])
pulumi.export("firstSubnetId", outputs["public_subnet_ids"][0])
