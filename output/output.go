// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/saucelabs/sypl/flag"
	"github.com/saucelabs/sypl/formatter"
	"github.com/saucelabs/sypl/internal/builtin"
	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/message"
	"github.com/saucelabs/sypl/processor"
	"github.com/saucelabs/sypl/shared"
	"github.com/saucelabs/sypl/status"
)

// Output process, and write the message to the defined writer. A writer is
// anything that implements io.Writer.
//
// Notes:
// - Any message with a `level` beyond `maxLevel` will not be written.
// - Messages are processed according to the order processors are added.
type output struct {
	// Golang's builtin logger.
	builtinLogger *builtin.Builtin

	// Formats the message.
	formatter formatter.IFormatter

	// Any message above the max level will not be written.
	maxLevel level.Level

	// Name of the processor.
	name string

	// Processors used to process the message.
	processors []processor.IProcessor

	// Status of the processor.
	status status.Status

	// Writer to write.
	writer io.Writer
}

// String interface implementation.
func (o output) String() string {
	return o.name
}

//////
// IMeta interface implementation.
//////

// GetName returns the processor name.
func (o *output) GetName() string {
	return o.name
}

// SetName sets the processor name.
func (o *output) SetName(name string) {
	o.name = name
}

// GetStatus returns the processor status.
func (o *output) GetStatus() status.Status {
	return o.status
}

// SetStatus sets the processor status.
func (o *output) SetStatus(s status.Status) {
	o.status = s
}

//////
// IOutput interface implementation.
//////

// GetBuiltinLogger returns the Golang's builtin logger.
func (o *output) GetBuiltinLogger() *builtin.Builtin {
	return o.builtinLogger
}

// SetBuiltinLogger sets the Golang's builtin logger.
func (o *output) SetBuiltinLogger(builtinLogger *builtin.Builtin) IOutput {
	o.builtinLogger = builtinLogger

	return o
}

// GetFormatter returns the formatter.
func (o *output) GetFormatter() formatter.IFormatter {
	return o.formatter
}

// SetFormatter sets the formatter.
func (o *output) SetFormatter(fmtr formatter.IFormatter) IOutput {
	o.formatter = fmtr

	return o
}

// GetMaxLevel returns the max level.
func (o *output) GetMaxLevel() level.Level {
	return o.maxLevel
}

// SetMaxLevel sets the max level.
func (o *output) SetMaxLevel(l level.Level) IOutput {
	o.maxLevel = l

	return o
}

// AddProcessors adds one or more processors.
func (o *output) AddProcessors(processors ...processor.IProcessor) IOutput {
	o.processors = append(o.processors, processors...)

	return o
}

// GetProcessor returns the registered processor by its name. If not found, will
// be nil.
func (o *output) GetProcessor(name string) processor.IProcessor {
	for _, p := range o.processors {
		if strings.EqualFold(p.GetName(), name) {
			return p
		}
	}

	return nil
}

// SetProcessors sets one or more processors.
func (o *output) GetProcessors() []processor.IProcessor {
	return o.processors
}

// GetProcessors returns registered processors.
func (o *output) SetProcessors(processors ...processor.IProcessor) IOutput {
	for _, processor := range processors {
		for i, p := range o.processors {
			if strings.EqualFold(p.GetName(), processor.GetName()) {
				o.processors[i] = processor
			}
		}
	}

	return o
}

// GetProcessorsNames returns the names of the registered processors.
func (o *output) GetProcessorsNames() []string {
	processorsNames := []string{}

	for _, processor := range o.processors {
		processorsNames = append(processorsNames, processor.GetName())
	}

	return processorsNames
}

// GetWriter returns the writer.
func (o *output) GetWriter() io.Writer {
	return o.writer
}

// SetWriter sets the writer.
func (o *output) SetWriter(w io.Writer) IOutput {
	o.writer = w

	return o
}

// Write the message to the defined output. In case of any error, it can be
// introspected, providing more information about the failure. The error will be
// the type of `ProcessingError`.
//nolint:nestif
func (o *output) Write(m message.IMessage) error {
	// Should allows to specify `Output`(s).
	processorsNames := o.GetProcessorsNames()

	if len(m.GetProcessorsNames()) > 0 {
		processorsNames = m.GetProcessorsNames()
	}

	m.SetProcessorsNames(processorsNames)

	// Strips the last line break, which allows the content to be
	// properly processed. It gets restore later, if any.
	m.Strip()

	// Executes processors in series.
	o.processProcessors(m, strings.Join(processorsNames, ","))

	// Should print the message - regardless of the level, if flagged
	// with `Force`.
	if m.GetFlag() == flag.Force || m.GetFlag() == flag.SkipAndForce {
		if err := o.write(m); err != nil {
			log.Println(shared.ErrorPrefix, err)

			return err
		}
	} else {
		// Debug capability.
		finalMaxLevel := o.GetMaxLevel()

		// Should only run if Debug env var is set.
		if os.Getenv(shared.DebugEnvVar) != "" {
			debug := m.GetDebugEnvVarRegexes()

			l, _, ok := debug.Level()

			if ok {
				finalMaxLevel = l
			}
		}

		// Should only print if message `level` isn't above `MaxLevel`.
		// Should only print if `level` isn't `None`.
		// Should only print if not flagged with `Mute`.
		if m.GetLevel() != level.None &&
			m.GetLevel() <= finalMaxLevel &&
			m.GetFlag() != flag.Mute {
			if err := o.write(m); err != nil {
				log.Println(shared.ErrorPrefix, err)

				return err
			}
		}
	}

	return nil
}

//////
// Helpers.
//////

// Processors logic of the Write method.
func (o *output) processProcessors(m message.IMessage, processorsNames string) {
	// Should not process if message is flagged with `Skip` or `SkipAndForce`.
	if m.GetFlag() != flag.Skip && m.GetFlag() != flag.SkipAndForce {
		for _, p := range o.processors {
			// Should only use enabled Processors, and named (listed) ones.
			//
			// Note: `Enabled` status is checked in the `Run` method.
			if strings.Contains(processorsNames, p.GetName()) {
				m.SetProcessorName(p.GetName())

				if err := p.Run(m); err != nil {
					log.Println(shared.ErrorPrefix,
						processor.NewProcessingError(m, err))
				}
			}
		}
	}
}

// DRY for the writing step.
func (o *output) write(m message.IMessage) error {
	// Should only format if any, and if not flagged.
	if o.GetFormatter() != nil &&
		m.GetFlag() != flag.Skip &&
		m.GetFlag() != flag.SkipAndForce {
		if err := o.GetFormatter().Run(m); err != nil {
			log.Println(shared.ErrorPrefix, processor.NewProcessingError(m, err))
		}
	}

	// Restore linebreak(s), if needed.
	m.Restore()

	// Write to writer.
	if err := o.GetBuiltinLogger().OutputBuiltin(
		builtin.DefaultCallDepth,
		m.GetContent().GetProcessed(),
	); err != nil {
		// It means application using Sypl was piped, but the pipe was broken so
		// nothing to do.
		if errors.Is(err, syscall.EPIPE) {
			return nil
		}

		// It the output passed to Sypl is already closed, nothing to do.
		if errors.Is(err, os.ErrClosed) {
			log.Printf(`%s Attempt to write to closed writer. Output: "%s". Error: "%v"`,
				shared.WarnPrefix,
				o.GetName(),
				err,
			)

			return nil
		}

		return fmt.Errorf(`output: "%s". error: "%w"`, o.GetName(), err)
	}

	return nil
}

//////
// Factory.
//////

// New is the Output factory.
func New(name string,
	maxLevel level.Level,
	w io.Writer,
	processors ...processor.IProcessor,
) IOutput {
	return &output{
		builtinLogger: builtin.NewBuiltin(w, "", 0),
		maxLevel:      maxLevel,

		name:       name,
		processors: processors,
		status:     status.Enabled,
		writer:     w,
	}
}
