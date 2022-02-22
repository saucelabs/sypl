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

	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/message"
	"github.com/saucelabs/sypl/options"
	"github.com/saucelabs/sypl/output"
	"github.com/saucelabs/sypl/processor"
	"github.com/saucelabs/sypl/shared"
	"github.com/spf13/afero"
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
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.None,
		maxLevel:  level.Trace,
		run: func(a args) string {
			buf, o := output.SafeBuffer(a.maxLevel, processor.PrefixBasedOnMask(shared.DefaultTimestampFormat))

			New(a.component).AddOutputs(o).Print(a.level, a.content)

			return buf.String()
		},
	}

	infoArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Debug,
		run: func(a args) string {
			buf, o := output.SafeBuffer(a.maxLevel, processor.PrefixBasedOnMask(shared.DefaultTimestampFormat))

			New(a.component).AddOutputs(o).Print(a.level, a.content)

			return buf.String()
		},
	}

	aboveArgs := args{
		level:    level.Trace,
		maxLevel: level.Debug,
		run: func(a args) string {
			buf, o := output.SafeBuffer(a.maxLevel, processor.PrefixBasedOnMask(shared.DefaultTimestampFormat))

			New(a.component).AddOutputs(o).Print(a.level, a.content)

			return buf.String()
		},
	}

	mutedArgs := args{
		level:    level.Info, // Will not be used.
		maxLevel: level.Trace,
		run: func(a args) string {
			buf, o := output.SafeBuffer(a.maxLevel, processor.MuteBasedOnLevel(level.Info, level.Warn))

			New(a.component).
				AddOutputs(o).
				Printf(level.Info, "%s", a.content).
				Printf(level.Info, "%s", a.content).
				Printf(level.Warn, "%s", a.content)

			return buf.String()
		},
	}

	fileArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
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
				shared.DefaultFileMode)
			if err != nil {
				t.Error("Failed to open virtal file", err)
			}

			defer f.Close()

			New(a.component).
				AddOutputs(output.FileBased("virtual", level.Debug, f, processor.Prefixer("Test Prefix - "))).
				Print(a.level, a.content)

			b, err := afero.ReadFile(appFs, filePath)
			if err != nil {
				t.Error("Failed to read virtal file", err)
			}

			return string(b)
		},
	}

	disableArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Trace,
		run: func(a args) string {
			buf, o := output.SafeBuffer(a.maxLevel, processor.PrefixBasedOnMask(shared.DefaultTimestampFormat))

			New(a.component).AddOutputs(o).PrintWithOptions(&options.Options{
				ProcessorsNames: []string{""},
			}, a.level, a.content)

			return buf.String()
		},
	}

	errorArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.None, // Will not be used.
		maxLevel:  level.Trace,
		run: func(a args) string {
			buf, o := output.SafeBuffer(a.maxLevel, processor.PrefixBasedOnMask(shared.DefaultTimestampFormat))

			New(a.component).AddOutputs(o).Errorf("%s", a.content)

			return buf.String()
		},
	}

	info2Args := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.None, // Will not be used.
		maxLevel:  level.Trace,
		run: func(a args) string {
			buf, o := output.SafeBuffer(a.maxLevel, processor.PrefixBasedOnMask(shared.DefaultTimestampFormat))

			New(a.component).AddOutputs(o).Infof("%s", a.content)

			return buf.String()
		},
	}

	warnArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.None, // Will not be used.
		maxLevel:  level.Trace,
		run: func(a args) string {
			buf, o := output.SafeBuffer(a.maxLevel, processor.PrefixBasedOnMask(shared.DefaultTimestampFormat))

			New(a.component).AddOutputs(o).Warnf("%s", a.content)

			return buf.String()
		},
	}

	debugArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.None, // Will not be used.
		maxLevel:  level.Trace,
		run: func(a args) string {
			buf, o := output.SafeBuffer(a.maxLevel, processor.PrefixBasedOnMask(shared.DefaultTimestampFormat))

			New(a.component).AddOutputs(o).Debugf("%s", a.content)

			return buf.String()
		},
	}

	trace2Args := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.None, // Will not be used.
		maxLevel:  level.Trace,
		run: func(a args) string {
			buf, o := output.SafeBuffer(a.maxLevel, processor.PrefixBasedOnMask(shared.DefaultTimestampFormat))

			New(a.component).AddOutputs(o).Tracef("%s", a.content)

			return buf.String()
		},
	}

	forceArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.Error, // Will not be used.
		maxLevel:  level.Fatal,
		run: func(a args) string {
			buf, o := output.SafeBuffer(a.maxLevel, processor.ForceBasedOnLevel(level.Error, level.Warn))

			New(a.component).AddOutputs(o).Printf(level.Error, "%s", a.content).Printf(level.Warn, "%s", a.content)

			return buf.String()
		},
	}

	printfArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Trace,
		run: func(a args) string {
			buf, o := output.SafeBuffer(a.maxLevel)

			New(a.component).AddOutputs(o).Printf(a.level, "%s", a.content)

			return buf.String()
		},
	}

	printfNewLineArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Trace,
		run: func(a args) string {
			buf, o := output.SafeBuffer(a.maxLevel)

			New(a.component).AddOutputs(o).Printf(a.level, "%s\n", a.content)

			return buf.String()
		},
	}

	printlnArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Trace,
		run: func(a args) string {
			buf, o := output.SafeBuffer(a.maxLevel)

			New(a.component).AddOutputs(o).Println(a.level, a.content)

			return buf.String()
		},
	}

	prefixBasedOnMaskExceptForLevelsArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.Info, // Will not be used.
		maxLevel:  level.Trace,
		run: func(a args) string {
			buf, o := output.SafeBuffer(
				a.maxLevel,
				processor.PrefixBasedOnMaskExceptForLevels(
					shared.DefaultTimestampFormat,
					level.Info,
					level.Warn,
				),
			)

			New(a.component).
				AddOutputs(o).
				Printf(level.Info, "%s", a.content).
				Printf(level.Warn, "%s", a.content)

			return buf.String()
		},
	}

	prefixBasedOnMaskExceptForLevelsDontArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Trace,
		run: func(a args) string {
			buf, o := output.SafeBuffer(
				a.maxLevel,
				processor.PrefixBasedOnMaskExceptForLevels(
					shared.DefaultTimestampFormat,
					level.Warn),
			)

			New(a.component).AddOutputs(o).Printf(a.level, "%s", a.content)

			return buf.String()
		},
	}

	printWithOptionsArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Trace,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutputs(output.New("buffer 1", a.maxLevel, bufWriter)).
				AddOutputs(output.New("buffer 2", a.maxLevel, bufWriter)).
				PrintWithOptions(&options.Options{
					OutputsNames: []string{"buffer 1"},
				}, a.level, shared.DefaultContentOutput)

			bufWriter.Flush()

			return buf.String()
		},
	}

	printWithOptionsDontArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Trace,
		run: func(a args) string {
			buf, o := output.SafeBuffer(a.maxLevel)

			New(a.component).AddOutputs(o).PrintWithOptions(&options.Options{
				OutputsNames: []string{"invalid"},
			}, a.level, shared.DefaultContentOutput)

			return buf.String()
		},
	}

	enableDisableOutputsArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Trace,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutputs(output.New("buffer 1", a.maxLevel, bufWriter)).
				AddOutputs(output.New("buffer 2", a.maxLevel, bufWriter)).
				PrintWithOptions(&options.Options{
					OutputsNames: []string{"buffer 2"},
				}, level.Info, shared.DefaultContentOutput)

			bufWriter.Flush()

			return buf.String()
		},
	}

	changeFirstCharCaseUpperArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Trace,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutputs(output.New("buffer 1", a.maxLevel, bufWriter, processor.ChangeFirstCharCase(processor.Uppercase))).
				Info(shared.DefaultContentOutput)

			bufWriter.Flush()

			return buf.String()
		},
	}

	changeFirstCharCaseLowerArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   "ContentTest",
		level:     level.Info,
		maxLevel:  level.Trace,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			New(a.component).
				AddOutputs(output.New("buffer 1", a.maxLevel, bufWriter, processor.ChangeFirstCharCase(processor.Lowercase))).
				Info(shared.DefaultContentOutput)

			bufWriter.Flush()

			return buf.String()
		},
	}

	nonChainedNewLoggerArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Trace,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			// Creates logger, and name it.
			testingLogger := New("Testing Logger 1")

			// Creates an `Output`. In this case, called Console that will print to
			// stdout and max print level @ Info.
			ConsoleToStdOut := output.New("Console", level.Info, bufWriter)

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

			// Prints: "My Prefix - Test message"
			testingLogger.Print(level.Info, "Test message")

			bufWriter.Flush()

			return buf.String()
		},
	}

	printflnArgs := args{
		component: shared.DefaultComponentNameOutput,
		content:   shared.DefaultContentOutput,
		level:     level.Info,
		maxLevel:  level.Trace,
		run: func(a args) string {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			// Creates logger, and name it.
			testingLogger := New("Testing Logger 1")

			// Creates an `Output`. In this case, called Buffer that will write
			// to the specified buffer, and max print level @ Info.
			BufferOutput := output.New("Buffer", level.Info, bufWriter)

			// Adds `Output` to logger.
			testingLogger.AddOutputs(BufferOutput)

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
				return "Test Prefix - " + shared.DefaultContentOutput
			},
		},
		{
			name: "Should print - Only prefix (Disabler)",
			args: disableArgs,
			want: func(a args) string {
				return shared.DefaultContentOutput
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
					"error",
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
					"info",
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
					"warn",
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
					"debug",
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
					"trace",
					a.content)
			},
		},
		{
			name: "Should print - Force",
			args: forceArgs,
			want: func(a args) string {
				return shared.DefaultContentOutput + shared.DefaultContentOutput
			},
		},
		{
			name: "Should print - Printf - No newline",
			args: printfArgs,
			want: func(a args) string {
				return shared.DefaultContentOutput
			},
		},
		{
			name: "Should print - Printf - Newline",
			args: printfNewLineArgs,
			want: func(a args) string {
				return shared.DefaultContentOutput + "\n"
			},
		},
		{
			name: "Should print - Println",
			args: printlnArgs,
			want: func(a args) string {
				return shared.DefaultContentOutput + "\n"
			},
		},
		{
			name: "Should print not prefixed - PrefixBasedOnMaskExceptForLevels",
			args: prefixBasedOnMaskExceptForLevelsArgs,
			want: func(a args) string {
				return shared.DefaultContentOutput + shared.DefaultContentOutput
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
				return shared.DefaultContentOutput
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
				return shared.DefaultContentOutput
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
				return shared.DefaultContentOutput
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
