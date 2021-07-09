// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/saucelabs/lumberjack/v3"
	"github.com/saucelabs/sypl/internal/builtin"
)

const defaultFileMode = 0644

// FileRotationOptions for the log file.
type FileRotationOptions lumberjack.Logger

// output defines and output that could be a console, a file, or anything that
// implements io.Writer.
//
// Notes:
// - Any message with a `level` beyond `maxLevel` will not be written.
// - Messages are processed according to the order `processors` are added.
type Output struct {
	Logger *builtin.Builtin

	enabled    bool
	maxLevel   Level
	name       string
	processors []*Processor
}

// AddProcessor adds a processor.
//
// Note: This method is chainable.
func (o *Output) AddProcessor(processor *Processor) *Output {
	o.processors = append(o.processors, processor)

	return o
}

// GetStatus returns if the output is enabled or disabled.
func (o *Output) GetStatus() bool {
	return o.enabled
}

// SetStatus allows to enable or disable an output.
func (o *Output) SetStatus(status bool) {
	o.enabled = status
}

// NewOutput creates a new `output`.
//
// Notes:
// - The created `output` is enabled by default.
// - processors can be added here, or later using the `AddProcessor` method.
// - This method is chainable.
func NewOutput(name string, maxLevel Level, writer io.Writer, processors ...*Processor) *Output {
	return &Output{
		Logger: builtin.NewBuiltin(writer, "", 0),

		enabled:    true,
		maxLevel:   maxLevel,
		name:       name,
		processors: processors,
	}
}

// OutputsNames extract the names of the given outputs.
func OutputsNames(outputs []*Output) string {
	outputNames := []string{}

	for _, output := range outputs {
		outputNames = append(outputNames, output.name)
	}

	return strings.Join(outputNames, ",")
}

//////
// Built-in outputs.
//////

// Console is a specialized `output` that outputs to the console (stdout).
func Console(maxLevel Level, processors ...*Processor) *Output {
	return NewOutput("Console", maxLevel, os.Stdout, processors...)
}

// FileBased is a specialized `output` that outputs to a file. If the usual, and
// common used "-" is used, it will behave as a Console writing to stdout and
// named "-".
func FileBased(name string, path string, maxLevel Level, writer io.Writer, processors ...*Processor) *Output {
	if path == "-" {
		return NewOutput("-", maxLevel, os.Stdout, processors...)
	}

	return NewOutput(name, maxLevel, writer, processors...)
}

// File is a specialized `output` that outputs to the specified file.
func File(path string, maxLevel Level, processors ...*Processor) *Output {
	f, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		defaultFileMode,
	)
	if err != nil {
		log.Fatalf("Failed to create/open %s: %s", path, err)
	}

	return FileBased("File", path, maxLevel, f, processors...)
}

// FileWithRotation is a specialized `output` that outputs to the specified
// file, with rotation.
func FileWithRotation(
	path string,
	maxLevel Level,
	options *FileRotationOptions,
	processors ...*Processor,
) *Output {
	rotation := &lumberjack.Logger{
		Filename: path,

		Compress:   options.Compress,
		MaxAge:     options.MaxAge,
		MaxBackups: options.MaxBackups,
		MaxBytes:   options.MaxBytes,
	}

	return FileBased("FileWithRotation", path, maxLevel, rotation, processors...)
}
