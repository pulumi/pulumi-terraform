import * as tf from "@pulumi/terraform";
import * as pulumi from "@pulumi/pulumi";

const config = new pulumi.Config()
const bucketName = config.require("bucketName")
const prefix = config.require("prefix")
const region = config.require("region")

const remotestate = new tf.state.RemoteStateReference("ossstate", {
   backendType: "oss",
   bucket: bucketName,
   prefix: prefix,
   key: "terraform.tfstate",
   region: region,
});

export const vpcId= remotestate.getOutput("vpc_id");
export const publicSubnetIds = remotestate.getOutput("public_subnet_ids");
export const bucketArn = remotestate.getOutput("bucket_arn");
