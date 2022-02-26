//go:build tools

package ambient

// Source: https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

import (
	// Required for go generate.
	_ "github.com/vburenin/ifacemaker"
)
