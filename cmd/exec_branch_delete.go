package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
)

func execGitBranchDelete(path string, branch LocalBranchName) error {
	app := "git"

	arg0 := "-C"
	arg1 := path
	arg2 := "branch"
	arg3 := "-D"
	arg4 := branch

	cmd := exec.Command(app, arg0, arg1, arg2, arg3, string(arg4))

	if Verbose {
		fmt.Println(cmd)
	}

	out, err := cmd.Output()

	printOutput(bytes.NewReader(out))

	return err
}
