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

type Status struct {
	color   color.Color
	message string
}

var CommittedChanges = Status{warnColor, "Changes to push"}
var NoChanges = Status{successColor, "Up to date"}
var NotVersioned = Status{errorColor, "Not versioned"}
var RemoteChanges = Status{warnColor, "Changes to pull"}
var UncommitedChanges = Status{warnColor, "Changes to commit"}
var UntrackedChanges = Status{warnColor, "Untracked changes"}

type RepositoryStatus struct {
	Versioned     bool
	VersionNumber string
	LocalBranch   string
	RemoteBranch  string
	LocalStatus   Status
	RemoteStatus  Status
	CommitsAhead  int
	CommitsBehind int
	Staged        int
	Unstaged      int
	Untracked     int
	Ignored       int
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

		fmt.Printf("%-7s %-50s %-30s %-30s %s\n", "STATUS", "REPOSITORY NAME", "BRANCH", "VERSION", "MESSAGE")

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
			repositoryDir := filepath.Join(rootDir, repositoryName)

			repositoryStatus, err := ExecuteGitStatus(repositoryDir)

			if err != nil {
				message := fmt.Sprintf("%-50s Unable to read git repository: %s", repositoryName, err.Error())
				printErrorMessage(message)
				continue
			}

			message := createMessage(repositoryStatus)

			formattedMessage := fmt.Sprintf("%-50s %-30s %-30s %s", repositoryName, repositoryStatus.LocalBranch, repositoryStatus.VersionNumber, message)
			printSuccessMessage(formattedMessage)
		}
	},
}

func ExecuteGitStatus(repositoryDir string) (RepositoryStatus, error) {

	if !isGitRepository(repositoryDir) {
		return RepositoryStatus{
			LocalStatus: NotVersioned,
		}, nil
	}

	out, err := execGitStatus(repositoryDir)

	if err != nil {
		return RepositoryStatus{}, err
	}

	return parseGitStatusOutput(out), nil
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

func parseGitStatusOutput(r io.Reader) RepositoryStatus {

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

	var remoteStatus Status

	switch {
	case ahead > 0:
		remoteStatus = CommittedChanges
	case behind > 0:
		remoteStatus = RemoteChanges
	default:
		remoteStatus = NoChanges
	}

	var statuses []FileStatus
	for s.Scan() {
		if len(s.Text()) < 1 {
			continue
		}

		statuses = append(statuses, parseFileLine(s.Text()))
	}

	staged, unstaged, untracked, ignored := calculateTotals(statuses)

	var localStatus Status

	switch {
	case staged+unstaged > 0:
		localStatus = UncommitedChanges
	case strict && untracked > 0:
		localStatus = UntrackedChanges
	default:
		localStatus = NoChanges
	}

	return RepositoryStatus{
		Versioned:     true,
		LocalBranch:   localBranch,
		RemoteBranch:  remoteBranch,
		LocalStatus:   localStatus,
		RemoteStatus:  remoteStatus,
		CommitsAhead:  ahead,
		CommitsBehind: behind,
		Staged:        staged,
		Unstaged:      unstaged,
		Untracked:     untracked,
		Ignored:       ignored,
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

func calculateTotals(fileStatuses []FileStatus) (staged, unstaged, untracked, ignored int) {
	for _, fs := range fileStatuses {
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
	return
}

func createMessage(repositoryStatus RepositoryStatus) string {
	var message string

	if repositoryStatus.LocalStatus == NotVersioned ||
		repositoryStatus.LocalStatus == repositoryStatus.RemoteStatus {
		return repositoryStatus.LocalStatus.color.Render(repositoryStatus.LocalStatus.message)
	}

	if repositoryStatus.LocalStatus != NoChanges {
		message = repositoryStatus.LocalStatus.message
	}

	if repositoryStatus.RemoteStatus != NoChanges {
		if message != "" {
			message += ", "
		}

		message = repositoryStatus.RemoteStatus.message
	}

	return warnColor.Render(message)
}

func init() {
	gitCmd.AddCommand(statusCmd)

	statusCmd.PersistentFlags().BoolVarP(&strict, "strict", "s", false, "treat untracked files as outstanding changes")
}
