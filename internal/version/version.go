package version

import "fmt"

// These values are injected at build time via -ldflags.
// Fallbacks are used if not set.
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

func String() string {
	return fmt.Sprintf(
		"portarbiter %s\ncommit: %s\nbuilt:  %s",
		Version,
		GitCommit,
		BuildDate,
	)
}

