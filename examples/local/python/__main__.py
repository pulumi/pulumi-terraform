import pulumi
import pulumi_terraform as terraform

outputs = terraform.local_state_reference(path="./terraform.0-12-24.tfstate")

pulumi.export("state", outputs.outputs)
