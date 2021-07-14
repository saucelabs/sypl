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
	"strings"
	"testing"
	"time"

	"github.com/saucelabs/sypl/level"
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
		level     level.Level
		maxLevel  level.Level

		run func(a args) string
	}

	noneArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     level.None,
		maxLevel:  level.Trace,
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
		level:     level.Info,
		maxLevel:  level.Debug,
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
		level:    level.Trace,
		maxLevel: level.Debug,
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
		level:    level.Info, // Will not be used.
		maxLevel: level.Trace,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter, MuteBasedOnLevel(level.Info, level.Warn))).
				Printf(level.Info, "%s", a.content).
				Printf(level.Warn, "%s", a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	fileArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		dir:       "/tmp",
		filename:  "test.log",
		level:     level.Info,
		maxLevel:  level.Debug,
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
				AddOutput(FileBased("virtual", filePath, level.Debug, f, Prefixer("Test Prefix - "))).
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
		level:     level.Info,
		maxLevel:  level.Trace,
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
		level:     level.None, // Will not be used.
		maxLevel:  level.Trace,
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
		level:     level.None, // Will not be used.
		maxLevel:  level.Trace,
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
		level:     level.None, // Will not be used.
		maxLevel:  level.Trace,
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
		level:     level.None, // Will not be used.
		maxLevel:  level.Trace,
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
		level:     level.None, // Will not be used.
		maxLevel:  level.Trace,
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
		level:     level.Error, // Will not be used.
		maxLevel:  level.Fatal,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter, ForceBasedOnLevel(level.Error, level.Warn))).
				Printf(level.Error, "%s", a.content).
				Printf(level.Warn, "%s", a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	printfArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Trace,
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
		level:     level.Info,
		maxLevel:  level.Trace,
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
		level:     level.Info,
		maxLevel:  level.Trace,
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
		level:     level.Info, // Will not be used.
		maxLevel:  level.Trace,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter,
					PrefixBasedOnMaskExceptForLevels(defaultTimestampFormat, level.Info, level.Warn),
				)).
				Printf(level.Info, "%s", a.content).
				Printf(level.Warn, "%s", a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	prefixBasedOnMaskExceptForLevelsDontArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Trace,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput(
					"buffer",
					a.maxLevel,
					bufWriter,
					PrefixBasedOnMaskExceptForLevels(defaultTimestampFormat, level.Warn)),
				).
				Printf(a.level, "%s", a.content)

			bufWriter.Flush()

			return buf.String()
		},
	}

	printWithOptionsArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Trace,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer 1", a.maxLevel, bufWriter)).
				AddOutput(NewOutput("buffer 2", a.maxLevel, bufWriter)).
				PrintWithOptions(&Options{
					OutputsNames: []string{"buffer 1"},
				}, a.level, defaultContentOutput)

			bufWriter.Flush()

			return buf.String()
		},
	}

	printWithOptionsDontArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Trace,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutput(NewOutput("buffer", a.maxLevel, bufWriter)).
				PrintWithOptions(&Options{
					OutputsNames: []string{"invalid"},
				}, a.level, defaultContentOutput)

			bufWriter.Flush()

			return buf.String()
		},
	}

	enableDisableOutputsArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Trace,
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
		level:     level.Info,
		maxLevel:  level.Trace,
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
		level:     level.Info,
		maxLevel:  level.Trace,
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
		level:     level.Info,
		maxLevel:  level.Trace,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			// Creates logger, and name it.
			testingLogger := New("Testing Logger 1")

			// Creates an `Output`. In this case, called Console that will print to
			// stdout and max print level @ Info.
			ConsoleToStdOut := NewOutput("Console", level.Info, bufWriter)

			// Creates a `Processor`. It will prefix all messages.
			Prefixer := func(prefix string) *Processor {
				return NewProcessor("Prefixer", func(message *Message) {
					message.SetProcessedContent(prefix + message.GetProcessedContent())
				})
			}

			// Adds `Processor` to `Output`.
			ConsoleToStdOut.AddProcessor(Prefixer(defaultPrefixValue))

			// Adds `Output` to logger.
			testingLogger.AddOutput(ConsoleToStdOut)

			// Prints: "My Prefix - Test message"
			testingLogger.Print(level.Info, "Test message")

			bufWriter.Flush()

			return buf.String()
		},
	}

	printflnArgs := args{
		component: defaultComponentNameOutput,
		content:   defaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Trace,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			// Creates logger, and name it.
			testingLogger := New("Testing Logger 1")

			// Creates an `Output`. In this case, called Buffer that will write
			// to the specified buffer, and max print level @ Info.
			BufferOutput := NewOutput("Buffer", level.Info, bufWriter)

			// Adds `Output` to logger.
			testingLogger.AddOutput(BufferOutput)

			testingLogger.
				Printlnf(level.Info, "%s %s", "element 1", "element 2")

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
					"Error",
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
					"Info",
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
					"Warn",
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
					"Debug",
					a.content)
			},
		},
		{
			name: "Should print - level.Trace level",
			args: trace2Args,
			want: func(a args) string {
				return fmt.Sprintf("%d [%d] [%s] [%s] %s",
					time.Now().Year(),
					os.Getpid(),
					a.component,
					"Trace",
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
					level.Info,
					a.content)
			},
		},
		{
			name: "Should print - printWithOptions",
			args: printWithOptionsArgs,
			want: func(a args) string {
				return defaultContentOutput
			},
		},
		{
			name: "Should not print - printWithOptions - name doesn't match",
			args: printWithOptionsDontArgs,
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
		{
			name: "Should print - printflnArgs",
			args: printflnArgs,
			want: func(a args) string {
				return "element 1 element 2\n"
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
	buf := new(strings.Builder)

	// Creates logger, and name it.
	testingLogger := New("Testing Logger")

	// Creates an `Output`. In this case, called "Console" that will print to
	// `stdout` and max print level @ `Info`.
	ConsoleToStdOut := NewOutput("Console", level.Info, buf)

	// Creates a `Processor`. It will prefix all messages.
	Prefixer := func(prefix string) *Processor {
		return NewProcessor("Prefixer", func(message *Message) {
			message.SetProcessedContent(prefix + message.GetProcessedContent())
		})
	}

	// Adds `Processor` to `Output`.
	ConsoleToStdOut.AddProcessor(Prefixer(defaultPrefixValue))

	// Adds `Output` to logger.
	testingLogger.AddOutput(ConsoleToStdOut)

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
	New("Testing Logger").
		// Creates two `Output`s. "Console" and "Error". "Console" will print to
		// `Fatal`, `Error`, and `Info`. "Error" will only print `Fatal`, and
		// `Error` levels.
		AddOutput(NewOutput("Console", level.Info, os.Stdout)).
		// Creates a `Processor`. It will prefix all messages. It will only
		// prefix messages for this specific `Output`, and @ `Error` level.
		AddOutput(NewOutput("Error", level.Error, os.Stderr).
			AddProcessor(func(prefix string) *Processor {
				return NewProcessor("Prefixer", func(message *Message) {
					if message.GetLevel() == level.Error {
						message.SetProcessedContent(prefix + message.GetProcessedContent())
					}
				})
			}(defaultPrefixValue))).
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
	New("Testing Logger").
		// Adds an `Output`. In this case, called "Console" that will print to
		// `stdout` and max print level @ `Info`.
		//
		// Adds a `Processor`. It will prefix all messages.
		AddOutput(Console(level.Info).AddProcessor(Prefixer(defaultPrefixValue))).
		// Prints: My Prefix - Test info message
		Infoln("Test info message")

	// output:
	// My Prefix - Test info message
}

// inlineUsingBuiltin same as `ChainedUsingBuiltin` but using inline form.
func ExampleNew_inlineUsingBuiltin() {
	New("Testing Logger", Console(level.Info).AddProcessor(Prefixer(defaultPrefixValue))).Infoln("Test info message")

	// output:
	// My Prefix - Test info message
}

// printWithOptions demonstrates `sypl` flexibility. `Options` enhances the
// usual `PrintX` methods allowing to specify flags, and tags.
func ExampleNew_printWithOptions() {
	// Creates logger, and name it.
	testingLogger := New("Testing Logger")

	// Creates 3 `Output`s, all called "Console" that will print to `stdout`, and
	// max print level @ `Info`.
	Console1ToStdOut := NewOutput("Console 1", level.Info, os.Stdout)
	Console2ToStdOut := NewOutput("Console 2", level.Info, os.Stdout)
	Console3ToStdOut := NewOutput("Console 3", level.Info, os.Stdout)

	// Creates a `Processor`. It will `prefix` all messages with the Output, and
	// Processor names.
	Prefixer := func() *Processor {
		return NewProcessor("Prefixer", func(message *Message) {
			prefix := fmt.Sprintf("Output: %s Processor: %s Content: ",
				message.GetOutput().GetName(),
				message.GetProcessor().GetName(),
			)

			message.SetProcessedContent(prefix + message.GetProcessedContent())
		})
	}

	// Creates a `Processor`. It will `suffix` all messages with the specified
	// `tag`.
	SuffixBasedOnTag := func(tag string) *Processor {
		return NewProcessor("SuffixBasedOnTag", func(message *Message) {
			if message.ContainTag(tag) {
				message.SetProcessedContent(message.GetProcessedContent() + " - My Suffix")
			}
		})
	}

	// Adds `Processor`s to `Output`s.
	Console1ToStdOut.
		AddProcessor(Prefixer()).
		AddProcessor(SuffixBasedOnTag("SuffixIt"))
	Console2ToStdOut.
		AddProcessor(Prefixer()).
		AddProcessor(SuffixBasedOnTag("SuffixIt"))
	Console3ToStdOut.
		AddProcessor(Prefixer()).
		AddProcessor(SuffixBasedOnTag("SuffixIt"))

	// Adds all `Output`s to logger.
	testingLogger.
		AddOutput(Console1ToStdOut).
		AddOutput(Console2ToStdOut).
		AddOutput(Console3ToStdOut)

	// Prints with prefix, without suffix.
	testingLogger.Println(level.Info, defaultContentOutput)

	// Prints with prefix, and suffix.
	testingLogger.PrintWithOptions(&Options{
		OutputsNames:    []string{"Console 1"},
		ProcessorsNames: []string{"Prefixer", "SuffixBasedOnTag"},
		Tags:            []string{"SuffixIt"},
	}, level.Info, defaultContentOutput)

	// output:
	// Output: Console 1 Processor: Prefixer Content: contentTest
	// Output: Console 2 Processor: Prefixer Content: contentTest
	// Output: Console 3 Processor: Prefixer Content: contentTest
	// Output: Console 1 Processor: Prefixer Content: contentTest - My Suffix
}

// PrintPretty example.
func ExampleNew_printPretty() {
	type TestType struct {
		Key1 string
		Key2 int
	}

	New("Testing Logger", Console(level.Info)).PrintPretty(&TestType{
		Key1: "text",
		Key2: 12,
	})

	// output:
	// {
	// 	"Key1": "text",
	// 	"Key2": 12
	// }
}
