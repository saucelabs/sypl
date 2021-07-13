package example

import (
	"bufio"
	"bytes"

	"github.com/saucelabs/sypl"
	"github.com/saucelabs/sypl/level"
)

// NonChainable is a non-chainable example of creating and setting up a new sypl
// logger. It writes to a custom buffer.
func NonChainable() string {
	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)

	// Creates logger, and name it.
	testingLogger := sypl.New("Testing Logger")

	// Creates an `Output`. In this case, called "Console" that will write to a
	// custom buffer, and max print level @ `Info`.
	ConsoleToStdOut := sypl.NewOutput("Console", level.Info, bufWriter)

	// Creates a `Processor`. It will prefix all messages.
	Prefixer := func(prefix string) *sypl.Processor {
		return sypl.NewProcessor("Prefixer", func(message *sypl.Message) {
			message.SetProcessedContent(prefix + message.GetProcessedContent())
		})
	}

	// Adds `Processor` to `Output`.
	ConsoleToStdOut.AddProcessor(Prefixer("My Prefix - "))

	// Adds `Output` to logger.
	testingLogger.AddOutput(ConsoleToStdOut)

	// Writes: "My Prefix - Test message"
	testingLogger.Print(level.Info, "Test info message")

	bufWriter.Flush()

	return buf.String()
}
