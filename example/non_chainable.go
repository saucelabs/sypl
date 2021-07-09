package example

import (
	"bufio"
	"bytes"

	"github.com/saucelabs/sypl"
)

// NonChainable is a non-chainable example of creating and setting up a new sypl
// logger. It writes to a custom buffer.
func NonChainable() string {
	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)

	// Creates logger, and name it.
	testingLogger := sypl.New("Testing Logger")

	// Creates an `Output`. In this case, called "Console" that will print to
	// `stdout` and max print level @ `INFO`.
	ConsoleToStdOut := sypl.NewOutput("Console", sypl.INFO, bufWriter)

	// Creates a `Processor`. It will prefix all messages.
	Prefixer := func(prefix string) *sypl.Processor {
		return sypl.NewProcessor("Prefixer", func(message *sypl.Message) {
			message.ContentProcessed = prefix + message.ContentProcessed
		})
	}

	// Adds `Processor` to `Output`.
	ConsoleToStdOut.AddProcessor(Prefixer("My Prefix - "))

	// Adds `Output` to logger.
	testingLogger.AddOutput(ConsoleToStdOut)

	// Prints: "My Prefix - Test message"
	testingLogger.Print(sypl.INFO, "Test info message")

	bufWriter.Flush()

	return buf.String()
}
