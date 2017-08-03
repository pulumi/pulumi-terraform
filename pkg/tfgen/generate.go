// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfgen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-fabric/pkg/diag"
	"github.com/pulumi/pulumi-fabric/pkg/tokens"
	"github.com/pulumi/pulumi-fabric/pkg/tools"
	"github.com/pulumi/pulumi-fabric/pkg/util/cmdutil"
	"github.com/pulumi/pulumi-fabric/pkg/util/contract"

	"github.com/pulumi/pulumi-terraform/pkg/tfbridge"
)

type generator struct {
}

func newGenerator() *generator {
	return &generator{}
}

const (
	tfgen              = "the Lumi Terraform Bridge (TFGEN) Tool"
	defaultOutDir      = "pack/"
	defaultOverlaysDir = "overlays/"
	maxWidth           = 120 // the ideal maximum width of the generated file.
)

// Generate creates Lumi packages out of one or more Terraform plugins.  It accepts a list of all of the input Terraform
// providers, already bound statically to the code (since we cannot obtain schema information dynamically), walks them
// and generates the Lumi code, and spews that code into the output directory.
func (g *generator) Generate(pkg string, provinfo tfbridge.ProviderInfo, outDir, overlaysDir string) error {
	// If outDir or overlaysDir are empty, default to pack/ in the pwd.
	if outDir == "" || overlaysDir == "" {
		p, err := os.Getwd()
		if err != nil {
			return err
		}
		if outDir == "" {
			outDir = filepath.Join(p, defaultOutDir)
		}
		if overlaysDir == "" {
			overlaysDir = filepath.Join(defaultOverlaysDir)
		}
	}

	// Now generate the provider code.
	return g.generateProvider(pkg, provinfo, outDir, overlaysDir)
}

// generateProvider creates a single standalone Lumi package for the given provider.
func (g *generator) generateProvider(pkg string, provinfo tfbridge.ProviderInfo, outDir, overlaysDir string) error {
	var files []string
	exports := make(map[string]string)               // a list of top-level exports.
	modules := make(map[string]string)               // a list of modules to export individually.
	submodules := make(map[string]map[string]string) // a map of sub-module name to exported members.

	// Ensure the output path exists.
	if err := tools.EnsureDir(outDir); err != nil {
		return err
	}

	// Place all configuration variables into a single config module.
	prov := provinfo.P
	if len(prov.Schema) > 0 {
		cfgfile, err := g.generateConfig(pkg, prov.Schema, provinfo.Config, outDir)
		if err != nil {
			return err
		}
		// ensure we export the config submodule and add its file to the project.
		submodules["config"] = map[string]string{
			"vars": cfgfile,
		}
		files = append(files, cfgfile)
	}

	// For each resource, create its own dedicated type and module export.
	var reserr error
	resmap := prov.ResourcesMap
	reshits := make(map[string]bool)
	for _, r := range stableResources(prov.ResourcesMap) {
		var resinfo tfbridge.ResourceInfo
		if provinfo.Resources != nil {
			if ri, has := provinfo.Resources[r]; has {
				resinfo = ri
				reshits[r] = true
			} else {
				// if this has a map, but this resource wasn't found, issue a warning.
				cmdutil.Diag().Warningf(
					diag.Message("Resource %v not found in provider map; using default naming"), r)
			}
		}
		result, err := g.generateResource(pkg, r, resmap[r], resinfo, outDir, outDir)
		if err != nil {
			// Keep track of the error, but keep going, so we can expose more at once.
			reserr = multierror.Append(reserr, err)
		}
		if result.Submod == "" {
			// if no sub-module, export flatly in our own index.
			exports[result.Name] = result.File
		} else {
			// otherwise, make sure to track this in the submodule so we can create and export it correctly.
			submod := submodules[result.Submod]
			if submod == nil {
				submod = make(map[string]string)
				submodules[result.Submod] = submod
			}
			submod[result.Name] = result.File
		}
		files = append(files, result.File)
	}
	if reserr != nil {
		return reserr
	}

	// Emit a warning if there is a map but some names didn't match.
	if provinfo.Resources != nil {
		var resnames []string
		for resname := range provinfo.Resources {
			resnames = append(resnames, resname)
		}
		sort.Strings(resnames)
		for _, resname := range resnames {
			if !reshits[resname] {
				cmdutil.Diag().Warningf(
					diag.Message("Resource %v (%v) wasn't found in the Terraform module; possible name mismatch?"),
					resname, provinfo.Resources[resname].Tok)
			}
		}
	}

	// Now go ahead and merge in any overlays into the modules if there are any.
	for _, overfile := range provinfo.Overlay.Files {
		// Copy the file into its place, and add it to the export and files list.
		from := filepath.Join(overlaysDir, overfile)
		to := filepath.Join(outDir, overfile)
		overname := removeExtension(to, ".ts")
		if _, has := exports[overname]; has {
			return errors.Errorf("Overlay file %v conflicts with a generated file", to)
		}
		if err := copyFile(from, to); err != nil {
			return err
		}
		exports[overname] = to
		files = append(files, to)
	}

	// Generate any submodules and add them to the export list.
	subs, extrafs, err := g.generateSubmodules(submodules, provinfo.Overlay.Modules, outDir, overlaysDir)
	if err != nil {
		return err
	}
	var subnames []string
	for sub := range subs {
		subnames = append(subnames, sub)
	}
	sort.Strings(subnames)
	for _, sub := range subnames {
		subf := subs[sub]
		if conflict, has := modules[sub]; has {
			cmdutil.Diag().Errorf(
				diag.Message("Conflicting submodule %v; exists for both %v and %v"), sub, conflict, subf)
		}
		modules[sub] = subf
		files = append(files, subf)
	}
	files = append(files, extrafs...)

	// Generate the index.ts file that reexports everything at the entrypoint.
	ixfile, err := g.generateIndex(exports, modules, outDir)
	if err != nil {
		return err
	}
	files = append(files, ixfile)

	// Generate all of the package metadata: Lumi.yaml, package.json, and tsconfig.json.
	if err := g.generatePackageMetadata(pkg, files, outDir); err != nil {
		return err
	}

	// Finally, emit the version information in a special VERSION file so we know where it came from.
	versionInfo := fmt.Sprintf("Generated by %s from:\n", os.Args[0])
	versionInfo += fmt.Sprintf("Repo: %v\n", provinfo.Git.Repo)
	if provinfo.Git.Tag != "" {
		versionInfo += fmt.Sprintf("Tag: %v\n", provinfo.Git.Tag)
	}
	if provinfo.Git.Commit != "" {
		versionInfo += fmt.Sprintf("Commit: %v\n", provinfo.Git.Commit)
	}
	versionInfo += "\n"
	return ioutil.WriteFile(filepath.Join(outDir, "VERSION"), []byte(versionInfo), 0600)
}

// copyFile is a stupid file copy routine.  It reads the file into memory to avoid messy OS-specific oddities.
func copyFile(from, to string) error {
	body, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(to, body, 0600)
}

// generateConfig takes a map of config variables and emits a config submodule to the given file.
func (g *generator) generateConfig(pkg string, cfg map[string]*schema.Schema,
	custom map[string]tfbridge.SchemaInfo, outDir string) (string, error) {
	// Sort the config variables to ensure they are emitted in a deterministic order.
	var cfgkeys []string
	for key := range cfg {
		cfgkeys = append(cfgkeys, key)
	}
	sort.Strings(cfgkeys)

	// Place a vars.ts file underneath the config/ submodule directory.
	confDir := filepath.Join(outDir, "config")

	// Ensure the config subdirectory exists.
	if err := tools.EnsureDir(confDir); err != nil {
		return "", err
	}

	// Open up the file and spew a standard "code-generated" warning header.
	file := filepath.Join(confDir, "vars.ts")
	w, err := tools.NewGenWriter(tfgen, file)
	if err != nil {
		return "", err
	}
	defer contract.IgnoreClose(w)
	w.EmitHeaderWarning()

	// First look for any custom types that will require any imports.
	if err := generateCustomImports(w, custom, pkg, outDir, confDir); err != nil {
		return "", err
	}

	// Now just emit a simple export for each variable.
	for _, key := range cfgkeys {
		// Generate a name and type to use for this key.
		prop, typ, err := g.propTyp("", key, cfg, custom[key], false)
		if err != nil {
			return "", err
		} else if prop != "" {
			g.generateComment(w, cfg[key].Description, "")
			w.Writefmtln("export let %v: %v;", prop, typ)
		}
	}
	w.Writefmtln("")

	// Ensure there weren't any custom fields that were unrecognized.
	for key := range custom {
		if _, has := cfg[key]; !has {
			cmdutil.Diag().Warningf(
				diag.Message("Custom config schema %v was not present in the Terraform metadata"), key)
		}
	}

	return file, nil
}

type resourceResult struct {
	Name   string // the resource name.
	File   string // the resource filename.
	Submod string // the submodule name, if any.
}

// generateResource generates a single module for the given resource.
func (g *generator) generateResource(pkg string, rawname string,
	res *schema.Resource, resinfo tfbridge.ResourceInfo, root, outDir string) (resourceResult, error) {
	// Transform the name as necessary.
	resname, filename := resourceName(pkg, rawname, resinfo)

	// Make a fully qualified file path that we will write to.
	file := filepath.Join(outDir, filename+".ts")

	// If the filename contains slashes, it is a sub-module, and we must ensure it exists.
	var submod string
	if slix := strings.Index(filename, "/"); slix != -1 {
		// Extract the module and file parts.
		submod = filename[:slix]
		if strings.Contains(filename[slix+1:], "/") {
			return resourceResult{},
				errors.Errorf("Modules nested more than one level deep not currently supported")
		}

		// Ensure the submodule directory exists.
		if err := tools.EnsureFileDir(file); err != nil {
			return resourceResult{}, err
		}
	}

	// Open up the file and spew a standard "code-generated" warning header.
	w, err := tools.NewGenWriter(tfgen, file)
	if err != nil {
		return resourceResult{}, err
	}
	defer contract.IgnoreClose(w)
	w.EmitHeaderWarning()

	// Now import the modules we need.
	w.Writefmtln("import * as lumi from \"@lumi/lumi\";")
	w.Writefmtln("import * as lumirt from \"@lumi/lumirt\";")
	w.Writefmtln("")

	// If there are imports required due to the custom schema info, emit them now.
	custom := resinfo.Fields
	if err := generateCustomImports(w, custom, pkg, outDir, filepath.Dir(file)); err != nil {
		return resourceResult{}, err
	}

	// Generate the resource class.
	var base string
	customNamed := len(resinfo.NameFields) > 0
	if customNamed {
		base = "lumi.Resource"
	} else {
		base = "lumi.NamedResource"
	}
	w.Writefmtln("export class %[1]v extends %v implements %[1]vArgs {", resname, base)

	// First, generate all instance properties.
	var finalerr error
	var inprops []string
	var inflags []string
	var intypes []string
	var schemas []*schema.Schema
	var customs []tfbridge.SchemaInfo
	if len(res.Schema) > 0 {
		for _, s := range stableSchemas(res.Schema) {
			if sch := res.Schema[s]; sch.Removed == "" {
				// Generate the property name, type, and flags; note that this is in the output position, hence the true.
				// TODO: figure out how to deal with sensitive fields.
				prop, outflags, typ, err := g.propFlagTyp(resname, s, res.Schema, custom[s], true /*out*/)
				if err != nil {
					// Keep going so we can accumulate as many errors as possible.
					err = errors.Errorf("%v:%v: %v", pkg, rawname, err)
					finalerr = multierror.Append(finalerr, err)
				} else if prop != "" && (customNamed || prop != tfbridge.NameProperty) {
					// Make a little comment in the code so it's easy to pick out output properties.
					inprop := inProperty(sch)
					var outcomment string
					if !inprop {
						outcomment = "/*out*/ "
					}

					w.Writefmtln("    public %vreadonly %v%v: %v;", outcomment, prop, outflags, typ)

					// Only keep track of input properties for purposes of initialization data structures.
					if inprop {
						// Regenerate the type and flags since optionals may be different in input positions.
						incust := custom[s]
						inflag := g.tfToJSFlags(sch, incust, false /*out*/)
						intype := g.tfToJSType(sch, incust, false /*out*/)
						inprops = append(inprops, prop)
						inflags = append(inflags, inflag)
						intypes = append(intypes, intype)
						schemas = append(schemas, sch)
						customs = append(customs, incust)
					}
				}
			}
		}
		w.Writefmtln("")
	}

	// Add the standard "factory" functions: get and query.  These are static and so must go before the constructor.
	w.Writefmtln("    public static get(id: lumi.ID): %v {", resname)
	w.Writefmtln("        return <any>undefined; // functionality provided by the runtime")
	w.Writefmtln("    }")
	w.Writefmtln("")
	w.Writefmtln("    public static query(q: any): %v[] {", resname)
	w.Writefmtln("        return <any>undefined; // functionality provided by the runtime")
	w.Writefmtln("    }")
	w.Writefmtln("")

	// Now create a constructor that chains supercalls and stores into properties.
	var argsflags string
	if len(inprops) == 0 {
		// If the number of input properties was zero, we make the args object optional.
		argsflags = "?"
	}
	if customNamed {
		w.Writefmtln("    constructor(args%v: %vArgs) {", argsflags, resname)
		w.Writefmtln("        super();")
	} else {
		w.Writefmtln("    constructor(name: string, args%v: %vArgs) {", argsflags, resname)
		w.Writefmtln("        super(name);")
	}
	var propsindent string
	if len(inprops) == 0 {
		propsindent = "    "
		w.Writefmtln("        if (args !== undefined) {")
	}
	for i, prop := range inprops {
		if !optionalProperty(schemas[i], customs[i], false) {
			w.Writefmtln("%v        if (lumirt.defaultIfComputed(args.%v, \"\") === undefined) {", propsindent, prop)
			w.Writefmtln("%v            throw new Error(\"Property argument '%v' is required, "+
				"but was missing\");", propsindent, prop)
			w.Writefmtln("%v        }", propsindent)
		}
		// Note that we must perform an <any> assignment because the types may not match perfectly.  Input arguments
		// may be missing certain optional computed properties that will be assigned by the provider upon cretion.
		w.Writefmtln("%v        this.%[2]v = <any>args.%[2]v;", propsindent, prop)
	}
	if len(inprops) == 0 {
		w.Writefmtln("        }")
	}
	w.Writefmtln("    }")
	w.Writefmtln("}")

	w.Writefmtln("")

	// Next, generate the args interface for this class.
	w.Writefmtln("export interface %vArgs {", resname)
	for i, prop := range inprops {
		g.generateComment(w, schemas[i].Description, "    ")
		w.Writefmtln("    readonly %v%v: %v;", prop, inflags[i], intypes[i])
	}
	w.Writefmtln("}")
	w.Writefmtln("")

	// Ensure there weren't any custom fields that were unrecognized.
	for key := range custom {
		if _, has := res.Schema[key]; !has {
			cmdutil.Diag().Warningf(
				diag.Message("Custom resource schema %v.%v was not present in the Terraform metadata"),
				resname, key)
		}
	}

	return resourceResult{
		Name:   resname,
		File:   file,
		Submod: submod,
	}, finalerr
}

// generateSubmodules creates a set of index files, if necessary, for the given submodules.  It returns a map of
// submodule name to the generated index file, so that a caller can be sure to re-export it as necessary.
func (g *generator) generateSubmodules(submodules map[string]map[string]string,
	overlays map[string]tfbridge.OverlayInfo,
	outDir string, overlaysDir string) (map[string]string, []string, error) {
	results := make(map[string]string) // the resulting module map.
	var extrafs []string               // the resulting "extra" files to include, if any.

	// Sort the submodules by name so that we emit in a deterministic order.
	var subs []string
	for submod := range submodules {
		subs = append(subs, submod)
	}
	sort.Strings(subs)

	// Now for each module, generate the requisite index.
	for _, sub := range subs {
		exports := submodules[sub]

		// If there are any overlays for this sub-module, copy them and add them.
		if overlays != nil {
			if over, has := overlays[sub]; has {
				for _, overfile := range over.Files {
					from := filepath.Join(overlaysDir, sub, overfile)
					to := filepath.Join(outDir, sub, overfile)
					overname := removeExtension(to, ".ts")
					if _, has := exports[overname]; has {
						return nil, nil, errors.Errorf("Overlay file %v conflicts with a generated file", to)
					}
					if err := copyFile(from, to); err != nil {
						return nil, nil, err
					}
					exports[overname] = to
					extrafs = append(extrafs, to)
				}
				if over.Modules != nil {
					cmdutil.Diag().Warningf(
						diag.Message("Modules more than one level deep not supported; sub-overlays for %v skipped"),
						sub)
				}
			}
		}

		index, err := g.generateIndex(exports, nil, filepath.Join(outDir, sub))
		if err != nil {
			return nil, nil, err
		}
		results[sub] = index
	}

	return results, extrafs, nil
}

// generateIndex creates a module index file for easy access to sub-modules and exports.
func (g *generator) generateIndex(exports, modules map[string]string, outDir string) (string, error) {
	// Open up the file and spew a standard "code-generated" warning header.
	file := filepath.Join(outDir, "index.ts")
	w, err := tools.NewGenWriter(tfgen, file)
	if err != nil {
		return "", err
	}
	defer contract.IgnoreClose(w)
	w.EmitHeaderWarning()

	// Import anything we will export as a sub-module, and then re-export it.
	if len(modules) > 0 {
		w.Writefmtln("// Export sub-modules:")
		var mods []string
		for mod := range modules {
			mods = append(mods, mod)
		}
		sort.Strings(mods)
		for _, mod := range mods {
			rel, err := relModule(outDir, modules[mod])
			if err != nil {
				return "", err
			}
			w.Writefmtln("import * as %v from \"%v\";", mod, rel)
		}
		w.Writefmt("export {")
		for i, mod := range mods {
			if i > 0 {
				w.Writefmt(", ")
			}
			w.Writefmt(mod)
		}
		w.Writefmtln("};")
		w.Writefmtln("")
	}

	// Export anything flatly that is a direct export rather than sub-module.
	if len(exports) > 0 {
		w.Writefmtln("// Export members:")
		var exps []string
		for exp := range exports {
			exps = append(exps, exp)
		}
		sort.Strings(exps)
		for _, exp := range exps {
			rel, err := relModule(outDir, exports[exp])
			if err != nil {
				return "", err
			}
			w.Writefmtln("export * from \"%v\";", rel)
		}
		w.Writefmtln("")
	}

	return file, nil
}

// relModule removes the path suffix from a module and makes it relative to the root path.
func relModule(root string, mod string) (string, error) {
	// Return the path as a relative path to the root, so that imports are relative.
	file, err := filepath.Rel(root, mod)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(file, ".") {
		file = "./" + file
	}
	return removeExtension(file, ".ts"), nil
}

// removeExtension removes the file extension, if any.
func removeExtension(file, ext string) string {
	if strings.HasSuffix(file, ext) {
		return file[:len(file)-len(ext)]
	}
	return file
}

// generatePackageMetadata generates all the non-code metadata required by a Lumi package.
func (g *generator) generatePackageMetadata(pkg string, files []string, outDir string) error {
	// There are three files to write out:
	//     1) Lumi.yaml: Lumi package information
	//     2) package.json: minimal NPM package metadata
	//     3) tsconfig.json: instructions for TypeScript compilation
	if err := g.generateLumiPackageMetadata(pkg, outDir); err != nil {
		return err
	}
	if err := g.generateNPMPackageMetadata(pkg, outDir); err != nil {
		return err
	}
	return g.generateTypeScriptProjectFile(pkg, files, outDir)
}

func (g *generator) generateLumiPackageMetadata(pkg string, outDir string) error {
	w, err := tools.NewGenWriter(tfgen, filepath.Join(outDir, "Lumi.yaml"))
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(w)
	w.Writefmtln("name: %v", pkg)
	w.Writefmtln("description: A Lumi resource provider for %v.", pkg)
	w.Writefmtln("dependencies:")
	w.Writefmtln("    lumi: \"*\"")
	w.Writefmtln("    lumirt: \"*\"")
	w.Writefmtln("")
	return nil
}

func (g *generator) generateNPMPackageMetadata(pkg string, outDir string) error {
	w, err := tools.NewGenWriter(tfgen, filepath.Join(outDir, "package.json"))
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(w)
	w.Writefmtln("{")
	w.Writefmtln("    \"name\": \"@lumi/%v\"", pkg)
	w.Writefmtln("}")
	w.Writefmtln("")
	return nil
}

func (g *generator) generateTypeScriptProjectFile(pkg string, files []string, outDir string) error {
	w, err := tools.NewGenWriter(tfgen, filepath.Join(outDir, "tsconfig.json"))
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(w)
	w.Writefmtln(`{
    "compilerOptions": {
        "outDir": ".lumi/bin",
        "target": "es6",
        "module": "commonjs",
        "moduleResolution": "node",
        "declaration": true,
        "sourceMap": true
    },
    "files": [`)
	for i, file := range files {
		var suffix string
		if i != len(files)-1 {
			suffix = ","
		}
		relfile, err := filepath.Rel(outDir, file)
		if err != nil {
			return err
		}
		w.Writefmtln("        \"%v\"%v", relfile, suffix)
	}
	w.Writefmtln(`    ]
}
`)
	return nil
}

func (g *generator) generateComment(w *tools.GenWriter, comment string, prefix string) {
	prefix += "// "
	if comment != "" {
		curr := 0
		w.Writefmt(prefix)
		for _, word := range strings.Fields(comment) {
			if curr > 0 {
				if curr+len(word)+1 > (maxWidth - len(prefix)) {
					curr = 0
					w.Writefmt("\n%v", prefix)
				} else {
					w.Writefmt(" ")
					curr++
				}
			}
			w.Writefmt(word)
			curr += len(word)
		}
		w.Writefmtln("")
	}
}

// inProperty checks whether the given property is supplied by the user (versus being always computed).
func inProperty(sch *schema.Schema) bool {
	return sch.Optional || sch.Required
}

// optionalProperty checks whether the given property is optional, either due to Terraform or an overlay.
func optionalProperty(sch *schema.Schema, custom tfbridge.SchemaInfo, out bool) bool {
	// If we're checking a property used in an output position, it isn't optional if it's computed.
	if out {
		return !sch.Computed && !custom.HasDefault() && sch.Optional
	}
	return sch.Optional || sch.Computed || custom.HasDefault()
}

// propFlagTyp returns the property name, flag, and type to use for a given property/field/schema element.  The out
// bit determines whether a property suitable for outputs is provided (e.g., it assumes compputeds have occurred).
func (g *generator) propFlagTyp(res string, key string, sch map[string]*schema.Schema,
	custom tfbridge.SchemaInfo, out bool) (string, string, string, error) {
	// Use the name override, if one exists, or use the standard name mangling otherwise.
	prop := custom.Name
	if prop == "" {
		var err error
		prop, err = propertyName(res, key)
		if err != nil {
			return "", "", "", err
		}
	}

	return prop, g.tfToJSFlags(sch[key], custom, out), g.tfToJSType(sch[key], custom, out), nil
}

// propTyp returns the property name and type, without flags, to use for a given property/field/schema element.  The
// out bit determines whether a property suitable for outputs is provided (e.g., it assumes compputeds have occurred).
func (g *generator) propTyp(res string, key string, sch map[string]*schema.Schema,
	custom tfbridge.SchemaInfo, out bool) (string, string, error) {
	// Use the name override, if one exists, or use the standard name mangling otherwise.
	prop := custom.Name
	if prop == "" {
		var err error
		prop, err = propertyName(res, key)
		if err != nil {
			return "", "", err
		}
	}

	return prop, g.tfToJSTypeFlags(sch[key], custom, out), nil
}

// tsToJSFlags returns the JavaScript flags for a given schema property.
func (g *generator) tfToJSFlags(sch *schema.Schema, custom tfbridge.SchemaInfo, out bool) string {
	if optionalProperty(sch, custom, out) {
		return "?"
	}
	return ""
}

// tfToJSType returns the JavaScript type name for a given schema property.
func (g *generator) tfToJSType(sch *schema.Schema, custom tfbridge.SchemaInfo, out bool) string {
	if custom.Type != "" {
		return string(custom.Type.Name())
	} else if custom.Asset != nil {
		return "lumi.asset." + string(custom.Asset.Type().Name())
	}
	var elemcust tfbridge.SchemaInfo
	if custom.Elem != nil {
		elemcust = *custom.Elem
	}
	return g.tfToJSValueType(sch.Type, sch.Elem, elemcust, out)
}

// tfToJSValueType returns the JavaScript type name for a given schema value type and element kind.
func (g *generator) tfToJSValueType(vt schema.ValueType, elem interface{},
	custom tfbridge.SchemaInfo, out bool) string {
	switch vt {
	case schema.TypeBool:
		return "boolean"
	case schema.TypeInt, schema.TypeFloat:
		return "number"
	case schema.TypeString:
		return "string"
	case schema.TypeList:
		return fmt.Sprintf("%v[]", g.tfElemToJSType(elem, custom, out))
	case schema.TypeMap:
		return fmt.Sprintf("{[key: string]: %v}", g.tfElemToJSType(elem, custom, out))
	case schema.TypeSet:
		// IDEA: we can't use ES6 sets here, because we're using values and not objects.  It would be possible to come
		//     up with a ValueSet of some sorts, but that depends on things like shallowEquals which is known to be
		//     brittle and implementation dependent.  For now, we will stick to arrays, and validate on the backend.
		return fmt.Sprintf("%v[]", g.tfElemToJSType(elem, custom, out))
	default:
		contract.Failf("Unrecognized schema type: %v", vt)
		return ""
	}
}

// tfElemToJSType returns the JavaScript type for a given schema element.  This element may be either a simple schema
// property or a complex structure.  In the case of a complex structure, this will expand to its nominal type.
func (g *generator) tfElemToJSType(elem interface{}, custom tfbridge.SchemaInfo, out bool) string {
	// If there is no element type specified, we will accept anything.
	if elem == nil {
		return "any"
	}

	switch e := elem.(type) {
	case schema.ValueType:
		return g.tfToJSValueType(e, nil, custom, out)
	case *schema.Schema:
		// A simple type, just return its type name.
		return g.tfToJSType(e, custom, out)
	case *schema.Resource:
		// A complex type, just expand to its nominal type name.
		// TODO: spill all complex structures in advance so that we don't have insane inline expansions.
		t := "{ "
		c := 0
		for _, s := range stableSchemas(e.Schema) {
			prop, flag, typ, err := g.propFlagTyp("", s, e.Schema, custom.Fields[s], out)
			contract.Assertf(err == nil, "No errors expected for non-resource properties")
			if prop != "" {
				if c > 0 {
					t += ", "
				}
				t += fmt.Sprintf("%v%v: %v", prop, flag, typ)
				c++
			}
		}
		return t + " }"
	default:
		contract.Failf("Unrecognized schema element type: %v", e)
		return ""
	}
}

// tfToJSTypeFlags returns the JavaScript type name for a given schema property, just like tfToJSType, except that if
// the schema is optional, we will emit an undefined union type (for non-field positions).
func (g *generator) tfToJSTypeFlags(sch *schema.Schema, custom tfbridge.SchemaInfo, out bool) string {
	ts := g.tfToJSType(sch, custom, out)
	if optionalProperty(sch, custom, out) {
		ts += " | undefined"
	}
	return ts
}

// generateCustomImports traverses a custom schema map, deeply, to figure out the set of imported names and files that
// will be required to access those names.  WARNING: this routine doesn't (yet) attempt to eliminate naming collisions.
func generateCustomImports(w *tools.GenWriter,
	infos map[string]tfbridge.SchemaInfo, pkg string, root string, curr string) error {
	imports := make(map[string][]string)
	if err := gatherCustomImports(infos, imports, pkg, root, curr); err != nil {
		return err
	}
	if len(imports) > 0 {
		var impfiles []string
		for impfile := range imports {
			impfiles = append(impfiles, impfile)
		}
		sort.Strings(impfiles)
		for _, impfile := range impfiles {
			w.Writefmt("import {")
			for i, impname := range imports[impfile] {
				if i > 0 {
					w.Writefmt(", ")
				}
				w.Writefmt(impname)
			}
			w.Writefmtln("} from \"%v\";", impfile)
		}
		w.Writefmtln("")
	}
	return nil
}

// gatherCustomImports gathers imports from an entire map of schema info.
func gatherCustomImports(infos map[string]tfbridge.SchemaInfo, imports map[string][]string,
	pkg string, root string, curr string) error {
	if infos != nil {
		for _, info := range infos {
			if err := gatherCustomImportsFrom(info, imports, pkg, root, curr); err != nil {
				return err
			}
		}
	}
	return nil
}

// gatherCustomImportsFrom gathers imports from a single schema info structure.
func gatherCustomImportsFrom(info tfbridge.SchemaInfo, imports map[string][]string,
	pkg string, root string, curr string) error {
	// If this property has a custom schema type, and it isn't "simple" (e.g., string, etc), then we need to
	// create a relative module import.  Note that we assume this is local to the current package!
	if info.Type != "" && !tokens.Token(info.Type).Simple() {
		haspkg := string(info.Type.Module().Package().Name())
		if haspkg != pkg {
			return errors.Errorf("Custom schema type %v was not in the current package %v", haspkg, pkg)
		}
		mod := info.Type.Module().Name()
		modfile := filepath.Join(root,
			strings.Replace(string(mod), tokens.TokenDelimiter, string(filepath.Separator), -1))
		relmod, err := relModule(curr, modfile)
		if err != nil {
			return err
		}
		imports[relmod] = append(imports[modfile], string(info.Type.Name()))
	}

	// If the property has an element type, recurse and propagate any results.
	if info.Elem != nil {
		if err := gatherCustomImportsFrom(*info.Elem, imports, pkg, root, curr); err != nil {
			return err
		}
	}

	// If the property has fields, then simply recurse and propagate any results, if any, to our map.
	if info.Fields != nil {
		if err := gatherCustomImports(info.Fields, imports, pkg, root, curr); err != nil {
			return err
		}
	}

	return nil
}

// resourceName translates a Terraform name into its Lumi name equivalent, plus a suggested filename.
func resourceName(pkg string, rawname string, resinfo tfbridge.ResourceInfo) (string, string) {
	if resinfo.Tok == "" {
		// default transformations.
		contract.Assert(strings.HasPrefix(rawname, pkg+"_"))
		name := rawname[len(pkg)+1:]                     // strip off the pkg prefix.
		return tfbridge.TerraformToLumiName(name, true), // PascalCase the resource name.
			tfbridge.TerraformToLumiName(name, false) // camelCase the filename.
	}
	// otherwise, a custom transformation exists; use it.
	return string(resinfo.Tok.Name()), string(resinfo.Tok.Module().Name())
}

// propertyName translates a Terraform underscore_cased_property_name into the JavaScript camelCasedPropertyName.
func propertyName(resname string, s string) (string, error) {
	if resname != "" && tfbridge.IsBuiltinLumiProperty(s) {
		return "",
			errors.Errorf("Property name '%v' on resource '%v' clashes with builtin name; "+
				"please rename it with an override", s, resname)
	}

	// BUGBUG: work around issue in the Elastic Transcoder where a field has a trailing ":".
	if strings.HasSuffix(s, ":") {
		s = s[:len(s)-1]
	}

	return tfbridge.TerraformToLumiName(s, false /*no to PascalCase; we want camelCase*/), nil
}

func stableResources(resources map[string]*schema.Resource) []string {
	var rs []string
	for r := range resources {
		rs = append(rs, r)
	}
	sort.Strings(rs)
	return rs
}

func stableSchemas(schemas map[string]*schema.Schema) []string {
	var ss []string
	for s := range schemas {
		ss = append(ss, s)
	}
	sort.Strings(ss)
	return ss
}
