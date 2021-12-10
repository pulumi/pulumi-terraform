module github.com/hashicorp/terraform/shim

go 1.15

require (
	github.com/golang/protobuf v1.5.2
	github.com/hashicorp/terraform v1.1.0
	github.com/hashicorp/terraform-svchost v0.0.0-20200729002733-f050f53b9734
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi/sdk/v3 v3.17.0
	github.com/stretchr/testify v1.7.0
	github.com/zclconf/go-cty v1.10.0
	google.golang.org/grpc v1.36.0
)

replace (
	github.com/hashicorp/consul/api => github.com/hashicorp/consul/api v1.9.1
	google.golang.org/grpc => google.golang.org/grpc v1.27.1
	k8s.io/client-go => k8s.io/client-go v0.21.3
)
