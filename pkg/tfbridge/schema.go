// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"fmt"

	"github.com/golang/glog"

	pbstruct "github.com/golang/protobuf/ptypes/struct"
	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/flatmap"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/resource/plugin"
	"github.com/pulumi/pulumi/pkg/util/contract"
)

// MakeTerraformInputs takes a property map plus custom schema info and does whatever is necessary to prepare it for
// use by Terraform.  Note that this function may have side effects, for instance if it is necessary to spill an asset
// to disk in order to create a name out of it.  Please take care not to call it superfluously!
func MakeTerraformInputs(res *PulumiResource, m resource.PropertyMap,
	tfs map[string]*schema.Schema, ps map[string]*SchemaInfo,
	defaults, useRawNames bool) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Enumerate the inputs provided and add them to the map using their Terraform names.
	for key, value := range m {
		// First translate the Pulumi property name to a Terraform name.
		name, tfi, psi := getInfoFromPulumiName(key, tfs, ps, useRawNames)
		contract.Assert(name != "")

		// And then translate the property value.
		v, err := MakeTerraformInput(res, name, value, tfi, psi, defaults, useRawNames)
		if err != nil {
			return nil, err
		}
		result[name] = v
	}

	// Now enumerate and propagate defaults if the corresponding values are still missing.
	for key, info := range ps {
		if v, has := result[key]; has {
			glog.V(9).Infof("Created Terraform input: %v = %v", key, v)
		} else if defaults && info.HasDefault() {
			if info.Default.Value != nil {
				result[key] = info.Default.Value
				glog.V(9).Infof("Created Terraform input: %v = %v (default)", key, result[key])
			} else if from := info.Default.From; from != nil {
				result[key] = from(res)
				glog.V(9).Infof("Created Terraform input: %v = %v (default from fnc)", key, result[key])
			} else {
				contract.Failf("Default missing Value or From")
			}
		} else {
			glog.V(9).Infof("Skipped Terraform input: %v (skipped or no defaults)", key)
		}
	}

	if glog.V(5) {
		for k, v := range result {
			glog.V(5).Infof("Terraform input %v = %v", k, v)
		}
	}

	return result, nil
}

// MakeTerraformInput takes a single property plus custom schema info and does whatever is necessary to prepare it for
// use by Terraform.  Note that this function may have side effects, for instance if it is necessary to spill an asset
// to disk in order to create a name out of it.  Please take care not to call it superfluously!
func MakeTerraformInput(res *PulumiResource, name string,
	v resource.PropertyValue, tfs *schema.Schema, ps *SchemaInfo, defaults, rawNames bool) (interface{}, error) {
	if v.IsNull() {
		return nil, nil
	} else if v.IsBool() {
		return v.BoolValue(), nil
	} else if v.IsNumber() {
		return int(v.NumberValue()), nil // convert floats to ints.
	} else if v.IsString() {
		return v.StringValue(), nil
	} else if v.IsArray() {
		// FIXME: marshal/unmarshal sets properly.
		var arr []interface{}
		for i, elem := range v.ArrayValue() {
			var etfs *schema.Schema
			if tfs != nil {
				if sch, issch := tfs.Elem.(*schema.Schema); issch {
					etfs = sch
				} else if _, isres := tfs.Elem.(*schema.Resource); isres {
					// The IsObject case below expects a schema whose `Elem` is
					// a Resource, so just pass the full List schema
					etfs = tfs
				}
			}
			var eps *SchemaInfo
			if ps != nil {
				eps = ps.Elem
			}
			e, err := MakeTerraformInput(res, fmt.Sprintf("%v[%v]", name, i), elem, etfs, eps, defaults, rawNames)
			if err != nil {
				return nil, err
			}
			arr = append(arr, e)
		}
		return arr, nil
	} else if v.IsAsset() {
		// We require that there be asset information, otherwise an error occurs.
		if ps == nil || ps.Asset == nil {
			return nil,
				errors.Errorf("Encountered an asset %v but asset translation instructions were missing", name)
		} else if !ps.Asset.IsAsset() {
			return nil,
				errors.Errorf("Invalid asset translation instructions for %v; expected an asset", name)
		}
		return ps.Asset.TranslateAsset(v.AssetValue())
	} else if v.IsArchive() {
		// We require that there be archive information, otherwise an error occurs.
		if ps == nil || ps.Asset == nil {
			return nil,
				errors.Errorf("Encountered an archive %v but asset translation instructions were missing", name)
		} else if !ps.Asset.IsArchive() {
			return nil,
				errors.Errorf("Invalid asset translation instructions for %v; expected an archive", name)
		}
		return ps.Asset.TranslateArchive(v.ArchiveValue())
	} else if v.IsObject() {
		var tfflds map[string]*schema.Schema
		if tfs != nil {
			if res, isres := tfs.Elem.(*schema.Resource); isres {
				tfflds = res.Schema
			}
		}
		var psflds map[string]*SchemaInfo
		if ps != nil {
			psflds = ps.Fields
		}
		return MakeTerraformInputs(res, v.ObjectValue(), tfflds, psflds, defaults, rawNames || useRawNames(tfs))
	} else if v.IsComputed() || v.IsOutput() {
		// If any variables are unknown, we need to mark them in the inputs so the config map treats it right.  This
		// requires the use of the special UnknownVariableValue sentinel in Terraform, which is how it internally stores
		// interpolated variables whose inputs are currently unknown.
		return config.UnknownVariableValue, nil
	}

	contract.Failf("Unexpected value marshaled: %v", v)
	return nil, nil
}

// MakeTerraformInputsFromRPC unmarshals an RPC payload of properties and turns the results into Terraform inputs.
func MakeTerraformInputsFromRPC(res *PulumiResource, m *pbstruct.Struct,
	tfs map[string]*schema.Schema, ps map[string]*SchemaInfo,
	allowUnknowns bool, defaults bool) (map[string]interface{}, error) {
	props, err := plugin.UnmarshalProperties(m,
		plugin.MarshalOptions{AllowUnknowns: allowUnknowns, SkipNulls: true})
	if err != nil {
		return nil, err
	}
	return MakeTerraformInputs(res, props, tfs, ps, defaults, false)
}

// MakeTerraformResult expands a Terraform-style flatmap into an expanded Pulumi resource property map.  This respects
// the property maps so that results end up with their correct Pulumi names when shipping back to the engine.
func MakeTerraformResult(props map[string]string,
	tfs map[string]*schema.Schema, ps map[string]*SchemaInfo) resource.PropertyMap {
	outs := make(map[string]interface{})
	for _, key := range flatmap.Map(props).Keys() {
		outs[key] = flatmap.Expand(props, key)
	}
	return MakeTerraformOutputs(outs, tfs, ps, false)
}

// MakeTerraformOutputs takes an expanded Terraform property map and returns a Pulumi equivalent.  This respects
// the property maps so that results end up with their correct Pulumi names when shipping back to the engine.
func MakeTerraformOutputs(outs map[string]interface{},
	tfs map[string]*schema.Schema, ps map[string]*SchemaInfo, rawNames bool) resource.PropertyMap {
	result := make(resource.PropertyMap)
	for key, value := range outs {
		// First do a lookup of the name/info.
		name, tfi, psi := getInfoFromTerraformName(key, tfs, ps, rawNames)
		contract.Assert(name != "")

		// Next perform a translation of the value accordingly.
		result[name] = MakeTerraformOutput(value, tfi, psi, rawNames)
	}

	if glog.V(5) {
		for k, v := range result {
			glog.V(5).Infof("Terraform output %v = %v", k, v)
		}
	}

	return result
}

// MakeTerraformOutput takes a single Terraform property and returns the Pulumi equivalent.
func MakeTerraformOutput(v interface{},
	tfs *schema.Schema, ps *SchemaInfo, rawNames bool) resource.PropertyValue {
	if v == nil {
		return resource.NewNullProperty()
	}
	switch t := v.(type) {
	case bool:
		return resource.NewBoolProperty(t)
	case int:
		return resource.NewNumberProperty(float64(t))
	case string:
		// If the string is the special unknown property sentinel, reflect back an unknown computed property.  Note that
		// Terraform doesn't carry the types along with it, so the best we can do is give back a computed string.
		if t == config.UnknownVariableValue {
			elem := resource.Computed{Element: resource.NewStringProperty("")}
			return resource.NewComputedProperty(elem)
		}
		// Else it's just a string.
		return resource.NewStringProperty(t)
	case []interface{}:
		var tfes *schema.Schema
		if tfs != nil {
			if sch, issch := tfs.Elem.(*schema.Schema); issch {
				tfes = sch
			} else if _, isres := tfs.Elem.(*schema.Resource); isres {
				// The map[string]interface{} case below expects a schema whose
				// `Elem` is a Resource, so just pass the full List schema
				tfes = tfs
			}
		}
		var pes *SchemaInfo
		if ps != nil {
			pes = ps.Elem
		}
		var arr []resource.PropertyValue
		for _, elem := range t {
			arr = append(arr, MakeTerraformOutput(elem, tfes, pes, rawNames))
		}
		return resource.NewArrayProperty(arr)
	case map[string]interface{}:
		var tfflds map[string]*schema.Schema
		if tfs != nil {
			if res, isres := tfs.Elem.(*schema.Resource); isres {
				tfflds = res.Schema
			}
		}
		var psflds map[string]*SchemaInfo
		if ps != nil {
			psflds = ps.Fields
		}
		obj := MakeTerraformOutputs(t, tfflds, psflds, rawNames || useRawNames(tfs))
		return resource.NewObjectProperty(obj)
	default:
		contract.Failf("Unexpected TF output property value: %v", v)
		return resource.NewNullProperty()
	}
}

// MakeTerraformConfig creates a Terraform config map, used in state and diff calculations, from a Pulumi property map.
func MakeTerraformConfig(res *PulumiResource, m resource.PropertyMap,
	tfs map[string]*schema.Schema, ps map[string]*SchemaInfo, defaults bool) (*terraform.ResourceConfig, error) {
	// Convert the resource bag into an untyped map, and then create the resource config object.
	inputs, err := MakeTerraformInputs(res, m, tfs, ps, defaults, false)
	if err != nil {
		return nil, err
	}
	return MakeTerraformConfigFromInputs(inputs)
}

// MakeTerraformConfigFromRPC creates a Terraform config map from a Pulumi RPC property map.
func MakeTerraformConfigFromRPC(res *PulumiResource, m *pbstruct.Struct,
	tfs map[string]*schema.Schema, ps map[string]*SchemaInfo,
	allowUnknowns, defaults bool) (*terraform.ResourceConfig, error) {
	props, err := plugin.UnmarshalProperties(m,
		plugin.MarshalOptions{AllowUnknowns: allowUnknowns, SkipNulls: true})
	if err != nil {
		return nil, err
	}
	return MakeTerraformConfig(res, props, tfs, ps, defaults)
}

// MakeTerraformConfigFromInputs creates a new Terraform configuration object from a set of Terraform inputs.
func MakeTerraformConfigFromInputs(inputs map[string]interface{}) (*terraform.ResourceConfig, error) {
	cfg, err := config.NewRawConfig(inputs)
	if err != nil {
		return nil, err
	}
	return terraform.NewResourceConfig(cfg), nil
}

// MakeTerraformAttributes converts a Pulumi property bag into its Terraform equivalent.  This requires
// flattening everything and serializing individual properties as strings.  This is a little awkward, but it's how
// Terraform represents resource properties (schemas are simply sugar on top).
func MakeTerraformAttributes(res *PulumiResource, m resource.PropertyMap,
	tfs map[string]*schema.Schema, ps map[string]*SchemaInfo, defaults bool) (map[string]string, error) {
	// Turn the resource properties into a map.  For the most part, this is a straight Mappable, but we use MapReplace
	// because we use float64s and Terraform uses ints, to represent numbers.
	inputs, err := MakeTerraformInputs(res, m, tfs, ps, defaults, false)
	if err != nil {
		return nil, err
	}
	return MakeTerraformAttributesFromInputs(inputs), nil
}

// MakeTerraformAttributesFromRPC unmarshals an RPC property map and calls through to MakeTerraformAttributes.
func MakeTerraformAttributesFromRPC(res *PulumiResource, m *pbstruct.Struct,
	tfs map[string]*schema.Schema, ps map[string]*SchemaInfo,
	allowUnknowns, defaults bool) (map[string]string, error) {
	props, err := plugin.UnmarshalProperties(m,
		plugin.MarshalOptions{AllowUnknowns: allowUnknowns, SkipNulls: true})
	if err != nil {
		return nil, err
	}
	return MakeTerraformAttributes(res, props, tfs, ps, defaults)
}

// MakeTerraformAttributesFromInputs creates a flat Terraform map from a structured set of Terraform inputs.
func MakeTerraformAttributesFromInputs(inputs map[string]interface{}) map[string]string {
	return flatmap.Flatten(inputs)
}

// MakeTerraformDiff takes a bag of old and new properties, and returns two things: the existing resource's state as
// an attribute map, alongside a Terraform diff for the old versus new state.  If there was no existing state, the
// returned attributes will be empty (because the resource doesn't yet exist).
func MakeTerraformDiff(old resource.PropertyMap, new resource.PropertyMap,
	tfs map[string]*schema.Schema, ps map[string]*SchemaInfo) (*terraform.InstanceState,
	*terraform.InstanceDiff, error) {
	// BUGBUG[pulumi/pulumi-terraform#22]: avoid spilling except for during creation.
	diff := make(map[string]*terraform.ResourceAttrDiff)
	// Add all new property values.
	if new != nil {
		inputs, err := MakeTerraformAttributes(nil, new, tfs, ps, false)
		if err != nil {
			return nil, nil, err
		}
		for p, v := range inputs {
			if diff[p] == nil {
				diff[p] = &terraform.ResourceAttrDiff{}
			}
			diff[p].New = v
		}
	}
	// Now add all old property values, provided they exist in new.
	existing := make(map[string]string)
	if old != nil {
		inputs, err := MakeTerraformAttributes(nil, old, tfs, ps, false)
		if err != nil {
			return nil, nil, err
		}
		for p, v := range inputs {
			if d, has := diff[p]; has {
				d.Old = v
			}
			existing[p] = v
		}
	}
	return &terraform.InstanceState{Attributes: existing},
		&terraform.InstanceDiff{Attributes: diff}, nil
}

// MakeTerraformDiffFromRPC takes RPC maps of old and new properties, unmarshals them, and calls into MakeTerraformDiff.
func MakeTerraformDiffFromRPC(old *pbstruct.Struct, new *pbstruct.Struct,
	tfs map[string]*schema.Schema, ps map[string]*SchemaInfo) (*terraform.InstanceState,
	*terraform.InstanceDiff, error) {
	var err error
	var oldprops resource.PropertyMap
	if old != nil {
		oldprops, err = plugin.UnmarshalProperties(old,
			plugin.MarshalOptions{SkipNulls: true})
		if err != nil {
			return nil, nil, err
		}
	}
	var newprops resource.PropertyMap
	if new != nil {
		newprops, err = plugin.UnmarshalProperties(new,
			plugin.MarshalOptions{AllowUnknowns: true, SkipNulls: true})
		if err != nil {
			return nil, nil, err
		}
	}
	return MakeTerraformDiff(oldprops, newprops, tfs, ps)
}

// useRawNames returns true if raw, unmangled names should be preserved.  This is only true for Terraform maps.
func useRawNames(tfs *schema.Schema) bool {
	return tfs != nil && tfs.Type == schema.TypeMap
}

// getInfoFromTerraformName does a map lookup to find the Pulumi name and schema info, if any.
func getInfoFromTerraformName(key string,
	tfs map[string]*schema.Schema, ps map[string]*SchemaInfo, rawName bool) (resource.PropertyKey,
	*schema.Schema, *SchemaInfo) {
	info := ps[key]
	var name string
	if info != nil {
		name = info.Name
	}
	if name == "" {
		if rawName {
			// If raw names are requested, simply preserve the Terraform name.
			name = key
		} else {
			// If no name override exists, use the default name mangling scheme.
			name = TerraformToPulumiName(key, false)
		}
	}
	return resource.PropertyKey(name), tfs[key], info
}

// getInfoFromPulumiName does a reverse map lookup to find the Terraform name and schema info for a Pulumi name, if any.
func getInfoFromPulumiName(key resource.PropertyKey,
	tfs map[string]*schema.Schema, ps map[string]*SchemaInfo, rawName bool) (string,
	*schema.Schema, *SchemaInfo) {
	// To do this, we will first look to see if there's a known custom schema that uses this name.  If yes, we
	// prefer to use that.  To do this, we must use a reverse lookup.  (In the future we may want to make a
	// lookaside map to avoid the traversal of this map.)  Otherwise, use the standard name mangling scheme.
	ks := string(key)
	for tfname, schinfo := range ps {
		if schinfo != nil && schinfo.Name == ks {
			return tfname, tfs[tfname], schinfo
		}
	}
	var name string
	if rawName {
		// If raw names are requested, they will not have been mangled, so preserve the name as-is.
		name = ks
	} else {
		// Otherwise, transform the Pulumi name to the Terraform name using the standard mangling scheme.
		name = PulumiToTerraformName(ks)
	}
	return name, tfs[name], ps[ks]
}
