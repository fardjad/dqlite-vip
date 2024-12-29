package version

// git_version is set by the build script
var gitVersion string = "unknown"
var staticVersion string = "0.0.0"

func Version() string {
	return staticVersion + " (" + gitVersion + ")"
}
