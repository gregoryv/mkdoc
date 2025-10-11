package main

import (
	_ "embed"
	"runtime/debug"
	"strings"
)

func Version() string {
	prefix := "## ["
	from := strings.Index(changelog, prefix) + len(prefix)
	to := from + strings.Index(changelog[from:], "]")
	return changelog[from:to]
}

// Revision returns e.g. fda2bca74b77956a2908a4172c2a82ed2104f804 or
// empty string if no revision is found.
func Revision(n int) string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			if n > 0 && len(setting.Value) > n {
				return setting.Value[:n]
			}
			return setting.Value
		}
	}
	return ""
}

//go:embed changelog.md
var changelog string
