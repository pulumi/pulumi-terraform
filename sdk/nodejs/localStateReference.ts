// *** WARNING: this file was generated by pulumi-language-nodejs. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "./utilities";

/**
 * Access state from the local filesystem.
 */
export function localStateReference(args?: LocalStateReferenceArgs, opts?: pulumi.InvokeOptions): Promise<LocalStateReferenceResult> {
    args = args || {};
    opts = pulumi.mergeOptions(utilities.resourceOptsDefaults(), opts || {});
    return pulumi.runtime.invoke("terraform:index:localStateReference", {
        "path": args.path,
        "workspaceDir": args.workspaceDir,
    }, opts);
}

export interface LocalStateReferenceArgs {
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
export interface LocalStateReferenceResult {
    /**
     * The outputs displayed from Terraform state.
     */
    readonly outputs: {[key: string]: any};
}
/**
 * Access state from the local filesystem.
 */
export function localStateReferenceOutput(args?: LocalStateReferenceOutputArgs, opts?: pulumi.InvokeOutputOptions): pulumi.Output<LocalStateReferenceResult> {
    args = args || {};
    opts = pulumi.mergeOptions(utilities.resourceOptsDefaults(), opts || {});
    return pulumi.runtime.invokeOutput("terraform:index:localStateReference", {
        "path": args.path,
        "workspaceDir": args.workspaceDir,
    }, opts);
}

export interface LocalStateReferenceOutputArgs {
    /**
     * The path to the tfstate file. This defaults to "terraform.tfstate" relative to the root module by default.
     */
    path?: pulumi.Input<string>;
    /**
     * The path to non-default workspaces.
     */
    workspaceDir?: pulumi.Input<string>;
}
