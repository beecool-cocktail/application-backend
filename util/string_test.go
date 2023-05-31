package util

import (
	"testing"
)

func TestConcatString(t *testing.T) {
	type args struct {
		strs []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Concatenating strings with spaces",
			args: args{
				strs: []string{"Hello", " ", "World", "!"},
			},
			want: "Hello World!",
		},
		{
			name: "Concatenating multiple words",
			args: args{
				strs: []string{"I", " ", "am", " ", "Golang"},
			},
			want: "I am Golang",
		},
		{
			name: "Concatenating empty strings",
			args: args {
				strs: []string{"", "", ""},
			},
			want: "",
		},
		{
			name: "Concatenating numeric strings",
			args: args{
				strs: []string{"123", "456", "789"},
			},
			want: "123456789",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConcatString(tt.args.strs...); got != tt.want {
				t.Errorf("ConcatString() = %v, want %v", got, tt.want)
			}
		})
	}
}
