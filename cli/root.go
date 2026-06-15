// Package cli assembles the qq command tree from the qq domain on top of the
// any-cli/kit framework.
package cli

import (
	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/qq-cli/qq"
)

// Build metadata, set via -ldflags at release time.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// NewApp assembles the kit application from the qq domain.
func NewApp() *kit.App {
	id := qq.BaseIdentity()
	id.Version = Version

	app := kit.New(id, kit.WithDefaults(qq.Defaults))
	qq.Register(app)

	return app
}
