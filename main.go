package main

import prockeeper "github.com/jiajiawang/prockeeper/src"

func main() {
	manager := &prockeeper.Manager{}
	manager.Run()
}
