import * as terraform from "@pulumi/terraform";

let outputs = terraform.localStateReferenceOutput({
  path: "./terraform.0-12-24.tfstate",
});

export const state = outputs.outputs;
