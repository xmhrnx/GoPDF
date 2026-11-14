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
	// Display is Data with annotation appearances drawn into the pages,
	// for rendering; it is Data itself when there are none.
	Display []byte
	// Signatures are the document's signature fields, checked.
	Signatures []SignatureField
	// Encrypted reports that the file was password protected; Data is
	// the decrypted copy.
	Encrypted bool
}

// Open parses data, decrypting it with password when it is protected.
// ErrPassword reports a missing or wrong password.
//
// Decrypting rewrites the file, so the signatures are checked on the
// original bytes.
func Open(data []byte, password string) (*Document, error) {
	plain, err := decrypt(data, password)
	if err != nil {
		return nil, err
	}
	doc, err := Load(plain)
	if err != nil || len(plain) == len(data) && string(plain) == string(data) {
		return doc, err
	}
	doc.Encrypted = true
	if sigs, err := ReadSignatures(data, doc.Pages, password); err == nil {
		doc.Signatures = sigs
	}
	return doc, nil
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
	display, err := displayCopy(data)
	if err != nil {
		display = data // still viewable, without annotation appearances
	}
	sigs, err := ReadSignatures(data, pages, "")
	if err != nil {
		sigs = nil // unreadable fields are shown as no signatures
	}
	return &Document{Data: data, Pages: pages, Meta: readMetadata(ctx), Display: display, Signatures: sigs}, nil
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
