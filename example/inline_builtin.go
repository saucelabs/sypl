package example

import (
	"github.com/saucelabs/sypl"
	"github.com/saucelabs/sypl/level"
)

// InlineBuiltin same as `ChainableBuiltin` but using the inline initialization.
// nolint:lll
func InlineBuiltin() {
	sypl.New("Testing Logger", sypl.Console(level.Info).AddProcessor(sypl.Prefixer("My Prefix - "))).Infoln("Test info message")
}
