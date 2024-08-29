package config

import (
	"fmt"
	"runtime/debug"
)

// These variables are initialized externally during the build. See the Makefile.
var GitCommit string
var GitLastTag string
var GitExactTag string
var BuildDate string

func GetVersion() string {
	if GitExactTag == "undefined" {
		GitExactTag = ""
	}

	if GitExactTag != "" {
		// we are exactly on a tag --> release version
		return GitLastTag
	}

	if GitLastTag != "" {
		// not exactly on a tag --> dev version
		return fmt.Sprintf("%s-dev-%.10s", GitLastTag, GitCommit)
	}

	// we don't have commit information, try golang build info
	if commit, dirty, err := getCommitAndDirty(); err == nil {
		if dirty {
			return fmt.Sprintf("dev-%.10s-dirty", commit)
		}
		return fmt.Sprintf("dev-%.10s", commit)
	}

	return "dev-unknown"
}

func getCommitAndDirty() (commit string, dirty bool, err error) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", false, fmt.Errorf("unable to read build info")
	}

	var commitFound bool

	// get the commit and modified status
	// (that is the flag for repository dirty or not)
	for _, kv := range info.Settings {
		switch kv.Key {
		case "vcs.revision":
			commit = kv.Value
			commitFound = true
		case "vcs.modified":
			if kv.Value == "true" {
				dirty = true
			}
		}
	}

	if !commitFound {
		return "", false, fmt.Errorf("no commit found")
	}

	return commit, dirty, nil
}
