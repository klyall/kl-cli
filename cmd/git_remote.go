/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/klyall/kl-cli/pkg/git"
	"github.com/klyall/kl-cli/pkg/output"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Runs 'git remote' across all sub-directories",
	Long:  `Runs 'git remote' across all sub-directories.`,
	Run: func(cmd *cobra.Command, args []string) {

		out := output.SStdOut{
			Out: os.Stdout,
		}

		gitRemote := git.Remote{
			Verbose:   Verbose,
			Outputter: out,
		}

		error := color.FgRed.Render

		// Find directories
		entries, err := os.ReadDir(WorkingDir)

		if err != nil {
			log.Fatal(err)
		}

		// Loop through directories
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			var message string

			repositoryName := entry.Name()

			repositoryDir := filepath.Join(WorkingDir, repositoryName)

			if isGitRepository(repositoryDir) {

				remote, err := gitRemote.Exec(repositoryDir)
				if err != nil {
					message := fmt.Sprintf("%-50s Unable to fetch git repository: %s", repositoryName, err.Error())
					out.Error(message)
					continue
				}

				switch {
				case remote.Fetch == "":
					message = out.RenderInfo("No remote")
				case remote.Fetch != remote.Push:
					message = out.RenderWarn("Remotes mismatch: %s (fetch) %S (push)", remote.Fetch, remote.Push)
				default:
					message = remote.Fetch
				}

			} else {
				message = error("Not versioned")
			}

			cliMessage := fmt.Sprintf("%-50s %s", repositoryName, message)
			out.Success(cliMessage)
		}
	},
}

func init() {
	gitCmd.AddCommand(remoteCmd)
}
