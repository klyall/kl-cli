package git

import (
	"bufio"
	"bytes"
	"github.com/klyall/kl-cli/pkg/output"
	"io"
	"os/exec"
	"strings"
)

type Branch struct {
	Verbose   bool
	Outputter output.Outputter
}

func (b Branch) ExecDelete(path string, branch LocalBranchName) error {
	app := "git"

	arg0 := "-C"
	arg1 := path
	arg2 := "branch"
	arg3 := "-D"
	arg4 := branch

	cmd := exec.Command(app, arg0, arg1, arg2, arg3, string(arg4))

	if b.Verbose {
		b.Outputter.Debug(cmd)
	}

	out, err := cmd.Output()

	if b.Verbose {
		b.Outputter.DebugBytes(out)
	}

	return err
}

func (b Branch) ExecRemote(path string) ([]RemoteBranchName, error) {
	app := "git"

	arg0 := "-C"
	arg1 := path
	arg2 := "branch"
	arg3 := "-r"

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)

	if b.Verbose {
		b.Outputter.Debug(cmd)
	}

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	branches := b.parseBranchRemoteOutput(bytes.NewReader(out))

	return branches, nil
}

func (b Branch) parseBranchRemoteOutput(r io.Reader) []RemoteBranchName {
	var branches []RemoteBranchName

	s := bufio.NewScanner(r)

	for s.Scan() {
		line := strings.TrimSpace(s.Text())

		if b.Verbose {
			b.Outputter.Debug(s.Text())
		}

		if line != "" {
			parts := strings.Split(line, " ")

			branches = append(branches, RemoteBranchName(parts[0]))
		}
	}

	return branches
}

func (b Branch) ExecVV(path string) ([]LocalBranch, error) {
	app := "git"

	arg0 := "-C"
	arg1 := path
	arg2 := "branch"
	arg3 := "-vv"

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)

	if b.Verbose {
		b.Outputter.Debug(cmd)
	}

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	branches := b.parseBranchVVOutput(bytes.NewReader(out))

	return branches, nil
}

func (b Branch) parseBranchVVOutput(r io.Reader) []LocalBranch {
	branches := []LocalBranch{}

	s := bufio.NewScanner(r)

	for s.Scan() {
		line := s.Text()

		if b.Verbose {
			b.Outputter.Debug(line)
		}

		if line != "" {
			branch := b.parseBranchVVLine(line)

			branches = append(branches, branch)
		}
	}

	return branches
}

func (b Branch) parseBranchVVLine(line string) LocalBranch {
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
