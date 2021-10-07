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
	"github.com/spf13/cobra"
)

// WebstorectlCommand represents the base command when called without any subcommands.
type WebstorectlCommand struct {
	*cobra.Command
}

// NewWebstorectlCommand returns an instance of the WebstorectlCommand.
func NewWebstorectlCommand() *WebstorectlCommand {
	c := &WebstorectlCommand{
		Command: &cobra.Command{
			Use:   "webstorectl",
			Short: "Manage webstore stuff like a boss",
			Long:  "Manage webstore stuff like a boss",
		},
	}

	c.addSubCommands()

	return c
}

// Run represents the main entry point into the command
// This is called by main.main() to execute the root command.
func (c *WebstorectlCommand) Run() {
	cobra.CheckErr(c.Execute())
}

// addSubCommands adds any additional subCommands to the root command.
func (c *WebstorectlCommand) addSubCommands() {
	c.newInitCommand()
	c.newGenerateCommand()
	//+kubebuilder:scaffold:operator-builder:subcommands
}
