// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package level

import (
	"testing"
)

func TestFromInt(t *testing.T) {
	type args struct {
		level int
	}
	tests := []struct {
		name string
		args args
		want Level
	}{
		{
			name: "Should work",
			args: args{
				level: 3,
			},
			want: Info,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromInt(tt.args.level); got != tt.want {
				t.Errorf("FromInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromString(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name string
		args args
		want Level
	}{
		{
			name: "Should work",
			args: args{
				level: "Info",
			},
			want: Info,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromString(tt.args.level); got != tt.want {
				t.Errorf("FromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelsToString(t *testing.T) {
	type args struct {
		levels []Level
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should work",
			args: args{
				levels: []Level{Info, Warn},
			},
			want: "Info,Warn",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LevelsToString(tt.args.levels); got != tt.want {
				t.Errorf("LevelsToString() = %v, want %v", got, tt.want)
			}
		})
	}
}