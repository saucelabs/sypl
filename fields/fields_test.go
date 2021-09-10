package fields

import (
	"reflect"
	"testing"
)

func TestCopy(t *testing.T) {
	type args struct {
		src Fields
		dst Fields
	}
	tests := []struct {
		name string
		args args
		want Fields
	}{
		{
			name: "Should work",
			args: args{
				src: Fields{"test": 2},
				dst: Fields{"test": 1},
			},
			want: Fields{"test": 2},
		},
		{
			name: "Should work",
			args: args{
				src: Fields{"test": 1},
				dst: Fields{"test": 2},
			},
			want: Fields{"test": 1},
		},
		{
			name: "Should work - only src - both init",
			args: args{
				src: Fields{"test": 1},
				dst: Fields{},
			},
			want: Fields{"test": 1},
		},
		{
			name: "Should work - only dst - both init",
			args: args{
				src: Fields{"test": 1},
				dst: Fields{},
			},
			want: Fields{"test": 1},
		},
		{
			name: "Should work - only src - dst nil",
			args: args{
				src: Fields{"test": 1, "test2": 2},
			},
			want: Fields{"test": 1, "test2": 2},
		},
		{
			name: "Should work - only dst - src nil",
			args: args{
				dst: Fields{"test": 1, "test2": 2},
			},
			want: nil,
		},
		{
			name: "Should work - no src and dst",
			args: args{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Copy(tt.args.src, tt.args.dst); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}
