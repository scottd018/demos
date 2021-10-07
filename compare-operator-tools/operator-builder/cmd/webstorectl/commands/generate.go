/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"sigs.k8s.io/yaml"

	appsv1alpha1 "github.com/scottd018/demos/apis/apps/v1alpha1"
	"github.com/scottd018/demos/apis/apps/v1alpha1/webstore"
)

type generateCommand struct {
	*cobra.Command
	workloadManifest string
}

// newGenerateCommand creates a new instance of the generate subcommand.
func (c *WebstorectlCommand) newGenerateCommand() {
	g := &generateCommand{}
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate child resource manifests from a workload's custom resource",
		Long:  "Generate child resource manifests from a workload's custom resource",
		RunE:  g.generate,
	}

	generateCmd.Flags().StringVarP(
		&g.workloadManifest,
		"workload-manifest",
		"w",
		"",
		"Filepath to the workload manifest to generate child resources for.",
	)
	generateCmd.MarkFlagRequired("workload-manifest")

	c.AddCommand(generateCmd)
}

// generate creates child resource manifests from a workload's custom resource.
func (g *generateCommand) generate(cmd *cobra.Command, args []string) error {

	filename, _ := filepath.Abs(g.workloadManifest)

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s, %w", filename, err)
	}

	var workload appsv1alpha1.WebStore

	err = yaml.Unmarshal(yamlFile, &workload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml %s into workload, %w", filename, err)
	}

	resourceObjects := make([]metav1.Object, len(webstore.CreateFuncs))

	for i, f := range webstore.CreateFuncs {
		resource, err := f(&workload)
		if err != nil {
			return err
		}

		resourceObjects[i] = resource
	}

	e := json.NewYAMLSerializer(json.DefaultMetaFactory, nil, nil)

	outputStream := os.Stdout

	for _, o := range resourceObjects {
		if _, err := outputStream.WriteString("---\n"); err != nil {
			return fmt.Errorf("failed to write output, %w", err)
		}

		if err := e.Encode(o.(runtime.Object), os.Stdout); err != nil {
			return fmt.Errorf("failed to write output, %w", err)
		}
	}

	return nil
}
