package ded

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIsBinaryFile(t *testing.T) {
	for _, tt := range []struct {
		path string
		want bool
	}{
		{"foo.png", true},
		{"foo.PNG", true},
		{"foo.txt", false},
		{"README", false},
	} {
		if got := isBinaryFilename(tt.path); !cmp.Equal(got, tt.want) {
			t.Errorf(
				"isBinaryFilename(%q)\ngot %v\nwant %v",
				tt.path,
				got,
				tt.want,
			)
		}
	}
}

func TestIsSCMPath(t *testing.T) {
	for _, tt := range []struct {
		path string
		want bool
	}{
		{"foo.png", false},
		{"foo/.git/whatever", true},
	} {
		if got := isSCMPath(tt.path); !cmp.Equal(got, tt.want) {
			t.Errorf(
				"isSCMPath(%q)\ngot %v\nwant %v",
				tt.path,
				got,
				tt.want,
			)
		}
	}
}
