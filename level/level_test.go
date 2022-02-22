// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package level

import (
	"errors"
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
		name       string
		args       args
		want       Level
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Should work - Valid",
			args: args{
				level: "Info",
			},
			want:       Info,
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "Should fail - Invalid",
			args: args{
				level: "Invalid",
			},
			want:       None,
			wantErr:    true,
			wantErrMsg: "invalid error level: Invalid. Available: none, fatal, error, info, warn, debug, trace",
		},
		{
			name:       "Should fail - empty",
			args:       args{},
			want:       None,
			wantErr:    true,
			wantErrMsg: "invalid error level: No level specified. Available: none, fatal, error, info, warn, debug, trace",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromString(tt.args.level)
			if got != tt.want {
				t.Errorf("FromString() = %v, want %v", got, tt.want)

				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error = %v, want %v", err, tt.wantErr)

				return
			}

			if tt.wantErr && (err != nil) {
				if err.Error() != tt.wantErrMsg {
					t.Errorf("Expected error = %v, want %v", err.Error(), tt.wantErrMsg)

					return
				}

				if !errors.Is(err, ErrInvalidLevel) {
					t.Errorf("Expected error to be ErrInvalidLevel")

					return
				}
			}
		})
	}
}

func TestMustFromString(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name string
		args args
		want Level
	}{
		{
			name: "Should work - Valid - Uppercased",
			args: args{
				level: "Info",
			},
			want: Info,
		},
		{
			name: "Should work - Valid - Lowercased",
			args: args{
				level: "info",
			},
			want: Info,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustFromString(tt.args.level); got != tt.want {
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
			name: "Should work - Valid",
			args: args{
				levels: []Level{Info, Warn},
			},
			want: "info,warn",
		},
		{
			name: "Should work - Unknown",
			args: args{
				levels: []Level{Level(10)},
			},
			want: "Unknown",
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
