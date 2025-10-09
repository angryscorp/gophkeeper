package buildinfo

import "fmt"

// BuildInfo holds version and build date metadata.
type BuildInfo struct {
	version string
	date    string
}

// New creates a new BuildInfo instance.
func New(version, date string) BuildInfo {
	return BuildInfo{
		version: version,
		date:    date,
	}
}

// String formats the build information. If fields are empty,
// "N/A" is shown instead.
func (buildInfo BuildInfo) String() string {
	if buildInfo.version == "" {
		buildInfo.version = "N/A"
	}

	if buildInfo.date == "" {
		buildInfo.date = "N/A"
	}

	return fmt.Sprintf("Build version: %s\nBuild date: %s\n", buildInfo.version, buildInfo.date)
}
