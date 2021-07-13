// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl

import (
	"testing"

	"github.com/saucelabs/sypl/flag"
	"github.com/saucelabs/sypl/level"
)

func Test_generateUUID(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Should work",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateUUID()

			if len(got) < 30 {
				t.Errorf("generateUUID() = %v", got)
			}
		})
	}
}

func TestNewMessage(t *testing.T) {
	type args struct {
		sypl      *Sypl
		output    *Output
		processor *Processor
		level     level.Level
		content   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should work - tags, flags",
			args: args{
				sypl:      nil,
				output:    nil,
				processor: nil,
				level:     level.Info,
				content:   "Test",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMessage(tt.args.sypl, tt.args.output, tt.args.processor, tt.args.level, tt.args.content)

			lenID := len(m.GetID())

			if lenID < 30 {
				t.Errorf("Got %d chars, expected %d chars", lenID, 30)
			}

			if m.GetFlag() != flag.None && len(m.GetTags()) != 0 {
				t.Errorf("Expected %s flag, and %d tags", flag.None.String(), 0)
			}

			m.SetFlag(flag.Force)
			m.AddTags("x", "y")
			if m.GetFlag() != flag.Force && m.GetTags()[0] != "x" && m.GetTags()[1] != "y" {
				t.Errorf("Expected %s flag, and %s and %s tags", flag.Force.String(), "x", "y")
			}

			m.SetFlag(flag.Mute)
			m.DeleteTag("x")
			if m.GetFlag() != flag.Mute && m.GetTags()[0] != "y" {
				t.Errorf("Expected %s flag, and %s tag", flag.Mute.String(), "y")
			}
		})
	}
}
