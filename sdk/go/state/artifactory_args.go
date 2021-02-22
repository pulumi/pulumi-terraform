package state

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"

	"github.com/pulumi/pulumi-terraform/sdk/v4/go/state/internal"
)

// ArtifactoryArgs specifies the configuration options for a Terraform Remote State
// stored in the Artifactory backend
type ArtifactoryArgs struct {
	// Username is the username to authenticate to Artifactory with.
	// Sourced from `ARTIFACTORY_USERNAME` environment var if unset.
	Username pulumi.StringPtrInput

	// Password is the password to authenticate to Artifactory with.
	// Sourced from `ARTIFACTORY_PASSWORD` environment var if unset.
	Password pulumi.StringPtrInput

	// Url is the base URL to artifactory, not the full repo and
	// subpath. However, it must include the path to the artifactory
	// installation - likely this will end in `/artifactory`.
	// Sourced from `ARTIFACTORY_URL` in the environment, if unset.
	Url pulumi.StringPtrInput

	// Repo is the repository name
	Repo pulumi.StringInput

	// Subpath is the path within the repository
	Subpath pulumi.StringInput

	// Workspace is the terraform workspace from which to read state
	Workspace pulumi.StringPtrInput
}

func (a *ArtifactoryArgs) toInternalArgs() pulumi.Input {
	return internal.ArtifactoryStateReferenceArgs{
		BackendType: pulumi.String("artifactory"),
		Username:    a.Username,
		Password:    a.Password,
		Url:         a.Url,
		Repo:        a.Repo,
		Subpath:     a.Subpath,
		Workspace:   a.Workspace,
	}
}

func (l *ArtifactoryArgs) validateArgs() error {
	if l.Repo == pulumi.String("") || l.Subpath == pulumi.String("") {
		return errors.New("`Repo` and `Subpath` are required parameters")
	}
	return nil
}
