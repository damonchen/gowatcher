package config

type Command struct {
	Type string
	Cmd  string
	Args []string
	Envs []string
}

type Config struct {
	IncludePaths  []string
	ExcludedPaths []string
	Command       []Command
}
