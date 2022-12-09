// Copyright 2016-2019, Pulumi Corporation.
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

package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/terraform/shim"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Provider struct {
	pulumirpc.UnimplementedResourceProviderServer

	version string
}

func NewProvider(ctx context.Context, host *provider.HostClient, version string) *Provider {
	log.SetOutput(NewTerraformLogRedirector(ctx, host))
	shim.InitTfBackend()

	p := &Provider{
		version: version,
	}
	return p
}

func validateAndExtractResourceType(urnValue string) (string, error) {
	urn := resource.URN(urnValue)
	resourceType := urn.Type()

	const resourceTypeRemoteStateReference = "terraform:state:RemoteStateReference"

	switch resourceType {
	case resourceTypeRemoteStateReference:
		return resourceTypeRemoteStateReference, nil
	default:
		return "", status.Error(codes.InvalidArgument, fmt.Sprintf("unknown resource type: %q", resourceType))
	}
}

func (*Provider) GetSchema(context.Context, *pulumirpc.GetSchemaRequest) (*pulumirpc.GetSchemaResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetSchema is not yet implemented")
}

func (*Provider) CheckConfig(context.Context, *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CheckConfig is not yet implemented")
}

func (*Provider) StreamInvoke(*pulumirpc.InvokeRequest, pulumirpc.ResourceProvider_StreamInvokeServer) error {
	return nil
}

func (*Provider) DiffConfig(context.Context, *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	return nil, status.Error(codes.Unimplemented, "DiffConfig is not yet implemented")
}

func (p *Provider) Configure(context.Context, *pulumirpc.ConfigureRequest) (*pulumirpc.ConfigureResponse, error) {
	return &pulumirpc.ConfigureResponse{}, nil
}

func (*Provider) Invoke(context.Context, *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Invoke is not yet implemented")
}

func (*Provider) Check(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Check is not yet implemented")
}

func (*Provider) Diff(context.Context, *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Diff is not yet implemented")
}

func (*Provider) Create(context.Context, *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Create is not yet implemented")
}

func (p *Provider) Construct(context.Context, *pulumirpc.ConstructRequest) (*pulumirpc.ConstructResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Construct is not yet implemented")
}

func (*Provider) Read(ctx context.Context, req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	if _, err := validateAndExtractResourceType(req.Urn); err != nil {
		return nil, err
	}

	return shim.RemoteStateReferenceRead(ctx, req)
}

func (*Provider) Update(context.Context, *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Update is not yet implemented")
}

func (*Provider) Delete(context.Context, *pulumirpc.DeleteRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "Delete is not yet implemented")
}

func (*Provider) Cancel(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (*Provider) Attach(context.Context, *pulumirpc.PluginAttach) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// Call dynamically executes a method in the provider associated with a component resource.
func (p *Provider) Call(ctx context.Context, req *pulumirpc.CallRequest) (*pulumirpc.CallResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Call is not yet implemented")
}

func (p *Provider) GetPluginInfo(context.Context, *empty.Empty) (*pulumirpc.PluginInfo, error) {
	return &pulumirpc.PluginInfo{
		Version: p.version,
	}, nil
}
