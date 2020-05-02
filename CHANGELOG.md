CHANGELOG
=========

## HEAD (Unreleased)
_(none)_

---

## v2.0.0 (2020-05-02)
* `RemoteStateReference` resources can now read states created with Terraform v0.12.24 and below
* Fix a bug with `RemoteStateReference`'s remote state backend on Python, where the workspace name was not  getting configured correctly (issue [#524](https://github.com/pulumi/pulumi-terraform/issues/524)).
* Added support for .NET with the `Pulumi.Terraform` NuGet package

## v1.1.0 (2019-10-04)
* Upgrade the Pulumi dependency requirements for NodeJS and Python SDKs

## v1.0.0 (2019-10-02)
* Use of the `RemoteStateReference` resource no longer results in a panic if the configured remote state cannot be accessed
* `RemoteStateReference` resources can now read states created with Terraform v0.12.9 and below
* Added support for Python with the `pulumi_terraform` package

## v0.18.1 (2019-05-16)
* Initial release of `@pulumi/terraform` with support for Node.js.
