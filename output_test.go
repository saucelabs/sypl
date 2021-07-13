// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/saucelabs/sypl/level"
)

func TestNewOutput(t *testing.T) {
	type args struct {
		maxLevel level.Level
		name     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should work",
			args: args{
				name:     "Buffer",
				maxLevel: level.Trace,
			},
			want: defaultPrefixValue + defaultContentOutput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			output := NewOutput(tt.args.name, tt.args.maxLevel, bufWriter, Prefixer(defaultPrefixValue))

			message := NewMessage(&Sypl{}, &Output{}, &Processor{}, level.Info, defaultContentOutput)
			message.SetSypl(nil)
			message.SetOutput(nil)
			message.SetProcessor(nil)

			if message.GetSypl() != nil &&
				message.GetOutput() != nil &&
				message.GetProcessor() != nil {
				t.Error("Message should not have sypl, output, and processor")
			}

			for _, processor := range output.processors {
				processor.Run(message)
			}

			if err := output.GetBuiltinLogger().OutputBuiltin(defaultCallDepth, message.GetProcessedContent()); err != nil {
				t.Errorf("Failed to log to output: %w", err)
			}

			bufWriter.Flush()

			if buf.String() != tt.want {
				t.Errorf("Got %v, want %v", buf.String(), tt.want)
			}
		})
	}
}

func TestOutput_GetStatus(t *testing.T) {
	type args struct {
		maxLevel level.Level
		name     string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should work",
			args: args{
				name:     "Buffer",
				maxLevel: level.Trace,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := Console(level.Trace)
			output.SetStatus(false)

			if output.GetStatus() != tt.want {
				t.Errorf("Got %v, want %v", output.GetStatus(), tt.want)
			}
		})
	}
}
