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
	"github.com/saucelabs/sypl/fields"
	"github.com/saucelabs/sypl/flag"
	"github.com/saucelabs/sypl/formatter"
	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/message"
	"github.com/saucelabs/sypl/options"
	"github.com/saucelabs/sypl/output"
	"github.com/saucelabs/sypl/processor"
	"github.com/saucelabs/sypl/safebuffer"
	"github.com/saucelabs/sypl/shared"
	"github.com/saucelabs/sypl/status"
)

//////
// Helpers
//////

// Checks if `src` has `texts`. If not, it prints a message.
//nolint:forbidigo
func stringContains(src string, texts ...string) bool {
	contains := true

	for _, text := range texts {
		if !strings.Contains(src, text) {
			contains = false

			fmt.Printf("Expected %s to contain %s\n", src, text)
		}
	}

	return contains
}

//////
// Examples.
//////

// NonChained is a non-chained example of creating, and setting up a `sypl`
// logger. It writes to a custom buffer.
func ExampleNew_notChained() {
	buf := new(strings.Builder)

	// Creates logger, and name it.
	testingLogger := sypl.New("Testing Logger")

	// Creates an `Output`. In this case, called "Console" that will print to
	// `stdout` and max print level @ `Info`.
	ConsoleToStdOut := output.New("Console", level.Info, buf)

	// Creates a `Processor`. It will prefix all messages.
	Prefixer := func(prefix string) processor.IProcessor {
		return processor.New("Prefixer", func(message message.IMessage) error {
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
		// Creates a `Processor`. It will prefix all messages. It will only
		// prefix messages for this specific `Output`, and @ `Error` level.
		AddOutputs(output.New("Error", level.Error, os.Stdout).
			AddProcessors(func(prefix string) processor.IProcessor {
				return processor.New("Prefixer", func(message message.IMessage) error {
					if message.GetLevel() == level.Error {
						message.GetContent().SetProcessed(prefix + message.GetContent().GetProcessed())
					}

					return nil
				})
			}(shared.DefaultPrefixValue))).
		// Prints:
		// My Prefix - Test error message
		Println(level.Error, "Test error message")

	// Output:
	// My Prefix - Test error message
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
	Console1ToStdOut := output.New("Buffer 1", level.Info, &c1buf)

	var c2buf safebuffer.Buffer
	Console2ToStdOut := output.New("Buffer 2", level.Info, &c2buf)

	var c3buf safebuffer.Buffer
	Console3ToStdOut := output.New("Buffer 3", level.Info, &c3buf)

	// Creates a `Processor`. It will `prefix` all messages with the Output, and
	// Processor names.
	Prefixer := func() processor.IProcessor {
		return processor.New("Prefixer", func(message message.IMessage) error {
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
		return processor.New("SuffixBasedOnTag", func(message message.IMessage) error {
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
		PrintlnWithOptions(&options.Options{
			Flag: flag.Force,
		}, level.Debug, shared.DefaultContentOutput).

		// Message will be processed, but not printed.
		PrintlnWithOptions(&options.Options{
			Flag: flag.Mute,
		}, level.Info, shared.DefaultContentOutput).

		// Message will not be processed, but printed.
		PrintlnWithOptions(&options.Options{
			Flag: flag.Skip,
		}, level.Info, shared.DefaultContentOutput).

		// Should not print - restricted by level.
		Debugln(shared.DefaultContentOutput).

		// SkipAndForce message will not be processed, but will be printed
		// independent of `Level` restrictions.
		PrintlnWithOptions(&options.Options{
			Flag: flag.SkipAndForce,
		}, level.Debug, shared.DefaultContentOutput).

		// Message will not be processed, neither printed.
		PrintlnWithOptions(&options.Options{
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

// Text formatter example.
func ExampleNew_textFormatter() {
	buf, o := output.SafeBuffer(level.Info)
	o.SetFormatter(formatter.Text())

	// Creates logger, and name it.
	sypl.New(shared.DefaultComponentNameOutput).
		AddOutputs(o).
		PrintlnWithOptions(&options.Options{
			Fields: fields.Fields{
				"field1": "value1",
				"field2": "value2",
				"field3": "value3",
			},
		}, level.Info, shared.DefaultContentOutput)

	s := buf.String()

	fmt.Print(
		strings.Contains(s, shared.DefaultContentOutput),
		strings.Contains(s, "field1=value1"),
		strings.Contains(s, "field2=value2"),
		strings.Contains(s, "field3=value3"),
		strings.Contains(s, "component="),
		strings.Contains(s, "level="),
		strings.Contains(s, "timestamp="),
	)

	// Prints:
	//
	// component=componentNameTest level=info field1=value1 field2=value2 field3=value3 timestamp=2021-08-10T22:50:36-07:00

	// output:
	// true true true true true true true
}

// JSON formatter example.
func ExampleNew_jsonFormatter() {
	buf, o := output.SafeBuffer(level.Info)
	o.SetFormatter(formatter.JSON())

	// Creates logger, and name it.
	sypl.New(shared.DefaultComponentNameOutput).
		AddOutputs(o).
		PrintWithOptions(&options.Options{
			Fields: fields.Fields{
				"field1": "value1",
				"field2": 1,
				"field3": true,
				"field4": []string{"1", "2"},
			},
		}, level.Info, shared.DefaultContentOutput)

	s := buf.String()

	fmt.Print(
		strings.Contains(s, `"component"`),
		strings.Contains(s, `"message"`),
		strings.Contains(s, `"field1"`),
		strings.Contains(s, `"field2"`),
		strings.Contains(s, `"field3"`),
		strings.Contains(s, `"field4"`),
		strings.Contains(s, `"level"`),
		strings.Contains(s, `"timestamp"`),
	)

	// Prints:
	//
	// {
	// 	"component": "componentNameTest",
	// 	"content": "contentTest",
	// 	"field1": "value1",
	// 	"field2": 1,
	// 	"field3": true,
	// 	"field4": [
	// 		"1",
	// 		"2"
	// 	],
	// 	"level": "info",
	// 	"timestamp": "2021-08-10T23:27:25-07:00"
	// }

	// output:
	// true true true true true true true true
}

// Simulates a problematic processor.
func ExampleNew_errorSimulator() {
	// Creates logger, and name it.
	sypl.New(shared.DefaultComponentNameOutput).
		AddOutputs(output.Console(level.Info, processor.ErrorSimulator("Test"))).
		Infoln(shared.DefaultContentOutput)

	// Prints:
	//
	// 2021/08/10 20:00:56 [sypl] [Error] Output: "Console" Processor: "ErrorSimulator" Error: "Test" Original Message: "contentTest"
}

// Child loggers example.
func ExampleNew_childLoggers() {
	// Creates logger, and name it.
	k8Logger := sypl.New("k8").
		AddOutputs(output.Console(level.Info).SetFormatter(formatter.Text()))

	k8Logger.Infoln("k8 connected")

	podLogger := k8Logger.New("pod")
	podLogger.Infoln("pod created")

	k8Logger.GetOutput("Console").SetStatus(status.Disabled)

	k8Logger.Infoln("k8 connected")
	podLogger.Infoln("pod created")

	// Prints:
	//
	// k8 connected component=k8 level=info timestamp=2021-08-11T09:24:13-07:00
	// pod created component=pod level=info timestamp=2021-08-11T09:24:13-07:00
}

// PrintMessagesToOutputs example.
func ExampleNew_printMessagesToOutputs() {
	// Creates logger, and name it.
	sypl.New("pod").AddOutputs(
		output.New("Console 1", level.Trace, os.Stdout).SetFormatter(formatter.Text()),
		output.New("Console 2", level.Trace, os.Stdout).SetFormatter(formatter.Text()),
		output.New("Console 3", level.Trace, os.Stdout).SetFormatter(formatter.Text()),
	).PrintMessagesToOutputs(
		sypl.MessageToOutput{OutputName: "Console 1", Level: level.Info, Content: "Test 1\n"},
		sypl.MessageToOutput{OutputName: "Console 1", Level: level.Debug, Content: "Test 2\n"},
		sypl.MessageToOutput{OutputName: "Console 2", Level: level.Info, Content: "Test 3\n"},
		sypl.MessageToOutput{OutputName: "Console 4", Level: level.Info, Content: "Test 4\n"},
	)

	// Prints:
	//
	// output=Console 2 level=Info message=Test 3
	// output=Console 1 level=Debug message=Test 2
	// output=Console 1 level=Info message=Test 1
}

// StdErr output example.
func ExampleNew_stdErrOutputExample() {
	// Creates logger, and name it.
	sypl.New(shared.DefaultComponentNameOutput, output.StdErr()).
		Infoln(shared.DefaultContentOutput).
		Errorln(shared.DefaultContentOutput)

	// Prints:
	//
	// contentTest
}

// NewDefault output example.
func ExampleNew_newDefault() {
	// Creates logger, and name it.
	sypl.NewDefault(shared.DefaultComponentNameOutput, level.Trace).
		Infoln(shared.DefaultContentOutput).
		Errorln("error message")

	// Prints:
	//
	// component=componentNameTest output=console level=info timestamp=2021-08-17T19:10:00-07:00 message=contentTest
	// component=componentNameTest output=console level=error timestamp=2021-08-17T19:10:00-07:00 message=error message
}

// PrintOnlyIfTagged output example.
func ExampleNew_printOnlyIfTagged() {
	// Creates logger, and name it.
	sypl.NewDefault(shared.DefaultComponentNameOutput, level.Trace, processor.PrintOnlyIfTagged("testTag")).
		Infoln(shared.DefaultContentOutput).
		PrintlnWithOptions(&options.Options{
			Tags: []string{"testTag"},
		}, level.Info, shared.DefaultContentOutput)

	// Prints:
	//
	// component=componentNameTest output=console level=info timestamp=2021-08-17T19:19:57-07:00 message=contentTest
}

// Updating outputs' max levels example.
func ExampleNew_updateOutputsMaxLevel() {
	// Creates logger, and name it.
	l := sypl.New(shared.DefaultComponentNameOutput).AddOutputs(
		output.New("Console 1", level.Info, os.Stdout),
		output.New("Console 2", level.Debug, os.Stdout),
		output.New("Console 3", level.Trace, os.Stdout),
	)

	l.PrintlnWithOptions(&options.Options{
		OutputsNames: []string{"Console 1"},
	}, level.Info, l.GetMaxLevel())

	l.SetMaxLevel(level.Info)

	l.PrintlnWithOptions(&options.Options{
		OutputsNames: []string{"Console 1"},
	}, level.Info, l.GetMaxLevel())

	// output:
	// map[Console 1:info Console 2:debug Console 3:trace]
	// map[Console 1:info Console 2:info Console 3:info]
}

// PrintMessagesToOutputsWithOptions example.
func ExampleNew_printMessagesToOutputsWithOptions() {
	// Creates logger, and name it.
	l := sypl.New(shared.DefaultComponentNameOutput).AddOutputs(
		output.New("Console 1", level.Trace, os.Stdout).SetFormatter(formatter.Text()),
		output.New("Console 2", level.Trace, os.Stdout).SetFormatter(formatter.Text()),
	)

	l.PrintMessagesToOutputsWithOptions(&options.Options{
		Fields: fields.Fields{"1": 2},
	},
		sypl.MessageToOutput{Content: fmt.Sprintln(shared.DefaultContentOutput), Level: level.Info, OutputName: "Console 1"},
		sypl.MessageToOutput{Content: fmt.Sprintln(shared.DefaultContentOutput), Level: level.Warn, OutputName: "Console 2"},
	)

	// Prints:
	//
	// component=componentNameTest output=console 2 level=warn timestamp=2021-08-19T11:48:34-07:00 message=contentTest 1=2
	// component=componentNameTest output=console 1 level=info timestamp=2021-08-19T11:48:34-07:00 message=contentTest 1=2
}

// PrintNewLine example.
func ExampleNew_printNewLine() {
	// Creates logger, and name it.
	l := sypl.New(shared.DefaultComponentNameOutput, output.Console(level.Info))

	l.Infoln(shared.DefaultContentOutput)
	l.PrintNewLine()
	l.Infoln(shared.DefaultContentOutput)

	// output:
	// contentTest
	//
	// contentTest
}

// Global fields example.
func ExampleNew_globalFields() {
	buf, o := output.SafeBuffer(level.Info)

	// Creates logger, and name it.
	l := sypl.New(shared.DefaultComponentNameOutput, o.SetFormatter(formatter.Text()))
	l.SetFields(fields.Fields{"a": 1})

	l.Infoln(shared.DefaultContentOutput) // a=1

	l.PrintlnWithOptions(&options.Options{
		Fields: fields.Fields{"a": 2, "b": 3},
	}, level.Info, shared.DefaultContentOutput) // a=2, b=3

	l.Infoln(shared.DefaultContentOutput) // a=1

	fmt.Println(stringContains(buf.String(), "contentTest", "a=1", "a=2", "b=3"))

	// output:
	// true
}

// Logging filtering, and debug capability example.
func ExampleNew_debugAndFilter() {
	// From any SYPL logger, bump all max levels to `info`
	// From any SYPL logger havong a `console` output, bump max levels to `debug`
	// From a logger named `pod`, for its output called `Console`, bump max levels to `trace`
	// From a logger named `pv`, for its output called `o1`, bump max levels to `trace`
	os.Setenv(shared.DebugEnvVar, "info,console:debug,pod:console:trace,pv:o1:trace")
	defer os.Unsetenv(shared.DebugEnvVar)

	// From any SYPL logger, only print the following ones.
	os.Setenv(shared.FilterEnvVar, "pod,svc,vs,np,cm,pv")
	defer os.Unsetenv(shared.FilterEnvVar)

	// Will print, max level bumped to `trace` by `pod:console:trace`.
	sypl.New("pod").AddOutputs(output.Console(level.Error)).Traceln("pod created")

	// Will print, max level bumped to `debug`.
	sypl.New("svc").AddOutputs(output.Console(level.Error)).Debugln("svc created")

	// Will print, max level bumped to `debug`.
	sypl.New("vs").AddOutputs(output.Console(level.Error)).Debugln("vs created")

	// Will print, max level bumped to `debug`.
	sypl.New("np").AddOutputs(output.Console(level.Error)).Debugln("np created")

	// Will print, max level bumped to `debug`.
	sypl.New("cm").AddOutputs(output.Console(level.Error)).Debugln("cm created")

	sypl.New("pv").AddOutputs(
		// Will print, max level bumped to `trace` by `pv:o1:trace`.
		output.New("o1", level.Error, os.Stdout),
		// Will not print, max level bumped to `info` but less than `debug`.
		output.New("o2", level.Error, os.Stdout),
	).Traceln("pv created")

	// Will not print, max level bumped to `debug` but it's filtered out.
	sypl.New("dp").AddOutputs(output.Console(level.Error)).Debugln("dp created")

	// output:
	// pod created
	// svc created
	// vs created
	// np created
	// cm created
	// pv created
}

// Sypl as io.Writer, including sub-loggers thru interfaces.
//nolint:gosimple
func ExampleNew_ioWriter() {
	buf, o := output.SafeBuffer(level.Trace)

	// Creates logger, and name it.
	l := sypl.New(shared.DefaultComponentNameOutput, o.SetFormatter(formatter.Text()))

	l.SetDefaultIoWriterLevel(level.Info)

	if _, err := l.Write([]byte(shared.DefaultContentOutput + "1 \n")); err != nil {
		fmt.Println(false)
	}

	l.SetDefaultIoWriterLevel(level.Warn)

	if _, err := l.Write([]byte(shared.DefaultContentOutput + "2 \n")); err != nil {
		fmt.Println(false)
	}

	var m sypl.ISypl

	m = l.New("sub")

	m.SetDefaultIoWriterLevel(level.Debug)

	if _, err := m.Write([]byte(shared.DefaultContentOutput + "3 \n")); err != nil {
		fmt.Println(false)
	}

	m.SetDefaultIoWriterLevel(level.Trace)

	if _, err := m.Write([]byte(shared.DefaultContentOutput + "4 \n")); err != nil {
		fmt.Println(false)
	}

	// Buf content should be...:
	// component=componentNameTest output=buffer level=info timestamp=2022-02-22T09:58:22-08:00 message=contentTest1
	// component=componentNameTest output=buffer level=warn timestamp=2022-02-22T09:58:22-08:00 message=contentTest2
	// component=sub output=buffer level=debug timestamp=2022-02-22T09:58:22-08:00 message=contentTest3
	// component=sub output=buffer level=trace timestamp=2022-02-22T09:58:22-08:00 message=contentTest4

	fmt.Println(stringContains(buf.String(),
		"componentNameTest",
		"sub",
		"contentTest1",
		"contentTest2",
		"contentTest3",
		"contentTest4",
		"info",
		"warn",
		"debug",
		"trace",
	))

	// output:
	// true
}
