package editor

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseRanges reads a page selection such as "1-3, 5, 8-" (one-based, as
// typed by the user) for a document of n pages and returns zero-based
// indexes in the order given.
func ParseRanges(spec string, n int) ([]int, error) {
	var pages []int
	for part := range strings.SplitSeq(spec, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		from, to, err := parseRange(part, n)
		if err != nil {
			return nil, err
		}
		for p := from; p <= to; p++ {
			pages = append(pages, p-1)
		}
	}
	if len(pages) == 0 {
		return nil, fmt.Errorf("enter pages such as 1-3, 5")
	}
	return pages, nil
}

func parseRange(part string, n int) (from, to int, err error) {
	lo, hi, isRange := strings.Cut(part, "-")
	if from, err = pageNumber(lo, 1, n); err != nil {
		return 0, 0, err
	}
	to = from
	if isRange {
		if to, err = pageNumber(hi, n, n); err != nil {
			return 0, 0, err
		}
	}
	if from > to {
		return 0, 0, fmt.Errorf("%q runs backwards", part)
	}
	return from, to, nil
}

func pageNumber(s string, empty, n int) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return empty, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 1 || v > n {
		return 0, fmt.Errorf("%q is not a page between 1 and %d", s, n)
	}
	return v, nil
}
