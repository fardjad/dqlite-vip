package version

// git_version is set by the build script
var git_version string = "unknown"
var version string = "0.1.0"

func Version() string {
	return version + " (" + git_version + ")"
}
