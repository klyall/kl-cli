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
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

var strict bool

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Runs 'git status' across all sub-directories",
	Long:  `A longer description that spans multiple lines `,
	Run: func(cmd *cobra.Command, args []string) {

		out := output.SStdOut{
			Out: os.Stdout,
		}

		gitStatus := git.Status{
			Verbose:   Verbose,
			Outputter: out,
			Strict:    strict,
		}

		fmt.Printf("%-7s %-50s %-30s %-30s %s\n", "STATUS", "REPOSITORY NAME", "BRANCH", "VERSION", "MESSAGE")

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
			repositoryDir := filepath.Join(WorkingDir, repositoryName)

			repositoryStatus, err := ExecuteGitStatus(repositoryDir, gitStatus)
			if err != nil {
				message := fmt.Sprintf("%-50s Unable to read git repository: %s", repositoryName, err.Error())
				out.Error(message)
				continue
			}

			message := createMessage(repositoryStatus, out)

			formattedMessage := fmt.Sprintf("%-50s %-30s %-30s %s", repositoryName, repositoryStatus.LocalBranch, repositoryStatus.VersionNumber, message)
			out.Success(formattedMessage)
		}
	},
}

func ExecuteGitStatus(repositoryDir string, gitStatus git.Status) (git.RepositoryStatus, error) {

	if !isGitRepository(repositoryDir) {
		return git.RepositoryStatus{
			LocalStatus: git.NotVersioned,
		}, nil
	}

	status, err := gitStatus.Exec(repositoryDir)
	if err != nil {
		return git.RepositoryStatus{}, err
	}

	return status, nil
}

func createMessage(repositoryStatus git.RepositoryStatus, out output.Outputter) string {
	var message string

	if repositoryStatus.LocalStatus == git.NotVersioned ||
		repositoryStatus.LocalStatus == repositoryStatus.RemoteStatus {
		return repositoryStatus.LocalStatus.Color.Render(repositoryStatus.LocalStatus.Message)
	}

	if repositoryStatus.LocalStatus != git.NoChanges {
		message = repositoryStatus.LocalStatus.Message
	}

	if repositoryStatus.RemoteStatus != git.NoChanges {
		if message != "" {
			message += ", "
		}

		message = repositoryStatus.RemoteStatus.Message
	}

	return out.RenderWarn(message)
}

func init() {
	gitCmd.AddCommand(statusCmd)

	statusCmd.PersistentFlags().BoolVarP(&strict, "strict", "s", false, "treat untracked files as outstanding changes")
}
