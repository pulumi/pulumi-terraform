CHANGELOG
=========

## HEAD (Unreleased)
_(none)_

---

## 3.4.0 (2020-10-26)
* Upgrade to Pulumi v2.12.0
* Upgrade to Terraform v0.13.5

## 3.3.0 (2020-10-01)
* Upgrade to Terraform v0.13.4

## 3.2.0 (2020-09-19)
* Upgrade to Terraform v0.13.2
* Upgrade to Pulumi v2.10.1

## 3.1.0 (2020-09-02)
* Upgrade to Terraform v0.13.2

## 3.0.0 (2020-08-28)
* Upgrade to Terraform v0.13.1
* Upgrade to Pulumi v2.9.1

## 2.5.0 (2020-08-10)
* Upgrade to Terraform v0.12.29

## 2.4.0 (2020-06-26)
* Upgrade to Terraform v0.12.28

## 2.3.0 (2020-05-27)
* Upgrade to Terraform v0.12.26

## 2.2.0 (2020-05-19)
* Upgrade to Terraform v0.12.25

## 2.1.0 (2020-05-15)
* Add support for Go with the `pulumi-terraform` Go SDK

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
