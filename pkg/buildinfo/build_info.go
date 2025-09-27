package buildinfo

import "fmt"

type BuildInfo struct {
	version string
	date    string
}

func New(version, date string) BuildInfo {
	return BuildInfo{
		version: version,
		date:    date,
	}
}

func (buildInfo BuildInfo) String() string {
	if buildInfo.version == "" {
		buildInfo.version = "N/A"
	}

	if buildInfo.date == "" {
		buildInfo.date = "N/A"
	}

	return fmt.Sprintf("Build version: %s\nBuild date: %s\n", buildInfo.version, buildInfo.date)
}
