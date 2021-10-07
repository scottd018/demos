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
	"os"

	"github.com/spf13/cobra"
)

const defaultManifest = `apiVersion: apps.acme.com/v1alpha1
kind: WebStore
metadata:
  name: webstore-sample
spec:
  webstoreImage: "nginx:1.17"
  serviceName: "webstore-svc"
  webStoreReplicas: 2
`

// newInitCommand creates a new instance of the init subcommand.
func (c *WebstorectlCommand) newInitCommand() {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Write a sample custom resource manifest for a workload to standard out",
		Long:  "Write a sample custom resource manifest for a workload to standard out",
		RunE: func(cmd *cobra.Command, args []string) error {
			outputStream := os.Stdout

			if _, err := outputStream.WriteString(defaultManifest); err != nil {
				return fmt.Errorf("failed to write to stdout, %w", err)
			}

			return nil
		},
	}

	c.AddCommand(initCmd)
}
