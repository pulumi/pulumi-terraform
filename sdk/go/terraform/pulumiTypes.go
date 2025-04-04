// Code generated by pulumi-language-go DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package terraform

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi-terraform/sdk/go/terraform/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var _ = internal.GetEnvOrDefault

type Workspace struct {
	Name   *string `pulumi:"name"`
	Prefix *string `pulumi:"prefix"`
}

// WorkspaceInput is an input type that accepts WorkspaceArgs and WorkspaceOutput values.
// You can construct a concrete instance of `WorkspaceInput` via:
//
//	WorkspaceArgs{...}
type WorkspaceInput interface {
	pulumi.Input

	ToWorkspaceOutput() WorkspaceOutput
	ToWorkspaceOutputWithContext(context.Context) WorkspaceOutput
}

type WorkspaceArgs struct {
	Name   pulumi.StringPtrInput `pulumi:"name"`
	Prefix pulumi.StringPtrInput `pulumi:"prefix"`
}

func (WorkspaceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*Workspace)(nil)).Elem()
}

func (i WorkspaceArgs) ToWorkspaceOutput() WorkspaceOutput {
	return i.ToWorkspaceOutputWithContext(context.Background())
}

func (i WorkspaceArgs) ToWorkspaceOutputWithContext(ctx context.Context) WorkspaceOutput {
	return pulumi.ToOutputWithContext(ctx, i).(WorkspaceOutput)
}

type WorkspaceOutput struct{ *pulumi.OutputState }

func (WorkspaceOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*Workspace)(nil)).Elem()
}

func (o WorkspaceOutput) ToWorkspaceOutput() WorkspaceOutput {
	return o
}

func (o WorkspaceOutput) ToWorkspaceOutputWithContext(ctx context.Context) WorkspaceOutput {
	return o
}

func (o WorkspaceOutput) Name() pulumi.StringPtrOutput {
	return o.ApplyT(func(v Workspace) *string { return v.Name }).(pulumi.StringPtrOutput)
}

func (o WorkspaceOutput) Prefix() pulumi.StringPtrOutput {
	return o.ApplyT(func(v Workspace) *string { return v.Prefix }).(pulumi.StringPtrOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*WorkspaceInput)(nil)).Elem(), WorkspaceArgs{})
	pulumi.RegisterOutputType(WorkspaceOutput{})
}
