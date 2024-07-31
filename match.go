package detector

import "net/http"

// Match holds a URL and it's associated filename and line/column.
type Match struct {
	Filename     string
	Line, Column int

	URL    string
	Result string
}

// Check makes a request to this Match's URL and checks the response status.
func (m *Match) Check() {
	resp, err := http.Get(m.URL)
	if err != nil {
		m.Result = err.Error()
		return
	}
	resp.Body.Close()
	m.Result = resp.Status
}
