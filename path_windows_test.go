package main

import (
	"testing"
)

func TestGetRealPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Test Windows Path",
			"C:\\Users\\appleboy\\test.txt",
			"/C/Users/appleboy/test.txt",
		},
	}
	for _, tt := range tests {
		if got := getRealPath(tt.args.path); got != tt.want {
			t.Errorf("%q. getRealPath() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
