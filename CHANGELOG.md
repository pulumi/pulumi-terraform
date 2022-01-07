CHANGELOG
=========

## Notice (2022-01-06)

*As of this notice, using CHANGELOG.md is DEPRECATED. We will be using [GitHub Releases](https://github.com/pulumi/pulumi-terraform/releases) for this repository*

## HEAD (Unreleased)
_(none)_

---

## 5.5.1 (2021-12-20)
* Upgrade to Terraform v1.1.2

## 5.5.0 (2021-12-10)
* Upgrade to Terraform v1.1.0

## 5.4.3 (2021-11-18)
* Upgrade to Terraform v1.0.11

## 5.4.2 (2021-11-01)
* Upgrade to Terraform v1.0.10

## 5.4.1 (2021-10-04)
* Upgrade to Terraform v1.0.8

## 5.4.0 (2021-09-27)
* Upgrade to Terraform v1.0.7

## 5.3.0 (2021-09-09)
* Upgrade to Terraform v1.0.6

## 5.2.1 (2021-08-17)
* Add support for S3-style supported backends to the S3 Backend

## 5.2.0 (2021-05-12)
* Upgrade to Terraform v0.15.3

## 5.1.0 (2021-04-30)
* Upgrade to Terraform v0.15.1

## 5.0.0 (2021-04-26)
* Upgrade to Terraform v0.15.0
* Upgrade to Pulumi v3.0  
  **Please Note:** This upgrade to Pulumi v3 means that we have introduced a major version bump in the provider

## 4.7.0 (2021-04-12)
* Upgrade to Terraform v0.14.10

## 4.6.0 (2021-03-30)
* Upgrade to Terraform v0.14.9

## 4.5.0 (2021-02-25)
* Upgrade to Terraform v0.14.7
* Add support for AliCloud OSS Backend type

## 4.4.0 (2021-02-08)
* Upgrade to Terraform v0.14.6

## 4.3.0 (2021-01-21)
* Upgrade to Terraform v0.14.5

## 4.2.0 (2021-01-08)
* Upgrade to Terraform v0.14.4

## 4.1.0 (2020-12-22)
* Upgrade to Terraform v0.14.3

## 4.0.0 (2020-12-07)
* Upgrade to Terraform v0.14.0

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
