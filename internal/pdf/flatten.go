package pdf

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"

	"gopdf/internal/annot"
)

// Flatten draws list into the page content of data, so that every viewer
// shows the markup. Annotations refer to pages by zero-based index.
func Flatten(data []byte, list []annot.Annotation) ([]byte, error) {
	if len(list) == 0 {
		return data, nil
	}
	ctx, err := readContext(data)
	if err != nil {
		return nil, err
	}
	byPage := map[int][]annot.Annotation{}
	for _, a := range list {
		byPage[a.Page] = append(byPage[a.Page], a)
	}
	fontRef, err := ctx.IndRefForNewObject(helveticaFont())
	if err != nil {
		return nil, err
	}
	for page, annots := range byPage {
		if err := flattenPage(ctx, page+1, annots, *fontRef); err != nil {
			return nil, fmt.Errorf("page %d: %w", page+1, err)
		}
	}
	return writeContext(ctx)
}

func helveticaFont() types.Dict {
	return types.Dict{
		"Type":     types.Name("Font"),
		"Subtype":  types.Name("Type1"),
		"BaseFont": types.Name("Helvetica"),
		"Encoding": types.Name("WinAnsiEncoding"),
	}
}

func flattenPage(ctx *model.Context, pageNr int, list []annot.Annotation, font types.IndirectRef) error {
	pageDict, _, inh, err := ctx.PageDict(context.Background(), pageNr, false)
	if err != nil {
		return err
	}
	res := copyDict(inh.Resources)
	fonts, err := subDict(ctx, res, "Font")
	if err != nil {
		return err
	}
	fonts[fontResName] = font
	states, err := subDict(ctx, res, "ExtGState")
	if err != nil {
		return err
	}
	gs := func(opacity float64) string {
		name := gsResPrefix + strconv.Itoa(int(opacity*100))
		states[name] = types.Dict{"Type": types.Name("ExtGState"), "CA": types.Float(opacity), "ca": types.Float(opacity)}
		return name
	}
	content := writeAnnotations(list, pageInfoFrom(inh).displayMatrix(), gs)
	pageDict["Resources"] = res
	return wrapContents(ctx, pageDict, content)
}

// copyDict returns a shallow copy of d, so adding entries never touches
// dictionaries shared with other pages.
func copyDict(d types.Dict) types.Dict {
	out := make(types.Dict, len(d))
	for k, v := range d {
		out[k] = v
	}
	return out
}

func subDict(ctx *model.Context, res types.Dict, key string) (types.Dict, error) {
	out := types.Dict{}
	if obj, ok := res[key]; ok {
		d, err := ctx.DereferenceDict(obj)
		if err != nil {
			return nil, err
		}
		out = copyDict(d)
	}
	res[key] = out
	return out, nil
}

// wrapContents isolates the page's existing content in q/Q so its
// graphics state can't leak into the markup appended after it.
func wrapContents(ctx *model.Context, pageDict types.Dict, content []byte) error {
	push, err := ctx.StreamDictIndRef([]byte("q\n"))
	if err != nil {
		return err
	}
	pop, err := ctx.StreamDictIndRef([]byte("\nQ\n"))
	if err != nil {
		return err
	}
	ours, err := ctx.StreamDictIndRef(content)
	if err != nil {
		return err
	}
	parts := types.Array{*push}
	if obj, ok := pageDict["Contents"]; ok {
		existing, err := contentParts(ctx, obj)
		if err != nil {
			return err
		}
		parts = append(parts, existing...)
	}
	pageDict["Contents"] = append(parts, *pop, *ours)
	return nil
}

func contentParts(ctx *model.Context, obj types.Object) (types.Array, error) {
	if ref, ok := obj.(types.IndirectRef); ok {
		d, err := ctx.Dereference(ref)
		if err != nil {
			return nil, err
		}
		if arr, ok := d.(types.Array); ok {
			return arr, nil
		}
		return types.Array{ref}, nil
	}
	if arr, ok := obj.(types.Array); ok {
		return arr, nil
	}
	return nil, fmt.Errorf("unexpected page contents %T", obj)
}
