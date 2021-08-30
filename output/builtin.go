// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"io"
	"log"
	"os"

	"github.com/saucelabs/lumberjack/v3"
	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/processor"
	"github.com/saucelabs/sypl/safebuffer"
	"github.com/saucelabs/sypl/shared"
)

// FileRotationOptions for the log file.
type FileRotationOptions lumberjack.Logger

// Handles the common used "-" making the output behave as a Console writing to
// stdout, and named "-".
func dashHandler(name, path string, maxLevel level.Level, processors ...processor.IProcessor) IOutput {
	if path == "-" {
		return New(name, maxLevel, os.Stdout, processors...)
	}

	return nil
}

//////
// Built-in outputs.
//////

// Console is a built-in `output` - named `Console`, that writes to `stdout`.
func Console(maxLevel level.Level, processors ...processor.IProcessor) IOutput {
	return New("Console", maxLevel, os.Stdout, processors...)
}

// StdErr is a built-in `output` - named `StdErr`, that only writes to `stderr`
// message @ Error level.
func StdErr(processors ...processor.IProcessor) IOutput {
	processors = append(processors, processor.PrintOnlyAtLevel(level.Fatal, level.Error))

	return New("StdErr", level.Error, os.Stderr, processors...)
}

// FileBased is a built-in `output`, that writes to the specified file.
// TODO: Remove `path`.
func FileBased(
	name string,
	maxLevel level.Level,
	writer io.Writer,
	processors ...processor.IProcessor,
) IOutput {
	return New(name, maxLevel, writer, processors...)
}

// File is a built-in `output` - named `File`, that writes to the specified file.
//
// Note: If the common used "-" is used, it will behave as a Console writing to
// stdout.
func File(path string, maxLevel level.Level, processors ...processor.IProcessor) IOutput {
	if o := dashHandler("File", path, maxLevel, processors...); o != nil {
		return o
	}

	f, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		shared.DefaultFileMode,
	)
	if err != nil {
		log.Fatalf("%s File Output: Failed to create/open %s: %s", shared.ErrorPrefix, path, err)
	}

	return FileBased("File", maxLevel, f, processors...)
}

// FileWithRotation is a built-in `output` - named `FileWithRotation`, that
// writes to the specified file, an is automatically rotated.
//
// Note: If the common used "-" is used, it will behave as a Console writing to
// stdout.
func FileWithRotation(
	path string,
	maxLevel level.Level,
	options *FileRotationOptions,
	processors ...processor.IProcessor,
) IOutput {
	if o := dashHandler("FileWithRotation", path, maxLevel, processors...); o != nil {
		return o
	}

	rotation := &lumberjack.Logger{
		Filename: path,

		Compress:   options.Compress,
		MaxAge:     options.MaxAge,
		MaxBackups: options.MaxBackups,
		MaxBytes:   options.MaxBytes,
	}

	return FileBased("FileWithRotation", maxLevel, rotation, processors...)
}

// SafeBuffer is a built-in `output` - named `Buffer`, that writes to the buffer.
func SafeBuffer(maxLevel level.Level, processors ...processor.IProcessor) (*safebuffer.Buffer, IOutput) {
	var buf safebuffer.Buffer

	o := New("Buffer", maxLevel, &buf, processors...)

	return &buf, o
}
