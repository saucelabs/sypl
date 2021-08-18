// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package processor

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/message"
	"github.com/saucelabs/sypl/shared"
	"github.com/saucelabs/sypl/status"
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
				timestamp: time.Now().Format(shared.DefaultTimestampFormat),
				component: shared.DefaultComponentNameOutput,
				level:     level.Trace,
			},
			want: fmt.Sprintf("%d [%d] [%s] [Trace] ",
				time.Now().Year(),
				os.Getpid(),
				shared.DefaultComponentNameOutput,
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
		message message.IMessage
		want    string
	}{
		{
			name: "Should work",
			args: args{
				prefix: shared.DefaultPrefixValue,
			},
			message: message.NewMessage(level.Info, shared.DefaultContentOutput),
			want:    shared.DefaultPrefixValue + shared.DefaultContentOutput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Prefixer(tt.args.prefix)
			if err := p.Run(tt.message); err != nil {
				t.Errorf("Run failed: %s", err)
			}

			if !strings.EqualFold(tt.message.GetContent().GetProcessed(), tt.want) {
				t.Errorf("Prefixer() = %v, want %v", tt.message.GetContent().GetProcessed(), tt.want)
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
		message message.IMessage
		want    string
	}{
		{
			name: "Should work",
			args: args{
				suffix: " - My Suffix",
			},
			message: message.NewMessage(level.Info, shared.DefaultContentOutput),
			want:    shared.DefaultContentOutput + " - My Suffix",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Suffixer(tt.args.suffix)
			if err := p.Run(tt.message); err != nil {
				t.Errorf("Run failed: %s", err)
			}

			if !strings.EqualFold(tt.message.GetContent().GetProcessed(), tt.want) {
				t.Errorf("Suffixer() = %v, want %v", tt.message.GetContent().GetProcessed(), tt.want)
			}
		})
	}
}

func TestNewProcessor(t *testing.T) {
	type args struct {
		name    string
		RunFunc RunFunc
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
				RunFunc: func(m message.IMessage) error {
					m.GetContent().SetProcessed(shared.DefaultPrefixValue + m.GetContent().GetProcessed())

					return nil
				},
			},
			want: shared.DefaultPrefixValue + shared.DefaultContentOutput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProcessor(tt.args.name, tt.args.RunFunc)

			m := message.NewMessage(level.Info, shared.DefaultContentOutput)

			if err := p.Run(m); err != nil {
				t.Errorf("Run failed: %s", err)
			}

			if m.GetContent().GetProcessed() != tt.want {
				t.Errorf("Got %v, want %v", m.GetContent().GetProcessed(), tt.want)
			}
		})
	}
}

func TestProcessor_SetStatus(t *testing.T) {
	type args struct {
		name    string
		RunFunc RunFunc
	}
	tests := []struct {
		name   string
		args   args
		status status.Status
		want   string
	}{
		{
			name: "Should work - No processing",
			args: args{
				name: "Prefixer",
				RunFunc: func(message message.IMessage) error {
					message.GetContent().SetProcessed(shared.DefaultPrefixValue + message.GetContent().GetProcessed())

					return nil
				},
			},
			status: status.Disabled,
			want:   shared.DefaultContentOutput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProcessor(tt.args.name, tt.args.RunFunc)

			m := message.NewMessage(level.Info, shared.DefaultContentOutput)

			p.SetStatus(tt.status)

			if err := p.Run(m); err != nil {
				t.Errorf("Run failed: %s", err)
			}

			if m.GetContent().GetProcessed() != tt.want {
				t.Errorf("Got %v, want %v", m.GetContent().GetProcessed(), tt.want)
			}
		})
	}
}

func TestProcessor_GetStatus(t *testing.T) {
	type args struct {
		name    string
		RunFunc RunFunc
	}
	tests := []struct {
		name   string
		args   args
		status status.Status
		want   status.Status
	}{
		{
			name: "Should work",
			args: args{
				name: "Prefixer",
				RunFunc: func(message message.IMessage) error {
					message.GetContent().SetProcessed(shared.DefaultPrefixValue + message.GetContent().GetProcessed())

					return nil
				},
			},
			status: status.Enabled,
			want:   status.Enabled,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProcessor(tt.args.name, tt.args.RunFunc)

			p.SetStatus(tt.status)

			if p.GetStatus() != tt.want {
				t.Errorf("Got %v, want %v", p.GetStatus(), tt.want)
			}
		})
	}
}
