package git

type LocalBranchName string
type RemoteBranchName string

type LocalBranch struct {
	LocalBranchName  LocalBranchName
	RemoteBranchName RemoteBranchName
	CurrentBranch    bool
}

type RepositoryStatus struct {
	Versioned     bool
	VersionNumber string
	LocalBranch   string
	RemoteBranch  string
	LocalStatus   StatusMessage
	RemoteStatus  StatusMessage
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

type RepositoryRemote struct {
	Fetch string
	Push  string
}
