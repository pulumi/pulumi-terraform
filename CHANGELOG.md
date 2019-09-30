# `@pulumi/terraform` CHANGELOG

This CHANGELOG details important changes made in each version of the
`terraform` provider, the `@pulumi/terraform` Node.js package and the
`pulumi_terraform` Python package.

## v0.18.4 (Unreleased)

- Initial release of `pulumi_terraform` for Python.
- `RemoteStateReference` resources can now read states created with Terraform
  0.12.6 and below.
- Use of the `RemoteStateReference` resource no longer results in a panic if
  the configured remote state cannot be accessed.

## v0.18.2 (Released May 28th, 2019)

- Improved the package `README` file to reflect usage of the
  `@pulumi/terraform` package rather than the Terraform bridge.

## v0.18.1 (Released May 16th, 2019)

- Initial release of `@pulumi/terraform` with support for Node.js.
