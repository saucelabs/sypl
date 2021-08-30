// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package processor

import (
	"errors"
	"testing"

	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/message"
)

func TestNewProcessingError(t *testing.T) {
	type args struct {
		m message.IMessage
		e error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should work",
			args: args{
				m: message.New(level.Info, "Message content"),
				e: errors.New("Error content"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewProcessingError(tt.args.m, tt.args.e); (err != nil) != tt.wantErr {
				t.Errorf("NewProcessingError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcessingError_Error(t *testing.T) {
	type fields struct {
		Cause         error
		Message       message.IMessage
		OutputName    string
		ProcessorName string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should work",
			fields: fields{
				Cause:         errors.New("Error content"),
				Message:       message.New(level.Info, "Message content"),
				OutputName:    "Test Output",
				ProcessorName: "Test Processor",
			},
			want: `Output: "Test Output" Processor: "Test Processor" Error: "Error content" Original Message: "Message content"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProcessingError{
				Cause:         tt.fields.Cause,
				Message:       tt.fields.Message,
				OutputName:    tt.fields.OutputName,
				ProcessorName: tt.fields.ProcessorName,
			}
			if got := p.Error(); got != tt.want {
				t.Errorf("ProcessingError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
