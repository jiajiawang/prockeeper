package prockeeper

import (
	"fmt"
	"os"
)

// CheckError ...
func CheckError(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}
}
