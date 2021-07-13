package example

import (
	"fmt"
	"os"

	"github.com/saucelabs/sypl"
	"github.com/saucelabs/sypl/level"
)

// PrintWithOptions is an example of printing specifying `Options`.
func PrintWithOptions() {
	// Creates logger, and name it.
	testingLogger := sypl.New("Testing Logger")

	// Creates 3 `Output`s, all called "Console" that will print to `stdout`, and
	// max print level @ `Info`.
	Console1ToStdOut := sypl.NewOutput("Console 1", level.Info, os.Stdout)
	Console2ToStdOut := sypl.NewOutput("Console 2", level.Info, os.Stdout)
	Console3ToStdOut := sypl.NewOutput("Console 3", level.Info, os.Stdout)

	// Creates a `Processor`. It will `prefix` all messages with the Output, and
	// Processor names.
	Prefixer := func() *sypl.Processor {
		return sypl.NewProcessor("Prefixer", func(message *sypl.Message) {
			prefix := fmt.Sprintf("Output: %s Processor: %s Content: ",
				message.GetOutput().GetName(),
				message.GetProcessor().GetName(),
			)

			message.SetProcessedContent(prefix + message.GetProcessedContent())
		})
	}

	// Creates a `Processor`. It will `suffix` all messages with the specified
	// `tag`.
	SuffixBasedOnTag := func(tag string) *sypl.Processor {
		return sypl.NewProcessor("SuffixBasedOnTag", func(message *sypl.Message) {
			if message.ContainTag(tag) {
				message.SetProcessedContent(message.GetProcessedContent() + " - My Suffix")
			}
		})
	}

	// Adds `Processor`s to `Output`s.
	Console1ToStdOut.
		AddProcessor(Prefixer()).
		AddProcessor(SuffixBasedOnTag("SuffixIt"))
	Console2ToStdOut.
		AddProcessor(Prefixer()).
		AddProcessor(SuffixBasedOnTag("SuffixIt"))
	Console3ToStdOut.
		AddProcessor(Prefixer()).
		AddProcessor(SuffixBasedOnTag("SuffixIt"))

	// Adds all `Output`s to logger.
	testingLogger.
		AddOutput(Console1ToStdOut).
		AddOutput(Console2ToStdOut).
		AddOutput(Console3ToStdOut)

	// Prints:
	// Output: Console 1 Processor: Prefixer Content: Test info message
	// Output: Console 2 Processor: Prefixer Content: Test info message
	// Output: Console 3 Processor: Prefixer Content: Test info message
	testingLogger.Println(level.Info, "Test info message")

	// Prints:
	// Output: Console 1 Processor: Prefixer Content: Test info message - My Suffix
	testingLogger.PrintWithOptions(&sypl.Options{
		OutputsNames:    []string{"Console 1"},
		ProcessorsNames: []string{"Prefixer", "SuffixBasedOnTag"},
		Tags:            []string{"SuffixIt"},
	}, level.Info, "Test info message")
}
