// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const defaultCallDepth = 2

// Printer defines possible printers.
type Printer interface {
	// Print prints the message.
	Print(level Level, args ...interface{}) *Sypl

	// Printf prints the message according with the specified format.
	Printf(level Level, format string, args ...interface{}) *Sypl

	// Print prints the message, also adding a new line and the end.
	Println(level Level, args ...interface{}) *Sypl

	// Fatal prints the message, and exit with os.Exit(1).
	Fatal(level Level, args ...interface{})

	// Fatalf prints the message according with the specified format, and exit
	// with os.Exit(1).
	Fatalf(level Level, format string, args ...interface{})

	// Fatalln prints the message, also adding a new line and the end, and exit
	// with os.Exit(1).
	Fatalln(level Level, args ...interface{})
}

// sypl definition. It's able to print messages according to the definition of
// each `output`.
type Sypl struct {
	name    string
	outputs []*Output
}

// AddOutput adds an `output` to logger.
//
// Note: This method is chainable.
func (cL *Sypl) AddOutput(output *Output) *Sypl {
	cL.outputs = append(cL.outputs, output)

	return cL
}

// Writes to the specified writer.
//
//
// Note: In case of any error, the standard log will be used to highlight the
// case, but IT WILL NOT STOP the application.
func (cL *Sypl) write(message *Message) {
	if err := message.output.Logger.OutputBuiltin(
		defaultCallDepth,
		message.ContentProcessed,
	); err != nil {
		log.Println("ERROR: Failed to log to output.", err)
	}
}

// Process a message according to logger's registered outputs, and its
// processors. If a list of outputs names is passed, it will only process and
// print if matches against the registered outputs, otherwise all registered
// outputs will be used. If not content if found, nothing is processed or
// printed.
func (cL *Sypl) process(level Level, content string, outputsNames ...string) *Sypl {
	if content == "" {
		return cL
	}

	for _, output := range cL.outputs {
		concatenatedNames := OutputsNames(cL.outputs)

		if len(outputsNames) > 0 {
			concatenatedNames = strings.Join(outputsNames, ",")
		}

		if strings.Contains(concatenatedNames, output.name) && output.enabled {
			// Content envelop.
			message := &Message{
				sypl: cL,

				ContentOriginal:  content,
				ContentProcessed: content,
				Level:            level,
				output:           output,
			}

			for _, processor := range output.processors {
				// Sets the `Processor` in use.
				message.processor = processor

				processor.Run(message)
			}

			if message.force {
				cL.write(message)
			}

			// Should only print if message `level` isn't above `MaxLevel`, it
			// has content and isn't muted.
			if message.Level != NONE &&
				message.Level <= output.maxLevel &&
				!message.mute {
				cL.write(message)
			}
		}
	}

	// Handles FATAL level. Here to ensure that message was written to all
	// outputs before exit.
	if content != "" && level == FATAL {
		os.Exit(1)
	}

	return cL
}

// Print prints the message (if has content).
func (cL *Sypl) Print(level Level, args ...interface{}) *Sypl {
	return cL.process(level, fmt.Sprint(args...))
}

// Printf prints the message (if has content) according with the specified
// format.
func (cL *Sypl) Printf(level Level, format string, args ...interface{}) *Sypl {
	return cL.process(level, fmt.Sprintf(format, args...))
}

// Println prints the message (if has content), also adding a new line to the
// end.
func (cL *Sypl) Println(level Level, args ...interface{}) *Sypl {
	return cL.process(level, fmt.Sprintln(args...))
}

// Fatal prints the message (if has content), and exit with os.Exit(1).
func (cL *Sypl) Fatal(args ...interface{}) {
	cL.process(FATAL, fmt.Sprint(args...))
}

// Fatalf prints the message (if has content) according with the specified
// format, and exit with os.Exit(1).
func (cL *Sypl) Fatalf(format string, args ...interface{}) {
	cL.process(FATAL, fmt.Sprintf(format, args...))
}

// Fatalln prints the message (if has content), also adding a new line and the
// end, and exit with os.Exit(1).
func (cL *Sypl) Fatalln(args ...interface{}) {
	cL.process(FATAL, fmt.Sprintln(args...))
}

//////
// Convenience methods.
//////

// Error prints the message (if has content) @ the ERROR level.
func (cL *Sypl) Error(args ...interface{}) *Sypl {
	return cL.process(ERROR, fmt.Sprint(args...))
}

// Errorf prints the message (if has content) according with the specified
// format @ the ERROR level.
func (cL *Sypl) Errorf(format string, args ...interface{}) *Sypl {
	return cL.process(ERROR, fmt.Sprintf(format, args...))
}

// Errorln prints the message (if has content), also adding a new line to the
// end @ the ERROR level.
func (cL *Sypl) Errorln(args ...interface{}) *Sypl {
	return cL.process(ERROR, fmt.Sprintln(args...))
}

// Info prints the message (if has content) @ the INFO level.
func (cL *Sypl) Info(args ...interface{}) *Sypl {
	return cL.process(INFO, fmt.Sprint(args...))
}

// Infof prints the message (if has content) according with the specified format
// @ the INFO level.
func (cL *Sypl) Infof(format string, args ...interface{}) *Sypl {
	return cL.process(INFO, fmt.Sprintf(format, args...))
}

// Infoln prints the message (if has content), also adding a new line to the end
// @ the INFO level.
func (cL *Sypl) Infoln(args ...interface{}) *Sypl {
	return cL.process(INFO, fmt.Sprintln(args...))
}

// Warn prints the message (if has content) @ the WARN level.
func (cL *Sypl) Warn(args ...interface{}) *Sypl {
	return cL.process(WARN, fmt.Sprint(args...))
}

// Warnf prints the message (if has content) according with the specified format
// @ the WARN level.
func (cL *Sypl) Warnf(format string, args ...interface{}) *Sypl {
	return cL.process(WARN, fmt.Sprintf(format, args...))
}

// Warnln prints the message (if has content), also adding a new line to the end
// @ the WARN level.
func (cL *Sypl) Warnln(args ...interface{}) *Sypl {
	return cL.process(WARN, fmt.Sprintln(args...))
}

// Debug prints the message (if has content) @ the DEBUG level.
func (cL *Sypl) Debug(args ...interface{}) *Sypl {
	return cL.process(DEBUG, fmt.Sprint(args...))
}

// Debugf prints the message (if has content) according with the specified
// format @ the DEBUG level.
func (cL *Sypl) Debugf(format string, args ...interface{}) *Sypl {
	return cL.process(DEBUG, fmt.Sprintf(format, args...))
}

// Debugln prints the message (if has content), also adding a new line to the
// end @ the DEBUG level.
func (cL *Sypl) Debugln(args ...interface{}) *Sypl {
	return cL.process(DEBUG, fmt.Sprintln(args...))
}

// Trace prints the message (if has content) @ the TRACE level.
func (cL *Sypl) Trace(args ...interface{}) *Sypl {
	return cL.process(TRACE, fmt.Sprint(args...))
}

// Tracef prints the message (if has content) according with the specified
// format @ the TRACE level.
func (cL *Sypl) Tracef(format string, args ...interface{}) *Sypl {
	return cL.process(TRACE, fmt.Sprintf(format, args...))
}

// Traceln prints the message (if has content), also adding a new line to the
// end @ the TRACE level.
func (cL *Sypl) Traceln(args ...interface{}) *Sypl {
	return cL.process(TRACE, fmt.Sprintln(args...))
}

// PrintByOutput prints the message (if has content) to the specified outputs.
func (cL *Sypl) PrintByOutput(outputsNames []string, level Level, args ...interface{}) *Sypl {
	return cL.process(level, fmt.Sprint(args...), outputsNames...)
}

// PrintfByOutput prints the message (if has content) according with the
// specified format, and to
// the specified outputs.
func (cL *Sypl) PrintfByOutput(outputsNames []string, level Level, format string, args ...interface{}) *Sypl {
	return cL.process(level, fmt.Sprintf(format, args...), outputsNames...)
}

// PrintlnByOutput prints the message (if has content), also adding a new line
// to the end, and to
// the specified outputs.
func (cL *Sypl) PrintlnByOutput(outputsNames []string, level Level, args ...interface{}) *Sypl {
	return cL.process(level, fmt.Sprintln(args...), outputsNames...)
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
