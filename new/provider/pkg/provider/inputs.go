// TODO: is this the same for all resources?
package provider

import "github.com/pulumi/pulumi-go-provider/infer"

// ResourceInputs are inputs common to resource CRUD operations.
type ResourceInputs struct {
	// The field tags are used to provide metadata on the schema representation.
	// pulumi:"optional" specifies that a field is optional. This must be a pointer.
	// provider:"replaceOnChanges" specifies that the resource will be replaced if the field changes.
	Triggers *[]any `pulumi:"triggers,optional" provider:"replaceOnChanges"`

	Create *string `pulumi:"create,optional"`
	Update *string `pulumi:"update,optional"`
	Delete *string `pulumi:"delete,optional"`
}

// TODO: how is this different for each resource?
// Annotate lets you provide descriptions and default values for fields and they will
// be visible in the provider's schema and the generated SDKs.
func (c *ResourceInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Triggers, `Trigger a resource replacement on changes to any of these values. The
trigger values can be of any type. If a value is different in the current update compared to the
previous update, the resource will be replaced, i.e., the "create" command will be re-run.
Please see the resource documentation for examples.`)
	a.Describe(&c.Create, "The command to run on create.")
	a.Describe(&c.Delete, `The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
Command resource from previous create or update steps.`)
	a.Describe(&c.Update, `The command to run on update, if empty, create will 
run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR 
are set to the stdout and stderr properties of the Command resource from previous 
create or update steps.`)
}
