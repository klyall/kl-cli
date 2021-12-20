package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

func execGitFetchPurge(path string) error {
	app := "git"

	arg0 := "-C"
	arg1 := path
	arg2 := "fetch"
	arg3 := "-p"

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)

	if Verbose {
		fmt.Println(cmd)
	}

	out, err := cmd.Output()

	printOutput(bytes.NewReader(out))

	return err
}

func printOutput(r io.Reader) {

	s := bufio.NewScanner(r)

	for s.Scan() {
		if Verbose {
			line := s.Text()

			if line != "" {
				fmt.Println(s.Text())
			}
		}
	}
}
