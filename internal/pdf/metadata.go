package pdf

import (
	"context"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// Metadata is the document information shown and edited in Properties.
type Metadata struct {
	Title    string
	Author   string
	Subject  string
	Keywords string
	Creator  string
	Producer string
}

func readMetadata(ctx *model.Context) Metadata {
	return Metadata{
		Title:    ctx.Title,
		Author:   ctx.Author,
		Subject:  ctx.Subject,
		Keywords: ctx.Keywords,
		Creator:  ctx.Creator,
		Producer: ctx.Producer,
	}
}

// SetMetadata writes the editable fields of m into the document.
func SetMetadata(data []byte, m Metadata) ([]byte, error) {
	ctx, err := readContext(data)
	if err != nil {
		return nil, err
	}
	props := map[string]string{
		"Title":    strings.TrimSpace(m.Title),
		"Author":   strings.TrimSpace(m.Author),
		"Subject":  strings.TrimSpace(m.Subject),
		"Keywords": strings.TrimSpace(m.Keywords),
	}
	if err := pdfcpu.PropertiesAdd(context.Background(), ctx, props); err != nil {
		return nil, err
	}
	return writeContext(ctx)
}
