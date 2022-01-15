package git

import (
	"bufio"
	"bytes"
	"github.com/klyall/kl-cli/pkg/output"
	"io"
	"os/exec"
	"strings"
)

type Remote struct {
	Verbose   bool
	Outputter output.Outputter
}

func (r Remote) Exec(path string) (RepositoryRemote, error) {
	app := "git"

	arg0 := "-C"
	arg1 := path
	arg2 := "remote"
	arg3 := "-v"

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)

	if r.Verbose {
		r.Outputter.Debug(cmd)
	}

	out, err := cmd.Output()
	if err != nil {
		return RepositoryRemote{}, err
	}

	return r.parseGitRemoteOutput(bytes.NewReader(out)), nil
}

func (r Remote) parseGitRemoteOutput(reader io.Reader) RepositoryRemote {

	s := bufio.NewScanner(reader)

	s.Scan()
	fetch := r.parseGitRemoteLine(s.Text())
	s.Scan()
	push := r.parseGitRemoteLine(s.Text())

	return RepositoryRemote{
		Fetch: fetch,
		Push:  push,
	}
}

func (r Remote) parseGitRemoteLine(line string) string {
	if len(line) == 0 {
		return ""
	}

	if r.Verbose {
		r.Outputter.Debug(line)
	}

	s := bufio.NewScanner(strings.NewReader(line))
	s.Split(bufio.ScanWords)

	s.Scan()
	s.Scan()

	return s.Text()
}
