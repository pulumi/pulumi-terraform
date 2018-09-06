// Copyright 2016-2018, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tfbridge

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/stretchr/testify/assert"
)

// TestTerraformInputs verifies that we translate Pulumi inputs into Terraform inputs.
func TestTerraformInputs(t *testing.T) {
	result, err := MakeTerraformInputs(
		nil, /*res*/
		nil, /*olds*/
		resource.NewPropertyMapFromMap(map[string]interface{}{
			"boolPropertyValue":   false,
			"numberPropertyValue": 42,
			"floatPropertyValue":  99.6767932,
			"stringo":             "ognirts",
			"arrayPropertyValue":  []interface{}{"an array"},
			"unknownArrayValue":   resource.Computed{Element: resource.NewStringProperty("")},
			"objectPropertyValue": map[string]interface{}{
				"propertyA": "a",
				"propertyB": true,
			},
			"mapPropertyValue": map[string]interface{}{
				"propertyA": "a",
				"propertyB": true,
				"propertyC": map[string]interface{}{
					"nestedPropertyA": true,
				},
			},
			"nestedResources": []map[string]interface{}{{
				"configuration": map[string]interface{}{
					"configurationValue": true,
				},
			}},
			"optionalConfig": map[string]interface{}{
				"someValue":      true,
				"someOtherValue": "a value",
			},
			"optionalConfigOther": map[string]interface{}{
				"someValue":      true,
				"someOtherValue": "a value",
			},
			"mapWithResourceElem": map[string]interface{}{
				"someValue": "a value",
			},
		}),
		map[string]*schema.Schema{
			// Type mapPropertyValue as a map so that keys aren't mangled in the usual way.
			"float_property_value": {Type: schema.TypeFloat},
			"unknown_array_value":  {Type: schema.TypeList},
			"map_property_value":   {Type: schema.TypeMap},
			"nested_resource": {
				Type:     schema.TypeList,
				MaxItems: 2,
				// Embed a `*schema.Resource` to validate that type directed
				// walk of the schema successfully walks inside Resources as well
				// as Schemas.
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration": {Type: schema.TypeMap},
					},
				},
			},
			"optional_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"some_value":       {Type: schema.TypeBool},
						"some_other_value": {Type: schema.TypeString},
					},
				},
			},
			"optional_config_other": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"some_value":       {Type: schema.TypeBool},
						"some_other_value": {Type: schema.TypeString},
					},
				},
			},
			"map_with_resource_elem": {
				Type: schema.TypeMap,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"some_value": {Type: schema.TypeString},
					},
				},
			},
		},
		map[string]*SchemaInfo{
			// Reverse map string_property_value to the stringo property.
			"string_property_value": {
				Name: "stringo",
			},
			"optional_config_other": {
				Name:        "optionalConfigOther",
				MaxItemsOne: boolPointer(true),
			},
		},
		nil,   /* assets */
		false, /*defaults*/
		false, /*useRawNames*/
	)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"bool_property_value":   false,
		"number_property_value": 42,
		"float_property_value":  99.6767932,
		"string_property_value": "ognirts",
		"array_property_value":  []interface{}{"an array"},
		"unknown_array_value":   []interface{}{config.UnknownVariableValue},
		"object_property_value": map[string]interface{}{
			"property_a": "a",
			"property_b": true,
		},
		"map_property_value": map[string]interface{}{
			"propertyA": "a",
			"propertyB": true,
			"propertyC": map[string]interface{}{
				"nestedPropertyA": true,
			},
		},
		"nested_resource": []interface{}{
			map[string]interface{}{
				"configuration": map[string]interface{}{
					"configurationValue": true,
				},
			},
		},
		"optional_config": []interface{}{
			map[string]interface{}{
				"some_value":       true,
				"some_other_value": "a value",
			},
		},
		"optional_config_other": []interface{}{
			map[string]interface{}{
				"some_value":       true,
				"some_other_value": "a value",
			},
		},
		"map_with_resource_elem": []interface{}{
			map[string]interface{}{
				"some_value": "a value",
			},
		},
	}, result)

	_, err = MakeTerraformInputs(
		nil, /*res*/
		nil, /*olds*/
		resource.NewPropertyMapFromMap(map[string]interface{}{
			"nilPropertyValue": nil,
		}),
		nil,   /* tfs */
		nil,   /* ps */
		nil,   /* assets */
		false, /*defaults*/
		false, /*useRawNames*/
	)
	assert.Error(t, err)
}

type MyString string

// TestTerraformOutputs verifies that we translate Terraform outputs into Pulumi outputs.
func TestTerraformOutputs(t *testing.T) {
	result := MakeTerraformOutputs(
		map[string]interface{}{
			"nil_property_value":       nil,
			"bool_property_value":      false,
			"number_property_value":    42,
			"float_property_value":     99.6767932,
			"string_property_value":    "ognirts",
			"my_string_property_value": MyString("ognirts"),
			"array_property_value":     []interface{}{"an array"},
			"object_property_value": map[string]interface{}{
				"property_a": "a",
				"property_b": true,
			},
			"map_property_value": map[string]interface{}{
				"propertyA": "a",
				"propertyB": true,
				"propertyC": map[string]interface{}{
					"nestedPropertyA": true,
				},
			},
			"nested_resource": []interface{}{
				map[string]interface{}{
					"configuration": map[string]interface{}{
						"configurationValue": true,
					},
				},
			},
			"optional_config": []interface{}{
				map[string]interface{}{
					"some_value":       true,
					"some_other_value": "a value",
				},
			},
			"optional_config_other": []interface{}{
				map[string]interface{}{
					"some_value":       true,
					"some_other_value": "a value",
				},
			},
		},
		map[string]*schema.Schema{
			// Type mapPropertyValue as a map so that keys aren't mangled in the usual way.
			"float_property_value":     {Type: schema.TypeFloat},
			"my_string_property_value": {Type: schema.TypeString},
			"map_property_value":       {Type: schema.TypeMap},
			"nested_resource": {
				Type:     schema.TypeList,
				MaxItems: 2,
				// Embed a `*schema.Resource` to validate that type directed
				// walk of the schema successfully walks inside Resources as well
				// as Schemas.
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration": {Type: schema.TypeMap},
					},
				},
			},
			"optional_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"some_value":       {Type: schema.TypeBool},
						"some_other_value": {Type: schema.TypeString},
					},
				},
			},
			"optional_config_other": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"some_value":       {Type: schema.TypeBool},
						"some_other_value": {Type: schema.TypeString},
					},
				},
			},
		},
		map[string]*SchemaInfo{
			// Reverse map string_property_value to the stringo property.
			"string_property_value": {
				Name: "stringo",
			},
			"optional_config_other": {
				Name:        "optionalConfigOther",
				MaxItemsOne: boolPointer(true),
			},
		},
		nil,   /* assets */
		false, /*useRawNames*/
	)
	assert.Equal(t, resource.NewPropertyMapFromMap(map[string]interface{}{
		"nilPropertyValue":      nil,
		"boolPropertyValue":     false,
		"numberPropertyValue":   42,
		"floatPropertyValue":    99.6767932,
		"stringo":               "ognirts",
		"myStringPropertyValue": "ognirts",
		"arrayPropertyValue":    []interface{}{"an array"},
		"objectPropertyValue": map[string]interface{}{
			"propertyA": "a",
			"propertyB": true,
		},
		"mapPropertyValue": map[string]interface{}{
			"propertyA": "a",
			"propertyB": true,
			"propertyC": map[string]interface{}{
				"nestedPropertyA": true,
			},
		},
		"nestedResources": []map[string]interface{}{{
			"configuration": map[string]interface{}{
				"configurationValue": true,
			},
		}},
		"optionalConfig": map[string]interface{}{
			"someValue":      true,
			"someOtherValue": "a value",
		},
		"optionalConfigOther": map[string]interface{}{
			"someValue":      true,
			"someOtherValue": "a value",
		},
	}), result)
}

func TestTerraformAttributes(t *testing.T) {
	result, err := MakeTerraformAttributesFromInputs(
		map[string]interface{}{
			"nil_property_value":    nil,
			"bool_property_value":   false,
			"number_property_value": 42,
			"float_property_value":  99.6767932,
			"string_property_value": "ognirts",
			"array_property_value":  []interface{}{"an array"},
			"object_property_value": map[string]interface{}{
				"property_a": "a",
				"property_b": true,
			},
			"map_property_value": map[string]interface{}{
				"propertyA": "a",
				"propertyB": true,
				"propertyC": map[string]interface{}{
					"nestedPropertyA": true,
				},
			},
			"nested_resources": []interface{}{
				map[string]interface{}{
					"configuration": map[string]interface{}{
						"configurationValue": true,
					},
				},
			},
			"set_property_value":            []interface{}{"set member 1", "set member 2"},
			"string_with_bad_interpolation": "some ${interpolated:value} with syntax errors",
		},
		map[string]*schema.Schema{
			"nil_property_value":    {Type: schema.TypeMap},
			"bool_property_value":   {Type: schema.TypeBool},
			"number_property_value": {Type: schema.TypeInt},
			"float_property_value":  {Type: schema.TypeFloat},
			"string_property_value": {Type: schema.TypeString},
			"array_property_value": {
				Type: schema.TypeList,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"object_property_value": {Type: schema.TypeMap},
			"map_property_value":    {Type: schema.TypeMap},
			"nested_resources": {
				Type:     schema.TypeList,
				MaxItems: 1,
				// Embed a `*schema.Resource` to validate that type directed
				// walk of the schema successfully walks inside Resources as well
				// as Schemas.
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration": {Type: schema.TypeMap},
					},
				},
			},
			"set_property_value": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"string_with_bad_interpolation": {Type: schema.TypeString},
		})

	assert.NoError(t, err)
	assert.Equal(t, result, map[string]string{
		"array_property_value.#":                              "1",
		"array_property_value.0":                              "an array",
		"bool_property_value":                                 "false",
		"float_property_value":                                "99.6767932",
		"map_property_value.%":                                "3",
		"map_property_value.propertyA":                        "a",
		"map_property_value.propertyB":                        "true",
		"map_property_value.propertyC.%":                      "1",
		"map_property_value.propertyC.nestedPropertyA":        "true",
		"nested_resources.#":                                  "1",
		"nested_resources.0.%":                                "1",
		"nested_resources.0.configuration.%":                  "1",
		"nested_resources.0.configuration.configurationValue": "true",
		"number_property_value":                               "42",
		"object_property_value.%":                             "2",
		"object_property_value.property_a":                    "a",
		"object_property_value.property_b":                    "true",
		"set_property_value.#":                                "2",
		"set_property_value.3618983862":                       "set member 2",
		"set_property_value.4237827189":                       "set member 1",
		"string_property_value":                               "ognirts",
		"string_with_bad_interpolation":                       "some ${interpolated:value} with syntax errors",
	})

	// MapFieldWriter has issues with values of TypeMap. Build a schema without such values s.t. we can test
	// MakeTerraformAttributes against the output of MapFieldWriter.
	sharedSchema := map[string]*schema.Schema{
		"bool_property_value":   {Type: schema.TypeBool},
		"number_property_value": {Type: schema.TypeInt},
		"float_property_value":  {Type: schema.TypeFloat},
		"string_property_value": {Type: schema.TypeString},
		"array_property_value": {
			Type: schema.TypeList,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
		"nested_resource_value": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"nested_set_property": {
						Type: schema.TypeSet,
						Elem: &schema.Schema{Type: schema.TypeString},
					},
					"nested_string_property": {Type: schema.TypeString},
				},
			},
		},
		"set_property_value": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
		"string_with_bad_interpolation": {Type: schema.TypeString},
	}
	sharedInputs := map[string]interface{}{
		"bool_property_value":   false,
		"number_property_value": 42,
		"float_property_value":  99.6767932,
		"string_property_value": "ognirts",
		"array_property_value":  []interface{}{"an array"},
		"nested_resource_value": map[string]interface{}{
			"nested_set_property":    []interface{}{"nested set member"},
			"nested_string_property": "value",
		},
		"set_property_value":            []interface{}{"set member 1", "set member 2"},
		"string_with_bad_interpolation": "some ${interpolated:value} with syntax errors",
	}

	// Build a TF attribute map using schema.MapFieldWriter.
	cfg, err := MakeTerraformConfigFromInputs(sharedInputs)
	assert.NoError(t, err)
	reader := &schema.ConfigFieldReader{Config: cfg, Schema: sharedSchema}
	writer := &schema.MapFieldWriter{Schema: sharedSchema}
	for k := range sharedInputs {
		f, ferr := reader.ReadField([]string{k})
		assert.NoError(t, ferr)

		err = writer.WriteField([]string{k}, f.Value)
		assert.NoError(t, err)
	}
	expected := writer.Map()

	// Build the same using MakeTerraformAttributesFromInputs.
	result, err = MakeTerraformAttributesFromInputs(sharedInputs, sharedSchema)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

// Test that meta-properties are correctly produced.
func TestMetaProperties(t *testing.T) {
	const resName = "example_resource"
	res := testTFProvider.ResourcesMap["example_resource"]

	info := &terraform.InstanceInfo{Type: resName}
	state := &terraform.InstanceState{ID: "0", Attributes: map[string]string{}, Meta: map[string]interface{}{}}
	read, err := testTFProvider.Refresh(info, state)
	assert.NoError(t, err)
	assert.NotNil(t, read)

	props := MakeTerraformResult(read, res.Schema, nil)
	assert.NotNil(t, props)

	attrs, meta, err := MakeTerraformAttributes(res, props, res.Schema, nil, false)
	assert.NoError(t, err)
	assert.NotNil(t, attrs)
	assert.NotNil(t, meta)

	assert.Equal(t, strconv.Itoa(res.SchemaVersion), meta["schema_version"])

	state.Attributes, state.Meta = attrs, meta
	read2, err := testTFProvider.Refresh(info, state)
	assert.NoError(t, err)
	assert.NotNil(t, read2)
	assert.Equal(t, read, read2)

	// Delete the resource's meta-property and ensure that we re-populate its schema version.
	delete(props, metaKey)

	attrs, meta, err = MakeTerraformAttributes(res, props, res.Schema, nil, false)
	assert.NoError(t, err)
	assert.NotNil(t, attrs)
	assert.NotNil(t, meta)

	assert.Equal(t, strconv.Itoa(res.SchemaVersion), meta["schema_version"])

	// Remove the resource's meta-attributes and ensure that we do not include them in the result.
	read2.Meta = map[string]interface{}{}
	props = MakeTerraformResult(read2, res.Schema, nil)
	assert.NotNil(t, props)
	assert.NotContains(t, props, metaKey)

	// Ensure that timeouts are populated and preserved.
	state.ID = ""
	cfg, err := config.NewRawConfig(map[string]interface{}{})
	assert.NoError(t, err)
	diff, err := testTFProvider.Diff(info, state, terraform.NewResourceConfig(cfg))
	assert.NoError(t, err)
	create, err := testTFProvider.Apply(info, state, diff)
	assert.NoError(t, err)

	props = MakeTerraformResult(create, res.Schema, nil)
	assert.NotNil(t, props)

	attrs, meta, err = MakeTerraformAttributes(res, props, res.Schema, nil, false)
	assert.NoError(t, err)
	assert.NotNil(t, attrs)
	assert.NotNil(t, meta)

	assert.Contains(t, meta, schema.TimeoutKey)
}

// Test that an unset list still generates a length attribute.
func TestEmptyListAttribute(t *testing.T) {
	result, err := MakeTerraformAttributesFromInputs(
		map[string]interface{}{},
		map[string]*schema.Schema{
			"list_property": {Type: schema.TypeList, Optional: true},
		})

	assert.NoError(t, err)
	assert.Equal(t, result, map[string]string{
		"list_property.#": "0",
	})
}

func TestDefaults(t *testing.T) {
	// Produce maps with the following properties, and then validate them:
	//     - aaa string; no defaults, no inputs => empty
	//     - bbb string; no defaults, input "BBB" => "BBB"
	//     - ccc string; TF default "CCC", no inputs => "CCC"
	//     - cc2 string; TF default "CC2" (func), no inputs => "CC2"
	//     - ddd string; TF default "TFD", input "DDD" => "DDD"
	//     - dd2 string; TF default "TD2" (func), input "DDD" => "DDD"
	//     - eee string; PS default "EEE", no inputs => "EEE"
	//     - ee2 string; PS default "EE2" (func), no inputs => "EE2"
	//     - fff string; PS default "PSF", input "FFF" => "FFF"
	//     - ff2 string; PS default "PF2", input "FFF" => "FFF"
	//     - ggg string; TF default "TFG", PS default "PSG", no inputs => "PSG" (PS wins)
	//     - hhh string; TF default "TFH", PS default "PSH", input "HHH" => "HHH"
	//     - iii string; old default "OLI", TF default "TFI", PS default "PSI", no input => "OLD"
	//     - jjj string: old input "OLJ", no defaults, no input => no merged input
	//     - lll: old default "OLL", TF default "TFL", no input => "OLL"
	//     - mmm: old default "OLM", PS default "PSM", no input => "OLM"
	//     - uuu: PS default "PSU", envvars w/o valiues => "PSU"
	//     - vvv: PS default 42, envvars with values => 1337
	//     - www: old default "OLW", deprecated, required, no input -> "OLW"
	//     - xxx: old default "OLX", deprecated, no input => nothing
	//     - yyy: TF default "TLY", deprecated, no input => nothing
	err := os.Setenv("PTFV2", "1337")
	assert.Nil(t, err)
	asset, err := resource.NewTextAsset("hello")
	assert.Nil(t, err)
	assets := make(AssetTable)
	tfs := map[string]*schema.Schema{
		"ccc": {Type: schema.TypeString, Default: "CCC"},
		"cc2": {Type: schema.TypeString, DefaultFunc: func() (interface{}, error) { return "CC2", nil }},
		"ddd": {Type: schema.TypeString, Default: "TFD"},
		"dd2": {Type: schema.TypeString, DefaultFunc: func() (interface{}, error) { return "TD2", nil }},
		"ggg": {Type: schema.TypeString, Default: "TFG"},
		"hhh": {Type: schema.TypeString, Default: "TFH"},
		"iii": {Type: schema.TypeString, Default: "TFI"},
		"jjj": {Type: schema.TypeString},
		"lll": {Type: schema.TypeString, Default: "TFL"},
		"mmm": {Type: schema.TypeString},
		"uuu": {Type: schema.TypeString},
		"vvv": {Type: schema.TypeInt},
		"www": {Type: schema.TypeString, Deprecated: "deprecated", Required: true},
		"xxx": {Type: schema.TypeString, Deprecated: "deprecated", Optional: true},
		"yyy": {Type: schema.TypeString, Default: "TLY", Deprecated: "deprecated", Optional: true},
		"zzz": {Type: schema.TypeString},
	}
	ps := map[string]*SchemaInfo{
		"eee": {Default: &DefaultInfo{Value: "EEE"}},
		"ee2": {Default: &DefaultInfo{From: func(res *PulumiResource) (interface{}, error) { return "EE2", nil }}},
		"fff": {Default: &DefaultInfo{Value: "PSF"}},
		"ff2": {Default: &DefaultInfo{From: func(res *PulumiResource) (interface{}, error) { return "PF2", nil }}},
		"ggg": {Default: &DefaultInfo{Value: "PSG"}},
		"hhh": {Default: &DefaultInfo{Value: "PSH"}},
		"iii": {Default: &DefaultInfo{Value: "PSI"}},
		"mmm": {Default: &DefaultInfo{Value: "PSM"}},
		"uuu": {Default: &DefaultInfo{Value: "PSU", EnvVars: []string{"PTFU", "PTFU2"}}},
		"vvv": {Default: &DefaultInfo{Value: 42, EnvVars: []string{"PTFV", "PTFV2"}}},
		"www": {Default: &DefaultInfo{Value: "PSW"}},
		"zzz": {Asset: &AssetTranslation{Kind: FileAsset}},
	}
	olds := resource.PropertyMap{
		"iii": resource.NewStringProperty("OLI"),
		"jjj": resource.NewStringProperty("OLJ"),
		"lll": resource.NewStringProperty("OLL"),
		"mmm": resource.NewStringProperty("OLM"),
		"www": resource.NewStringProperty("OLW"),
		"xxx": resource.NewStringProperty("OLX"),
	}
	props := resource.PropertyMap{
		"bbb": resource.NewStringProperty("BBB"),
		"ddd": resource.NewStringProperty("DDD"),
		"dd2": resource.NewStringProperty("DDD"),
		"fff": resource.NewStringProperty("FFF"),
		"ff2": resource.NewStringProperty("FFF"),
		"hhh": resource.NewStringProperty("HHH"),
		"zzz": resource.NewAssetProperty(asset),
	}
	inputs, err := MakeTerraformInputs(nil, olds, props, tfs, ps, assets, true, false)
	assert.NoError(t, err)
	outputs := MakeTerraformOutputs(inputs, tfs, ps, assets, false)
	assert.Equal(t, resource.NewPropertyMapFromMap(map[string]interface{}{
		"bbb": "BBB",
		"ccc": "CCC",
		"cc2": "CC2",
		"ddd": "DDD",
		"dd2": "DDD",
		"eee": "EEE",
		"ee2": "EE2",
		"fff": "FFF",
		"ff2": "FFF",
		"ggg": "PSG",
		"hhh": "HHH",
		"iii": "OLI",
		"lll": "OLL",
		"mmm": "OLM",
		"uuu": "PSU",
		"vvv": 1337,
		"www": "OLW",
		"zzz": asset,
	}), outputs)
}

func TestComputedAsset(t *testing.T) {
	assets := make(AssetTable)
	tfs := map[string]*schema.Schema{
		"zzz": {Type: schema.TypeString},
	}
	ps := map[string]*SchemaInfo{
		"zzz": {Asset: &AssetTranslation{Kind: FileAsset}},
	}
	olds := resource.PropertyMap{}
	props := resource.PropertyMap{
		"zzz": resource.NewStringProperty(config.UnknownVariableValue),
	}
	inputs, err := MakeTerraformInputs(nil, olds, props, tfs, ps, assets, true, false)
	assert.NoError(t, err)
	outputs := MakeTerraformOutputs(inputs, tfs, ps, assets, false)
	assert.Equal(t, resource.PropertyMap{
		"zzz": resource.PropertyValue{V: resource.Computed{Element: resource.PropertyValue{V: ""}}},
	}, outputs)
}

func TestInvalidAsset(t *testing.T) {
	assets := make(AssetTable)
	tfs := map[string]*schema.Schema{
		"zzz": {Type: schema.TypeString},
	}
	ps := map[string]*SchemaInfo{
		"zzz": {Asset: &AssetTranslation{Kind: FileAsset}},
	}
	olds := resource.PropertyMap{}
	props := resource.PropertyMap{
		"zzz": resource.NewStringProperty("invalid"),
	}
	inputs, err := MakeTerraformInputs(nil, olds, props, tfs, ps, assets, true, false)
	assert.NoError(t, err)
	assert.Panics(t, func() {
		MakeTerraformOutputs(inputs, tfs, ps, assets, false)
	})
}

func boolPointer(b bool) *bool {
	return &b
}

func TestCustomTransforms(t *testing.T) {
	doc := map[string]interface{}{
		"a": 99,
		"b": false,
	}
	tfs := &schema.Schema{Type: schema.TypeString}
	psi := &SchemaInfo{Transform: TransformJSONDocument}
	v1, err := MakeTerraformInput(
		nil, "v", resource.PropertyValue{}, resource.NewObjectProperty(resource.NewPropertyMapFromMap(doc)),
		tfs, psi, nil, false, false)
	assert.NoError(t, err)
	if !assert.Equal(t, `{"a":99,"b":false}`, v1) {
		assert.Equal(t, `{"b":false,"a":99}`, v1)
	}
	v2, err := MakeTerraformInput(
		nil, "v", resource.PropertyValue{}, resource.NewStringProperty(`{"a":99,"b":false}`),
		tfs, psi, nil, false, false)
	assert.NoError(t, err)
	assert.Equal(t, `{"a":99,"b":false}`, v2)
}
