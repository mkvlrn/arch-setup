package start

// Config describes data from config.json.
type Config struct {
	BasePackages []string `json:"basePackages"`
	MainPackages []string `json:"mainPackages"`

	RepoHTTP string `json:"repoHttp"`
	RepoSSH  string `json:"repoSsh"`

	MirrorListPath  string `json:"mirrorListPath"`
	MirrorListCheck string `json:"mirrorListCheck"`

	XdgMkDir []string `json:"xdgMkDir"`
	XdgRmRf  []string `json:"xdgRmRf"`

	Username string
	HomeDir  string
	RepoDir  string
	TempDir  string
}
