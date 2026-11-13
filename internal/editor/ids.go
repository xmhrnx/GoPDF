package editor

import "sync/atomic"

var (
	lastID      atomic.Int64
	lastVersion atomic.Int64
)

// NewID returns a process-unique annotation ID.
func NewID() int64 { return lastID.Add(1) }

// nextVersion returns a process-unique document version.
func nextVersion() int { return int(lastVersion.Add(1)) }
