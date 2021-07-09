// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl

import (
	"bufio"
	"bytes"
	"testing"
)

func TestNewOutput(t *testing.T) {
	type args struct {
		maxLevel Level
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
				maxLevel: TRACE,
			},
			want: defaultPrefixValue + defaultContentOutput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			bufWriter := bufio.NewWriter(&buf)

			output := NewOutput(tt.args.name, tt.args.maxLevel, bufWriter, Prefixer(defaultPrefixValue))

			message := &Message{
				ContentOriginal:  defaultContentOutput,
				ContentProcessed: defaultContentOutput,
				Level:            INFO,
			}

			for _, processor := range output.processors {
				processor.Run(message)
			}

			if err := output.Logger.OutputBuiltin(defaultCallDepth, message.ContentProcessed); err != nil {
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
		maxLevel Level
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
				maxLevel: TRACE,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := Console(TRACE)
			output.SetStatus(false)

			if output.GetStatus() != tt.want {
				t.Errorf("Got %v, want %v", output.GetStatus(), tt.want)
			}
		})
	}
}
