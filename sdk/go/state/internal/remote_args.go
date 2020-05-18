package internal

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type RemoteBackendStateReferenceArgs struct {
	BackendType pulumi.StringPtrInput

	Organization pulumi.StringInput
	Hostname     pulumi.StringPtrInput
	Token        pulumi.StringPtrInput
	Workspaces   WorkspaceArgs
}

type WorkspaceArgs struct {
	Name   pulumi.StringPtrInput
	Prefix pulumi.StringPtrInput
}

type remoteBackendStateReferenceArgs struct {
	BackendType *string `pulumi:"backendType"`

	Organization string        `pulumi:"organization"`
	Hostname     *string       `pulumi:"hostname"`
	Token        *string       `pulumi:"token"`
	Workspaces   workspaceArgs `pulumi:"workspaces"`
}

type workspaceArgs struct {
	Name   *string `pulumi:"name"`
	Prefix *string `pulumi:"prefix"`
}

func (RemoteBackendStateReferenceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*remoteBackendStateReferenceArgs)(nil)).Elem()
}

func (WorkspaceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*workspaceArgs)(nil)).Elem()
}
