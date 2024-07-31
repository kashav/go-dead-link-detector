package detector

import (
	"bufio"
	"regexp"
	"strings"

	"mvdan.cc/xurls"
)

// Parse parses all links from a single file and returns a list of Matches.
func Parse(filename, raw, ignoreExpr string) []Match {
	var ignore = func(url string) bool {
		return len(ignoreExpr) > 0 && regexp.MustCompile(ignoreExpr).MatchString(url)
	}

	var (
		line   string
		lineNo int

		matches = make([]Match, 0)
		scanner = bufio.NewScanner(strings.NewReader(raw))
	)

	for scanner.Scan() {
		lineNo++
		line = scanner.Text()
		if ms := xurls.Strict.FindAllStringIndex(line, -1); ms != nil {
			for _, m := range ms {
				if ignore(line[m[0]:m[1]]) {
					continue
				}
				matches = append(matches, Match{
					Filename: filename,
					Line:     lineNo,
					Column:   m[0],
					URL:      line[m[0]:m[1]],
				})
			}
		}
	}

	return matches
}
