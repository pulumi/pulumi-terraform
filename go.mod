module github.com/pulumi/pulumi-terraform

go 1.12

require (
	github.com/Azure/azure-amqp-common-go v1.1.4 // indirect
	github.com/blang/semver v3.5.1+incompatible
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/gedex/inflector v0.0.0-20170307190818-16278e9db813
	github.com/gliderlabs/ssh v0.1.3 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.2
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/hcl2 v0.0.0-20190821123243-0c888d1241f6
	github.com/hashicorp/terraform v0.12.7
	github.com/mitchellh/copystructure v1.0.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.1.0 // indirect
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/pulumi/pulumi v1.6.0
	github.com/reconquest/loreley v0.0.0-20160708080500-2ab6b7470a54 // indirect
	github.com/sabhiram/go-gitignore v0.0.0-20180611051255-d3107576ba94 // indirect
	github.com/spf13/cobra v0.0.3
	github.com/stretchr/testify v1.3.1-0.20190311161405-34c6fa2dc709
	github.com/tidwall/pretty v0.0.0-20190325153808-1166b9ac2b65 // indirect
	github.com/uber-go/atomic v1.3.2 // indirect
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	github.com/zclconf/go-cty v1.0.1-0.20190708163926-19588f92a98f
	go.mongodb.org/mongo-driver v1.0.1 // indirect
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	google.golang.org/grpc v1.22.0
)

replace (
	git.apache.org/thrift.git => github.com/apache/thrift v0.12.0
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v12.4.3+incompatible
	github.com/Nvveen/Gotty => github.com/ijc25/Gotty v0.0.0-20170406111628-a8b993ba6abd
	github.com/hashicorp/vault => github.com/hashicorp/vault v1.2.0
)
