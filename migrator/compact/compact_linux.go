//go:build linux
// +build linux

package compact

import (
	"syscall"
)

const mmapFlags = syscall.MAP_POPULATE
