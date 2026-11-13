package pdf

import (
	"context"
	"io"
	"strconv"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// selection turns zero-based page indexes into pdfcpu's one-based
// selection strings.
func selection(pages []int) []string {
	s := make([]string, len(pages))
	for i, p := range pages {
		s[i] = strconv.Itoa(p + 1)
	}
	return s
}

// Rotate turns pages by deg degrees clockwise (a multiple of 90).
func Rotate(data []byte, pages []int, deg int) ([]byte, error) {
	return transform(data, func(c context.Context, rs io.ReadSeeker, w io.Writer, conf *model.Configuration) error {
		return api.Rotate(c, rs, w, deg, selection(pages), conf)
	})
}

// Arrange builds a document from the pages listed in order, by zero-based
// index. Pages may repeat (duplicate) or be left out (delete).
func Arrange(data []byte, order []int) ([]byte, error) {
	return transform(data, func(c context.Context, rs io.ReadSeeker, w io.Writer, conf *model.Configuration) error {
		return api.Collect(c, rs, w, selection(order), conf)
	})
}

// InsertBlank adds an empty page of size w x h points after page index
// after; after -1 inserts at the start.
func InsertBlank(data []byte, after int, w, h float64) ([]byte, error) {
	before, at := false, after
	if after < 0 {
		before, at = true, 0
	}
	pageConf := &pdfcpu.PageConfiguration{PageDim: &types.Dim{Width: w, Height: h}}
	return transform(data, func(c context.Context, rs io.ReadSeeker, w io.Writer, conf *model.Configuration) error {
		return api.InsertPages(c, rs, w, selection([]int{at}), before, pageConf, conf)
	})
}
