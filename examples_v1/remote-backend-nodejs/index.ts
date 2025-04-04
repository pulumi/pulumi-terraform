import * as tf from "@pulumi/terraform";
import * as pulumi from "@pulumi/pulumi";

const config = new pulumi.Config()
const token = config.requireSecret("tfeToken")
const organization = config.require("organization")
const workspaceName = config.require("workspaceName")

const remotestate = new tf.state.RemoteStateReference("remote-backend-state", {
   backendType: "remote",
   token: token,
   organization: organization,
   workspaces: {
      name: workspaceName
   }
});

export const password = remotestate.getOutput("password");
