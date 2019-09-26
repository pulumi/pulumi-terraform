/**
 * This code is copied from package github.com/hashicorp/terraform/flatmap, and is distributed
 * under the terms of the Mozilla Public License v2.
 *
 * Copyright 2014-2019 HashiCorp Inc.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
 * If a copy of the MPL was not distributed with this file, You can obtain one at
 * http://mozilla.org/MPL/2.0/.
 */

package flatmap

import (
	"fmt"
	"reflect"
)

// Flatten takes a structure and turns into a flat map[string]string.
//
// Within the "thing" parameter, only primitive values are allowed. Structs are
// not supported. Therefore, it can only be slices, maps, primitives, and
// any combination of those together.
//
// See the tests for examples of what inputs are turned into.
func Flatten(thing map[string]interface{}) Map {
	result := make(map[string]string)

	for k, raw := range thing {
		flatten(result, k, reflect.ValueOf(raw))
	}

	return Map(result)
}

func flatten(result map[string]string, prefix string, v reflect.Value) {
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			result[prefix] = "true"
		} else {
			result[prefix] = "false"
		}
	case reflect.Int:
		result[prefix] = fmt.Sprintf("%d", v.Int())
	case reflect.Map:
		flattenMap(result, prefix, v)
	case reflect.Slice:
		flattenSlice(result, prefix, v)
	case reflect.String:
		result[prefix] = v.String()
	default:
		panic(fmt.Sprintf("Unknown: %s", v))
	}
}

func flattenMap(result map[string]string, prefix string, v reflect.Value) {
	for _, k := range v.MapKeys() {
		if k.Kind() == reflect.Interface {
			k = k.Elem()
		}

		if k.Kind() != reflect.String {
			panic(fmt.Sprintf("%s: map key is not string: %s", prefix, k))
		}

		flatten(result, fmt.Sprintf("%s.%s", prefix, k.String()), v.MapIndex(k))
	}
}

func flattenSlice(result map[string]string, prefix string, v reflect.Value) {
	prefix = prefix + "."

	result[prefix+"#"] = fmt.Sprintf("%d", v.Len())
	for i := 0; i < v.Len(); i++ {
		flatten(result, fmt.Sprintf("%s%d", prefix, i), v.Index(i))
	}
}
