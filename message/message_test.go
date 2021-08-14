// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/saucelabs/sypl/flag"
	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/shared"
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
		level   level.Level
		content string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should work - tags, flags",
			args: args{
				level:   level.Info,
				content: "Test",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMessage(tt.args.level, tt.args.content)

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

func TestCopy(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Should work",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMessage(level.Info, shared.DefaultContentOutput)

			if res := deep.Equal(Copy(m), m); len(res) > 0 {
				t.Log("Expected:", shared.Prettify(m))
				t.Log("Got:", shared.Prettify(Copy(m)))
				t.Error("Diff:", res)
			}
		})
	}
}

func Test_message_strip(t *testing.T) {
	tests := []struct {
		name                string
		m                   IMessage
		wantChars           []string
		wantContent         string
		wantContentOriginal string
		wantSize            int
	}{
		{
			name:                "Should work - \\n",
			m:                   NewMessage(level.Info, "Test 1\n"),
			wantChars:           []string{"\n"},
			wantContent:         "Test 1",
			wantContentOriginal: "Test 1\n",
			wantSize:            1,
		},
		{
			name:                "Should work - Println",
			m:                   NewMessage(level.Info, fmt.Sprintln("Test 1")),
			wantChars:           []string{"\n"},
			wantContent:         "Test 1",
			wantContentOriginal: "Test 1\n",
			wantSize:            1,
		},
		{
			name:                "Should work - Printf",
			m:                   NewMessage(level.Info, fmt.Sprintf("%s\n", "Test 1")),
			wantChars:           []string{"\n"},
			wantContent:         "Test 1",
			wantContentOriginal: "Test 1\n",
			wantSize:            1,
		},
		{
			name:                "Should work - \\r",
			m:                   NewMessage(level.Info, "Test 1\r"),
			wantChars:           []string{"\r"},
			wantContent:         "Test 1",
			wantContentOriginal: "Test 1\r",
			wantSize:            1,
		},
		{
			name:                "Should work - many (\\n\\r\\n)",
			m:                   NewMessage(level.Info, "Test 1\n\r\n"),
			wantChars:           []string{"\n", "\r", "\n"},
			wantContent:         "Test 1",
			wantContentOriginal: "Test 1\n\r\n",
			wantSize:            3,
		},
		{
			name:                "Should work - none",
			m:                   NewMessage(level.Info, "Test 1"),
			wantChars:           []string{},
			wantContent:         "Test 1",
			wantContentOriginal: "Test 1",
			wantSize:            0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Strip()

			if tt.wantSize != len(tt.m.getLineBreaker().ControlChars) {
				t.Errorf("Strip got %v expected %v", tt.wantSize, len(tt.m.getLineBreaker().ControlChars))
			}

			if d := deep.Equal(tt.m.getLineBreaker().ControlChars, tt.wantChars); len(d) > 0 {
				t.Errorf("Strip %+v", d)
			}

			if tt.wantContent != tt.m.GetContent().GetProcessed() {
				t.Errorf("Strip got %v expected %v", tt.wantContent, tt.m.GetContent().GetProcessed())
			}

			tt.m.Restore()

			if tt.m.GetContent().GetProcessed() != tt.wantContentOriginal {
				t.Errorf("Restore got %v expected %v", tt.m.GetContent().GetProcessed(), tt.wantContentOriginal)
			}
		})
	}
}
