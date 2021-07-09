// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/afero"
)

const (
	defaultComponentNameOutput = "componentNameTest"
	defaultContentOutput       = "contentTest"
	defaultPrefixValue         = "My Prefix - "
	defaultTimestampFormat     = "2006"
)

func TestNew(t *testing.T) {
	type args struct {
		component string
		content   string
		dir       string
		filename  string
		level     Level
		maxLevel  Level

		run func(a args) string
	}

	noneArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     NONE,
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter, PrefixBasedOnMask(defaultTimestampFormat))).
				Print(a.level, a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	infoArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     INFO,
		maxLevel:  DEBUG,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter, PrefixBasedOnMask(defaultTimestampFormat))).
				Print(a.level, a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	aboveArgs := args{
		level:    TRACE,
		maxLevel: DEBUG,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter, PrefixBasedOnMask(defaultTimestampFormat))).
				Print(a.level, a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	mutedArgs := args{
		level:    INFO, // Will not be used.
		maxLevel: TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter, MuteBasedOnLevel(INFO, WARN))).
				Printf(INFO, "%s", a.content).
				Printf(WARN, "%s", a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	fileArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		dir:       "/tmp",
		filename:  "test.log",
		level:     INFO,
		maxLevel:  DEBUG,
		run: func(a args) string {
			filePath := filepath.Join(a.dir, a.filename)

			appFs := afero.NewMemMapFs()
			f, err := appFs.OpenFile(
				filePath,
				os.O_APPEND|os.O_CREATE|os.O_WRONLY,
				defaultFileMode)
			if err != nil {
				t.Error("Failed to open virtal file", err)
			}

			defer f.Close()

			New(a.component).
				AddOutput(FileBased("virtual", filePath, DEBUG, f, Prefixer("Test Prefix - "))).
				Print(a.level, a.content)

			b, err := afero.ReadFile(appFs, filePath)
			if err != nil {
				t.Error("Failed to read virtal file", err)
			}

			return string(b)
		},
	}

	disableArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     INFO,
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter).
					AddProcessor(EnableDisableProcessors(false, "Prefixer", "Suffixer")).
					AddProcessor(Prefixer("Testing prefix - ")).
					AddProcessor(Suffixer(" - Testing suffix"))).
				Print(a.level, a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	errorArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     NONE, // Will not be used.
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter, PrefixBasedOnMask(defaultTimestampFormat))).
				Errorf("%s", a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	info2Args := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     NONE, // Will not be used.
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter, PrefixBasedOnMask(defaultTimestampFormat))).
				Infof("%s", a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	warnArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     NONE, // Will not be used.
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter, PrefixBasedOnMask(defaultTimestampFormat))).
				Warnf("%s", a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	debugArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     NONE, // Will not be used.
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter, PrefixBasedOnMask(defaultTimestampFormat))).
				Debugf("%s", a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	trace2Args := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     NONE, // Will not be used.
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter, PrefixBasedOnMask(defaultTimestampFormat))).
				Tracef("%s", a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	forceArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     ERROR, // Will not be used.
		maxLevel:  FATAL,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter, ForceBasedOnLevel(ERROR, WARN))).
				Printf(ERROR, "%s", a.content).
				Printf(WARN, "%s", a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	printfArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     INFO,
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter)).
				Printf(a.level, "%s", a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	printfNewLineArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     INFO,
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter)).
				Printf(a.level, "%s\n", a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	printlnArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     INFO,
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter)).
				Println(a.level, a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	prefixBasedOnMaskExceptForLevelsArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     INFO, // Will not be used.
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter,
					PrefixBasedOnMaskExceptForLevels(defaultTimestampFormat, INFO, WARN),
				)).
				Printf(INFO, "%s", a.content).
				Printf(WARN, "%s", a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	prefixBasedOnMaskExceptForLevelsDontArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     INFO,
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput(
					"buffer",
					a.maxLevel,
					bufWriter,
					PrefixBasedOnMaskExceptForLevels(defaultTimestampFormat, WARN)),
				).
				Printf(a.level, "%s", a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	printByOutputArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     INFO,
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer 1", a.maxLevel, bufWriter)).
				AddOutput(NewOutput("buffer 2", a.maxLevel, bufWriter)).
				PrintByOutput([]string{"buffer 1"}, a.level, defaultContentOutput)

			bufWriter.Flush()

			return buf.String()
		},
	}

	printByOutputDontArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     INFO,
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter)).
				PrintByOutput([]string{"invalid"}, a.level, defaultContentOutput)

			bufWriter.Flush()

			return buf.String()
		},
	}

	enableDisableOutputsArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     INFO,
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer 1", a.maxLevel, bufWriter, EnableDisableOutputs(false, "buffer 2"))).
				AddOutput(NewOutput("buffer 2", a.maxLevel, bufWriter)).
				Info(defaultContentOutput)

			bufWriter.Flush()

			return buf.String()
		},
	}

	changeFirstCharCaseUpperArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     INFO,
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer 1", a.maxLevel, bufWriter, ChangeFirstCharCase(Uppercase))).
				Info(defaultContentOutput)

			bufWriter.Flush()

			return buf.String()
		},
	}

	changeFirstCharCaseLowerArgs := args{
		component: defaultComponentNameOutput,
		content:   "ContentTest",
		level:     INFO,
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer 1", a.maxLevel, bufWriter, ChangeFirstCharCase(Lowercase))).
				Info(defaultContentOutput)

			bufWriter.Flush()

			return buf.String()
		},
	}

	nonChainedNewLoggerArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     INFO,
		maxLevel:  TRACE,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			// Creates logger, and name it.
			testingLogger := New("Testing Logger 1")

			// Creates an `Output`. In this case, called Console that will print to
			// stdout and max print level @ INFO.
			ConsoleToStdOut := NewOutput("Console", INFO, bufWriter)

			// Creates a `Processor`. It will prefix all messages.
			Prefixer := func(prefix string) *Processor {
				return NewProcessor("Prefixer", func(message *Message) {
					message.ContentProcessed = prefix + message.ContentProcessed
				})
			}

			// Adds `Processor` to `Output`.
			ConsoleToStdOut.AddProcessor(Prefixer(defaultPrefixValue))

			// Adds `Output` to logger.
			testingLogger.AddOutput(ConsoleToStdOut)

			// Prints: "My Prefix - Test message"
			testingLogger.Print(INFO, "Test message")

			bufWriter.Flush()

			return buf.String()
		},
	}

	tests := []struct {
		name string
		args args
		want func(a args) string
	}{
		{
			name: "Should not print - None",
			args: noneArgs,
			want: func(a args) string {
				return ""
			},
		},
		{
			name: "Should print - Masked Prefix",
			args: infoArgs,
			want: func(a args) string {
				return fmt.Sprintf("%d [%d] [%s] [%s] %s",
					time.Now().Year(),
					os.Getpid(),
					a.component,
					a.level,
					a.content)
			},
		},
		{
			name: "Should not print - Above MaxLevel",
			args: aboveArgs,
			want: func(a args) string {
				return ""
			},
		},
		{
			name: "Should not print - Muted",
			args: mutedArgs,
			want: func(a args) string {
				return ""
			},
		},
		{
			name: "Should print - File based",
			args: fileArgs,
			want: func(a args) string {
				return "Test Prefix - " + defaultContentOutput
			},
		},
		{
			name: "Should print - Only prefix (Disabler)",
			args: disableArgs,
			want: func(a args) string {
				return defaultContentOutput
			},
		},
		{
			name: "Should print - Error level",
			args: errorArgs,
			want: func(a args) string {
				return fmt.Sprintf("%d [%d] [%s] [%s] %s",
					time.Now().Year(),
					os.Getpid(),
					a.component,
					"ERROR",
					a.content)
			},
		},
		{
			name: "Should print - Info level",
			args: info2Args,
			want: func(a args) string {
				return fmt.Sprintf("%d [%d] [%s] [%s] %s",
					time.Now().Year(),
					os.Getpid(),
					a.component,
					"INFO",
					a.content)
			},
		},
		{
			name: "Should print - Warn level",
			args: warnArgs,
			want: func(a args) string {
				return fmt.Sprintf("%d [%d] [%s] [%s] %s",
					time.Now().Year(),
					os.Getpid(),
					a.component,
					"WARN",
					a.content)
			},
		},
		{
			name: "Should print - Debug level",
			args: debugArgs,
			want: func(a args) string {
				return fmt.Sprintf("%d [%d] [%s] [%s] %s",
					time.Now().Year(),
					os.Getpid(),
					a.component,
					"DEBUG",
					a.content)
			},
		},
		{
			name: "Should print - Trace level",
			args: trace2Args,
			want: func(a args) string {
				return fmt.Sprintf("%d [%d] [%s] [%s] %s",
					time.Now().Year(),
					os.Getpid(),
					a.component,
					"TRACE",
					a.content)
			},
		},
		{
			name: "Should print - Force",
			args: forceArgs,
			want: func(a args) string {
				return defaultContentOutput + defaultContentOutput
			},
		},
		{
			name: "Should print - Printf - No newline",
			args: printfArgs,
			want: func(a args) string {
				return defaultContentOutput
			},
		},
		{
			name: "Should print - Printf - Newline",
			args: printfNewLineArgs,
			want: func(a args) string {
				return defaultContentOutput + "\n"
			},
		},
		{
			name: "Should print - Println",
			args: printlnArgs,
			want: func(a args) string {
				return defaultContentOutput + "\n"
			},
		},
		{
			name: "Should print not prefixed - PrefixBasedOnMaskExceptForLevels",
			args: prefixBasedOnMaskExceptForLevelsArgs,
			want: func(a args) string {
				return defaultContentOutput + defaultContentOutput
			},
		},
		{
			name: "Should print prefixed - PrefixBasedOnMaskExceptForLevels",
			args: prefixBasedOnMaskExceptForLevelsDontArgs,
			want: func(a args) string {
				return fmt.Sprintf("%d [%d] [%s] [%s] %s",
					time.Now().Year(),
					os.Getpid(),
					a.component,
					INFO,
					a.content)
			},
		},
		{
			name: "Should print - printByOutput",
			args: printByOutputArgs,
			want: func(a args) string {
				return defaultContentOutput
			},
		},
		{
			name: "Should not print - printByOutput - name doesn't match",
			args: printByOutputDontArgs,
			want: func(a args) string {
				return ""
			},
		},
		{
			name: "Should print - enableDisableOutputsArgs",
			args: enableDisableOutputsArgs,
			want: func(a args) string {
				return defaultContentOutput
			},
		},
		{
			name: "Should print - changeFirstCharCaseUpperArgs",
			args: changeFirstCharCaseUpperArgs,
			want: func(a args) string {
				return "ContentTest"
			},
		},
		{
			name: "Should print - changeFirstCharCaseLowerArgs",
			args: changeFirstCharCaseLowerArgs,
			want: func(a args) string {
				return defaultContentOutput
			},
		},
		{
			name: "Should print - nonChainedNewLoggerArgs",
			args: nonChainedNewLoggerArgs,
			want: func(a args) string {
				return "My Prefix - Test message"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := tt.args.run(tt.args)
			want := tt.want(tt.args)

			if message != want {
				t.Errorf("Got %v, want %v", message, want)
			}
		})
	}
}

// NonChained is a non-chained example of creating, and setting up a `sypl`
// logger. It writes to a custom buffer.
func ExampleNew_notChained() {
	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)

	// Creates logger, and name it.
	testingLogger := New("Testing Logger")

	// Creates an `Output`. In this case, called "Console" that will print to
	// `stdout` and max print level @ `INFO`.
	ConsoleToStdOut := NewOutput("Console", INFO, bufWriter)

	// Creates a `Processor`. It will prefix all messages.
	Prefixer := func(prefix string) *Processor {
		return NewProcessor("Prefixer", func(message *Message) {
			message.ContentProcessed = prefix + message.ContentProcessed
		})
	}

	// Adds `Processor` to `Output`.
	ConsoleToStdOut.AddProcessor(Prefixer(defaultPrefixValue))

	// Adds `Output` to logger.
	testingLogger.AddOutput(ConsoleToStdOut)

	// Prints: "My Prefix - Test message"
	testingLogger.Print(INFO, "Test info message")

	bufWriter.Flush()

	fmt.Println(buf.String())

	// Output:
	// My Prefix - Test info message
}

// Chained is the chained example of creating, and setting up a `sypl` logger.
// It writes to both `stdout` and `stderr`.
func ExampleNew_chained() {
	// Creates logger, and name it.
	New("Testing Logger").
		// Creates two `Output`s. "Console" and "Error". "Console" will print to
		// `Fatal`, `Error`, and `Info`. "Error" will only print `Fatal`, and
		// `Error` levels.
		AddOutput(NewOutput("Console", INFO, os.Stdout)).
		// Creates a `Processor`. It will prefix all messages. It will only
		// prefix messages for this specific `Output`, and @ `ERROR` level.
		AddOutput(NewOutput("Error", ERROR, os.Stderr).
			AddProcessor(func(prefix string) *Processor {
				return NewProcessor("Prefixer", func(message *Message) {
					if message.Level == ERROR {
						message.ContentProcessed = prefix + message.ContentProcessed
					}
				})
			}(defaultPrefixValue))).
		// Prints:
		// Test info message
		Println(INFO, "Test info message").
		// Prints:
		// Test error message
		// My Prefix - Test error message
		Println(ERROR, "Test error message")

	// Output:
	// My Prefix - Test error message

	// Note: Go "example" parser only captured the last message.
}

// ChainedUsingBuiltin is the chained example of creating, and setting up a
// `sypl` logger using built-in `Output`, and `Processor`. It writes to
// `stdout`, and `stderr`.
func ExampleNew_chainedUsingBuiltin() {
	// Creates logger, and name it.
	New("Testing Logger").
		// Adds an `Output`. In this case, called "Console" that will print to
		// `stdout` and max print level @ `INFO`.
		//
		// Adds a `Processor`. It will prefix all messages.
		AddOutput(Console(INFO).AddProcessor(Prefixer(defaultPrefixValue))).
		// Prints: My Prefix - Test info message
		Infoln("Test info message")

	// output:
	// My Prefix - Test info message
}

// inlineUsingBuiltin same as `ChainedUsingBuiltin` but using inline form.
func ExampleNew_inlineUsingBuiltin() {
	New("Testing Logger", Console(INFO).AddProcessor(Prefixer(defaultPrefixValue))).Infoln("Test info message")

	// output:
	// My Prefix - Test info message
}
