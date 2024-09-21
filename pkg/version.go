package sfid

import "github.com/gnames/gnlib/ent/gnvers"

var (
	// Version provides the current version of `sfid`.
	Version = "v0.0.1"

	// Build provides a timestamp when `sfid` was compiled.
	Build = "N/A"
)

// GetVersion returns BHLnames version and build information.
func GetVersion() gnvers.Version {
	return gnvers.Version{Version: Version, Build: Build}
}
