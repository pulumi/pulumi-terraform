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
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gedex/inflector"

	"github.com/golang/glog"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"

	"github.com/pulumi/pulumi/pkg/tools"
	"github.com/pulumi/pulumi/pkg/util/contract"

	"github.com/pulumi/pulumi-terraform/pkg/tfbridge"
)

// newDotnetGenerator returns a language generator that understands how to produce C# packages.
func newDotnetGenerator(pkg, version string, info tfbridge.ProviderInfo, overlaysDir, outDir string) langGenerator {
	return &dotnetGenerator{
		pkg:         pkg,
		version:     version,
		info:        info,
		overlaysDir: overlaysDir,
		outDir:      outDir,
	}
}

type dotnetGenerator struct {
	pkg         string
	version     string
	info        tfbridge.ProviderInfo
	overlaysDir string
	outDir      string
}

// commentChars returns the comment characters to use for single-line comments.
func (g *dotnetGenerator) commentChars() string {
	return "//"
}

// moduleDir returns the directory for the given module.
func (g *dotnetGenerator) moduleDir(mod *module) string {
	dir := g.outDir
	if mod != nil && mod.name != "" {
		dir = filepath.Join(dir, csName(mod.name))
	}
	return dir
}

// openWriter opens a writer for the given module and file name, emitting the standard header automatically.
func (g *dotnetGenerator) openWriter(mod *module, name string, needsSDK bool) (*tools.GenWriter, error) {
	dir := g.moduleDir(mod)
	file := filepath.Join(dir, name)
	w, err := tools.NewGenWriter(tfgen, file)
	if err != nil {
		return nil, err
	}

	// Emit a standard warning header ("do not edit", etc).
	w.EmitHeaderWarning(g.commentChars())

	// If needed, emit the standard Pulumi SDK import statement.
	if needsSDK {
		g.emitSDKImport(mod, w)
	}

	return w, nil
}

func (g *dotnetGenerator) emitSDKImport(mod *module, w *tools.GenWriter) {
	w.Writefmtln("using System;")
	w.Writefmtln("using System.Collections.Generic;")
	w.Writefmtln("using System.Linq;")
	w.Writefmtln("")
}

func (g *dotnetGenerator) emitPackage(pack *pkg) error {
	namespace := "Pulumi." + csName(pack.name)
	err := g.emitModules(pack.modules, namespace)
	if err != nil {
		return err
	}

	if err := g.emitUtilities(namespace); err != nil {
		return err
	}

	return g.emitPackageMetadata(pack)
}

func (g *dotnetGenerator) emitUtilities(namespace string) error {
	w, err := g.openWriter(nil, "Utilities.cs", false)
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(w)

	w.Writefmt(dotnetUtilitiesFile, namespace)
	return nil
}

func (g *dotnetGenerator) emitModules(mmap moduleMap, namespace string) error {
	for _, mod := range mmap.values() {
		if err := g.emitModule(mod, namespace); err != nil {
			return nil
		}
	}
	return nil
}

func (g *dotnetGenerator) emitModule(mod *module, namespace string) error {
	glog.V(3).Infof("emitModule(%s)", mod.name)

	dir := g.moduleDir(mod)
	if err := tools.EnsureDir(dir); err != nil {
		return errors.Wrapf(err, "creating module directory")
	}

	namespace = fmt.Sprintf("%s.%s", namespace, csName(mod.name))

	for _, member := range mod.members {
		err := g.emitModuleMember(mod, namespace, member)
		if err != nil {
			return errors.Wrapf(err, "emitting module %s member %s", mod.name, member.Name())
		}
	}

	return nil
}

func (g *dotnetGenerator) emitModuleMember(mod *module, namespace string, member moduleMember) error {
	glog.V(3).Infof("emitModuleMember(%s, %s)", mod, member.Name())

	switch t := member.(type) {
	case *resourceType:
		return g.emitResourceType(mod, namespace, t)
	case *resourceFunc:
		return g.emitResourceFunc(mod, namespace, t)
	case *variable:
		contract.Assertf(mod.config(),
			"only expected top-level variables in config module (%s is not one)", mod.name)
		// skip the variable, we will process it later.
		return nil
	case *overlayFile:
		return g.emitOverlay(mod, t)
	default:
		contract.Failf("unexpected member type: %v", reflect.TypeOf(member))
		return nil
	}
}

func (g *dotnetGenerator) emitDocComment(w *tools.GenWriter, comment, prefix string) {
	if comment != "" {
		var written bool
		lines := strings.Split(comment, "\n")
		for i, docLine := range lines {
			// Break if we get to the last line and it's empty
			if i == len(lines)-1 && strings.TrimSpace(docLine) == "" {
				break
			}

			// If the first line, start a doc comment.
			if i == 0 {
				w.Writefmtln(`%s/// <summary>`, prefix)
				written = true
			}

			// Print the line of documentation
			w.Writefmtln("%s/// %s", prefix, docLine)
		}
		if written {
			w.Writefmtln(`%s/// </summary>`, prefix)
		}
	}
}

func (g *dotnetGenerator) emitRawDocComment(w *tools.GenWriter, comment, prefix string) {
	if comment != "" {
		curr := 0
		for _, word := range strings.Fields(comment) {
			if curr > 0 {
				if curr+len(word)+1 > (maxWidth - len(prefix)) {
					curr = 0
					w.Writefmtln("")
					w.Writefmt("%s/// ", prefix)
				} else {
					w.Writefmt(" ")
					curr++
				}
			} else {
				w.Writefmtln(`%s/// <summary>`, prefix)
				w.Writefmt("%s/// ", prefix)
			}
			w.Writefmt(word)
			curr += len(word)
		}
		w.Writefmtln("")
		w.Writefmtln(`%s/// </summary>`, prefix)
	}
}

func (g *dotnetGenerator) emitDoc(w *tools.GenWriter, doc, rawdoc, prefix string) {
	if doc != "" {
		g.emitDocComment(w, doc, prefix)
	} else if rawdoc != "" {
		g.emitRawDocComment(w, rawdoc, prefix)
	}
}

func (g *dotnetGenerator) emitSubstructures(w *tools.GenWriter, class, key string, sch *schema.Schema, input, io bool) {
	switch sch.Type {
	case schema.TypeString:
	case schema.TypeInt:
	case schema.TypeFloat:
	case schema.TypeBool:
		return
	}

	switch e := sch.Elem.(type) {
	case *schema.Schema:
		g.emitSubstructures(w, class, key, e, input, io)
	case *schema.Resource:
		name := csStructureName(class, key, input)
		interfaceType := "Pulumi.IProtobuf"
		if io {
			interfaceType = "Pulumi.IIOProtobuf"
		}
		w.Writefmtln("	public sealed class %s : %s {", name, interfaceType)
		for _, s := range stableSchemas(e.Schema) {
			sch := e.Schema[s]
			w.Writefmtln("		public %s %s { get; set; }", csType(class, s, sch, input, io), csName(s))
		}
		w.Writefmtln("")

		valueType := "Google.Protobuf.WellKnownTypes.Value"
		if io {
			valueType = fmt.Sprintf("Pulumi.IO<%s>", valueType)
		}
		w.Writefmtln("		public %s ToProtobuf() {", valueType)
		w.Writefmtln("			return Pulumi.Protobuf.ToProtobuf(")
		for i, s := range stableSchemas(e.Schema) {
			if i != 0 {
				w.Writefmtln(",")
			}
			sch := e.Schema[s]
			expr := g.exprToProtobuf(csName(s), sch, 0)
			w.Writefmt("				new KeyValuePair<string, %s>(\"%s\", %s)", valueType, s, expr)
		}
		w.Writefmtln(");")
		w.Writefmtln("		} // ToProtobuf")
		w.Writefmtln("")

		if !input {
			w.Writefmtln("		public static %s FromProtobuf(Google.Protobuf.WellKnownTypes.Value value) {", name)
			w.Writefmtln("			var obj = value.StructValue;")
			w.Writefmtln("			return new %s() {", name)
			for _, s := range stableSchemas(e.Schema) {
				sch := e.Schema[s]
				expr := g.exprFromProtobuf(class, key, fmt.Sprintf("obj.Fields[\"%s\"]", s), sch, 0)
				w.Writefmtln("				%s = %s,", csName(s), expr)
			}
			w.Writefmtln("			};")
			w.Writefmtln("		} // FromProtobuf")
			w.Writefmtln("")
		}

		w.Writefmtln("	} // %s", name)
		w.Writefmtln("")
	}
}

func (g *dotnetGenerator) emitResourceType(mod *module, namespace string, res *resourceType) error {
	// Create a resource file for this resource's module.
	name := csName(res.name)
	w, err := g.openWriter(mod, name+".cs", true)
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(w)

	w.Writefmtln("namespace %s {", namespace)
	w.Writefmtln("")

	// Emit sub-structures
	for _, prop := range res.inprops {
		g.emitSubstructures(w, name, prop.name, prop.schema, true, true)
	}
	for _, prop := range res.outprops {
		g.emitSubstructures(w, name, prop.name, prop.schema, false, true)
	}

	// Emit the argument struct
	args := ""
	if res.argst != nil {
		args = g.emitStruct(w, name, res.argst, true, true, "	")
	}

	g.emitDoc(w, res.doc, "", "	")
	w.Writefmtln("	public class %s : Pulumi.CustomResource {", name)

	// Emit the properties
	g.emitProperties(w, name, res.outprops, false, true, "		")

	// Emit the ctor
	w.Writefmt("		public %s(string name", name)
	if args != "" {
		w.Writefmt(", %s args, ", args)
	}
	w.Writefmtln("Pulumi.ResourceOptions opts = default(Pulumi.ResourceOptions))")
	if args != "" {
		w.Writefmtln("			: base(\"%s\", name, SerialiseArgs(args), opts) {", res.info.Tok)
	} else {
		w.Writefmtln("			: base(\"%s\", name, null, opts) {", res.info.Tok)
	}

	// Emit the output transformers
	for _, arg := range res.inprops {
		expr := g.exprFromProtobuf(name, csName(arg.name), "item", arg.schema, 1)
		w.Writefmtln("			%s = Outputs[\"%s\"].Select(item => %s);", csName(arg.name), arg.name, expr)
	}
	for _, arg := range res.outprops {
		expr := g.exprFromProtobuf(name, csName(arg.name), "item", arg.schema, 1)
		w.Writefmtln("			%s = Outputs[\"%s\"].Select(item => %s);", csName(arg.name), arg.name, expr)
	}

	w.Writefmtln("		} // ctor")
	w.Writefmtln("")

	if args != "" {
		// Transform the resource arguments into Values.
		w.Writefmt("		private static Dictionary")
		w.Writefmt("<string, Pulumi.IO<Google.Protobuf.WellKnownTypes.Value>>")
		w.Writefmtln(" SerialiseArgs(%s args) {", args)
		w.Writefmtln("			var props = new Dictionary<string, Pulumi.IO<Google.Protobuf.WellKnownTypes.Value>>();")
		ins := make(map[string]bool)
		for _, prop := range res.inprops {
			ins[prop.name] = true
			expr := g.exprToProtobuf(fmt.Sprintf("args.%s", csName(prop.name)), prop.schema, 0)
			w.Writefmtln("			props[\"%s\"] = %s;", prop.name, expr)
		}
		for _, prop := range res.outprops {
			if !ins[prop.name] {
				w.Writefmtln("			props[\"%s\"] = null; //out", prop.name)
			}
		}
		w.Writefmtln("			return props;")
		w.Writefmtln("		} // SerialiseArgs")
		w.Writefmtln("")
	}

	w.Writefmtln("	} // %s", name)
	w.Writefmtln("} // %s", namespace)
	return nil
}

func (g *dotnetGenerator) emitProperties(w *tools.GenWriter,
	class string,
	props []*variable,
	input, io bool,
	prefix string) {

	for _, prop := range props {
		g.emitDoc(w, prop.doc, prop.rawdoc, prefix)
		w.Writefmt(prefix)
		typ := csType(class, prop.name, prop.schema, input, io)
		name := csName(prop.name)
		w.Writefmtln("public %s %s { get; set; }", typ, name)
		w.Writefmtln("")
	}
}

func (g *dotnetGenerator) emitStruct(w *tools.GenWriter,
	class string,
	argst *plainOldType,
	input, io bool,
	prefix string) string {

	g.emitDoc(w, argst.doc, "", prefix)

	name := csName(argst.name)
	w.Writefmtln("%spublic struct %s {", prefix, name)

	g.emitProperties(w, class, argst.props, input, io, prefix+"	")

	w.Writefmtln("%s} // %s", prefix, name)
	w.Writefmtln("")
	return name
}

func (g *dotnetGenerator) exprFromProtobuf(class, key, expr string, sch *schema.Schema, depth int) string {
	item := "item"
	if depth != 0 {
		item = fmt.Sprintf("item%d", depth)
	}

	switch sch.Type {
	case schema.TypeString:
		expr = fmt.Sprintf("Protobuf.ToString(%s)", expr)
	case schema.TypeInt:
		expr = fmt.Sprintf("Protobuf.ToInt(%s)", expr)
	case schema.TypeFloat:
		expr = fmt.Sprintf("Protobuf.ToDouble(%s)", expr)
	case schema.TypeBool:
		expr = fmt.Sprintf("Protobuf.ToBool(%s)", expr)
	case schema.TypeList:
		if _, ok := sch.Elem.(*schema.Resource); ok {
			name := csStructureName(class, key, false)
			if tfbridge.IsMaxItemsOne(sch, nil) {
				expr = fmt.Sprintf("%s.FromProtobuf(%s)", name, expr)
			} else {
				expr = fmt.Sprintf("Protobuf.ToList(%s, %s => %s.FromProtobuf(%s))", expr, item, name, item)
			}
		} else if elemSchema, ok := sch.Elem.(*schema.Schema); ok {
			if tfbridge.IsMaxItemsOne(sch, nil) {
				expr = g.exprFromProtobuf(class, key, expr, elemSchema, depth)
			} else {
				subexpr := g.exprFromProtobuf(class, key, item, elemSchema, depth+1)
				expr = fmt.Sprintf("Protobuf.ToList(%s, %s => %s)", expr, item, subexpr)
			}
		}
	case schema.TypeMap:
		expr = fmt.Sprintf("Protobuf.ToMap(%s)", expr)
	case schema.TypeSet:
		if _, ok := sch.Elem.(*schema.Resource); ok {
			name := csStructureName(class, key, false)
			if tfbridge.IsMaxItemsOne(sch, nil) {
				expr = fmt.Sprintf("%s.FromProtobuf(%s)", name, expr)
			} else {
				expr = fmt.Sprintf("Protobuf.ToList(%s, %s => %s.FromProtobuf(%s))", expr, item, name, item)
			}
		} else if elemSchema, ok := sch.Elem.(*schema.Schema); ok {
			if tfbridge.IsMaxItemsOne(sch, nil) {
				expr = g.exprFromProtobuf(class, key, expr, elemSchema, depth)
			} else {
				subexpr := g.exprFromProtobuf(class, key, item, elemSchema, depth+1)
				expr = fmt.Sprintf("Protobuf.ToList(%s, %s => %s)", expr, item, subexpr)
			}
		}
	}
	return expr
}

func (g *dotnetGenerator) exprToProtobuf(expr string, sch *schema.Schema, depth int) string {
	item := "item"
	if depth != 0 {
		item = fmt.Sprintf("item%d", depth)
	}

	switch sch.Type {
	case schema.TypeString:
		expr = fmt.Sprintf("Protobuf.ToProtobuf(%s)", expr)
	case schema.TypeInt:
		expr = fmt.Sprintf("Protobuf.ToProtobuf(%s)", expr)
	case schema.TypeFloat:
		expr = fmt.Sprintf("Protobuf.ToProtobuf(%s)", expr)
	case schema.TypeBool:
		expr = fmt.Sprintf("Protobuf.ToProtobuf(%s)", expr)
	case schema.TypeList:
		if _, ok := sch.Elem.(*schema.Resource); ok {
			if tfbridge.IsMaxItemsOne(sch, nil) {
				expr = fmt.Sprintf("Protobuf.ToProtobuf(%s)", expr)
			} else {
				expr = fmt.Sprintf("Protobuf.ToProtobuf(%s, %s => Protobuf.ToProtobuf(%s))", expr, item, item)
			}
		} else if elemSchema, ok := sch.Elem.(*schema.Schema); ok {
			if tfbridge.IsMaxItemsOne(sch, nil) {
				expr = g.exprToProtobuf(expr, elemSchema, depth)
			} else {
				subexpr := g.exprToProtobuf(item, elemSchema, depth+1)
				expr = fmt.Sprintf("Protobuf.ToProtobuf(%s, %s => %s)", expr, item, subexpr)
			}
		}
	case schema.TypeMap:
		expr = fmt.Sprintf("Protobuf.ToProtobuf(%s)", expr)
	case schema.TypeSet:
		if _, ok := sch.Elem.(*schema.Resource); ok {
			if tfbridge.IsMaxItemsOne(sch, nil) {
				expr = fmt.Sprintf("Protobuf.ToProtobuf(%s)", expr)
			} else {
				expr = fmt.Sprintf("Protobuf.ToProtobuf(%s, %s => Protobuf.ToProtobuf(%s))", expr, item, item)
			}
		} else if elemSchema, ok := sch.Elem.(*schema.Schema); ok {
			if tfbridge.IsMaxItemsOne(sch, nil) {
				expr = g.exprToProtobuf(expr, elemSchema, depth)
			} else {
				subexpr := g.exprToProtobuf(item, elemSchema, depth+1)
				expr = fmt.Sprintf("Protobuf.ToProtobuf(%s, %s => %s)", expr, item, subexpr)
			}
		}
	}

	return expr
}

func (g *dotnetGenerator) emitResourceFunc(mod *module, namespace string, fun *resourceFunc) error {
	name := csName(fun.name)
	w, err := g.openWriter(mod, name+".cs", true)
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(w)

	w.Writefmtln("namespace %s {", namespace)
	w.Writefmtln("")

	// Emit sub-structures
	for _, prop := range fun.args {
		g.emitSubstructures(w, name, prop.name, prop.schema, true, false)
	}
	for _, prop := range fun.rets {
		g.emitSubstructures(w, name, prop.name, prop.schema, false, false)
	}

	// Emit the argument struct, if needed
	args := ""
	if fun.argst != nil {
		args = g.emitStruct(w, name, fun.argst, true, false, "	")
	}

	// Emit the result struct, if needed
	rets := ""
	if fun.retst != nil {
		rets = g.emitStruct(w, name, fun.retst, false, false, "	")
	}

	modName := csName(mod.name) + "Module"
	w.Writefmtln("	public static partial class %s {", modName)

	g.emitDoc(w, fun.doc, "", "		")
	if rets != "" {
		w.Writefmt("		public static System.Threading.Tasks.Task<%s>", rets)
	} else {
		w.Writefmt("		public static System.Threading.Tasks.Task")
	}
	w.Writefmt(" %s(", name)
	if args != "" {
		w.Writefmt("%s args, ", args)
	}
	w.Writefmtln("Pulumi.InvokeOptions opts = default(Pulumi.InvokeOptions)) {")

	// Copy the function arguments into a Struct.
	w.Writefmtln("			var invokeArgs = new Google.Protobuf.WellKnownTypes.Struct();")
	for _, arg := range fun.args {
		expr := g.exprToProtobuf(fmt.Sprintf("args.%s", csName(arg.name)), arg.schema, 0)
		w.Writefmt("			")
		w.Writefmtln("invokeArgs.Fields[\"%s\"] = %s;", arg.name, expr)

	}
	w.Writefmtln("")

	// Now simply invoke the runtime function with the arguments.
	w.Writefmtln("			var task = Pulumi.Runtime.InvokeAsync(\"%s\", invokeArgs, opts);", fun.info.Tok)
	w.Writefmtln("")

	// And copy the results to an object
	w.Writefmtln("			return task.ContinueWith(response => {")
	w.Writefmtln("				var protobuf = response.Result;")
	w.Writefmtln("				var result = new %s();", rets)
	for _, prop := range fun.rets {
		expr := g.exprFromProtobuf(name, csName(prop.name), fmt.Sprintf("protobuf.Fields[\"%s\"]", prop.name), prop.schema, 0)

		// Now perform the assignment
		w.Writefmtln("				if (protobuf.Fields.ContainsKey(\"%s\")) {", prop.name)
		w.Writefmtln("					result.%s = %s;", csName(prop.name), expr)
		w.Writefmtln("				}")
	}
	w.Writefmtln("				return result;")
	w.Writefmtln("			});")
	w.Writefmtln("		} // %s", name)
	w.Writefmtln("")
	w.Writefmtln("	} // %s", modName)
	w.Writefmtln("} // %s", namespace)

	return nil
}

func (g *dotnetGenerator) emitOverlay(mod *module, overlay *overlayFile) error {
	// Copy the file from the overlays directory to the destination.
	dir := g.moduleDir(mod)
	dst := filepath.Join(dir, overlay.name)
	if err := copyFile(overlay.src, dst); err != nil {
		return err
	}
	return nil
}

// emitPackageMetadata generates all the non-code metadata required by a Pulumi package.
func (g *dotnetGenerator) emitPackageMetadata(pack *pkg) error {
	// The generator already emitted Pulumi.yaml, so that just leaves the `Pulumi.csproj`.

	projectName := "Pulumi." + csName(pack.name) + ".csproj"
	w, err := tools.NewGenWriter(tfgen, filepath.Join(g.outDir, projectName))
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(w)

	// Emit a standard warning header ("do not edit", etc).
	w.Writefmtln("<!--")
	w.EmitHeaderWarning("")
	w.Writefmtln("-->")

	// First remove the 'v'
	version := g.version[1:]
	// Now split any pre-release suffix off
	parts := strings.Split(version, "-")
	// If we're a pre-release version add a pre-release wildcard suffix
	if len(parts) > 1 {
		version = parts[0] + "-*"
	}

	// Now create a standard csproj package from the metadata.
	w.Writefmtln("<Project Sdk=\"Microsoft.NET.Sdk\">")
	w.Writefmtln("	<PropertyGroup>")
	w.Writefmtln("		<TargetFramework>netstandard2.0</TargetFramework>")
	w.Writefmtln("	</PropertyGroup>")
	w.Writefmtln("	<ItemGroup>")
	//w.Writefmtln("		<ProjectReference Include=\"../../../pulumi/sdk/dotnet/Pulumi/Pulumi.csproj\" />")
	w.Writefmtln("		<PackageReference Include=\"Pulumi\" Version=\"%s\"/>", version)
	w.Writefmtln("	</ItemGroup>")
	w.Writefmtln("</Project>")

	return nil
}

func csType(class, key string, sch *schema.Schema, input, io bool) string {

	typeExpr := ""
	switch sch.Type {
	case schema.TypeBool:
		typeExpr = fmt.Sprintf("bool")
	case schema.TypeInt:
		typeExpr = fmt.Sprintf("int")
	case schema.TypeFloat:
		typeExpr = fmt.Sprintf("double")
	case schema.TypeString:
		typeExpr = fmt.Sprintf("string")
	case schema.TypeList:
		subtype := ""
		if _, ok := sch.Elem.(*schema.Resource); ok {
			subtype = csStructureName(class, key, input)
			if input && io && !tfbridge.IsMaxItemsOne(sch, nil) {
				subtype = fmt.Sprintf("Pulumi.IO<%s>", subtype)
			}
		} else if elemSchema, ok := sch.Elem.(*schema.Schema); ok {
			subtype = csType(class, key, elemSchema, input, io && input)
		}
		if tfbridge.IsMaxItemsOne(sch, nil) {
			typeExpr = subtype
		} else {
			typeExpr = fmt.Sprintf("%s[]", subtype)
		}
	case schema.TypeMap:
		if subSchema, ok := sch.Elem.(*schema.Schema); ok {
			if subSchema.Type == schema.TypeString {
				typeExpr = "System.Collections.Generic.Dictionary<string, string>"
			} else {
				panic(fmt.Sprintf("%s was a TypeMap but it's subschema wasn't a string: %+v", key, sch))
			}
		} else if _, ok := sch.Elem.(*schema.Resource); ok {
			subtype := csStructureName(class, key, input)
			typeExpr = fmt.Sprintf("System.Collections.Generic.Dictionary<string, %s>", subtype)
		} else if sch.Elem == nil {
			typeExpr = "System.Collections.Generic.Dictionary<string, string>"			
		} else {
			panic(fmt.Sprintf("%s was a TypeMap but it's Elem field wasn't understood:\n%+v\n%+v", key, sch, sch.Elem))
		}
	case schema.TypeSet:
		subtype := ""
		if _, ok := sch.Elem.(*schema.Resource); ok {
			subtype = csStructureName(class, key, input)
			if input && io && !tfbridge.IsMaxItemsOne(sch, nil) {
				subtype = fmt.Sprintf("Pulumi.IO<%s>", subtype)
			}
		} else if elemSchema, ok := sch.Elem.(*schema.Schema); ok {
			subtype = csType(class, key, elemSchema, input, io && input)
		}
		if tfbridge.IsMaxItemsOne(sch, nil) {
			typeExpr = subtype
		} else {
			typeExpr = fmt.Sprintf("%s[]", subtype)
		}
	}

	if io {
		return fmt.Sprintf("Pulumi.IO<%s>", typeExpr)
	}
	return typeExpr
}

func csName(name string) string {
	return tfbridge.TerraformToPulumiName(name, nil, true)
}

// Used for substructure names. We only have the terrafrom key to look at and it's nearly always a plural form
func csStructureName(class, name string, input bool) string {
	argSuffix := ""
	if input {
		argSuffix = "Args"
	}
	name = fmt.Sprintf("%s%s%s", class, argSuffix, tfbridge.TerraformToPulumiName(name, nil, true))
	return inflector.Singularize(name)
}
