package config

type Task struct {
	Commands []string `yaml:"commands"`
}

type Service struct {
	Name string `yaml:"name"`
	Tasks map[string]Task `yaml:"tasks"`
}

// Config is the top-level, root config file
type Config struct {
	Services []string `yaml:"services"`
	RegisteredServices map[string]*Service `yaml:"-"`
	Changes []string `yaml:"-"`
}
