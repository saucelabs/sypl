// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl

import (
	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/message"
	"github.com/saucelabs/sypl/meta"
	"github.com/saucelabs/sypl/options"
	"github.com/saucelabs/sypl/output"
)

// IBasePrinter specifies the foundation for other printers.
type IBasePrinter interface {
	// PrintMessage prints messages. It's a powerful option because it gives
	// full-control over the message. Use `NewMessage` to create the message.
	PrintMessage(messages ...message.IMessage) ISypl

	// PrintWithOptions is a more flexible way of printing, allowing to specify
	// a few message's options. For full-control over the message is possible
	// via `PrintMessage`.
	PrintWithOptions(o *options.Options, l level.Level, args ...interface{}) ISypl

	// PrintfWithOptions prints according with the specified format. It's a more
	// flexible way of printing, allowing to specify a few message's options.
	// For full-control over the message is possible via `PrintMessage`.
	PrintfWithOptions(o *options.Options, l level.Level, format string, args ...interface{}) ISypl

	// PrintfWithOptions prints according with the specified format, also adding
	// a new line to the end. It's a more flexible way of printing, allowing to
	// specify a few message's options. For full-control over the message is
	// possible via `PrintMessage`.
	PrintlnfWithOptions(o *options.Options, l level.Level, format string, args ...interface{}) ISypl

	// PrintfWithOptions prints, also adding a new line to the end. It's a more
	// flexible way of printing, allowing to specify a few message's options.
	// For full-control over the message is possible via `PrintMessage`.
	PrintlnWithOptions(o *options.Options, l level.Level, args ...interface{}) ISypl
}

// IBasicPrinter specifies the basic printers.
type IBasicPrinter interface {
	// Print just prints.
	Print(l level.Level, args ...interface{}) ISypl

	// Printf prints according with the specified format.
	Printf(l level.Level, format string, args ...interface{}) ISypl

	// Printlnf prints according with the specified format, also adding a new
	// line to the end.
	Printlnf(l level.Level, format string, args ...interface{}) ISypl

	// Println prints, also adding a new line to the end.
	Println(l level.Level, args ...interface{}) ISypl
}

// IConvenientPrinter specifies convenient printers.
type IConvenientPrinter interface {
	// PrintPretty prints data structures as JSON text.
	//
	// Notes:
	// - Only exported fields of the data structure will be printed.
	// - Message isn't processed.
	PrintPretty(l level.Level, data interface{}) ISypl

	// PrintlnPretty prints data structures as JSON text, also adding a new line
	// to the end.
	//
	// Notes:
	// - Only exported fields of the data structure will be printed.
	// - Message isn't processed.
	PrintlnPretty(l level.Level, data interface{}) ISypl

	// PrintMessagerPerOutput allows you to concurrently print messages, each
	// one, at the specified level and to the specified output.
	//
	// Note: If the named output doesn't exits, the message will not be printed.
	PrintMessagesToOutputs(messagesToOutputs ...MessageToOutput) ISypl
}

// ILeveledPrinter specifies the leveled printers.
type ILeveledPrinter interface {
	// Fatal prints, and exit with os.Exit(1).
	Fatal(args ...interface{}) ISypl

	// Fatalf prints according with the format, and exit with os.Exit(1).
	Fatalf(format string, args ...interface{}) ISypl

	// Fatallnf prints according with the format, also adding a new line to the
	// end, and exit with os.Exit(1).
	Fatallnf(format string, args ...interface{}) ISypl

	// Fatalln prints, also adding a new line and the end, and exit with
	// os.Exit(1).
	Fatalln(args ...interface{}) ISypl

	// Error prints @ the Error level.
	Error(args ...interface{}) ISypl

	// Errorf prints according with the format @ the Error level.
	Errorf(format string, args ...interface{}) ISypl

	// Errorlnf prints according with the format @ the Error level, also adding
	// a new line to the end.
	Errorlnf(format string, args ...interface{}) ISypl

	// Errorln prints, also adding a new line to the end @ the Error level.
	Errorln(args ...interface{}) ISypl

	// Serror prints like Error, and returns an error with the non-processed
	// content.
	Serror(args ...interface{}) error

	// Serrorf prints like Errorf, and returns an error with the non-processed
	// content.
	Serrorf(format string, args ...interface{}) error

	// Serrorlnf prints like Errorlnf, and returns an error with the
	// non-processed content.
	Serrorlnf(format string, args ...interface{}) error

	// Serrorln prints like Errorln, and returns an error with the non-processed
	// content.
	Serrorln(args ...interface{}) error

	// Info prints @ the Info level.
	Info(args ...interface{}) ISypl

	// Infof prints according with the specified format @ the Info level.
	Infof(format string, args ...interface{}) ISypl

	// Infolnf prints according with the specified format @ the Info level, also
	// adding a new line to the end.
	Infolnf(format string, args ...interface{}) ISypl

	// Infoln prints, also adding a new line to the end @ the Info level.
	Infoln(args ...interface{}) ISypl

	// Warn prints @ the Warn level.
	Warn(args ...interface{}) ISypl

	// Warnf prints according with the specified format @ the Warn level.
	Warnf(format string, args ...interface{}) ISypl

	// Warnlnf prints according with the specified format @ the Warn level, also
	// adding a new line to the end.
	Warnlnf(format string, args ...interface{}) ISypl

	// Warnln prints, also adding a new line to the end @ the Warn level.
	Warnln(args ...interface{}) ISypl

	// Debug prints @ the Debug level.
	Debug(args ...interface{}) ISypl

	// Debugf prints according with the specified format @ the Debug level.
	Debugf(format string, args ...interface{}) ISypl

	// Debuglnf prints according with the specified format @ the Debug level,
	// also adding a new line to the end.
	Debuglnf(format string, args ...interface{}) ISypl

	// Debugln prints, also adding a new line to the end @ the Debug level.
	Debugln(args ...interface{}) ISypl

	// Trace prints @ the Trace level.
	Trace(args ...interface{}) ISypl

	// Tracef prints according with the specified format @ the Trace level.
	Tracef(format string, args ...interface{}) ISypl

	// Tracelnf prints according with the specified format @ the Trace level,
	// also adding a new line to the end.
	Tracelnf(format string, args ...interface{}) ISypl

	// Traceln prints, also adding a new line to the end @ the Trace level.
	Traceln(args ...interface{}) ISypl
}

// IPrinters is all available printers.
type IPrinters interface {
	IBasePrinter
	IBasicPrinter
	IConvenientPrinter
	ILeveledPrinter
}

// ISypl specified what a Sypl logger does.
type ISypl interface {
	meta.IMeta
	IPrinters

	// String interface.
	String() string

	// AddOutputs adds one or more outputs.
	AddOutputs(outputs ...output.IOutput) ISypl

	// GetOutput returns the registered output by its name. If not found, will be nil.
	GetOutput(name string) output.IOutput

	// SetOutputs sets one or more outputs. Use to update output(s).
	SetOutputs(outputs ...output.IOutput)

	// GetOutputs returns registered outputs.
	GetOutputs() []output.IOutput

	// GetOutputsNames returns the names of the registered outputs.
	GetOutputsNames() []string

	// New creates a child logger.
	New(name string) ISypl

	// Process messages, per output, and process accordingly.
	process(messages ...message.IMessage)
}
