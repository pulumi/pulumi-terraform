import pulumi
import pulumi_terraform as terraform

outputs = terraform.state.get_local_reference(path="./terraform.0-12-24.tfstate")

pulumi.export("state", outputs.outputs)
