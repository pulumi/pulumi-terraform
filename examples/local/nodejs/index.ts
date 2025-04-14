import { state as tf_state } from "@pulumi/terraform";

let outputs = tf_state.getLocalReferenceOutput({
  path: "./terraform.0-12-24.tfstate",
}).outputs;

export const state = outputs;
export const bucketArn = outputs["bucket_arn"];
export const firstSubnetId = outputs["public_subnet_ids"].apply(x => x[0]);
