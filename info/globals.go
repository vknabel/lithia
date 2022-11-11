package info

var Version = "dev"
var Commit string
var Date string
var BuiltBy = "dev"
var Debug bool

func init() {
	if BuiltBy == "dev" {
		Debug = true
	}
}
