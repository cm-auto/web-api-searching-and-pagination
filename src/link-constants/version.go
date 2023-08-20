package linkconstants

var version string = ""
func GetVersion() string{
	if version == ""{
		panic("version has not been specified!")
	}
	return version
}
