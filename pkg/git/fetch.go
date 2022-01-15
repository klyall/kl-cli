package git

import (
	"github.com/klyall/kl-cli/pkg/output"
	"os/exec"
)

type Fetch struct {
	Verbose   bool
	Outputter output.Outputter
}

func (f Fetch) ExecWithPurge(path string) error {
	app := "git"

	arg0 := "-C"
	arg1 := path
	arg2 := "fetch"
	arg3 := "-p"

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)

	if f.Verbose {
		f.Outputter.Debug(cmd)
	}

	out, err := cmd.Output()

	if f.Verbose {
		f.Outputter.DebugBytes(out)
	}

	return err
}

func (f Fetch) Exec(path string) error {
	app := "git"

	arg0 := "-C"
	arg1 := path
	arg2 := "fetch"

	cmd := exec.Command(app, arg0, arg1, arg2)

	if f.Verbose {
		f.Outputter.Debug(cmd)
	}

	out, err := cmd.Output()

	if f.Verbose {
		f.Outputter.DebugBytes(out)
	}

	return err
}
