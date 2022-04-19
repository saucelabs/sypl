// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.


package sypl

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/saucelabs/sypl/level"
)

const (
	_programmerErrorTemplate = "You've found a bug in sypl! Please file a bug at " +
		"https://github.com/saucelabs/sypl/issues/new and reference this error: %v"
)


// RedirectStdLog redirects output from the standard library's package-global
// logger to the supplied logger at level.Info. Since zap already handles caller
// annotations, timestamps, etc., it automatically disables the standard
// library's annotations and prefixing.
//
// It returns a function to restore the original prefix and flags and reset the
// standard library's output to os.Stderr.
func RedirectStdLog(sypl *Sypl) func() {
	f, err := redirectStdLogAt(sypl, level.Info)
	if err != nil {
		// Passing InfoLevel to redirectStdLogAt should always work
		panic(fmt.Sprintf(_programmerErrorTemplate, err))
	}
	return f
}

// RedirectStdLogAt redirects output from the standard library's package-global
// logger to the supplied logger at the specified level. Since sypl already
// handles caller annotations, timestamps, etc., it automatically disables the
// standard library's annotations and prefixing.
//
// It returns a function to restore the original prefix and flags and reset the
// standard library's output to os.Stderr.
func RedirectStdLogAt(sypl *Sypl, lvl level.Level) (func(), error) {
	return redirectStdLogAt(sypl, lvl)
}

func redirectStdLogAt(sypl *Sypl, lvl level.Level) (func(), error) {
	flags := log.Flags()
	prefix := log.Prefix()
	log.SetFlags(0)
	log.SetPrefix("")
	logFunc, err := levelToFunc(sypl, lvl)
	if err != nil {
		return nil, err
	}
	log.SetOutput(&loggerWriter{logFunc})
	return func() {
		log.SetFlags(flags)
		log.SetPrefix(prefix)
		log.SetOutput(os.Stderr)
	}, nil
}

func levelToFunc(sypl *Sypl, lvl level.Level) (func(...interface{}) ISypl, error) {
	switch lvl {
	case level.Trace:
		return sypl.Traceln, nil
	case level.Debug:
		return sypl.Debugln, nil
	case level.Info:
		return sypl.Infoln, nil
	case level.Warn:
		return sypl.Warnln, nil
	case level.Error:
		return sypl.Errorln, nil
	case level.Fatal:
		return sypl.Fatalln, nil
	}
	return nil, fmt.Errorf("Unrecognized log level: %q", lvl)
}

type loggerWriter struct {
	logFunc func(...interface{}) ISypl
}

func (l *loggerWriter) Write(p []byte) (int, error) {
	p = bytes.TrimSpace(p)
	l.logFunc(string(p))
	return len(p), nil
}
