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
	"regexp"
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

type RepositoryStatus struct {
	LocalBranch   string
	RemoteBranch  string
	CommitsAhead  int
	CommitsBehind int
	FilesStatus   []FileStatus
}

type FileStatus struct {
	Text      string
	Staged    bool
	Unstaged  bool
	Untracked bool
	Ignored   bool
}

var strict bool

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Runs 'git status' across all sub-directories",
	Long:  `A longer description that spans multiple lines `,
	Run: func(cmd *cobra.Command, args []string) {

		error := color.FgRed.Render
		warn := color.FgYellow.Render
		// info := color.FgCyan.Render
		// pass := color.FgCyan.Render
		success := color.FgCyan.Render

		rootDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		// Find directories
		entries, err := os.ReadDir(rootDir)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%-7s %-50s %-30s %-30s %s\n", "STATUS", "REPOSITORY NAME", "BRANCH", "VERSION", "MESSAGE")

		// Loop through directories
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			repositoryName := entry.Name()
			var branch, version, message string

			repositoryDir := filepath.Join(rootDir, repositoryName)

			if isGitRepository(repositoryDir) {
				// fmt.Println(subdirectory)

				out, err := execGitStatus(repositoryDir)
				if err != nil {
					log.Fatalf("unable to read git repository status : %s", err.Error())
				}

				status := parseOutput(out)
				branch = status.LocalBranch
				version = "Unknown"
				// version = status.RemoteBranch

				var staged, unstaged, untracked, ignored int

				// fmt.Println("Branch is ", status.Branch)
				for _, fs := range status.FilesStatus {
					if Verbose {
						fmt.Println(fs.Text)
					}

					if fs.Staged {
						staged++
					}
					if fs.Unstaged {
						unstaged++
					}
					if fs.Untracked {
						untracked++
					}
					if fs.Ignored {
						ignored++
					}
				}

				switch {
				// case len(status.FilesStatus) == 0:
				// message = success("Up to date")
				case staged+unstaged > 0:
					message = "Changes to commit"
				case strict && untracked > 0:
					message = "Untracked changes"
				}

				switch {
				case status.CommitsAhead > 0:
					if message != "" {
						message += ", "
					}

					message += warn("Changes to push")
				case status.CommitsBehind > 0:
					if message != "" {
						message += ", "
					}

					message += warn("Changes to pull")
				}

				if message == "" {
					message = success("Up to date")
				} else {
					message = warn(message)
				}

			} else {
				message = error("Not versioned")
			}

			cliMessage := fmt.Sprintf("%-50s %-30s %-30s %s", repositoryName, branch, version, message)
			printSuccessMessage(cliMessage)
		}
	},
}

func isGitRepository(dir string) bool {
	gitDir := filepath.Join(dir, ".git")

	_, err := os.ReadDir(gitDir)

	return err == nil
}

// execGitStatus read the git status of the repository located at path
func execGitStatus(path string) (io.Reader, error) {
	app := "git"

	arg0 := "-C"
	arg1 := path
	arg2 := "status"
	arg3 := "-s"
	arg4 := "-b"
	arg5 := "--porcelain"

	cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4, arg5)

	if Verbose {
		fmt.Println(cmd)
	}

	out, err := cmd.Output()

	return bytes.NewReader(out), err
}

//Parse parses a git status output command
//It is compatible with the short version of the git status command
func parseOutput(r io.Reader) RepositoryStatus {

	s := bufio.NewScanner(r)

	var localBranch, remoteBranch string
	var ahead, behind int

	//Extract branch name
	for s.Scan() {
		//Skip any empty line
		if len(s.Text()) < 1 {
			continue
		}

		localBranch, remoteBranch, ahead, behind = parseBranchLine(s.Text())
		break
	}

	var statuses []FileStatus
	for s.Scan() {
		if len(s.Text()) < 1 {
			continue
		}

		statuses = append(statuses, parseFileLine(s.Text()))
	}

	return RepositoryStatus{
		LocalBranch:   localBranch,
		RemoteBranch:  remoteBranch,
		CommitsAhead:  ahead,
		CommitsBehind: behind,
		FilesStatus:   statuses,
	}
}

func parseBranchLine(input string) (string, string, int, int) {
	if Verbose {
		fmt.Println(input)
	}
	// Example line:
	//## develop...origin/develop [ahead 1, behind 18]

	s := bufio.NewScanner(strings.NewReader(input))
	s.Split(bufio.ScanWords)

	//check if input is a status branch line output
	s.Scan()

	if s.Text() != "##" {
		return "", "", 0, 0
	}

	//read next word and return the branch name(s)
	s.Scan()
	b := strings.Split(s.Text(), "...")

	localBranch := b[0]
	remoteBranch := ""

	if len(b) > 1 {
		remoteBranch = b[1]
	}

	var ahead = parseCommitsAhead(input)
	var behind = parseCommitsBehind(input)

	return localBranch, remoteBranch, ahead, behind
}

func parseCommitsAhead(input string) int {
	var r = regexp.MustCompile(`ahead ([0-9]+)`)

	ahead := r.FindStringSubmatch(input)
	if ahead == nil {
		return 0
	}

	i, err := strconv.Atoi(ahead[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return i
}

func parseCommitsBehind(input string) int {
	var r = regexp.MustCompile(`behind ([0-9]+)`)

	behind := r.FindStringSubmatch(input)
	if behind == nil {
		return 0
	}

	i, err := strconv.Atoi(behind[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return i
}

func parseFileLine(input string) FileStatus {
	staged := input[0] != ' ' && input[0] != '?' && input[0] != '!'
	unstaged := input[1] != ' ' && input[1] != '?' && input[1] != '!'
	untracked := strings.HasPrefix(input, "??")
	ignored := strings.HasPrefix(input, "!!")

	return FileStatus{
		Text:      input,
		Staged:    staged,
		Unstaged:  unstaged,
		Untracked: untracked,
		Ignored:   ignored,
	}
}

func init() {
	gitCmd.AddCommand(statusCmd)

	statusCmd.PersistentFlags().BoolVarP(&strict, "strict", "s", false, "treat untracked files as outstanding changes")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
