package state

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"

	"github.com/pulumi/pulumi-terraform/sdk/v3/go/state/internal"
)

// SwiftStateArgs specifies the configuration options for a Terraform Remote State
// stored in the swift backend.
type SwiftStateArgs struct {
	// AuthUrl is the Identity authentication URL. Sourced from `OS_AUTH_URL` in the environment, if unset.
	AuthUrl pulumi.StringInput

	// Container is the name of the container in which the Terraform state file is stored.
	Container pulumi.StringInput

	// UserName is the username with which to log in. Sourced from `OS_USERNAME` in the environment, if unset
	UserName pulumi.StringPtrInput

	// UserID is the user ID with which to log in. Sourced from `OS_USER_ID` in the environment, if unset
	UserID pulumi.StringPtrInput

	// Password is the password with which to log in. Sourced from `OS_PASSWORD` in the environment, if unset
	Password pulumi.StringPtrInput

	// Token is the access token with which to log in in stead of a username and password. Sourced from
	// `OS_AUTH_TOKEN` in the environment, if unset.
	Token pulumi.StringPtrInput

	// RegionName is the region in which the state file is stored. Sourced from `OS_REGION_NAME`, if unset
	RegionName pulumi.StringInput

	// TenantID is the ID of the tenant (for identity v2) or project (identity v3) which which to log in.
	// Sourced from `OS_TENANT_ID` or `OS_PROJECT_ID` in the environment, if unset.
	TenantID pulumi.StringPtrInput

	// TenantName is the name of the tenant (for identity v2) or project (identity v3) which which to log in.
	// Sourced from `OS_TENANT_NAME` or `OS_PROJECT_NAME` in the environment, if unset.
	TenantName pulumi.StringPtrInput

	// DomainName is the name of the domain to scope the log in to (identity v3). Sourced from
	// `OS_USER_DOMAIN_NAME`, `OS_PROJECT_DOMAIN_NAME` or `OS_DOMAIN_NAME` in the environment, if unset.
	DomainName pulumi.StringPtrInput

	// DomainID is the ID of the domain to scope the log in to (identity v3). Sourced from `OS_USER_DOMAIN_ID`,
	// `OS_PROJECT_DOMAIN_ID` or `OS_DOMAIN_ID` in the environment, if unset.
	DomainID pulumi.StringPtrInput

	// Insecure is whether to disable verification of the server TLS certificate. Sourced from
	// `OS_INSECURE` in the environment, if unset.
	Insecure pulumi.BoolPtrInput

	// CACertFile is the path to a CA root certificate for verifying the server TLS certificate. Sourced from
	// `OS_CACERT` in the environment, if unset.
	CACertFile pulumi.StringPtrInput

	// Cert is the path to a client certificate for TLS client authentication. Sourced from `OS_CERT`
	// in the environment, if unset.
	Cert pulumi.StringPtrInput

	// Key is the path to the private key corresponding to the client certificate for TLS client
	// authentication. Sourced from `OS_KEY` in the environment, if unset.
	Key pulumi.StringPtrInput
}

func (l *SwiftStateArgs) toInternalArgs() pulumi.Input {
	return internal.SwiftStateReferenceArgs{
		BackendType: pulumi.String("swift"),
		AuthUrl:     l.AuthUrl,
		Container:   l.Container,
		UserName:    l.UserName,
		UserID:      l.UserID,
		Password:    l.Password,
		Token:       l.Token,
		RegionName:  l.RegionName,
		TenantID:    l.TenantID,
		TenantName:  l.TenantName,
		DomainID:    l.DomainID,
		DomainName:  l.DomainName,
		Insecure:    l.Insecure,
		CACertFile:  l.CACertFile,
		Cert:        l.Cert,
		Key:         l.Key,
	}
}

func (l *SwiftStateArgs) validateArgs() error {
	if l.AuthUrl == pulumi.String("") || l.Container == pulumi.String("") ||
		l.RegionName == pulumi.String("") {
		return errors.New("`AuthUrl`, `Container` and `RegionName` are required parameters")
	}
	return nil
}
