package version

var (
	// Version is the main version number that is being run at the moment.
	Version = "1.41.1" // x-release-please-version

	// Tag is A pre-release marker for the Version. If this is ""
	// (empty string) then it means that it is a final release. Otherwise, this
	// is a pre-release such as "dev" (in development), "beta", "rc1", etc.
	//
	// For releases, GoReleaser will automatically set this to an empty string.
	Tag = "dev"

	// FullVersion is the full version string, including the prerelease tag.
	// This is dynamically generated based on the Version and Tag variables.
	FullVersion = func() string {
		s := Version
		if Tag != "" {
			s += "-" + Tag
		}
		return s
	}()
)
