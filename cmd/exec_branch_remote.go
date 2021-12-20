package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type RemoteBranchName string

func execGitBranchRemote(path string) ([]RemoteBranchName, error) {
	app := "git"

	arg0 := "-C"
	arg1 := path
	arg2 := "branch"
	arg3 := "-r"

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)

	if Verbose {
		fmt.Println(cmd)
	}

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	branches := parseBranchRemoteOutput(bytes.NewReader(out))

	return branches, nil
}

func parseBranchRemoteOutput(r io.Reader) []RemoteBranchName {
	branches := []RemoteBranchName{}

	s := bufio.NewScanner(r)

	for s.Scan() {
		line := strings.TrimSpace(s.Text())

		if Verbose {
			fmt.Println(s.Text())
		}

		if line != "" {
			parts := strings.Split(line, " ")

			branches = append(branches, RemoteBranchName(parts[0]))
		}
	}

	return branches
}
