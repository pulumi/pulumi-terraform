import { state as tf_state } from "@pulumi/terraform";

let outputs = tf_state.getLocalReferenceOutput({
  path: "./terraform.0-12-24.tfstate",
});

export const state = outputs.outputs;
