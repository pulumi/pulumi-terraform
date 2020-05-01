module github.com/pulumi/pulumi-terraform/provider

go 1.14

require (
	github.com/golang/protobuf v1.4.0
	github.com/hashicorp/terraform v0.12.24
	github.com/hashicorp/terraform-svchost v0.0.0-20191011084731-65d371908596
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi/pkg/v2 v2.1.0
	github.com/pulumi/pulumi/sdk/v2 v2.1.0
	github.com/stretchr/testify v1.5.1
	github.com/zclconf/go-cty v1.4.0
	golang.org/x/net v0.0.0-20200425230154-ff2c4b7c35a0
	google.golang.org/grpc v1.29.1
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v12.4.3+incompatible
	github.com/ugorji/go => github.com/ugorji/go v0.0.0-20180813092308-00b869d2f4a5
)
