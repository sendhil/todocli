package todocli

// Config is a struct representing Configuration parameters for todocli
type Config struct {
	Path     string            `yaml:"path"`
	Mappings map[string]string `yaml:"mappings"`
}
