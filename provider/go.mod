module github.com/pulumi/pulumi-terraform/provider/v3

go 1.14

require (
	github.com/golang/protobuf v1.4.2
	github.com/hashicorp/terraform v0.13.5
	github.com/hashicorp/terraform-svchost v0.0.0-20191011084731-65d371908596
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi/pkg/v2 v2.12.0
	github.com/pulumi/pulumi/sdk/v2 v2.12.0
	github.com/stretchr/testify v1.6.1
	github.com/zclconf/go-cty v1.5.1
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9
	google.golang.org/grpc v1.29.1
	k8s.io/client-go v11.0.0+incompatible // indirect
)

replace (
	github.com/ugorji/go => github.com/ugorji/go v0.0.0-20180813092308-00b869d2f4a5
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
)
