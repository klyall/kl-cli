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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Runs 'git fetch' across all sub-directories",
	Long:  `Runs 'git fetch' across all sub-directories.`,
	Run: func(cmd *cobra.Command, args []string) {

		error := color.FgRed.Render

		rootDir, err := os.Getwd()

		if err != nil {
			log.Fatal(err)
		}

		// Find directories
		entries, err := os.ReadDir(rootDir)
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

			repositoryDir := filepath.Join(rootDir, repositoryName)

			if isGitRepository(repositoryDir) {

				out, err := execGitFetch(repositoryDir)

				if err != nil {
					message := fmt.Sprintf("%-50s Unable to fetch git repository: %s", repositoryName, err.Error())
					printErrorMessage(message)
					continue
				}

				parseGitFetchOutput(out)

				message = info("Fetch complete")
			} else {
				message = error("Not versioned")
			}

			cliMessage := fmt.Sprintf("%-50s %s", repositoryName, message)
			printSuccessMessage(cliMessage)
		}
	},
}

func execGitFetch(path string) (io.Reader, error) {
	app := "git"

	arg0 := "-C"
	arg1 := path
	arg2 := "fetch"

	cmd := exec.Command(app, arg0, arg1, arg2)

	if Verbose {
		fmt.Println(cmd)
	}

	out, err := cmd.Output()

	return bytes.NewReader(out), err
}

func parseGitFetchOutput(r io.Reader) {

	s := bufio.NewScanner(r)

	if Verbose {
		s.Scan()
		line := s.Text()

		if line != "" {
			fmt.Println(s.Text())
		}
	}
}

func init() {
	gitCmd.AddCommand(fetchCmd)
}