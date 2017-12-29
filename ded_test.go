package ded

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPath_NoMatches(t *testing.T) {
	for _, tt := range []struct {
		filename, raw, ignoreExpr string
		want                      []Match
	}{
		{"empty", "", "", []Match{}},
		{"file1", "text", "", []Match{}},
		{"file2", "https://example.com", "example", []Match{}},
		{"file3", "https://example.com", "https://example.com", []Match{}},
		{"file4", "text\nhttps://example.com", "https://example.com", []Match{}},
		{"file5", "https://example.com\nhttps://example.com", "https://example.com", []Match{}},
		{"file6", "https://example.com?a=1", "https://example.com", []Match{}},
	} {
		if got := Parse(tt.filename, tt.raw, tt.ignoreExpr); !cmp.Equal(got, tt.want) {
			t.Errorf(
				"Parse(%q, %q, %q)\ngot %v\nwant %v",
				tt.filename,
				tt.raw,
				tt.ignoreExpr,
				got,
				tt.want,
			)
		}
	}
}

func TestPath_SingleMatch(t *testing.T) {
	for _, tt := range []struct {
		filename, raw, ignoreExpr string
		want                      []Match
	}{
		{"file1", "https://example.com", "", []Match{{"file1", 1, 0, "https://example.com", ""}}},
		{"file2", "text\nhttps://example.com", "", []Match{{"file2", 2, 0, "https://example.com", ""}}},
		{"file3", `  https://example.com`, "", []Match{{"file3", 1, 2, "https://example.com", ""}}},
		{"file4", "https://example.com\nhttps://google.com", "https://google.com", []Match{{"file4", 1, 0, "https://example.com", ""}}},
		{"file5", `https://example.com  https://google.com`, "https://google.com", []Match{{"file5", 1, 0, "https://example.com", ""}}},
		{"file6", `https://google.com  https://example.com`, "https://google.com", []Match{{"file6", 1, 20, "https://example.com", ""}}},
		{"file7", `https://example.com?a=https://example.com`, "", []Match{{"file7", 1, 0, "https://example.com?a=https://example.com", ""}}},
	} {
		if got := Parse(tt.filename, tt.raw, tt.ignoreExpr); !cmp.Equal(got, tt.want) {
			t.Errorf(
				"Parse(%q, %q, %q)\ngot %v\nwant %v",
				tt.filename,
				tt.raw,
				tt.ignoreExpr,
				got,
				tt.want,
			)
		}
	}
}

func TestPath_MultipleMatches(t *testing.T) {
	for _, tt := range []struct {
		filename, raw, ignoreExpr string
		want                      []Match
	}{
		{"file1", `https://example.com  https://example.com`, "", []Match{
			{"file1", 1, 0, "https://example.com", ""},
			{"file1", 1, 21, "https://example.com", ""}}},
		{"file2", `https://example.com  https://example.com  https://google.com`, "", []Match{
			{"file2", 1, 0, "https://example.com", ""},
			{"file2", 1, 21, "https://example.com", ""},
			{"file2", 1, 42, "https://google.com", ""}}},
		{
			"file3",
			`https://example.com  https://example.com  https://google.com`,
			"https://google.com",
			[]Match{
				{"file3", 1, 0, "https://example.com", ""},
				{"file3", 1, 21, "https://example.com", ""}}},
		{
			"file4",
			`https://example.com
https://example.com
https://google.com`,
			"",
			[]Match{
				{"file4", 1, 0, "https://example.com", ""},
				{"file4", 2, 0, "https://example.com", ""},
				{"file4", 3, 0, "https://google.com", ""},
			},
		},
		{
			"file5",
			`text text https://example.com text
text https://example.com?a=1 text`,
			"",
			[]Match{
				{"file5", 1, 10, "https://example.com", ""},
				{"file5", 2, 5, "https://example.com?a=1", ""},
			},
		},
	} {
		if got := Parse(tt.filename, tt.raw, tt.ignoreExpr); !cmp.Equal(got, tt.want) {
			t.Errorf(
				"Parse(%q, %q, %q)\ngot %v\nwant %v",
				tt.filename,
				tt.raw,
				tt.ignoreExpr,
				got,
				tt.want,
			)
		}
	}
}
