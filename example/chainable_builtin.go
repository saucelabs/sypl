package example

import (
	"github.com/saucelabs/sypl"
)

// ChainableBuiltin is a chainable example of creating and setting up a new sypl
// logger using buil-in outputs and processors. It writes to `stdout` and
// `stderr`.
func ChainableBuiltin() {
	// Creates logger, and name it.
	sypl.New("Testing Logger").
		// Adds an `Output`. In this case, called "Console" that will print to
		// `stdout` and max print level @ `INFO`.
		//
		// Adds a `Processor`. It will prefix all messages.
		AddOutput(sypl.Console(sypl.INFO).AddProcessor(sypl.Prefixer("My Prefix - "))).
		// Prints: My Prefix - Test info message
		Infoln("Test info message")
}
