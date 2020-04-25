package prockeeper

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config ...
type Config struct {
	Services []*struct {
		Name    string
		Command string
	}
}

// ParseConfig ...
func ParseConfig() *Config {
	file, err := ioutil.ReadFile("./prockeeper.yml")
	CheckError(err)

	config := &Config{}

	err = yaml.Unmarshal(file, config)
	CheckError(err)
	return config
}
