package todocli

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/go-yaml/yaml"
)

var cachedConfig *Config

func getConfig() *Config {
	if cachedConfig != nil {
		return cachedConfig
	}

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	configPath := fmt.Sprintf("%s/.todocli.yaml", usr.HomeDir)

	if _, err := os.Stat(configPath); err != nil {
		panic("Can't find .todocli.yaml")
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	config := Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	cachedConfig = &config

	return &config
}
