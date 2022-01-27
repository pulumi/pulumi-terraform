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
	"github.com/pulumi/pulumi/sdk/v3/proto/go"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ResourceTypeRemoteStateReference = "terraform:state:RemoteStateReference"
	ResourceTypeModule               = "terraform:tf:Module"
)

type Provider struct {
	host    *provider.HostClient
	version string
}

func NewProvider(ctx context.Context, host *provider.HostClient, version string) *Provider {
	log.SetOutput(NewTerraformLogRedirector(ctx, host))
	shim.InitTfBackend()

	return &Provider{
		host:    host,
		version: version,
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
	// TODO: anything to configure? step 1
	return &pulumirpc.ConfigureResponse{}, nil
}

func (*Provider) Invoke(context.Context, *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Invoke is not yet implemented")
}

func (*Provider) Check(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	// TODO: validate module inputs, etc. step 2
	// TODO: any way we can make the TF internal resources look like children somehow?
	return &pulumirpc.CheckResponse{
		Inputs: req.News,
	}, nil
}

func (*Provider) Diff(context.Context, *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	// TODO: diff the previous/new state, etc. step 3
	return &pulumirpc.DiffResponse{}, nil
}

func (p *Provider) Create(ctx context.Context, req *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	switch t := resource.URN(req.Urn).Type(); t {
	case ResourceTypeModule:
		return p.createModuleResource(ctx, req)
	default:
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unknown resource type: %q", t))
	}
}

func (p *Provider) Construct(context.Context, *pulumirpc.ConstructRequest) (*pulumirpc.ConstructResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Construct is not yet implemented")
}

func (*Provider) Read(ctx context.Context, req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	switch t := resource.URN(req.Urn).Type(); t {
	case ResourceTypeRemoteStateReference:
		return shim.RemoteStateReferenceRead(ctx, req)
	default:
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unknown resource type: %q", t))
	}
}

func (p *Provider) Update(ctx context.Context, req *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	switch t := resource.URN(req.Urn).Type(); t {
	case ResourceTypeModule:
		return p.updateModuleResource(ctx, req)
	default:
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unknown resource type: %q", t))
	}
}

func (p *Provider) Delete(ctx context.Context, req *pulumirpc.DeleteRequest) (*empty.Empty, error) {
	switch t := resource.URN(req.Urn).Type(); t {
	case ResourceTypeModule:
		return p.destroyModuleResource(ctx, req)
	default:
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unknown resource type: %q", t))
	}
}

func (*Provider) Cancel(context.Context, *empty.Empty) (*empty.Empty, error) {
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
