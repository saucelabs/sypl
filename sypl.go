// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/saucelabs/sypl/flag"
	"github.com/saucelabs/sypl/level"
)

const defaultCallDepth = 2

// Printer defines possible printers.
type Printer interface {
	// PrintWithOptions prints the message with `Options`.
	PrintWithOptions(options *Options, level level.Level, args ...interface{}) *Sypl

	// PrintfWithOptions prints the message, with `Options` according with the
	// specified format.
	PrintfWithOptions(options *Options, level level.Level, format string, args ...interface{}) *Sypl

	// PrintlnWithOptions prints the message (if has content), with `Options`, also
	// adding a new line to the end.
	PrintlnWithOptions(options *Options, level level.Level, args ...interface{}) *Sypl
}

// Options extends `PrintX` capabilities.
type Options struct {
	// Flag define behaviors.
	Flag flag.Flag

	// OutputsNames name of the outputs that should be used to print.
	OutputsNames []string

	// ProcessorsNames name of the processors that should be used.
	ProcessorsNames []string

	// Tags are indicators consumed by `Output`s and `Processor`s.
	Tags []string
}

// NewDefaultOptions creates a new set of options base on default values.
func NewDefaultOptions() *Options {
	return &Options{
		Flag:            flag.None,
		OutputsNames:    []string{},
		ProcessorsNames: []string{},
		Tags:            []string{},
	}
}

// sypl definition. It's able to print messages according to the definition of
// each `output`.
type Sypl struct {
	name    string
	outputs []*Output
}

// GetName returns the sypl name.
func (sypl *Sypl) GetName() string {
	return sypl.name
}

// AddOutput adds an `output` to logger.
//
// Note: This method is chainable.
func (sypl *Sypl) AddOutput(output *Output) *Sypl {
	sypl.outputs = append(sypl.outputs, output)

	return sypl
}

// GetOutput returns the specified output by its index.
func (sypl *Sypl) GetOutput(i int) *Output {
	if i < 0 || i > len(sypl.outputs)-1 {
		return nil
	}

	return sypl.outputs[i]
}

// GetOutputs returns registered outputs.
func (sypl *Sypl) GetOutputs() []*Output {
	return sypl.outputs
}

// Writes to the specified writer.
//
//
// Note: In case of any error, the standard log will be used to highlight the
// case, but IT WILL NOT STOP the application.
func (sypl *Sypl) write(message *Message) {
	if err := message.GetOutput().GetBuiltinLogger().OutputBuiltin(
		defaultCallDepth,
		message.GetProcessedContent(),
	); err != nil {
		log.Println("[sypl] [Error] write: Failed to write to output.", err)
	}
}

// Processor logic of the `process` method.
func (sypl *Sypl) processProcessor(
	output *Output,
	message *Message,
	processorsNames string,
) {
	// Should not process if message is flagged with `Skip` or `SkipAndForce`.
	if message.GetFlag() != flag.Skip || message.GetFlag() != flag.SkipAndForce {
		for _, processor := range output.processors {
			// Should only use named (listed) ones.
			// Should only use `enabled` `Processor`s, see logic in
			// `.Run` method.
			if strings.Contains(processorsNames, processor.GetName()) {
				message.SetProcessor(processor)

				processor.Run(message)
			}
		}
	}
}

// Output logic of the `process` method.
func (sypl *Sypl) processOutput(
	options *Options,
	lvl level.Level,
	content string,
	outputsNames string,
) {
	for _, output := range sypl.outputs {
		// Should only use `enabled` `Output`(s), and named (listed) ones.
		if strings.Contains(outputsNames, output.name) && output.enabled {
			// Message is isolated per `Output`.
			message := NewMessage(sypl, output, nil, lvl, content)
			message.SetFlag(options.Flag)
			message.AddTags(options.Tags...)

			processorsNames := ProcessorsNames(output.processors)

			// Should allows to specify `Processor`(s).
			if len(options.ProcessorsNames) > 0 {
				processorsNames = strings.Join(options.ProcessorsNames, ",")
			}

			sypl.processProcessor(output, message, processorsNames)

			// Should print the message - regardless of the level, if flagged
			// with `Force`.
			if message.GetFlag() == flag.Force {
				sypl.write(message)
			}

			// Should only print if message `level` isn't above `MaxLevel`.
			// Should only print if `level` isn't `None`.
			// Should only print if not flagged with `Mute`.
			if message.GetLevel() != level.None &&
				message.GetLevel() <= output.GetMaxLevel() &&
				message.GetFlag() != flag.Mute {
				sypl.write(message)
			}
		}
	}
}

// Process a message according to logger's registered outputs, and its
// processors. If a list of outputs names is passed, it will only process and
// print if matches against the registered outputs, otherwise all registered
// outputs will be used. If not content if found, nothing is processed or
// printed.
func (sypl *Sypl) process(options *Options, lvl level.Level, content string) *Sypl {
	// Do nothing if message as no context, or flagged with `SkipAndMute`.
	if content == "" || options.Flag == flag.SkipAndMute {
		return sypl
	}

	outputsNames := OutputsNames(sypl.outputs)

	// Should allows to specify `Output`(s).
	if len(options.OutputsNames) > 0 {
		outputsNames = strings.Join(options.OutputsNames, ",")
	}

	sypl.processOutput(options, lvl, content, outputsNames)

	// Should exit if `level` is `Fatal`.
	if lvl == level.Fatal {
		os.Exit(1)
	}

	return sypl
}

// prettify encodes data returning its JSON-stringified version.
func prettify(data interface{}) string {
	buf := new(bytes.Buffer)

	enc := json.NewEncoder(buf)
	enc.SetIndent("", "\t")

	if err := enc.Encode(data); err != nil {
		log.Println("[sypl] [Error] prettify: Failed to encode data")

		return ""
	}

	return buf.String()
}

//////
// Interface Implementation.
//////

// PrintWithOptions prints the message with `Options`.
func (sypl *Sypl) PrintWithOptions(options *Options, level level.Level, args ...interface{}) *Sypl {
	return sypl.process(options, level, fmt.Sprint(args...))
}

// PrintfWithOptions prints the message, with `Options` according with the
// specified format.
func (sypl *Sypl) PrintfWithOptions(options *Options, level level.Level, format string, args ...interface{}) *Sypl {
	return sypl.process(options, level, fmt.Sprintf(format, args...))
}

// PrintlnWithOptions prints the message (if has content), with `Options`, also
// adding a new line to the end.
func (sypl *Sypl) PrintlnWithOptions(options *Options, level level.Level, args ...interface{}) *Sypl {
	return sypl.process(options, level, fmt.Sprintln(args...))
}

//////
// Base methods.
//////

// Print prints the message (if has content).
func (sypl *Sypl) Print(level level.Level, args ...interface{}) *Sypl {
	return sypl.PrintWithOptions(NewDefaultOptions(), level, args...)
}

// Printf prints the message (if has content) according with the specified
// format.
func (sypl *Sypl) Printf(level level.Level, format string, args ...interface{}) *Sypl {
	return sypl.PrintfWithOptions(NewDefaultOptions(), level, format, args...)
}

// Printlnf prints the message (if has content) according with the specified
// format, also adding a new line to the end.
func (sypl *Sypl) Printlnf(level level.Level, format string, args ...interface{}) *Sypl {
	format += "\n"

	return sypl.PrintfWithOptions(NewDefaultOptions(), level, format, args...)
}

// Println prints the message (if has content), also adding a new line to the
// end.
func (sypl *Sypl) Println(level level.Level, args ...interface{}) *Sypl {
	return sypl.PrintlnWithOptions(NewDefaultOptions(), level, args...)
}

// PrintPretty prints data structures as JSON text.
//
// Note: Message will not be processed, but will printed independent of `Level`
// restrictions.
func (sypl *Sypl) PrintPretty(data interface{}) *Sypl {
	return sypl.PrintWithOptions(&Options{
		Flag: flag.SkipAndForce,
	}, level.Info, prettify(data))
}

// PrintlnPretty prints data structures as JSON text, also adding a new line to
// the end.
//
// Note: Message will not be processed, but will printed independent of `Level`
// restrictions.
func (sypl *Sypl) PrintlnPretty(data interface{}) *Sypl {
	return sypl.PrintlnWithOptions(&Options{
		Flag: flag.SkipAndForce,
	}, level.Info, prettify(data))
}

//////
// Convenient methods.
//////

// Error prints the message (if has content) @ the Error level.
func (sypl *Sypl) Error(args ...interface{}) *Sypl {
	return sypl.Print(level.Error, args...)
}

// Errorf prints the message (if has content) according with the specified
// format @ the Error level.
func (sypl *Sypl) Errorf(format string, args ...interface{}) *Sypl {
	return sypl.Printf(level.Error, format, args...)
}

// Errorlnf prints the message (if has content) according with the specified
// format @ the Error level, also adding a new line to the end.
func (sypl *Sypl) Errorlnf(format string, args ...interface{}) *Sypl {
	return sypl.Printlnf(level.Error, format, args...)
}

// Errorln prints the message (if has content), also adding a new line to the
// end @ the Error level.
func (sypl *Sypl) Errorln(args ...interface{}) *Sypl {
	return sypl.Println(level.Error, args...)
}

// Fatal prints the message (if has content), and exit with os.Exit(1).
func (sypl *Sypl) Fatal(args ...interface{}) {
	sypl.Print(level.Fatal, args...)
}

// Fatalf prints the message (if has content) according with the specified
// format, and exit with os.Exit(1).
func (sypl *Sypl) Fatalf(format string, args ...interface{}) {
	sypl.Printf(level.Fatal, format, args...)
}

// Fatallnf prints the message (if has content) according with the specified
// format, also adding a new line to the end, and exit with os.Exit(1).
func (sypl *Sypl) Fatallnf(format string, args ...interface{}) {
	sypl.Printlnf(level.Fatal, format, args...)
}

// Fatalln prints the message (if has content), also adding a new line and the
// end, and exit with os.Exit(1).
func (sypl *Sypl) Fatalln(args ...interface{}) {
	sypl.Println(level.Fatal, args...)
}

// Info prints the message (if has content) @ the Info level.
func (sypl *Sypl) Info(args ...interface{}) *Sypl {
	return sypl.Print(level.Info, args...)
}

// Infof prints the message (if has content) according with the specified format
// @ the Info level.
func (sypl *Sypl) Infof(format string, args ...interface{}) *Sypl {
	return sypl.Printf(level.Info, format, args...)
}

// Infolnf prints the message (if has content) according with the specified
// format @ the Info level, also adding a new line to the end.
func (sypl *Sypl) Infolnf(format string, args ...interface{}) *Sypl {
	return sypl.Printlnf(level.Info, format, args...)
}

// Infoln prints the message (if has content), also adding a new line to the end
// @ the Info level.
func (sypl *Sypl) Infoln(args ...interface{}) *Sypl {
	return sypl.Println(level.Info, args...)
}

// Warn prints the message (if has content) @ the Warn level.
func (sypl *Sypl) Warn(args ...interface{}) *Sypl {
	return sypl.Print(level.Warn, args...)
}

// Warnf prints the message (if has content) according with the specified format
// @ the Warn level.
func (sypl *Sypl) Warnf(format string, args ...interface{}) *Sypl {
	return sypl.Printf(level.Warn, format, args...)
}

// Warnlnf prints the message (if has content) according with the specified
// format @ the Warn level, also adding a new line to the end.
func (sypl *Sypl) Warnlnf(format string, args ...interface{}) *Sypl {
	return sypl.Printlnf(level.Warn, format, args...)
}

// Warnln prints the message (if has content), also adding a new line to the end
// @ the Warn level.
func (sypl *Sypl) Warnln(args ...interface{}) *Sypl {
	return sypl.Println(level.Warn, args...)
}

// Debug prints the message (if has content) @ the Debug level.
func (sypl *Sypl) Debug(args ...interface{}) *Sypl {
	return sypl.Print(level.Debug, args...)
}

// Debugf prints the message (if has content) according with the specified
// format @ the Debug level.
func (sypl *Sypl) Debugf(format string, args ...interface{}) *Sypl {
	return sypl.Printf(level.Debug, format, args...)
}

// Debuglnf prints the message (if has content) according with the specified
// format @ the Debug level, also adding a new line to the end.
func (sypl *Sypl) Debuglnf(format string, args ...interface{}) *Sypl {
	return sypl.Printlnf(level.Debug, format, args...)
}

// Debugln prints the message (if has content), also adding a new line to the
// end @ the Debug level.
func (sypl *Sypl) Debugln(args ...interface{}) *Sypl {
	return sypl.Println(level.Debug, args...)
}

// Trace prints the message (if has content) @ the Trace level.
func (sypl *Sypl) Trace(args ...interface{}) *Sypl {
	return sypl.Print(level.Trace, args...)
}

// Tracef prints the message (if has content) according with the specified
// format @ the Trace level.
func (sypl *Sypl) Tracef(format string, args ...interface{}) *Sypl {
	return sypl.Printf(level.Trace, format, args...)
}

// Tracelnf prints the message (if has content) according with the specified
// format @ the Trace level, also adding a new line to the end.
func (sypl *Sypl) Tracelnf(format string, args ...interface{}) *Sypl {
	return sypl.Printlnf(level.Trace, format, args...)
}

// Traceln prints the message (if has content), also adding a new line to the
// end @ the Trace level.
func (sypl *Sypl) Traceln(args ...interface{}) *Sypl {
	return sypl.Println(level.Trace, args...)
}

// New creates a new custom logger.
//
// Notes: Outputs can be added here, or later using the `AddOutput` method.
func New(name string, outputs ...*Output) *Sypl {
	return &Sypl{
		name:    name,
		outputs: outputs,
	}
}
