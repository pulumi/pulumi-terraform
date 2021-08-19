module github.com/pulumi/pulumi-terraform/provider/v5

go 1.16

require (
	github.com/golang/protobuf v1.5.2
	github.com/hashicorp/terraform/shim v0.0.0
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi/pkg/v3 v3.10.3
	github.com/pulumi/pulumi/sdk/v3 v3.10.3
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20210505214959-0714010a04ed
	google.golang.org/grpc v1.37.0
)

replace (
	github.com/hashicorp/consul => github.com/hashicorp/consul/api v1.8.1
	github.com/hashicorp/terraform/shim => ./shim
	github.com/spf13/cobra => github.com/spf13/cobra v1.1.3
	github.com/spf13/viper => github.com/spf13/viper v1.7.1
	google.golang.org/grpc => google.golang.org/grpc v1.27.1
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
)
