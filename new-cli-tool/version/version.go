package version

import "fmt"

var (
	semVer    string
	buildDate string
)

func ExpressVersion() string {
	version := fmt.Sprintf("semVer: %s, buildDate: %s", semVer, buildDate)
	return version
}

func GetSemVer() string {
	return semVer
}
