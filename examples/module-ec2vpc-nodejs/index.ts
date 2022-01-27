import * as aws from "@pulumi/aws/config/vars";
import * as pulumi from "@pulumi/pulumi";
import * as tf from "@pulumi/terraform";

// Input variable definitions.
const cfg = new pulumi.Config();
const vpcName = cfg.get("vpc_name") ?? "example-vpc";
const vpcCidr = cfg.get("vpc_cidr") ?? "10.0.0.0/16";
const vpcAzs = cfg.getObject<string[]>("vpc_azs") ?? ["us-west-2a", "us-west-2b", "us-west-2c"];
const vpcPrivateSubnets = cfg.getObject<string[]>("vpc_private_subnets") ?? ["10.0.1.0/24", "10.0.2.0/24"];
const vpcPublicSubnets = cfg.getObject<string[]>("vpc_public_subnets") ?? ["10.0.101.0/24", "10.0.102.0/24"];
const vpcEnableNatGateway = cfg.getBoolean("vpc_enable_nat_gateway") ?? true;
const vpcTags = cfg.getObject<Record<string, string>>("vpc_tags") ?? { "Pulumi": "true", "Environment": "dev" };

// Provision the modules: first a VPC and then EC2 instances within that VPC.
const vpc = new tf.Module("vpc", {
    source: "terraform-aws-modules/vpc/aws",
    version: "2.21.0",
    providers: { aws },
    inputs: {
        name: vpcName,
        cidr: vpcCidr,

        azs: vpcAzs,
        private_subnets: vpcPrivateSubnets,
        public_subnets: vpcPublicSubnets,

        enable_nat_gateway: vpcEnableNatGateway,

        tags: vpcTags,
    },
});
const ec2Instances = new tf.Module("ec2_instances", {
    source: "terraform-aws-modules/ec2-instance/aws",
    version: "2.12.0",
    providers: { aws },
    inputs: {
        name: "my-ec2-cluster",
        instance_count: 2,

        ami: "ami-0c5204531f799e0c6",
        instance_type: "t2.micro",
        vpc_security_group_ids: [ vpc.getOutput<string>("default_security_group_id") ],
        subnet_id: vpc.getOutput<string[]>("public_subnets")[0],

        tags: {
            "Pulumi": "true",
            "Environment": "dev",
        },
    },
});

// Output variable definitions.
module.exports = {
    // IDs of the VPC's public subnets.
    vpc_public_subnets: vpc.outputs["public_subnets"],
    // Public IP addresses of EC2 instances.
    ec2_instance_public_ipds: ec2Instances.outputs["public_ip"],
};
