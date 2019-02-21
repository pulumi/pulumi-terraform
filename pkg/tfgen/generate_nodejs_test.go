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

package tfgen

import (
	"testing"

	"github.com/pulumi/pulumi-terraform/pkg/tfbridge"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pulumi/pulumi/pkg/tokens"
	"github.com/stretchr/testify/assert"
)

type typeTest struct {
	schema         *schema.Schema
	info           *tfbridge.SchemaInfo
	expectedOutput string
	expectedInput  string
}

var tsTypeTests = []typeTest{
	{
		// Bool Schema
		schema:         &schema.Schema{Type: schema.TypeBool},
		expectedOutput: "boolean",
		expectedInput:  "boolean",
	},
	{
		// Int Schema
		schema:         &schema.Schema{Type: schema.TypeInt},
		expectedOutput: "number",
		expectedInput:  "number",
	},
	{
		// Float Schema
		schema:         &schema.Schema{Type: schema.TypeFloat},
		expectedOutput: "number",
		expectedInput:  "number",
	},
	{
		// String Schema
		schema:         &schema.Schema{Type: schema.TypeString},
		expectedOutput: "string",
		expectedInput:  "string",
	},
	{
		// Basic Set Schema
		schema: &schema.Schema{
			Type: schema.TypeSet,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
		expectedOutput: "string[]",
		expectedInput:  "string[]",
	},
	{
		// Basic List Schema
		schema: &schema.Schema{
			Type: schema.TypeList,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
		expectedOutput: "string[]",
		expectedInput:  "string[]",
	},
	{
		// Basic Map Schema
		schema: &schema.Schema{
			Type: schema.TypeMap,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
		expectedOutput: "{[key: string]: string}",
		expectedInput:  "{[key: string]: string}",
	},
	{
		// Resource Map Schema
		schema: &schema.Schema{
			Type: schema.TypeMap,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"foo": {Type: schema.TypeString},
				},
			},
		},
		expectedOutput: "{ foo: string }",
		expectedInput:  "{ foo: string }",
	},
	{
		// Basic alt types
		info: &tfbridge.SchemaInfo{
			Type:     "string",
			AltTypes: []tokens.Type{"Foo"},
		},
		expectedOutput: "string",
		expectedInput:  "string | Foo",
	},
}

func Test_TsTypes(t *testing.T) {
	for _, test := range tsTypeTests {
		v := &variable{
			name:   "foo",
			schema: test.schema,
			info:   test.info,
			opt:    true,
		}

		// Output
		v.out = true
		assert.Equal(t, test.expectedOutput, tsType(v, false))

		// Input
		v.out = false
		assert.Equal(t, test.expectedInput, tsType(v, false))
	}
}

func Test_Issue130(t *testing.T) {
	schema := &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}

	assert.Equal(t, "string", tsType(&variable{
		name:   "condition",
		schema: schema,
		out:    true,
	}, false))

	assert.Equal(t, "string", tsType(&variable{
		name:   "condition",
		schema: schema,
		out:    false,
	}, false))
}
