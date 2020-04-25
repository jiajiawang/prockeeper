package prockeeper

// CheckError ...
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}
