package example

import (
	"os"

	"github.com/saucelabs/sypl"
)

// Chainable is a chainable example of creating and setting up a new sypl
// logger. It writes to `stdout` and `stderr`.
func Chainable() {
	// Creates logger, and name it.
	sypl.New("Testing Logger").
		// Creates two `Output`s. "Console" and "Error". "Console" will print to
		// `Fatal`, `Error`, and `Info`. "Error" will only print `Fatal`, and
		// `Error` levels.
		AddOutput(sypl.NewOutput("Console", sypl.INFO, os.Stderr)).
		// Creates a `Processor`. It will prefix all messages. It will only
		// prefix messages for this specific `Output`, and @ `ERROR` level.
		AddOutput(sypl.NewOutput("Error", sypl.ERROR, os.Stdout).
			AddProcessor(func(prefix string) *sypl.Processor {
				return sypl.NewProcessor("Prefixer", func(message *sypl.Message) {
					if message.Level == sypl.ERROR {
						message.ContentProcessed = prefix + message.ContentProcessed
					}
				})
			}("My Prefix - "))).
		// Prints: Test info message
		Println(sypl.INFO, "Test info message").
		// Prints:
		// Test error message
		// My Prefix - Test error message
		Println(sypl.ERROR, "Test error message")
}
