package pdf

import (
	"strings"
	"unicode/utf8"
)

// snippetRadius is how many runes of context surround a match.
const snippetRadius = 36

// Match is one hit of a text search.
type Match struct {
	Page    int
	Snippet string
}

// FindInText returns up to limit matches of query (case-insensitive) in
// text, with surrounding context collapsed onto one line.
func FindInText(page int, text, query string, limit int) []Match {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return nil
	}
	lower := strings.ToLower(text)
	var out []Match
	for from := 0; len(out) < limit; {
		i := strings.Index(lower[from:], q)
		if i < 0 {
			break
		}
		at := from + i
		out = append(out, Match{Page: page, Snippet: snippet(text, at, len(q))})
		from = at + len(q)
	}
	return out
}

func snippet(text string, at, n int) string {
	start, end := at, at+n
	for i := 0; i < snippetRadius && start > 0; i++ {
		_, size := utf8.DecodeLastRuneInString(text[:start])
		start -= size
	}
	for i := 0; i < snippetRadius && end < len(text); i++ {
		_, size := utf8.DecodeRuneInString(text[end:])
		end += size
	}
	s := strings.Join(strings.Fields(text[start:end]), " ")
	if start > 0 {
		s = "…" + s
	}
	if end < len(text) {
		s += "…"
	}
	return s
}
