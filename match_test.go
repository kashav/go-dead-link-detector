package ded

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMatch(t *testing.T) {
	for _, tt := range []struct {
		URL  string
		want string
	}{
		{"https://example.com", "200 OK"},
		{"https://github.com/oasfasidjai", "404 Not Found"},
		{"https://kashavmadan.com", "Get https://kashavmadan.com: dial tcp: lookup kashavmadan.com: no such host"},
	} {
		m := Match{URL: tt.URL}
		m.Check()
		if got := m.Result; !cmp.Equal(got, tt.want) {
			t.Errorf("Match{URL: %q}.Check()\ngot %v\nwant %v", m.URL, got, tt.want)
		}
	}
}
