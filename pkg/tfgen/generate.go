// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfgen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/diag"
	"github.com/pulumi/pulumi/pkg/tokens"
	"github.com/pulumi/pulumi/pkg/tools"
	"github.com/pulumi/pulumi/pkg/util/cmdutil"
	"github.com/pulumi/pulumi/pkg/util/contract"

	"github.com/pulumi/pulumi-terraform/pkg/tfbridge"
)

type generator struct {
	pkg      string                // the TF package name (e.g. `aws`)
	version  string                // the package version.
	provinfo tfbridge.ProviderInfo // the provider info for customizing code generation
}

func newGenerator(pkg string, version string, provinfo tfbridge.ProviderInfo) *generator {
	return &generator{
		pkg:      pkg,
		version:  version,
		provinfo: provinfo,
	}
}

const (
	tfgen              = "the Pulumi Terraform Bridge (TFGEN) Tool"
	defaultOutDir      = "pack/"
	defaultOverlaysDir = "overlays/"
	maxWidth           = 120 // the ideal maximum width of the generated file.
)

// Generate creates Pulumi packages out of one or more Terraform plugins.  It accepts a list of all of the input
// Terraform providers, already bound statically to the code (since we cannot obtain schema information dynamically),
// walks them and generates the Pulumi code, and spews that code into the output directory.
func (g *generator) Generate(outDir, overlaysDir string) error {
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
	return g.generateProvider(g.provinfo, outDir, overlaysDir)
}

// generateProvider creates a single standalone Pulumi package for the given provider.
func (g *generator) generateProvider(provinfo tfbridge.ProviderInfo, outDir, overlaysDir string) error {
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
		cfgfile, err := g.generateConfig(prov.Schema, provinfo.Config, outDir)
		if err != nil {
			return err
		}
		// ensure we export the config submodule and add its file to the project.
		submodules["config"] = map[string]string{
			"vars": cfgfile,
		}
		files = append(files, cfgfile)
	}

	var pendingExports []genResult

	// Generate all resources.
	if len(prov.ResourcesMap) > 0 {
		resources, err := g.generateResources(prov.ResourcesMap, provinfo.Resources, outDir)
		if err != nil {
			return err
		}
		pendingExports = append(pendingExports, resources...)
	}

	// Place all data sources into a single data module.
	if len(prov.DataSourcesMap) > 0 {
		dataSources, err := g.generateDataSources(prov.DataSourcesMap, provinfo.DataSources, outDir)
		if err != nil {
			return err
		}
		pendingExports = append(pendingExports, dataSources...)
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

	// Make sure all exports are added to the proper lists.
	for _, export := range pendingExports {
		if export.Submod == "" || export.Submod == "index" {
			// if no sub-module, export flatly in our own index.
			exports[export.Name] = export.File
		} else {
			// otherwise, make sure to track this in the submodule so we can create and export it correctly.
			submod := submodules[export.Submod]
			if submod == nil {
				submod = make(map[string]string)
				submodules[export.Submod] = submod
			}
			submod[export.Name] = export.File
		}
		files = append(files, export.File)
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

	// Generate all of the package metadata: Pulumi.yaml, package.json, and tsconfig.json.
	err = g.generatePackageMetadata(files, outDir, provinfo.Overlay)
	if err != nil {
		return err
	}

	// Finally, emit the version information in a special VERSION file so we know where it came from.
	gitinfo, err := getGitInfo(provinfo.Name)
	if err != nil {
		return err
	}
	versionInfo := fmt.Sprintf("Generated by %s from:\n", os.Args[0])
	versionInfo += fmt.Sprintf("Repo: %v\n", gitinfo.Repo)
	if gitinfo.Tag != "" {
		versionInfo += fmt.Sprintf("Tag: %v\n", gitinfo.Tag)
	}
	if gitinfo.Commit != "" {
		versionInfo += fmt.Sprintf("Commit: %v\n", gitinfo.Commit)
	}
	versionInfo += "\n"
	return ioutil.WriteFile(filepath.Join(outDir, "VERSION"), []byte(versionInfo), 0600)
}

// copyFile is a stupid file copy routine.  It reads the file into memory to avoid messy OS-specific oddities.
func copyFile(from, to string) error {
	err := os.MkdirAll(path.Dir(to), 0700)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(to, body, 0600)
}

// generateConfig takes a map of config variables and emits a config submodule to the given file.
func (g *generator) generateConfig(cfg map[string]*schema.Schema,
	custom map[string]*tfbridge.SchemaInfo, outDir string) (string, error) {
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

	// We'll need the Pulumi SDK.
	w.Writefmtln("import * as pulumi from \"pulumi\";")
	w.Writefmtln("")

	// First look for any custom types that will require any imports.
	if err := generateCustomImports(w, custom, g.pkg, outDir, confDir); err != nil {
		return "", err
	}

	// Create a config bag for this package.
	w.Writefmtln("let _config = new pulumi.Config(\"%v:config\");", g.pkg)
	w.Writefmtln("")

	// Now just emit a simple export for each variable.
	for _, key := range cfgkeys {
		// Generate a name and type to use for this key.
		prop, typ, err := g.propTyp(key, cfg[key], custom[key], true /*out*/)
		if err != nil {
			return "", err
		} else if prop != "" {
			var getfunc string
			if optionalProperty(cfg[key], custom[key], false) {
				getfunc = "get"
			} else {
				getfunc = "require"
			}
			if cfg[key].Type != schema.TypeString {
				// Only try to parse a JSON object if the config isn't a straight string.
				getfunc = fmt.Sprintf("%sObject<%s>", getfunc, typ)
			}
			var anycast string
			if custom[key] != nil && custom[key].Type != "" {
				// If there's a custom type, we need to inject a cast to silence the compiler.
				anycast = "<any>"
			}
			g.generateCommentAdjustLines(w, cfg[key].Description, "")
			w.Writefmtln("export let %[1]v: %[2]v = %[3]s_config.%[4]v(\"%[1]v\");", prop, typ, anycast, getfunc)
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

func (g *generator) generateResources(resmap map[string]*schema.Resource,
	custom map[string]*tfbridge.ResourceInfo, outDir string) ([]genResult, error) {
	// For each resource, create its own dedicated type and module export.
	var results []genResult
	var reserr error
	reshits := make(map[string]bool)
	for _, r := range stableResources(resmap) {
		var resinfo *tfbridge.ResourceInfo
		if resmap != nil {
			resinfo = custom[r]
		}
		if resinfo == nil {
			// if this resource was missing, issue a warning and skip it.
			cmdutil.Diag().Warningf(
				diag.Message("Resource %v not found in provider map; skipping"), r)
			continue
		}
		reshits[r] = true
		result, err := g.generateResource(r, resmap[r], resinfo, outDir, outDir)
		if err != nil {
			// Keep track of the error, but keep going, so we can expose more at once.
			reserr = multierror.Append(reserr, err)
		} else {
			results = append(results, result)
		}
	}
	if reserr != nil {
		return nil, reserr
	}

	// Emit a warning if there is a map but some names didn't match.
	var resnames []string
	for resname := range custom {
		resnames = append(resnames, resname)
	}
	sort.Strings(resnames)
	for _, resname := range resnames {
		if !reshits[resname] {
			cmdutil.Diag().Warningf(
				diag.Message("Resource %v (%v) wasn't found in the Terraform module; possible name mismatch?"),
				resname, custom[resname].Tok)
		}
	}

	return results, nil
}

type genResult struct {
	Name   string // the export's name.
	File   string // the export's filename.
	Submod string // the submodule name, if any.
}

// generateResource generates a single module for the given resource.
func (g *generator) generateResource(rawname string,
	res *schema.Resource, resinfo *tfbridge.ResourceInfo, root, outDir string) (genResult, error) {
	// Transform the name as necessary.
	resname, filename := resourceName(g.pkg, rawname, resinfo)

	// Make a fully qualified file path that we will write to.
	file := filepath.Join(outDir, filename+".ts")

	// If the filename contains slashes, it is a sub-module, and we must ensure it exists.
	var submod string
	if slix := strings.Index(filename, "/"); slix != -1 {
		// Extract the module and file parts.
		submod = filename[:slix]
		if strings.Contains(filename[slix+1:], "/") {
			return genResult{},
				errors.Errorf("Modules nested more than one level deep not currently supported")
		}

		// Ensure the submodule directory exists.
		if err := tools.EnsureFileDir(file); err != nil {
			return genResult{}, err
		}
	}

	// Open up the file and spew a standard "code-generated" warning header.
	w, err := tools.NewGenWriter(tfgen, file)
	if err != nil {
		return genResult{}, err
	}
	defer contract.IgnoreClose(w)
	w.EmitHeaderWarning()

	// Now import the modules we need.
	w.Writefmtln("import * as pulumi from \"pulumi\";")
	w.Writefmtln("")

	// If there are imports required due to the custom schema info, emit them now.
	custom := resinfo.Fields
	if err = generateCustomImports(w, custom, g.pkg, outDir, filepath.Dir(file)); err != nil {
		return genResult{}, err
	}

	// Collect documentation information
	parsedDocs, err := getDocsForPackage(g.pkg, ResourceDocs, rawname, resinfo.Docs)
	if err != nil {
		return genResult{}, err
	}

	// Write the TypeDoc/JSDoc for the resource class
	if parsedDocs.Description != "" {
		g.generateComment(w, parsedDocs.Description, "")
	}

	// Generate the resource class.
	w.Writefmtln("export class %s extends pulumi.CustomResource {", resname)

	// First, generate all instance properties.
	var finalerr error
	var inprops []string
	var reqprops int
	var outprops []string
	var inflags []string
	var intypes []string
	var inrawnames []string
	var schemas []*schema.Schema
	var customs []*tfbridge.SchemaInfo
	if len(res.Schema) > 0 {
		for _, s := range stableSchemas(res.Schema) {
			if sch := res.Schema[s]; sch.Removed == "" {
				// Generate the property name, type, and flags; note that this is in the output position, hence the true.
				// TODO: figure out how to deal with sensitive fields.
				prop, outflags, typ, err := g.propFlagTyp(s, sch, custom[s], true /*out*/)
				if err != nil {
					// Keep going so we can accumulate as many errors as possible.
					err = errors.Errorf("%v:%v: %v", g.pkg, rawname, err)
					finalerr = multierror.Append(finalerr, err)
				} else if prop != "" {
					// Make a little comment in the code so it's easy to pick out output properties.
					inprop := inProperty(sch)
					var outcomment string
					if !inprop {
						outcomment = "/*out*/ "
					}
					// Emit documentation for the property if available
					if argDoc, ok := parsedDocs.Arguments[s]; ok {
						g.generateComment(w, argDoc, "    ")
					} else if attrDoc, ok := parsedDocs.Attributes[s]; ok {
						g.generateComment(w, attrDoc, "    ")
					}
					// Emit the property as a property; it has to carry undefined because of planning.
					w.Writefmtln("    public %vreadonly %v%v: pulumi.Computed<%v>;",
						outcomment, prop, outflags, typ)

					// Only keep track of input properties for purposes of initialization data structures.
					if inprop {
						// Regenerate the type and flags since optionals may be different in input positions.
						incust := custom[s]
						inflag := g.tfToJSFlags(sch, incust, false /*out*/)
						intype := g.tfToJSType(sch, incust, false /*out*/)
						inprops = append(inprops, prop)
						inflags = append(inflags, inflag)
						intypes = append(intypes, intype)
						inrawnames = append(inrawnames, s)
						schemas = append(schemas, sch)
						customs = append(customs, incust)
						if !optionalProperty(sch, incust, false) {
							reqprops++
						}
					} else {
						// Remember output properties because we still want to "zero-initialize" them as properties.
						outprops = append(outprops, prop)
					}
				}
			}
		}
		w.Writefmtln("")
	}

	// Now create a constructor that chains supercalls and stores into properties.
	w.Writefmtln("    /**")
	w.Writefmtln("     * Create a %s resource with the given unique name, arguments, and optional additional", resname)
	w.Writefmtln("     * resource dependencies.")
	w.Writefmtln("     *")
	w.Writefmtln("     * @param urnName A _unique_ name for this %s instance", resname)
	w.Writefmtln("     * @param args A collection of arguments for creating this %s instance", resname)
	w.Writefmtln("     * @param parent An optional parent resource to which this resource belongs")
	w.Writefmtln("     * @param dependsOn A optional array of additional resources this instance depends on")
	w.Writefmtln("     */")
	var argsflags string
	if reqprops == 0 {
		// If the number of input properties was zero, we make the args object optional.
		argsflags = "?"
	}
	w.Writefmtln("    constructor(urnName: string, args%v: %vArgs, opts?: pulumi.ResourceOptions) {",
		argsflags, resname)
	if reqprops == 0 {
		// If the property arg isn't required, zero-init it if it wasn't actually passed in.
		w.Writefmtln("        args = args || {};")
	}

	// First, validate all required arguments.
	for i, prop := range inprops {
		if !optionalProperty(schemas[i], customs[i], false) {
			w.Writefmtln("        if (args.%v === undefined) {", prop)
			w.Writefmtln("            throw new Error(\"Missing required property '%v'\");", prop)
			w.Writefmtln("        }")
		}
	}

	// Now invoke the super constructor with the type, name, and a property map.
	w.Writefmtln("        super(\"%s\", urnName, {", resinfo.Tok)
	for _, prop := range inprops {
		w.Writefmtln("            \"%[1]s\": args.%[1]s,", prop)
	}
	for _, prop := range outprops {
		w.Writefmtln("            \"%s\": undefined,", prop)
	}
	w.Writefmtln("        }, opts);")

	w.Writefmtln("    }")
	w.Writefmtln("}")
	w.Writefmtln("")

	// Next, generate the args interface for this class.
	w.Writefmtln("/**")
	w.Writefmtln(" * The set of arguments for constructing a %s resource.", resname)
	w.Writefmtln(" */")
	w.Writefmtln("export interface %vArgs {", resname)
	for i, prop := range inprops {
		if argDoc, ok := parsedDocs.Arguments[inrawnames[i]]; ok {
			g.generateComment(w, argDoc, "    ")
		} else {
			g.generateCommentAdjustLines(w, schemas[i].Description, "    ")
		}
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

	return genResult{
		Name:   resname,
		File:   file,
		Submod: submod,
	}, finalerr
}

func (g *generator) generateDataSources(sources map[string]*schema.Resource,
	custom map[string]*tfbridge.DataSourceInfo, outDir string) ([]genResult, error) {
	// Sort and enumerate all variables in a deterministic order.
	var srckeys []string
	for key := range sources {
		srckeys = append(srckeys, key)
	}
	sort.Strings(srckeys)

	// For each data source, create its own dedicated function and module export.
	var results []genResult
	var dserr error
	dshits := make(map[string]bool)
	for _, ds := range srckeys {
		var dsinfo *tfbridge.DataSourceInfo
		if sources != nil {
			dsinfo = custom[ds]
		}
		if dsinfo == nil {
			// if this data source was missing, issue a warning and skip it.
			cmdutil.Diag().Warningf(
				diag.Message("Data source %v not found in provider map; skipping"), ds)
			continue
		}
		dshits[ds] = true
		result, err := g.generateDataSource(ds, sources[ds], dsinfo, outDir, outDir)
		if err != nil {
			// Keep track of the error, but keep going, so we can expose more at once.
			dserr = multierror.Append(dserr, err)
		} else {
			results = append(results, result)
		}
	}
	if dserr != nil {
		return nil, dserr
	}

	// Emit a warning if there is a map but some names didn't match.
	var dsnames []string
	for dsname := range custom {
		dsnames = append(dsnames, dsname)
	}
	sort.Strings(dsnames)
	for _, dsname := range dsnames {
		if !dshits[dsname] {
			cmdutil.Diag().Warningf(
				diag.Message("Data source %v (%v) wasn't found in the Terraform module; possible name mismatch?"),
				dsname, custom[dsname].Tok)
		}
	}

	return results, nil
}

// generateDataSource generates a single module for the given data source function.
func (g *generator) generateDataSource(rawname string,
	ds *schema.Resource, dsinfo *tfbridge.DataSourceInfo, root, outDir string) (genResult, error) {
	// Transform the name as necessary.
	dsname, filename := dataSourceName(g.pkg, rawname, dsinfo)

	// Make a fully qualified file path that we will write to.
	file := filepath.Join(outDir, filename+".ts")

	// If the filename contains slashes, it is a sub-module, and we must ensure it exists.
	var submod string
	if slix := strings.Index(filename, "/"); slix != -1 {
		// Extract the module and file parts.
		submod = filename[:slix]
		if strings.Contains(filename[slix+1:], "/") {
			return genResult{},
				errors.Errorf("Modules nested more than one level deep not currently supported")
		}

		// Ensure the submodule directory exists.
		if err := tools.EnsureFileDir(file); err != nil {
			return genResult{}, err
		}
	}

	// Open up the file and spew a standard "code-generated" warning header.
	w, err := tools.NewGenWriter(tfgen, file)
	if err != nil {
		return genResult{}, err
	}
	defer contract.IgnoreClose(w)
	w.EmitHeaderWarning()

	// We'll need the Pulumi SDK.
	w.Writefmtln("import * as pulumi from \"pulumi\";")
	w.Writefmtln("")

	// Collect documentation information for this data source.
	parsedDocs, err := getDocsForPackage(g.pkg, DataSourceDocs, rawname, dsinfo.Docs)
	if err != nil {
		return genResult{}, err
	}

	// Write the TypeDoc/JSDoc for the data source function.
	if parsedDocs.Description != "" {
		g.generateComment(w, parsedDocs.Description, "")
	}

	// Sort the args and return properties so we are ready to go.
	args := ds.Schema
	var argkeys []string
	for arg := range args {
		argkeys = append(argkeys, arg)
	}
	sort.Strings(argkeys)

	// See if arguments for this function are optional.
	argc := 0
	reqc := 0
	optflag := "?"
	for _, arg := range argkeys {
		if inProperty(args[arg]) {
			argc++
			if !optionalProperty(args[arg], dsinfo.Fields[arg], false) {
				reqc++
				optflag = ""
			}
		}
	}

	// Now, emit the function signature.
	firstch, firstsz := utf8.DecodeRuneInString(dsname)
	dstype := string(unicode.ToUpper(firstch)) + dsname[firstsz:]
	var argsig string
	if argc > 0 {
		argsig = fmt.Sprintf("args%v: %vArgs", optflag, dstype)
	}
	w.Writefmtln("export function %v(%v): Promise<%vResult> {", dsname, argsig, dstype)
	if argc > 0 && reqc == 0 {
		w.Writefmtln("    args = args || {};")
	}
	w.Writefmtln("    return pulumi.runtime.invoke(\"%s\", {", dsinfo.Tok)
	for _, arg := range argkeys {
		if inProperty(args[arg]) {
			name, err := g.propName(arg, dsinfo.Fields[arg])
			if err != nil {
				return genResult{}, err
			}
			w.Writefmtln("        \"%[1]s\": args.%[1]s,", name)
		}
	}
	w.Writefmtln("    });")
	w.Writefmtln("}")
	w.Writefmtln("")

	// Emit the arguments interface (used as input).
	if argc > 0 {
		w.Writefmtln("/**")
		w.Writefmtln(" * A collection of arguments for invoking %s.", dsname)
		w.Writefmtln(" */")
		w.Writefmtln("export interface %vArgs {", dstype)
		for _, arg := range argkeys {
			// Only emit input properties in the arguments data structure.
			if inProperty(args[arg]) {
				// Emit documentation for the property if available
				if argDoc, ok := parsedDocs.Arguments[arg]; ok {
					g.generateComment(w, argDoc, "    ")
				}
				prop, flags, typ, err := g.propFlagTyp(arg, args[arg], dsinfo.Fields[arg], false /*out*/)
				if err != nil {
					return genResult{}, err
				}
				w.Writefmtln("    %v%v: %v;", prop, flags, typ)
			}
		}
		w.Writefmtln("}")
		w.Writefmtln("")
	}

	// Emit the result interface (used for return values).
	w.Writefmtln("/**")
	w.Writefmtln(" * A collection of values returned by %s.", dsname)
	w.Writefmtln(" */")
	w.Writefmtln("export interface %vResult {", dstype)
	for _, arg := range argkeys {
		// Only emit computed properties in the resulting return data structure.
		if args[arg].Computed {
			// Emit documentation for the property if available
			if attrDoc, ok := parsedDocs.Attributes[arg]; ok {
				g.generateComment(w, attrDoc, "    ")
			}
			prop, flags, typ, err := g.propFlagTyp(arg, args[arg], dsinfo.Fields[arg], true /*out*/)
			if err != nil {
				return genResult{}, err
			}
			w.Writefmtln("    %v%v: %v;", prop, flags, typ)
		}
	}
	w.Writefmtln("}")
	w.Writefmtln("")

	return genResult{
		Name:   dsname,
		File:   file,
		Submod: submod,
	}, nil
}

// generateSubmodules creates a set of index files, if necessary, for the given submodules.  It returns a map of
// submodule name to the generated index file, so that a caller can be sure to re-export it as necessary.
func (g *generator) generateSubmodules(submodules map[string]map[string]string,
	overlays map[string]*tfbridge.OverlayInfo,
	outDir string, overlaysDir string) (map[string]string, []string, error) {
	results := make(map[string]string) // the resulting module map.
	var extrafs []string               // the resulting "extra" files to include, if any.

	// Sort the submodules by name so that we emit in a deterministic order.
	subnames := stableSubmodules(submodules, overlays)

	// Now for each module, generate the requisite index.
	for _, sub := range subnames {
		exports, ok := submodules[sub]
		if !ok {
			exports = map[string]string{}
		}

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

// generatePackageMetadata generates all the non-code metadata required by a Pulumi package.
func (g *generator) generatePackageMetadata(files []string, outDir string,
	overlay *tfbridge.OverlayInfo) error {
	// There are three files to write out:
	//     1) Pulumi.yaml: Pulumi package information
	//     2) package.json: minimal NPM package metadata
	//     3) tsconfig.json: instructions for TypeScript compilation
	if err := g.generatePulumiPackageMetadata(outDir); err != nil {
		return err
	}
	if err := g.generateNPMPackageMetadata(outDir, overlay); err != nil {
		return err
	}
	return g.generateTypeScriptProjectFile(files, outDir)
}

func (g *generator) generatePulumiPackageMetadata(outDir string) error {
	w, err := tools.NewGenWriter(tfgen, filepath.Join(outDir, "Pulumi.yaml"))
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(w)
	w.Writefmtln("name: %v", g.pkg)
	w.Writefmtln("description: A Pulumi Fabric resource provider for %v.", g.pkg)
	w.Writefmtln("language: nodejs")
	w.Writefmtln("")
	return nil
}

func (g *generator) generateNPMPackageMetadata(outDir string, overlay *tfbridge.OverlayInfo) error {
	w, err := tools.NewGenWriter(tfgen, filepath.Join(outDir, "package.json"))
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(w)
	w.Writefmtln(`{`)
	w.Writefmtln(`    "name": "@pulumi/%v",`, g.pkg)
	w.Writefmtln(`    "version": "%s",`, g.version)
	w.Writefmtln(`    "scripts": {`)
	w.Writefmtln(`        "build": "tsc"`)
	w.Writefmtln(`    },`)
	if len(overlay.Dependencies) > 0 {
		w.Writefmtln(`    "dependencies": {`)
		var deps []string
		for dep := range overlay.Dependencies {
			deps = append(deps, dep)
		}
		sort.Strings(deps)
		for i, dep := range deps {
			var comma string
			if i != len(deps)-1 {
				comma = ","
			}
			w.Writefmtln(`         "%s": "%s"%s`, dep, overlay.Dependencies[dep], comma)
		}
		w.Writefmtln(`    },`)
	}
	w.Writefmtln(`    "devDependencies": {`)
	if len(overlay.DevDependencies) > 0 {
		var deps []string
		for dep := range overlay.DevDependencies {
			deps = append(deps, dep)
		}
		sort.Strings(deps)
		for _, dep := range deps {
			w.Writefmtln(`        "%s": "%s",`, dep, overlay.DevDependencies[dep])
		}
	}
	w.Writefmtln(`        "typescript": "^2.5.2"`)
	w.Writefmtln(`    },`)
	w.Writefmtln(`    "peerDependencies": {`)
	if len(overlay.PeerDependencies) > 0 {
		var deps []string
		for dep := range overlay.PeerDependencies {
			deps = append(deps, dep)
		}
		sort.Strings(deps)
		for _, dep := range deps {
			w.Writefmtln(`        "%s": "%s",`, dep, overlay.PeerDependencies[dep])
		}
	}
	w.Writefmtln(`        "pulumi": "*"`)
	w.Writefmtln(`    }`)
	w.Writefmtln(`}`)
	return nil
}

func (g *generator) generateTypeScriptProjectFile(files []string, outDir string) error {
	w, err := tools.NewGenWriter(tfgen, filepath.Join(outDir, "tsconfig.json"))
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(w)
	w.Writefmtln(`{
    "compilerOptions": {
        "outDir": "bin",
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

// sanitizeForDocComment ensures that no `*/` sequence appears in the string, to avoid
// accidentally closing the comment block.
func sanitizeForDocComment(str string) string {
	return strings.Replace(str, "*/", "*&#47;", -1)
}

func (g *generator) generateCommentAdjustLines(w *tools.GenWriter, comment string, prefix string) {
	if comment != "" {
		curr := 0
		w.Writefmtln("%v/**", prefix)
		w.Writefmt("%v * ", prefix)
		for _, word := range strings.Fields(comment) {
			word = sanitizeForDocComment(word)
			if curr > 0 {
				if curr+len(word)+1 > (maxWidth - len(prefix)) {
					curr = 0
					w.Writefmt("\n%v * ", prefix)
				} else {
					w.Writefmt(" ")
					curr++
				}
			}
			w.Writefmt(word)
			curr += len(word)
		}
		w.Writefmtln("")
		w.Writefmtln("%v */", prefix)
	}
}

func (g *generator) generateComment(w *tools.GenWriter, doc string, prefix string) {
	if doc != "" {
		lines := strings.Split(doc, "\n")
		w.Writefmtln("%v/**", prefix)
		for i, docLine := range lines {
			docLine = sanitizeForDocComment(docLine)
			// Break if we get to the last line and it's empty
			if i == len(lines)-1 && strings.TrimSpace(docLine) == "" {
				break
			}
			// Print the line of documentation
			w.Writefmtln("%v * %s", prefix, docLine)
		}
		w.Writefmtln("%v */", prefix)
	}
}

// inProperty checks whether the given property is supplied by the user (versus being always computed).
func inProperty(sch *schema.Schema) bool {
	return sch.Optional || sch.Required
}

// optionalProperty checks whether the given property is optional, either due to Terraform or an overlay.
func optionalProperty(sch *schema.Schema, custom *tfbridge.SchemaInfo, out bool) bool {
	// If we're checking a property used in an output position, it isn't optional if it's computed.
	customDefault := custom != nil && custom.HasDefault()
	if out {
		return sch.Optional && !sch.Computed && !customDefault
	}
	return sch.Optional || sch.Computed || customDefault
}

func (g *generator) propName(key string, custom *tfbridge.SchemaInfo) (string, error) {
	// Use the name override, if one exists, or use the standard name mangling otherwise.
	var prop string
	if custom != nil {
		prop = custom.Name
	}
	if prop == "" {
		var err error
		prop, err = propertyName(key)
		if err != nil {
			return "", err
		}
	}
	return prop, nil
}

// propFlagTyp returns the property name, flag, and type to use for a given property/field/schema element.  The out
// bit determines whether a property suitable for outputs is provided (e.g., it assumes compputeds have occurred).
func (g *generator) propFlagTyp(name string, sch *schema.Schema,
	custom *tfbridge.SchemaInfo, out bool) (string, string, string, error) {
	prop, err := g.propName(name, custom)
	if err != nil {
		return "", "", "", err
	}
	flags := g.tfToJSFlags(sch, custom, out)
	typ := g.tfToJSType(sch, custom, out)
	return prop, flags, typ, nil
}

// propTyp returns the property name and type, without flags, to use for a given property/field/schema element.  The
// out bit determines whether a property suitable for outputs is provided (e.g., it assumes compputeds have occurred).
func (g *generator) propTyp(name string, sch *schema.Schema,
	custom *tfbridge.SchemaInfo, out bool) (string, string, error) {
	prop, err := g.propName(name, custom)
	if err != nil {
		return "", "", err
	}
	typ := g.tfToJSTypeFlags(sch, custom, out)
	return prop, typ, nil
}

// tsToJSFlags returns the JavaScript flags for a given schema property.
func (g *generator) tfToJSFlags(sch *schema.Schema, custom *tfbridge.SchemaInfo, out bool) string {
	if optionalProperty(sch, custom, out) {
		return "?"
	}
	return ""
}

// tfToJSType returns the JavaScript type name for a given schema property.
func (g *generator) tfToJSType(sch *schema.Schema, custom *tfbridge.SchemaInfo, out bool) string {
	var elem *tfbridge.SchemaInfo
	if custom != nil {
		if custom.Type != "" {
			t := string(custom.Type.Name())
			if len(custom.AltTypes) > 0 {
				for _, at := range custom.AltTypes {
					t = fmt.Sprintf("%s | %s", t, at.Name())
				}
			}
			if !out {
				t = fmt.Sprintf("pulumi.ComputedValue<%s>", t)
			}
			return t
		} else if custom.Asset != nil {
			return "pulumi.asset." + custom.Asset.Type()
		}
		elem = custom.Elem
	}
	return g.tfToJSValueType(sch.Type, sch.Elem, elem, out)
}

// tfToJSValueType returns the JavaScript type name for a given schema value type and element kind.
func (g *generator) tfToJSValueType(vt schema.ValueType, elem interface{},
	custom *tfbridge.SchemaInfo, out bool) string {
	// First figure out the raw type.
	var t string
	var array bool
	switch vt {
	case schema.TypeBool:
		t = "boolean"
	case schema.TypeInt, schema.TypeFloat:
		t = "number"
	case schema.TypeString:
		t = "string"
	case schema.TypeList:
		t = g.tfElemToJSType(elem, custom, out)
		array = true
	case schema.TypeMap:
		t = fmt.Sprintf("{[key: string]: %v}", g.tfElemToJSType(elem, custom, out))
	case schema.TypeSet:
		// IDEA: we can't use ES6 sets here, because we're using values and not objects.  It would be possible to come
		//     up with a ValueSet of some sorts, but that depends on things like shallowEquals which is known to be
		//     brittle and implementation dependent.  For now, we will stick to arrays, and validate on the backend.
		t = g.tfElemToJSType(elem, custom, out)
		array = true
	default:
		contract.Failf("Unrecognized schema type: %v", vt)
	}

	// Now, if it is an input property value, it must be wrapped in a ComputedValue<T>.
	if !out {
		t = fmt.Sprintf("pulumi.ComputedValue<%s>", t)
	}

	// Finally make sure arrays are arrays; this must be done after the above, so we get a ComputedValue<T>[],
	// and not a ComputedValue<T[]>, which would constrain the ability to flexibly construct them.
	if array {
		t = fmt.Sprintf("%s[]", t)
	}

	return t
}

// tfElemToJSType returns the JavaScript type for a given schema element.  This element may be either a simple schema
// property or a complex structure.  In the case of a complex structure, this will expand to its nominal type.
func (g *generator) tfElemToJSType(elem interface{}, custom *tfbridge.SchemaInfo, out bool) string {
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
			var fldinfo *tfbridge.SchemaInfo
			if custom != nil {
				fldinfo = custom.Fields[s]
			}
			prop, flag, typ, err := g.propFlagTyp(s, e.Schema[s], fldinfo, out)
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
func (g *generator) tfToJSTypeFlags(sch *schema.Schema, custom *tfbridge.SchemaInfo, out bool) string {
	ts := g.tfToJSType(sch, custom, out)
	if optionalProperty(sch, custom, out) {
		ts += " | undefined"
	}
	return ts
}

// generateCustomImports traverses a custom schema map, deeply, to figure out the set of imported names and files that
// will be required to access those names.  WARNING: this routine doesn't (yet) attempt to eliminate naming collisions.
func generateCustomImports(w *tools.GenWriter,
	infos map[string]*tfbridge.SchemaInfo, pkg string, root string, curr string) error {
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
func gatherCustomImports(infos map[string]*tfbridge.SchemaInfo, imports map[string][]string,
	pkg string, root string, curr string) error {
	for _, info := range infos {
		if err := gatherCustomImportsFrom(info, imports, pkg, root, curr); err != nil {
			return err
		}
	}
	return nil
}

// gatherCustomImportsFrom gathers imports from a single schema info structure.
func gatherCustomImportsFrom(info *tfbridge.SchemaInfo, imports map[string][]string,
	pkg string, root string, curr string) error {
	if info != nil {
		// If this property has custom schema types that aren't "simple" (e.g., string, etc), then we need to
		// create a relative module import.  Note that we assume this is local to the current package!
		var custty []tokens.Type
		if info.Type != "" {
			custty = append(custty, info.Type)
			custty = append(custty, info.AltTypes...)
		}
		for _, ct := range custty {
			if !tokens.Token(ct).Simple() {
				haspkg := string(ct.Module().Package().Name())
				if haspkg != pkg {
					return errors.Errorf("Custom schema type %v was not in the current package %v", haspkg, pkg)
				}
				mod := ct.Module().Name()
				modfile := filepath.Join(root,
					strings.Replace(string(mod), tokens.TokenDelimiter, string(filepath.Separator), -1))
				relmod, err := relModule(curr, modfile)
				if err != nil {
					return err
				}
				imports[relmod] = append(imports[modfile], string(ct.Name()))
			}
		}

		// If the property has an element type, recurse and propagate any results.
		if info.Elem != nil {
			if err := gatherCustomImportsFrom(info.Elem, imports, pkg, root, curr); err != nil {
				return err
			}
		}

		// If the property has fields, then simply recurse and propagate any results, if any, to our map.
		if info.Fields != nil {
			if err := gatherCustomImports(info.Fields, imports, pkg, root, curr); err != nil {
				return err
			}
		}
	}

	return nil
}

// dataSourceName translates a Terraform name into its Pulumi name equivalent.
func dataSourceName(pkg string, rawname string, dsinfo *tfbridge.DataSourceInfo) (string, string) {
	if dsinfo == nil || dsinfo.Tok == "" {
		// default transformations.
		name := withoutPackageName(pkg, rawname)            // strip off the pkg prefix.
		return tfbridge.TerraformToPulumiName(name, false), // camelCase the data source name.
			tfbridge.TerraformToPulumiName(name, false) // camelCase the filename.
	}
	// otherwise, a custom transformation exists; use it.
	return string(dsinfo.Tok.Name()), string(dsinfo.Tok.Module().Name())
}

// resourceName translates a Terraform name into its Pulumi name equivalent, plus a suggested filename.
func resourceName(pkg string, rawname string, resinfo *tfbridge.ResourceInfo) (string, string) {
	if resinfo == nil || resinfo.Tok == "" {
		// default transformations.
		name := withoutPackageName(pkg, rawname)           // strip off the pkg prefix.
		return tfbridge.TerraformToPulumiName(name, true), // PascalCase the resource name.
			tfbridge.TerraformToPulumiName(name, false) // camelCase the filename.
	}
	// otherwise, a custom transformation exists; use it.
	return string(resinfo.Tok.Name()), string(resinfo.Tok.Module().Name())
}

// withoutPackageName strips off the package prefix from a raw name.
func withoutPackageName(pkg string, rawname string) string {
	contract.Assert(strings.HasPrefix(rawname, pkg+"_"))
	name := rawname[len(pkg)+1:] // strip off the pkg prefix.
	return name
}

// propertyName translates a Terraform underscore_cased_property_name into the JavaScript camelCasedPropertyName.
func propertyName(s string) (string, error) {
	// BUGBUG: work around issue in the Elastic Transcoder where a field has a trailing ":".
	if strings.HasSuffix(s, ":") {
		s = s[:len(s)-1]
	}

	return tfbridge.TerraformToPulumiName(s, false /*no to PascalCase; we want camelCase*/), nil
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

func stableSubmodules(submodules map[string]map[string]string, overlays map[string]*tfbridge.OverlayInfo) []string {
	subMap := map[string]bool{}
	for submod := range submodules {
		subMap[submod] = true
	}
	for submod := range overlays {
		subMap[submod] = true
	}
	var subs []string
	for submod := range subMap {
		subs = append(subs, submod)
	}
	sort.Strings(subs)
	return subs
}
