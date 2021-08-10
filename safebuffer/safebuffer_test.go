// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package safebuffer

import (
	"fmt"
	"testing"
)

const defaultContentOutput = "contentTest"

func TestBuffer_String(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "Should work",
			arg:  defaultContentOutput,
			want: defaultContentOutput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf Buffer

			fmt.Fprint(&buf, tt.arg)

			buf.Reset()

			fmt.Fprint(&buf, tt.arg)

			if buf.String() != tt.want {
				t.Errorf("Expected %s got %s", defaultContentOutput, buf.String())
			}
		})
	}
}
