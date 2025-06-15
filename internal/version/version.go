// Package version provides version information and build details for the application.
package version

import (
	"fmt"
	"runtime"
)

const (
	// ShortCommitHashLength defines the length for shortened commit hashes
	ShortCommitHashLength = 7
)

// Build-time variables set by linker flags
var (
	Version     = "dev"
	Commit      = "unknown"
	Date        = "unknown"
	BuiltBy     = "unknown"
	BuildNumber = "0"
)

// GetVersion returns the complete version string
func GetVersion() string {
	if BuildNumber != "0" && BuildNumber != "" {
		return fmt.Sprintf("%s (build %s)", Version, BuildNumber)
	}
	return Version
}

// GetFullVersionInfo returns detailed version information
func GetFullVersionInfo() string {
	return fmt.Sprintf(`go-mdfmt version %s
Build commit: %s
Build date: %s
Built by: %s
Go version: %s
OS/Arch: %s/%s`,
		GetVersion(),
		Commit,
		Date,
		BuiltBy,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}

// BuildInfo contains build information
type BuildInfo struct {
	Version   string
	Commit    string
	Date      string
	BuiltBy   string
	GoVersion string
	Platform  string
}

// GetBuildInfo returns the current build information
func GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		BuiltBy:   BuiltBy,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns a formatted version string
func (b BuildInfo) String() string {
	if b.Version == "dev" {
		return fmt.Sprintf("mdfmt %s (%s) built with %s on %s",
			b.Version, b.Commit, b.GoVersion, b.Platform)
	}
	return fmt.Sprintf("mdfmt %s built on %s with %s",
		b.Version, b.Date, b.GoVersion)
}

// Short returns a short version string with version and commit information.
func (b *BuildInfo) Short() string {
	if b.Commit == "" {
		return b.Version
	}
	if len(b.Commit) > ShortCommitHashLength {
		return fmt.Sprintf("%s (%s)", b.Version, b.Commit[:ShortCommitHashLength])
	}
	return fmt.Sprintf("%s (%s)", b.Version, b.Commit)
}
