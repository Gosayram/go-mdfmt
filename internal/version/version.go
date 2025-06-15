// Package version provides version information and build details for the application.
package version

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	// ShortCommitHashLength defines the length for shortened commit hashes
	ShortCommitHashLength = 7
	// UnknownValue represents unknown build information
	UnknownValue = "unknown"
)

// Build-time variables set by linker flags
var (
	Version     = "dev"
	Commit      = UnknownValue
	Date        = UnknownValue
	BuiltBy     = UnknownValue
	BuildNumber = "0"
)

// GetVersion returns the complete version string
func GetVersion() string {
	if BuildNumber != "0" && BuildNumber != "" {
		return fmt.Sprintf("%s (build %s)", Version, BuildNumber)
	}
	return Version
}

// GetFullVersionInfo returns detailed version information in modern format
func GetFullVersionInfo() string {
	var parts []string

	// Main version line
	versionLine := fmt.Sprintf("mdfmt %s", GetVersion())
	if Commit != UnknownValue && Commit != "" {
		shortCommit := Commit
		if len(Commit) > ShortCommitHashLength {
			shortCommit = Commit[:ShortCommitHashLength]
		}
		versionLine += fmt.Sprintf(" (%s)", shortCommit)
	}
	parts = append(parts, versionLine)

	// Build info line
	var buildInfo []string
	if Date != UnknownValue && Date != "" {
		// Format date more nicely
		formattedDate := strings.ReplaceAll(Date, "_", " ")
		buildInfo = append(buildInfo, fmt.Sprintf("built %s", formattedDate))
	}
	if BuiltBy != UnknownValue && BuiltBy != "" {
		buildInfo = append(buildInfo, fmt.Sprintf("by %s", BuiltBy))
	}
	buildInfo = append(buildInfo,
		fmt.Sprintf("with %s", runtime.Version()),
		fmt.Sprintf("for %s/%s", runtime.GOOS, runtime.GOARCH))

	if len(buildInfo) > 0 {
		parts = append(parts, strings.Join(buildInfo, " "))
	}

	return strings.Join(parts, "\n")
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
