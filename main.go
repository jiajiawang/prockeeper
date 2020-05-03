package main

import (
	"flag"

	prockeeper "github.com/jiajiawang/prockeeper/src"
)

func init() {
	var help bool
	flag.BoolVar(&help, "help", false, "Show usage")
	flag.Parse()
	if help {
		prockeeper.Usage()
	}
}

func main() {
	manager := &prockeeper.Manager{}
	manager.Run()
}
