package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

// unknown is a sentinel we substitute in place whenever an unknown value is encountered. This should
// only happen during a preview but lets us still interact with TF and use Pulumi to manage dependencies.
// TODO: we don't actually use this right now so unknowns will screw everything up.
const unknown = "74D93920-ED26-11E3-AC10-0800200C9A66"

// TODO: there's got to be definitions for these already in the Terraform repo somewhere.
type terraformFile struct {
	Provider map[string]terraformFileProviderBlock `json:"provider,omitempty"`
	Module   map[string]terraformFileModuleBlock   `json:"module,omitempty"`
	Output   map[string]terraformFileOutputBlock   `json:"output,omitempty"`
}
type terraformFileProviderBlock map[string]interface{}
type terraformFileModuleBlock map[string]interface{}
type terraformFileOutputBlock map[string]interface{}

// pulumiNameToTfName uses the standard convention to map a Pulumi name to a Terraform name, without
// requiring a Terraform schema.
func pulumiNameToTfName(tfName string) string {
	var result string
	for i, c := range tfName {
		if c >= 'A' && c <= 'Z' {
			// if upper case, add an underscore (if it's not #1), and then the lower case version.
			if i != 0 {
				result += "_"
			}
			result += string(unicode.ToLower(c))
		} else {
			result += string(c)
		}
	}

	return result
}

func (p *Provider) createModuleResource(ctx context.Context, req *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	// Unmarshal the inputs.
	props, err := plugin.UnmarshalProperties(req.Properties, plugin.MarshalOptions{})
	if err != nil {
		return nil, err
	}

	// Now apply the state updates; since this is a create, pass an empty state file.
	urn := resource.URN(req.Urn)
	id, outputs, err := p.applyModuleResource(ctx, urn, props.Mappable(), "", false)
	if err != nil {
		return nil, err
	}

	// Serialize the output properties: 1) a bag of module outputs and 2) base64-encoded TF state.
	propOuts := resource.NewPropertyMapFromMap(outputs)
	resultProperties, err := plugin.MarshalProperties(propOuts, plugin.MarshalOptions{})
	if err != nil {
		return nil, err
	}

	// Create an ID out of the unique bits: source+version+name, and return the results.
	return &pulumirpc.CreateResponse{
		Id:         id,
		Properties: resultProperties,
	}, nil
}

func (p *Provider) updateModuleResource(ctx context.Context, req *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	// Unmarshal the old properties so we can fetch the state.
	olds, err := plugin.UnmarshalProperties(req.Olds, plugin.MarshalOptions{})
	if err != nil {
		return nil, err
	}
	oldState, _ := olds.Mappable()["encodedState"].(string)
	if oldState != "" {
		oldStateDec, _ := base64.StdEncoding.DecodeString(oldState)
		oldState = string(oldStateDec)
	}

	// Unmarshal the new inputs and use them.
	props, err := plugin.UnmarshalProperties(req.News, plugin.MarshalOptions{})
	if err != nil {
		return nil, err
	}

	// Now apply the state updates; since this is a create, passing the prior state.
	urn := resource.URN(req.Urn)
	_, outputs, err := p.applyModuleResource(ctx, urn, props.Mappable(), oldState, false)
	if err != nil {
		return nil, err
	}

	// Serialize the output properties: 1) a bag of module outputs and 2) base64-encoded TF state.
	propOuts := resource.NewPropertyMapFromMap(outputs)
	resultProperties, err := plugin.MarshalProperties(propOuts, plugin.MarshalOptions{})
	if err != nil {
		return nil, err
	}

	// Return the resulting properties returned by the TF apply.
	return &pulumirpc.UpdateResponse{Properties: resultProperties}, nil
}

func (p *Provider) destroyModuleResource(ctx context.Context, req *pulumirpc.DeleteRequest) (*empty.Empty, error) {
	// Unmarshal the old properties so we can fetch the state.
	// TODO: do we even need to generate a terraform program on destroy?
	props, err := plugin.UnmarshalProperties(req.Properties, plugin.MarshalOptions{})
	if err != nil {
		return nil, err
	}
	oldState, _ := props.Mappable()["encodedState"].(string)
	if oldState != "" {
		oldStateDec, _ := base64.StdEncoding.DecodeString(oldState)
		oldState = string(oldStateDec)
	}

	// Now run Terraform to destroy the resource.
	urn := resource.URN(req.Urn)
	if _, _, err = p.applyModuleResource(ctx, urn, props.Mappable(), oldState, true); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (p *Provider) applyModuleResource(ctx context.Context, urn resource.URN,
	props map[string]interface{}, state string, destroy bool) (string, map[string]interface{}, error) {
	// To create or update a module resource, we will:
	//     - Initialize a temporary TF workspace.
	//     - Generate a simple HCL.JSON file with the module declaration.
	//     - Run terraform apply.
	//     - Serialize the resulting TF state file.
	//     - Return the results in a Pulumi-friendly format.

	name := string(urn.Name())
	source := props["source"].(string)
	version := props["version"].(string)
	providers := props["providers"].(map[string]interface{})
	inputs := props["inputs"].(map[string]interface{})

	// Make a terraform file, populate it with the provider and module resources, and serialize to JSON.
	tf := &terraformFile{
		Provider: make(map[string]terraformFileProviderBlock),
		Module:   make(map[string]terraformFileModuleBlock),
		Output:   make(map[string]terraformFileOutputBlock),
	}
	for pkey, pval := range providers {
		// HACKHACK: we've made it simple to just pass provider configuration by referencing the
		// Pulumi config bag in the program. Ideally we could just reach out to the engine here
		// and fetch it for all relevant providers. Because of this though we translate the name.
		tf.Provider[pkey] = make(map[string]interface{})
		for k, v := range pval.(map[string]interface{}) {
			tf.Provider[pkey][pulumiNameToTfName(k)] = v
		}
	}
	tf.Module[name] = terraformFileModuleBlock(inputs)
	tf.Module[name]["source"] = source
	tf.Module[name]["version"] = version

	// We need to map the module as an output so we can read and return its output properties from the TF state.
	tf.Output["outs"] = map[string]interface{}{
		"value": fmt.Sprintf("${module.%s}", name),
	}

	// Marshal the program to JSON and get ready to emit it and run a Terraform command.
	tfJSON, err := json.Marshal(tf)
	if err != nil {
		return "", nil, err
	}

	moduleOuts, state, err := p.runTerraform(ctx, urn, string(tfJSON), state, destroy)
	if err != nil {
		return "", nil, err
	}
	outputs := map[string]interface{}{
		"source":       source,
		"version":      version,
		"providers":    providers,
		"inputs":       inputs,
		"outputs":      moduleOuts,
		"encodedState": base64.StdEncoding.EncodeToString([]byte(state)),
	}

	// Manufacture an ID out of the unique bits of the resource: source+version+name.
	return fmt.Sprintf("%s@%s#%s", source, version, name), outputs, nil
}

// runTerraform executes the given HCL JSON program with the given state, and returns
// the resulting error, output variables, and state. The input state should be the empty string
// for the first time a program is being run.
func (p *Provider) runTerraform(ctx context.Context, urn resource.URN,
	program string, state string, destroy bool) (map[string]interface{}, string, error) {
	// Create a temporary directory to contain the TF program.
	// TODO: this is going to initialize the same thing over and over again, download plugins, etc.
	// Should we be caching things somehow?
	dir, err := os.MkdirTemp("", "pulumi-tf-module")
	if err != nil {
		return nil, state, err
	}

	// Spit out the program itself.
	tfmainPath := filepath.Join(dir, "main.tf.json")
	if err = ioutil.WriteFile(tfmainPath, []byte(program), 0644); err != nil {
		return nil, state, err
	}

	// Write out the state file if there is one.
	tfstatePath := filepath.Join(dir, "terraform.tfstate")
	if state != "" {
		if err = ioutil.WriteFile(tfstatePath, []byte(state), 0644); err != nil {
			return nil, state, err
		}
	}

	// logCmdOutput logs a specific stream of output from the given comamnd. It returns
	// a channel that's closed when that output has completed.
	logCmdOutput := func(pipeFunc func() (io.ReadCloser, error)) chan bool {
		done := make(chan bool)
		go func() {
			pipe, pipeErr := pipeFunc()
			if pipeErr != nil {
				panic(pipeErr)
			}
			reader := bufio.NewReader(pipe)
			for {
				msg, err := reader.ReadString('\n')
				if strings.TrimSpace(msg) != "" {
					// TODO: should be ephemeral, but it doesn't seem to be exposed.
					// BUGBUG: should be passing 'urn' here but it seems to go missing if we do ...
					p.host.Log(ctx, diag.Info, "", msg)
				}
				if err != nil {
					break
				}
			}
			close(done)
		}()
		return done
	}

	// runCmdAndWait runs the command, streams stdout/stderr, and returns any errors.
	runCmdAndWait := func(cmd string, args ...string) error {
		c := exec.Command(cmd, args...)
		c.Dir = dir
		cStdoutDone := logCmdOutput(c.StdoutPipe)
		cStderrDone := logCmdOutput(c.StderrPipe)
		if err = c.Start(); err != nil {
			return err
		}
		<-cStdoutDone
		<-cStderrDone
		return c.Wait()
	}

	// Run the initialization.
	if err = runCmdAndWait("terraform", "init"); err != nil {
		return nil, state, err
	}

	// Run the desired command, apply or destroy, as appropriate.
	var verb string
	if destroy {
		verb = "destroy"
	} else {
		verb = "apply"
	}
	if err = runCmdAndWait("terraform", verb, "-auto-approve"); err != nil {
		return nil, state, err
	}

	// Afterwards, read in the state file and return it.
	newState, err := ioutil.ReadFile(tfstatePath)
	if err != nil {
		return nil, state, err
	}

	// Deserialize the state so we can fetch the module output properties to return.
	var stateStruct terraformState
	if err = json.Unmarshal(newState, &stateStruct); err != nil {
		return nil, string(newState), err
	}

	return stateStruct.Outputs.Outs.Value, string(newState), nil
}

type terraformState struct {
	Outputs terraformStateOutputs `json:"outputs"`
}

type terraformStateOutputs struct {
	Outs terraformStateOutput `json:"outs"`
}

type terraformStateOutput struct {
	Value map[string]interface{} `json:"value"`
}
