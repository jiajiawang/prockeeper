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
		Dir     string
	}
}

// ParseConfig ...
func ParseConfig(path string) *Config {
	file, err := ioutil.ReadFile(path)
	CheckError(err)

	config := &Config{}

	err = yaml.Unmarshal(file, config)
	CheckError(err)
	return config
}
