package config

// Config describes data from config.json.
type Config struct {
	Env     Env
	Repo    Repo
	Machine Machine
	Mise    Mise   `json:"mise"`
	Pacman  Pacman `json:"pacman"`
	Yay     Yay    `json:"yay"`
	Xdg     Xdg    `json:"xdg"`
	GetNF   GetNF  `json:"getnf"`
}

// Env describes the environment configuration.
type Env struct {
	CI bool
}

// Repo describes the repository configuration.
type Repo struct {
	HTTP string `json:"http"`
	SSH  string `json:"ssh"`
}

// Machine describes the machine configuration.
type Machine struct {
	Username string
	HomeDir  string
	RepoDir  string
	TempDir  string
}

// Mise describes the mise configuration.
type Mise struct {
	Tools    []string   `json:"tools"`
	Settings [][]string `json:"settings"`
}

// Pacman describes the pacman configuration.
type Pacman struct {
	Install   []string `json:"install"`
	Uninstall []string `json:"uninstall"`
}

// Yay describes the yay configuration.
type Yay struct {
	MirrorListPath  string   `json:"mirrorListPath"`
	MirrorListCheck string   `json:"mirrorListCheck"`
	Packages        []string `json:"packages"`
}

// Xdg describes the XDG configuration.
type Xdg struct {
	MkDir []string `json:"mkdir"`
	RmRf  []string `json:"rmrf"`
}

// GetNF describes the getnf configuration.
type GetNF []string
