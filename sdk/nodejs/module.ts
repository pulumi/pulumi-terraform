// Copyright 2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import * as pulumi from "@pulumi/pulumi";

export interface ModuleArgs {
    /**
     * This module's name in the Terraform Registry.
     */
    source: string;
    /**
     * This module's version in The Terraform Registry.
     */
    version: string;
    /**
     * A bag of provider configurations.
     * HACKHACK: ideally we wouldn't need this and could do it automatically!
     */
    providers: {[key: string]: any};
    /**
     * A weakly typed bag of inputs to pass to the module construction.
     */
    inputs: {[key: string]: any};
}

/**
 * Enables provisioning a Terraform module directly as an opaque Pulumi resource.
 */
export class Module extends pulumi.CustomResource {
    /**
     * The root outputs of the Terraform Module.
     */
    public readonly outputs: pulumi.Output<{ [name: string]: any }>;
    /**
     * The base64-encoded TF state (blob) for this Terraform Module. This is largely
     * an implementation detail which can be ignored.
     */
    public readonly serializedState: pulumi.Output<string>;

    /**
     * Create a Module resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the remote state reference.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args: ModuleArgs, opts?: pulumi.CustomResourceOptions) {
        super("terraform:tf:Module", name, {
            outputs: undefined,
            serializedState: undefined,
            ...args
        }, opts);
    }

    /**
     * Fetches the value of a root output from the Terraform Module.
     *
     * @param name The name of the output to fetch. The name is formatted exactly as per
     * the "output" block in the Terraform Module's configuration.
     */
    public getOutput<T = any>(name: pulumi.Input<string>): pulumi.Output<T> {
        return pulumi.all([pulumi.output(name), this.outputs]).apply(([n, os]) => os[n] as T);
    }
}

