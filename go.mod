module github.com/pulumi/pulumi-terraform

go 1.12

require (
	cloud.google.com/go/logging v1.0.0 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/hashicorp/terraform v0.12.9
	github.com/hashicorp/vault v1.2.3 // indirect
	github.com/pkg/errors v0.8.1
	github.com/pulumi/pulumi v1.2.0
	github.com/stretchr/testify v1.4.0
	github.com/zclconf/go-cty v1.1.0
	golang.org/x/net v0.0.0-20190926025831-c00fd9afed17
	google.golang.org/grpc v1.24.0
)

replace (
	github.com/Azure/go-autorest/tracing => github.com/Azure/go-autorest v13.0.1+incompatible
	github.com/hashicorp/consul => github.com/hashicorp/consul v1.6.1
	github.com/hashicorp/vault => github.com/hashicorp/vault v1.2.3
	github.com/ugorji/go/codec => github.com/ugorji/go v1.1.4
)
