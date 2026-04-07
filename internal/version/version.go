package version

import "runtime/debug"

var (
	// version is a semver version (https://semver.org).
	version = "1.62.1" // x-releaser-pleaser-version

	// versionPrerelease is a semver version pre-release identifier (https://semver.org).
	//
	// For final releases, we set this to an empty string.
	versionPrerelease = "dev"

	// Version of the hcloud CLI.
	Version = func() string {
		if versionPrerelease != "" {
			return version + "-" + versionPrerelease
		}
		return version
	}()

	// Can be set by goreleaser because debug.ReadBuildInfo() is not available for goreleaser builds
	commit = ""

	// Commit is the latest full commit hash during build time
	Commit = func() string {
		if commit != "" {
			return commit
		}
		return getSettingsValue("vcs.revision", "unknown")
	}()

	// Can be set by goreleaser because debug.ReadBuildInfo() is not available for goreleaser builds
	commitDate = ""

	// CommitDate is the timestamp of the latest commit during build time in RFC3339
	CommitDate = func() string {
		if commitDate != "" {
			return commitDate
		}
		return getSettingsValue("vcs.time", "unknown")
	}()

	// Can be set by goreleaser because debug.ReadBuildInfo() is not available for goreleaser builds
	modified = ""

	// Modified specifies whether the git worktree was dirty during build time
	Modified = func() bool {
		if modified != "" {
			return modified == "true"
		}
		return getSettingsValue("vcs.modified", "false") == "true"
	}()

	// used for getSettingsValue
	info, ok = debug.ReadBuildInfo()
)

// getSettingsValue is a helper for getting values from debug.ReadBuildInfo()
// This is only a fallback for builds that do not use goreleaser, since goreleaser
// usually injects the above variables using ldflags. debug.ReadBuildInfo() will be
// used for example when installing using 'go install'.
func getSettingsValue(key, def string) string {
	if !ok {
		return def
	}
	for _, setting := range info.Settings {
		if setting.Key == key {
			return setting.Value
		}
	}
	return def
}
