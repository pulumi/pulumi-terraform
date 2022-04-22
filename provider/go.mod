module github.com/pulumi/pulumi-terraform/provider/v5

go 1.16

require (
	github.com/golang/protobuf v1.5.2
	github.com/hashicorp/terraform/shim v0.0.0
	github.com/pulumi/pulumi/pkg/v3 v3.25.0
	github.com/pulumi/pulumi/sdk/v3 v3.25.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f
	google.golang.org/grpc v1.37.0
)

replace (
	github.com/hashicorp/terraform/shim => ./shim
	google.golang.org/grpc => google.golang.org/grpc v1.27.1
	k8s.io/client-go => k8s.io/client-go v0.23.4
)
