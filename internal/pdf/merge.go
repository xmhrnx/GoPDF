package pdf

import (
	"bytes"
	"context"
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// Merge appends the pages of each document in others to data, in order.
func Merge(data []byte, others ...[]byte) ([]byte, error) {
	readers := make([]io.ReadSeeker, 0, len(others)+1)
	readers = append(readers, bytes.NewReader(data))
	for _, o := range others {
		readers = append(readers, bytes.NewReader(o))
	}
	var out bytes.Buffer
	if err := api.MergeRaw(context.Background(), readers, &out, false, newConf()); err != nil {
		return nil, friendly(err)
	}
	return out.Bytes(), nil
}
