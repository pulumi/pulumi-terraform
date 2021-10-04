module github.com/hashicorp/terraform/shim

go 1.15

require (
	github.com/golang/protobuf v1.4.3
	github.com/hashicorp/terraform v1.0.8
	github.com/hashicorp/terraform-svchost v0.0.0-20200729002733-f050f53b9734
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi/sdk/v3 v3.10.3
	github.com/stretchr/testify v1.6.1
	github.com/zclconf/go-cty v1.9.1
)

replace (
	google.golang.org/grpc => google.golang.org/grpc v1.27.1
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
)