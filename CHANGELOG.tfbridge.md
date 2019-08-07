# Pulumi-Terraform Bridge CHANGELOG

This CHANGELOG details changes made in the Terraform Bridge components of this repository - that is, `tfbridge`, `tfgen` and associated packages.

A CHANGELOG for the `pulumi-terraform` provider containing the `RemoteStateReference` resource, and associated SDKs, is in CHANGELOG.md.

## v0.18.4 (Unreleased)

- Terraform-based providers can now communicate detailed information about the difference between a resource's desired and actual state during a Pulumi update.
- Add the ability to inject CustomTimeouts into the InstanceDiff during a pulumi update.
- Better error message for missing required fields with default config ([#400](https://github.com/pulumi/pulumi-terraform/issues/400)).
- Change how Tfgen deals with package classes that are named Index to make them index_.ts
- Protect against panic in provider Create with InstanceState Meta initialization
- Use of the `RemoteStateReference` resource no longer results in a panic if the configured remote state cannot be accessed.
- Allow a provider to depend on a specific version of TypeScript.
- Allow users to specify a specific provider version.
- Add the ability to deprecate resources and datasources.
- Emit an appropriate user warning when Pulumi binary not found in Python setup.py.
- Add support for suppressing differences between the desired and actual state of a resource via the `ignoreChanges` property.
- Fix a bug that caused the recalculation of defaults for values that are normalized in resource state.
- The Python SDK generated for a provider now supports synchronous invokes.

## v0.18.3 (Released June 20, 2019)

- Fixed a bug that caused unnecessary changes if the first operation after upgrading a bridged provider was a `pulumi refresh`.
- Fixed a bug that caused maps with keys containing a '.' character to be incorrectly treated as containing nested maps when deserializing Terraform attributes.

### Improvements

- Automatically generate `isInstance` type guards for implementations of `Resource`.
- `TransformJSONDocument` now accepts arrays (in addition to maps).

