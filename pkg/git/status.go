package git

import (
	"bufio"
	"bytes"
	"github.com/gookit/color"
	"github.com/klyall/kl-cli/pkg/output"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type StatusMessage struct {
	Color   color.Color
	Message string
}

var CommittedChanges = StatusMessage{output.WarnColor, "Changes to push"}
var NoChanges = StatusMessage{output.SuccessColor, "Up to date"}
var NotVersioned = StatusMessage{output.ErrorColor, "Not versioned"}
var RemoteChanges = StatusMessage{output.WarnColor, "Changes to pull"}
var UncommittedChanges = StatusMessage{output.WarnColor, "Changes to commit"}
var UntrackedChanges = StatusMessage{output.WarnColor, "Untracked changes"}

type Status struct {
	Verbose   bool
	Outputter output.Outputter
	Strict    bool
}

func (s Status) Exec(path string) (RepositoryStatus, error) {
	app := "git"

	arg0 := "-C"
	arg1 := path
	arg2 := "status"
	arg3 := "-s"
	arg4 := "-b"
	arg5 := "--porcelain"

	cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4, arg5)

	if s.Verbose {
		s.Outputter.Debug(cmd)
	}

	out, err := cmd.Output()
	if err != nil {
		return RepositoryStatus{}, err
	}

	return s.parseGitStatusOutput(bytes.NewReader(out)), nil
}

func (s Status) parseGitStatusOutput(r io.Reader) RepositoryStatus {
	scanner := bufio.NewScanner(r)

	var localBranch, remoteBranch string
	var ahead, behind int

	//Extract branch name
	for scanner.Scan() {
		//Skip any empty line
		if len(scanner.Text()) < 1 {
			continue
		}

		localBranch, remoteBranch, ahead, behind = s.parseBranchLine(scanner.Text())
		break
	}

	var remoteStatus StatusMessage

	switch {
	case ahead > 0:
		remoteStatus = CommittedChanges
	case behind > 0:
		remoteStatus = RemoteChanges
	default:
		remoteStatus = NoChanges
	}

	var statuses []FileStatus
	for scanner.Scan() {
		if len(scanner.Text()) < 1 {
			continue
		}

		statuses = append(statuses, parseFileLine(scanner.Text()))
	}

	staged, unstaged, untracked, ignored := s.calculateTotals(statuses)

	var localStatus StatusMessage

	switch {
	case staged+unstaged > 0:
		localStatus = UncommittedChanges
	case s.Strict && untracked > 0:
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

func (s Status) parseBranchLine(input string) (string, string, int, int) {
	if s.Verbose {
		s.Outputter.Debug(input)
	}
	// Example line:
	//## develop...origin/develop [ahead 1, behind 18]

	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)

	//check if input is a status branch line output
	scanner.Scan()

	if scanner.Text() != "##" {
		return "", "", 0, 0
	}

	//read next word and return the branch name(s)
	scanner.Scan()
	b := strings.Split(scanner.Text(), "...")

	localBranch := b[0]
	remoteBranch := ""

	if len(b) > 1 {
		remoteBranch = b[1]
	}

	var ahead = s.parseCommitsAhead(input)
	var behind = s.parseCommitsBehind(input)

	return localBranch, remoteBranch, ahead, behind
}

func (s Status) parseCommitsAhead(input string) int {
	var r = regexp.MustCompile(`ahead ([0-9]+)`)

	ahead := r.FindStringSubmatch(input)
	if ahead == nil {
		return 0
	}

	i, err := strconv.Atoi(ahead[1])
	if err != nil {
		s.Outputter.Error(err)
		os.Exit(2)
	}

	return i
}

func (s Status) parseCommitsBehind(input string) int {
	var r = regexp.MustCompile(`behind ([0-9]+)`)

	behind := r.FindStringSubmatch(input)
	if behind == nil {
		return 0
	}

	i, err := strconv.Atoi(behind[1])
	if err != nil {
		s.Outputter.Error(err)
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

func (s Status) calculateTotals(fileStatuses []FileStatus) (staged, unstaged, untracked, ignored int) {
	for _, fs := range fileStatuses {
		if s.Verbose {
			s.Outputter.Debug(fs.Text)
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
