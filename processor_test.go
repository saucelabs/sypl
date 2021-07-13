// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/saucelabs/sypl/content"
	"github.com/saucelabs/sypl/level"
)

func Test_generateDefaultPrefix(t *testing.T) {
	type args struct {
		timestamp string
		component string
		level     level.Level
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should work",
			args: args{
				timestamp: time.Now().Format(defaultTimestampFormat),
				component: defaultComponentNameOutput,
				level:     level.Trace,
			},
			want: fmt.Sprintf("%d [%d] [%s] [Trace] ",
				time.Now().Year(),
				os.Getpid(),
				defaultComponentNameOutput,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateDefaultPrefix(tt.args.timestamp, tt.args.component, tt.args.level); got != tt.want {
				t.Errorf("generateDefaultPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrefixer(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name    string
		args    args
		message *Message
		want    string
	}{
		{
			name: "Should work",
			args: args{
				prefix: defaultPrefixValue,
			},
			message: NewMessage(nil, nil, nil, level.Info, defaultContentOutput),
			want:    defaultPrefixValue + defaultContentOutput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Prefixer(tt.args.prefix)
			p.Run(tt.message)

			if !strings.EqualFold(tt.message.GetProcessedContent(), tt.want) {
				t.Errorf("Prefixer() = %v, want %v", tt.message.GetProcessedContent(), tt.want)
			}
		})
	}
}

func TestSuffixer(t *testing.T) {
	type args struct {
		suffix string
	}
	tests := []struct {
		name    string
		args    args
		message *Message
		want    string
	}{
		{
			name: "Should work",
			args: args{
				suffix: " - My Suffix",
			},
			message: NewMessage(nil, nil, nil, level.Info, defaultContentOutput),
			want:    defaultContentOutput + " - My Suffix",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Suffixer(tt.args.suffix)
			p.Run(tt.message)

			if !strings.EqualFold(tt.message.GetProcessedContent(), tt.want) {
				t.Errorf("Suffixer() = %v, want %v", tt.message.GetProcessedContent(), tt.want)
			}
		})
	}
}

func TestNewProcessor(t *testing.T) {
	type args struct {
		name          string
		processorFunc ProcessorFunc
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should work",
			args: args{
				name: "Prefixer",
				processorFunc: func(message *Message) {
					message.SetProcessedContent(defaultPrefixValue + message.GetProcessedContent())
				},
			},
			want: defaultPrefixValue + defaultContentOutput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProcessor(tt.args.name, tt.args.processorFunc)

			m := &Message{
				content: content.NewContent(defaultContentOutput),
			}

			p.Run(m)

			if m.GetProcessedContent() != tt.want {
				t.Errorf("Got %v, want %v", m.GetProcessedContent(), tt.want)
			}
		})
	}
}

func TestProcessor_SetStatus(t *testing.T) {
	type args struct {
		name          string
		processorFunc ProcessorFunc
	}
	tests := []struct {
		name   string
		args   args
		status bool
		want   string
	}{
		{
			name: "Should work",
			args: args{
				name: "Prefixer",
				processorFunc: func(message *Message) {
					message.SetProcessedContent(defaultPrefixValue + message.GetProcessedContent())
				},
			},
			status: true,
			want:   defaultPrefixValue + defaultContentOutput,
		},
		{
			name: "Should work - No processing",
			args: args{
				name: "Prefixer",
				processorFunc: func(message *Message) {
					message.SetProcessedContent(defaultPrefixValue + message.GetProcessedContent())
				},
			},
			status: false,
			want:   defaultContentOutput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProcessor(tt.args.name, tt.args.processorFunc)

			m := &Message{
				content: content.NewContent(defaultContentOutput),
			}

			p.SetStatus(tt.status)

			p.Run(m)

			if m.GetProcessedContent() != tt.want {
				t.Errorf("Got %v, want %v", m.GetProcessedContent(), tt.want)
			}
		})
	}
}

func TestProcessor_GetStatus(t *testing.T) {
	type args struct {
		name          string
		processorFunc ProcessorFunc
	}
	tests := []struct {
		name   string
		args   args
		status bool
		want   bool
	}{
		{
			name: "Should work",
			args: args{
				name: "Prefixer",
				processorFunc: func(message *Message) {
					message.SetProcessedContent(defaultPrefixValue + message.GetProcessedContent())
				},
			},
			status: true,
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProcessor(tt.args.name, tt.args.processorFunc)

			p.SetStatus(tt.status)

			if p.GetStatus() != tt.want {
				t.Errorf("Got %v, want %v", p.GetStatus(), tt.want)
			}
		})
	}
}
