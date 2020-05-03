package prockeeper

import (
	"flag"
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

var configFile string

func init() {
	flag.StringVar(&configFile, "c", "./prockeeper.yml", "config file")
	flag.Parse()
}

// ParseConfig ...
func ParseConfig() *Config {
	file, err := ioutil.ReadFile(configFile)
	CheckError(err)

	config := &Config{}

	err = yaml.Unmarshal(file, config)
	CheckError(err)
	return config
}
