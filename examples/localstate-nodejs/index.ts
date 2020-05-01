import * as tf from "@pulumi/terraform";
import * as pulumi from "@pulumi/pulumi";
import * as path from "path";

const config = new pulumi.Config()
const statefile = config.require("statefile")

const remotestate = new tf.state.RemoteStateReference("localstate", {
   backendType: "local",
   path: path.join(__dirname, statefile),
});

export const vpcId= remotestate.getOutput("vpc_id");
export const publicSubnetIds = remotestate.getOutput("public_subnet_ids");
export const bucketArn = remotestate.getOutput("bucket_arn");
