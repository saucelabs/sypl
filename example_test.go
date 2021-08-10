// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl_test

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/saucelabs/sypl"
	"github.com/saucelabs/sypl/flag"
	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/message"
	"github.com/saucelabs/sypl/options"
	"github.com/saucelabs/sypl/output"
	"github.com/saucelabs/sypl/processor"
	"github.com/saucelabs/sypl/safebuffer"
	"github.com/saucelabs/sypl/shared"
)

// NonChained is a non-chained example of creating, and setting up a `sypl`
// logger. It writes to a custom buffer.
func ExampleNew_notChained() {
	buf := new(strings.Builder)

	// Creates logger, and name it.
	testingLogger := sypl.New("Testing Logger")

	// Creates an `Output`. In this case, called "Console" that will print to
	// `stdout` and max print level @ `Info`.
	ConsoleToStdOut := output.NewOutput("Console", level.Info, buf)

	// Creates a `Processor`. It will prefix all messages.
	Prefixer := func(prefix string) processor.IProcessor {
		return processor.NewProcessor("Prefixer", func(message message.IMessage) error {
			message.GetContent().SetProcessed(prefix + message.GetContent().GetProcessed())

			return nil
		})
	}

	// Adds `Processor` to `Output`.
	ConsoleToStdOut.AddProcessors(Prefixer(shared.DefaultPrefixValue))

	// Adds `Output` to logger.
	testingLogger.AddOutputs(ConsoleToStdOut)

	// Writes: "My Prefix - Test message"
	testingLogger.Print(level.Info, "Test info message")

	fmt.Println(buf.String())

	// Output:
	// My Prefix - Test info message
}

// Chained is the chained example of creating, and setting up a `sypl` logger.
// It writes to both `stdout` and `stderr`.
func ExampleNew_chained() {
	// Creates logger, and name it.
	sypl.New("Testing Logger").
		// Creates two `Output`s. "Console" and "Error". "Console" will print to
		// `Fatal`, `Error`, and `Info`. "Error" will only print `Fatal`, and
		// `Error` levels.
		AddOutputs(output.NewOutput("Console", level.Info, os.Stdout)).
		// Creates a `Processor`. It will prefix all messages. It will only
		// prefix messages for this specific `Output`, and @ `Error` level.
		AddOutputs(output.NewOutput("Error", level.Error, os.Stderr).
			AddProcessors(func(prefix string) processor.IProcessor {
				return processor.NewProcessor("Prefixer", func(message message.IMessage) error {
					if message.GetLevel() == level.Error {
						message.GetContent().SetProcessed(prefix + message.GetContent().GetProcessed())
					}

					return nil
				})
			}(shared.DefaultPrefixValue))).
		// Prints:
		// Test info message
		Println(level.Info, "Test info message").
		// Prints:
		// Test error message
		// My Prefix - Test error message
		Println(level.Error, "Test error message")

	// Output:
	// My Prefix - Test error message

	// Note: Go "example" parser only captured the last message.
}

// ChainedUsingBuiltin is the chained example of creating, and setting up a
// `sypl` logger using built-in `Output`, and `Processor`. It writes to
// `stdout`, and `stderr`.
func ExampleNew_chainedUsingBuiltin() {
	// Creates logger, and name it.
	sypl.New("Testing Logger").
		// Adds an `Output`. In this case, called "Console" that will print to
		// `stdout` and max print level @ `Info`.
		//
		// Adds a `Processor`. It will prefix all messages.
		AddOutputs(output.Console(level.Info).AddProcessors(processor.Prefixer(shared.DefaultPrefixValue))).
		// Prints: My Prefix - Test info message
		Infoln("Test info message")

	// output:
	// My Prefix - Test info message
}

// inlineUsingBuiltin same as `ChainedUsingBuiltin` but using inline form.
func ExampleNew_inlineUsingBuiltin() {
	sypl.New("Testing Logger", output.Console(level.Info).
		AddProcessors(processor.Prefixer(shared.DefaultPrefixValue))).
		Infoln("Test info message")

	// output:
	// My Prefix - Test info message
}

// printWithOptions demonstrates `sypl` flexibility. `Options` enhances the
// usual `PrintX` methods allowing to specify flags, and tags.
func ExampleNew_printWithOptions() {
	// Creates logger, and name it.
	testingLogger := sypl.New("Testing Logger")

	// Creates 3 `Output`s, all called "Console" that will print to `stdout`, and
	// max print level @ `Info`.
	var c1buf safebuffer.Buffer
	Console1ToStdOut := output.NewOutput("Buffer 1", level.Info, &c1buf)

	var c2buf safebuffer.Buffer
	Console2ToStdOut := output.NewOutput("Buffer 2", level.Info, &c2buf)

	var c3buf safebuffer.Buffer
	Console3ToStdOut := output.NewOutput("Buffer 3", level.Info, &c3buf)

	// Creates a `Processor`. It will `prefix` all messages with the Output, and
	// Processor names.
	Prefixer := func() processor.IProcessor {
		return processor.NewProcessor("Prefixer", func(message message.IMessage) error {
			prefix := fmt.Sprintf("Output: %s Processor: %s Content: ",
				message.GetOutputName(),
				message.GetProcessorName(),
			)

			message.GetContent().SetProcessed(prefix + message.GetContent().GetProcessed())

			return nil
		})
	}

	// Creates a `Processor`. It will `suffix` all messages with the specified
	// `tag`.
	SuffixBasedOnTag := func(tag string) processor.IProcessor {
		return processor.NewProcessor("SuffixBasedOnTag", func(message message.IMessage) error {
			if message.ContainTag(tag) {
				message.GetContent().SetProcessed(message.GetContent().GetProcessed() + " - My Suffix")
			}

			return nil
		})
	}

	// Adds `Processor`s to `Output`s.
	Console1ToStdOut.AddProcessors(Prefixer(), SuffixBasedOnTag("SuffixIt"))
	Console2ToStdOut.AddProcessors(Prefixer(), SuffixBasedOnTag("SuffixIt"))
	Console3ToStdOut.AddProcessors(Prefixer(), SuffixBasedOnTag("SuffixIt"))

	// Adds all `Output`s to logger.
	testingLogger.AddOutputs(Console1ToStdOut, Console2ToStdOut, Console3ToStdOut)

	// Prints with prefix, without suffix.
	testingLogger.Print(level.Info, shared.DefaultContentOutput)

	fmt.Println(strings.EqualFold(c1buf.String(), "Output: Buffer 1 Processor: Prefixer Content: contentTest"))
	fmt.Println(strings.EqualFold(c2buf.String(), "Output: Buffer 2 Processor: Prefixer Content: contentTest"))
	fmt.Println(strings.EqualFold(c3buf.String(), "Output: Buffer 3 Processor: Prefixer Content: contentTest"))

	c1buf.Reset()
	c2buf.Reset()
	c3buf.Reset()

	// Prints with prefix, and suffix.
	testingLogger.PrintWithOptions(&options.Options{
		OutputsNames:    []string{"Buffer 1"},
		ProcessorsNames: []string{"Prefixer", "SuffixBasedOnTag"},
		Tags:            []string{"SuffixIt"},
	}, level.Info, shared.DefaultContentOutput)

	fmt.Println(strings.EqualFold(c1buf.String(), "Output: Buffer 1 Processor: Prefixer Content: contentTest - My Suffix"))

	// fmt.Print(c1buf.String(), c2buf.String(), c3buf.String())

	// output:
	// true
	// true
	// true
	// true
}

// PrintPretty example.
func ExampleNew_printPretty() {
	type TestType struct {
		Key1 string
		Key2 int
	}

	sypl.New("Testing Logger", output.Console(level.Info)).PrintPretty(level.Info, &TestType{
		Key1: "text",
		Key2: 12,
	})

	// output:
	// {
	// 	"Key1": "text",
	// 	"Key2": 12
	// }
}

// Flags example.
func ExampleNew_flags() {
	// Creates logger, and name it.
	sypl.New("Testing Logger", output.Console(level.Info, processor.Prefixer(shared.DefaultPrefixValue))).
		// Message will be processed, and printed independent of `Level`
		// restrictions.
		PrintWithOptionsln(&options.Options{
			Flag: flag.Force,
		}, level.Debug, shared.DefaultContentOutput).

		// Message will be processed, but not printed.
		PrintWithOptionsln(&options.Options{
			Flag: flag.Mute,
		}, level.Info, shared.DefaultContentOutput).

		// Message will not be processed, but printed.
		PrintWithOptionsln(&options.Options{
			Flag: flag.Skip,
		}, level.Info, shared.DefaultContentOutput).

		// Should not print - restricted by level.
		Debugln(shared.DefaultContentOutput).

		// SkipAndForce message will not be processed, but will be printed
		// independent of `Level` restrictions.
		PrintWithOptionsln(&options.Options{
			Flag: flag.SkipAndForce,
		}, level.Debug, shared.DefaultContentOutput).

		// Message will not be processed, neither printed.
		PrintWithOptionsln(&options.Options{
			Flag: flag.SkipAndMute,
		}, level.Debug, shared.DefaultContentOutput)

	// output:
	// My Prefix - contentTest
	// contentTest
	// contentTest
}

// Serror{f|lnf|ln} example.
//nolint:goerr113
func ExampleNew_serrorX() {
	// Creates logger, and name it.
	testingLogger := sypl.New("Testing Logger", output.Console(level.Info))

	sErrorResult := testingLogger.Serror(shared.DefaultContentOutput)

	errExample := errors.New("Failed to reach something")
	sErrorfResult := testingLogger.Serrorf("Failed to do something, %s", errExample)
	sErrorlnfResult := testingLogger.Serrorlnf("Failed to do something, %s", errExample)
	sErrorlnResult := testingLogger.Serrorln(shared.DefaultContentOutput)

	fmt.Print(
		sErrorResult.Error() == shared.DefaultContentOutput,
		sErrorfResult.Error() == "Failed to do something, Failed to reach something",
		sErrorlnfResult.Error() == "Failed to do something, Failed to reach something"+"\n",
		sErrorlnResult.Error() == shared.DefaultContentOutput+"\n",
	)

	// output:
	// contentTestFailed to do something, Failed to reach somethingFailed to do something, Failed to reach something
	// contentTest
	// true true true true
}
