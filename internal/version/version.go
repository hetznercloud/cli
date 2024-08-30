package version

var (
	// version is a semver version (https://semver.org).
	version = "1.47.0" // x-release-please-version

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
)
