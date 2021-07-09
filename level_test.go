// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl

import (
	"testing"
)

func TestLevelFromInt(t *testing.T) {
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
			want: INFO,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LevelFromInt(tt.args.level); got != tt.want {
				t.Errorf("LevelFromInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelFromString(t *testing.T) {
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
				level: "info",
			},
			want: INFO,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LevelFromString(tt.args.level); got != tt.want {
				t.Errorf("LevelFromString() = %v, want %v", got, tt.want)
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
				levels: []Level{INFO, WARN},
			},
			want: "INFO,WARN",
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
