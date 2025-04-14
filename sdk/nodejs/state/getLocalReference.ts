// *** WARNING: this file was generated by pulumi-language-nodejs. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "../utilities";

/**
 * Access state from the local filesystem.
 */
export function getLocalReference(args?: GetLocalReferenceArgs, opts?: pulumi.InvokeOptions): Promise<GetLocalReferenceResult> {
    args = args || {};
    opts = pulumi.mergeOptions(utilities.resourceOptsDefaults(), opts || {});
    return pulumi.runtime.invoke("terraform:state:getLocalReference", {
        "path": args.path,
        "workspaceDir": args.workspaceDir,
    }, opts);
}

export interface GetLocalReferenceArgs {
    /**
     * The path to the tfstate file. This defaults to "terraform.tfstate" relative to the root module by default.
     */
    path?: string;
    /**
     * The path to non-default workspaces.
     */
    workspaceDir?: string;
}

/**
 * The result of fetching from a Terraform state store.
 */
export interface GetLocalReferenceResult {
    /**
     * The outputs displayed from Terraform state.
     */
    readonly outputs: {[key: string]: any};
}
/**
 * Access state from the local filesystem.
 */
export function getLocalReferenceOutput(args?: GetLocalReferenceOutputArgs, opts?: pulumi.InvokeOutputOptions): pulumi.Output<GetLocalReferenceResult> {
    args = args || {};
    opts = pulumi.mergeOptions(utilities.resourceOptsDefaults(), opts || {});
    return pulumi.runtime.invokeOutput("terraform:state:getLocalReference", {
        "path": args.path,
        "workspaceDir": args.workspaceDir,
    }, opts);
}

export interface GetLocalReferenceOutputArgs {
    /**
     * The path to the tfstate file. This defaults to "terraform.tfstate" relative to the root module by default.
     */
    path?: pulumi.Input<string>;
    /**
     * The path to non-default workspaces.
     */
    workspaceDir?: pulumi.Input<string>;
}
