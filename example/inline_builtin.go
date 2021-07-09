package example

import (
	"github.com/saucelabs/sypl"
)

// InlineBuiltin same as `ChainableBuiltin` but using the inline initialization.
// nolint:lll
func InlineBuiltin() {
	sypl.New("Testing Logger", sypl.Console(sypl.INFO).AddProcessor(sypl.Prefixer("My Prefix - "))).Infoln("Test info message")
}
