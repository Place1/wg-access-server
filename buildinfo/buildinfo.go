package buildinfo

import (
    "strings"
   _ "embed"
)

//go:generate sh get_version.sh
//go:embed version.txt
var VersionRaw string

//go:generate sh get_commit.sh
//go:embed commit.txt
var CommitHashRaw string

func Version() string {
    return strings.TrimSpace(VersionRaw)
}

func CommitHash() string {
    return strings.TrimSpace(CommitHashRaw)
}

func ShortCommitHash() string {
    if 7 < len(CommitHash()) {
        return CommitHash()[0:7]
    }
    return CommitHash()
}
