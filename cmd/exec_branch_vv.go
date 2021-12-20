package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type LocalBranchName string

type LocalBranch struct {
	LocalBranchName  LocalBranchName
	RemoteBranchName RemoteBranchName
	CurrentBranch    bool
}

func execGitBranchVV(path string) ([]LocalBranch, error) {
	app := "git"

	arg0 := "-C"
	arg1 := path
	arg2 := "branch"
	arg3 := "-vv"

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)

	if Verbose {
		fmt.Println(cmd)
	}

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	branches := parseBranchVVOutput(bytes.NewReader(out))

	return branches, nil
}

func parseBranchVVOutput(r io.Reader) []LocalBranch {
	branches := []LocalBranch{}

	s := bufio.NewScanner(r)

	for s.Scan() {
		line := s.Text()

		if Verbose {
			fmt.Println(line)
		}

		if line != "" {
			branch := parseBranchVVLine(line)

			branches = append(branches, branch)
		}
	}

	return branches
}

func parseBranchVVLine(line string) LocalBranch {
	currentBranch := line[0] == '*'

	parts := strings.Split(strings.TrimSpace(line[1:]), " ")

	var remoteBranch string

	if start := strings.Index(line, "["); start != -1 {

		end := strings.Index(line, "]")
		remoteStatus := line[start+1 : end]

		if i := strings.Index(remoteStatus, ":"); i != -1 {
			remoteBranch = remoteStatus[:i]
		} else {
			remoteBranch = remoteStatus
		}
	}

	branch := LocalBranch{
		LocalBranchName:  LocalBranchName(parts[0]),
		RemoteBranchName: RemoteBranchName(remoteBranch),
		CurrentBranch:    currentBranch,
	}

	return branch
}
