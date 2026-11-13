package pdf

import (
	"bytes"
	"context"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// Document is a parsed, unencrypted PDF. It is immutable: editing
// functions take Data and return new bytes for a new Document.
type Document struct {
	Data  []byte
	Pages []PageInfo
	Meta  Metadata
}

// Open parses data, decrypting it with password when it is protected.
// ErrPassword reports a missing or wrong password.
func Open(data []byte, password string) (*Document, error) {
	plain, err := decrypt(data, password)
	if err != nil {
		return nil, err
	}
	return Load(plain)
}

// Load parses unencrypted PDF bytes.
func Load(data []byte) (*Document, error) {
	ctx, err := readContext(data)
	if err != nil {
		return nil, err
	}
	pages, err := readPageInfos(ctx)
	if err != nil {
		return nil, err
	}
	return &Document{Data: data, Pages: pages, Meta: readMetadata(ctx)}, nil
}

// PageCount returns the number of pages.
func (d *Document) PageCount() int { return len(d.Pages) }

func decrypt(data []byte, password string) ([]byte, error) {
	conf := newConf()
	conf.UserPW, conf.OwnerPW = password, password
	ctx, err := api.ReadAndValidate(context.Background(), bytes.NewReader(data), conf)
	if err != nil {
		return nil, friendly(err)
	}
	if ctx.E == nil {
		return data, nil
	}
	var out bytes.Buffer
	if err := api.Decrypt(context.Background(), bytes.NewReader(data), &out, conf); err != nil {
		return nil, friendly(err)
	}
	return out.Bytes(), nil
}
