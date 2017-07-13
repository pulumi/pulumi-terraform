// *** WARNING: this file was generated by the Lumi Terraform Bridge (TFGEN) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as lumi from "@lumi/lumi";

export class AvailabilitySet extends lumi.NamedResource implements AvailabilitySetArgs {
    public readonly location: string;
    public readonly managed?: boolean;
    public readonly _name: string;
    public readonly platformFaultDomainCount?: number;
    public readonly platformUpdateDomainCount?: number;
    public readonly resourceGroupName: string;
    public readonly tags?: {[key: string]: any};

    constructor(name: string, args: AvailabilitySetArgs) {
        super(name);
        this.location = args.location;
        this.managed = args.managed;
        this._name = args._name;
        this.platformFaultDomainCount = args.platformFaultDomainCount;
        this.platformUpdateDomainCount = args.platformUpdateDomainCount;
        this.resourceGroupName = args.resourceGroupName;
        this.tags = args.tags;
    }
}

export interface AvailabilitySetArgs {
    readonly location: string;
    readonly managed?: boolean;
    readonly _name: string;
    readonly platformFaultDomainCount?: number;
    readonly platformUpdateDomainCount?: number;
    readonly resourceGroupName: string;
    readonly tags?: {[key: string]: any};
}
