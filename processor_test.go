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
)

func Test_generateDefaultPrefix(t *testing.T) {
	type args struct {
		timestampFormat string
		component       string
		level           Level
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should work",
			args: args{
				timestampFormat: defaultTimestampFormat,
				component:       defaultComponentNameOutput,
				level:           TRACE,
			},
			want: fmt.Sprintf("%d [%d] [%s] [TRACE] ",
				time.Now().Year(),
				os.Getpid(),
				defaultComponentNameOutput,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateDefaultPrefix(tt.args.timestampFormat, tt.args.component, tt.args.level); got != tt.want {
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
			message: &Message{
				ContentProcessed: defaultContentOutput,
			},
			want: defaultPrefixValue + defaultContentOutput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Prefixer(tt.args.prefix)
			p.Run(tt.message)

			if !strings.EqualFold(tt.message.ContentProcessed, tt.want) {
				t.Errorf("Prefixer() = %v, want %v", tt.message.ContentProcessed, tt.want)
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
					message.ContentProcessed = defaultPrefixValue + message.ContentProcessed
				},
			},
			want: defaultPrefixValue + defaultContentOutput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProcessor(tt.args.name, tt.args.processorFunc)

			m := &Message{
				ContentOriginal:  defaultContentOutput,
				ContentProcessed: defaultContentOutput,
			}

			p.Run(m)

			if m.ContentProcessed != tt.want {
				t.Errorf("Got %v, want %v", m.ContentProcessed, tt.want)
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
					message.ContentProcessed = defaultPrefixValue + message.ContentProcessed
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
					message.ContentProcessed = defaultPrefixValue + message.ContentProcessed
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
				ContentOriginal:  defaultContentOutput,
				ContentProcessed: defaultContentOutput,
			}

			p.SetStatus(tt.status)

			p.Run(m)

			if m.ContentProcessed != tt.want {
				t.Errorf("Got %v, want %v", m.ContentProcessed, tt.want)
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
					message.ContentProcessed = defaultPrefixValue + message.ContentProcessed
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
