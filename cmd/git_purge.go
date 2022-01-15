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
	"github.com/klyall/kl-cli/pkg/git"
	"github.com/klyall/kl-cli/pkg/output"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var dryRun bool

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Runs git purge across all sub-directories",
	Long:  `Removes all local branches that no longer have a valid remote branch.`,
	Run: func(cmd *cobra.Command, args []string) {

		out := output.SStdOut{
			Out: os.Stdout,
		}

		gitFetch := git.Fetch{
			Verbose:   Verbose,
			Outputter: out,
		}

		gitBranch := git.Branch{
			Verbose:   Verbose,
			Outputter: out,
		}

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

			repositoryName := entry.Name()
			var message string

			repositoryDir := filepath.Join(WorkingDir, repositoryName)

			if !isGitRepository(repositoryDir) {
				continue
			}

			err := gitFetch.ExecWithPurge(repositoryDir)
			if err != nil {
				message := fmt.Sprintf("%-50s Unable to fetch git repository: %s", repositoryName, err.Error())
				out.Error(message)
				continue
			}

			remoteBranches, err := gitBranch.ExecRemote(repositoryDir)
			if err != nil {
				message := fmt.Sprintf("%-50s Unable to retieve remote branches for repository: %s", repositoryName, err.Error())
				out.Error(message)
				continue
			}

			localBranches, err := gitBranch.ExecVV(repositoryDir)
			if err != nil {
				message := fmt.Sprintf("%-50s Unable to retieve branches for repository: %s", repositoryName, err.Error())
				out.Error(message)
				continue
			}

			for _, lb := range localBranches {
				if lb.RemoteBranchName != "" && !contains(remoteBranches, lb.RemoteBranchName) {
					if lb.CurrentBranch {
						msg := fmt.Sprintf("Unable to delete current branch '%s'", lb.LocalBranchName)
						out.Error(formatMessage(repositoryName, msg))
					} else if dryRun {
						msg := fmt.Sprintf("Dry Run: %s branch will be deleted", lb.LocalBranchName)
						out.Warn(formatMessage(repositoryName, msg))
					} else {
						err := gitBranch.ExecDelete(repositoryDir, lb.LocalBranchName)
						if err != nil {
							msg := fmt.Sprintf("Unable to delete local branch '%s': %s", lb.LocalBranchName, err.Error())
							message := formatMessage(repositoryName, msg)
							out.Error(message)
						} else {
							msg := fmt.Sprintf("%s branch deleted", lb.LocalBranchName)
							out.Success(formatMessage(repositoryName, msg))
						}
					}
				}
			}

			if len(remoteBranches) == 0 {
				message = out.RenderInfo("No remote")
			} else if dryRun {
				message = out.RenderInfo("Purge Dry Run")
			} else {
				message = out.RenderInfo("Purged")
			}

			out.Success(formatMessage(repositoryName, message))
		}
	},
}

func formatMessage(repositoryName, message string) string {
	return fmt.Sprintf("%-50s %s", repositoryName, message)
}

func contains(r []git.RemoteBranchName, branch git.RemoteBranchName) bool {
	for _, b := range r {
		if b == branch {
			return true
		}
	}
	return false
}

func init() {
	gitCmd.AddCommand(purgeCmd)

	purgeCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "show what would be done, without making any changes.")
}
