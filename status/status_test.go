// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package status

import "testing"

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name string
		f    Status
		want string
	}{
		{
			name: "Should work - Enabled",
			f:    Enabled,
			want: "Enabled",
		},
		{
			name: "Should work - Disabled",
			f:    Disabled,
			want: "Disabled",
		},
		{
			name: "Should work - Unknown",
			f:    Status(10),
			want: "Unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.String(); got != tt.want {
				t.Errorf("Status.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
