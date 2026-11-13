// Package pdf reads, edits and renders PDF files. Documents are immutable
// byte slices: every edit returns new bytes, which keeps undo cheap and
// lets rendering run on another goroutine without locks on the model.
package pdf

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// ErrPassword means the document is encrypted and the password given is
// missing or wrong.
var ErrPassword = errors.New("this document is protected by a password")

func init() {
	// Keep pdfcpu from creating a configuration folder in the user's home.
	api.DisableConfigDir()
}

func newConf() *model.Configuration {
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed
	return conf
}

// transform runs a pdfcpu stream operation over data and returns its output.
func transform(data []byte, op func(c context.Context, rs io.ReadSeeker, w io.Writer, conf *model.Configuration) error) ([]byte, error) {
	var out bytes.Buffer
	if err := op(context.Background(), bytes.NewReader(data), &out, newConf()); err != nil {
		return nil, friendly(err)
	}
	return out.Bytes(), nil
}

// readContext parses data into a pdfcpu context for low-level edits.
func readContext(data []byte) (*model.Context, error) {
	ctx, err := api.ReadAndValidate(context.Background(), bytes.NewReader(data), newConf())
	if err != nil {
		return nil, friendly(err)
	}
	return ctx, nil
}

// writeContext serializes ctx back into PDF bytes.
func writeContext(ctx *model.Context) ([]byte, error) {
	var out bytes.Buffer
	if err := api.WriteContext(context.Background(), ctx, &out); err != nil {
		return nil, friendly(err)
	}
	return out.Bytes(), nil
}

func friendly(err error) error {
	if errors.Is(err, pdfcpu.ErrWrongPassword) || errors.Is(err, pdfcpu.ErrOwnerPasswordRequired) {
		return ErrPassword
	}
	return err
}
