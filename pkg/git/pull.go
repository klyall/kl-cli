package git

import (
	"github.com/klyall/kl-cli/pkg/output"
	"os/exec"
)

type Pull struct {
	Verbose   bool
	Outputter output.Outputter
}

func (p Pull) Exec(path string) error {
	app := "git"

	arg0 := "-C"
	arg1 := path
	arg2 := "pull"

	cmd := exec.Command(app, arg0, arg1, arg2)

	if p.Verbose {
		p.Outputter.Debug(cmd)
	}

	out, err := cmd.Output()

	if p.Verbose {
		p.Outputter.DebugBytes(out)
	}

	return err
}
