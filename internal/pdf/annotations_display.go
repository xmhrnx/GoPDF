package pdf

import (
	"context"
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"

	"gopdf/internal/geom"
)

// Annotation flags (PDF 32000 12.5.3) that keep an annotation off screen.
const (
	annotHidden = 1 << 1
	annotNoView = 1 << 5
)

// displayCopy returns data with every annotation's normal appearance
// (signature fields, filled form fields, comments) drawn into its page's
// content, for renderers that only draw page content. It returns data
// itself when no page has an annotation to draw. The copy is for display
// only; it is never saved.
func displayCopy(data []byte) ([]byte, error) {
	ctx, err := readContext(data)
	if err != nil {
		return nil, err
	}
	drawn := false
	for n := 1; n <= ctx.PageCount; n++ {
		ok, err := drawAnnotations(ctx, n)
		if err != nil {
			return nil, fmt.Errorf("page %d: %w", n, err)
		}
		drawn = drawn || ok
	}
	if !drawn {
		return data, nil
	}
	return writeContext(ctx)
}

// drawAnnotations appends the appearances of page n's annotations to its
// content, reporting whether there were any.
func drawAnnotations(ctx *model.Context, n int) (bool, error) {
	pageDict, _, inh, err := ctx.PageDict(context.Background(), n, false)
	if err != nil {
		return false, err
	}
	annots, err := ctx.DereferenceArray(pageDict["Annots"])
	if err != nil || len(annots) == 0 {
		return false, err
	}
	res := copyDict(inh.Resources)
	xobjects, err := subDict(ctx, res, "XObject")
	if err != nil {
		return false, err
	}
	var cb contentBuilder
	for i, obj := range annots {
		form, m, ok := appearance(ctx, obj)
		if !ok {
			continue
		}
		name := fmt.Sprintf("GoPDFAnnot%d", i)
		xobjects[name] = form
		cb.op("q %s %s %s %s %s %s cm /%s Do Q", num(m[0]), num(m[1]), num(m[2]), num(m[3]), num(m[4]), num(m[5]), name)
	}
	if cb.b.Len() == 0 {
		return false, nil
	}
	pageDict["Resources"] = res
	return true, wrapContents(ctx, pageDict, cb.bytes())
}

// appearance returns the normal appearance of annotation obj and the
// matrix that places it on the page (PDF 32000 12.5.5), or false when
// it has none to show.
func appearance(ctx *model.Context, obj types.Object) (types.Object, [6]float64, bool) {
	var m [6]float64
	annot, err := ctx.DereferenceDict(obj)
	if err != nil || annot == nil {
		return nil, m, false
	}
	if f, ok := annot["F"].(types.Integer); ok && int(f)&(annotHidden|annotNoView) != 0 {
		return nil, m, false
	}
	if name, ok := annot["Subtype"].(types.Name); ok && name == "Popup" {
		return nil, m, false
	}
	ap, err := ctx.DereferenceDict(annot["AP"])
	if err != nil || ap == nil {
		return nil, m, false
	}
	formRef := ap["N"]
	normal, err := ctx.Dereference(formRef)
	if err != nil {
		return nil, m, false
	}
	if states, ok := normal.(types.Dict); ok {
		// A checkbox or radio button: the appearance of its current state.
		state, _ := annot["AS"].(types.Name)
		formRef = states[string(state)]
		if normal, err = ctx.Dereference(formRef); err != nil {
			return nil, m, false
		}
	}
	sd, ok := normal.(types.StreamDict)
	if !ok {
		return nil, m, false
	}
	rect, ok := numbers(ctx, annot["Rect"], 4)
	bbox, ok2 := numbers(ctx, sd.Dict["BBox"], 4)
	if !ok || !ok2 {
		return nil, m, false
	}
	formMatrix := [6]float64{1, 0, 0, 1, 0, 0}
	if fm, ok := numbers(ctx, sd.Dict["Matrix"], 6); ok {
		copy(formMatrix[:], fm)
	}
	// The form's bounding box, transformed by its matrix, is fitted to
	// the annotation's rectangle.
	box := mapRect(formMatrix, geom.RectFromPoints(geom.Point{X: bbox[0], Y: bbox[1]}, geom.Point{X: bbox[2], Y: bbox[3]}))
	target := geom.RectFromPoints(geom.Point{X: rect[0], Y: rect[1]}, geom.Point{X: rect[2], Y: rect[3]})
	if box.Dx() == 0 || box.Dy() == 0 || target.Dx() == 0 || target.Dy() == 0 {
		return nil, m, false
	}
	sx, sy := target.Dx()/box.Dx(), target.Dy()/box.Dy()
	fit := [6]float64{sx, 0, 0, sy, target.Min.X - box.Min.X*sx, target.Min.Y - box.Min.Y*sy}
	if _, indirect := formRef.(types.IndirectRef); !indirect {
		return nil, m, false
	}
	return formRef, multiply(formMatrix, fit), true
}

// multiply returns the matrix that applies a, then b.
func multiply(a, b [6]float64) [6]float64 {
	return [6]float64{
		a[0]*b[0] + a[1]*b[2], a[0]*b[1] + a[1]*b[3],
		a[2]*b[0] + a[3]*b[2], a[2]*b[1] + a[3]*b[3],
		a[4]*b[0] + a[5]*b[2] + b[4], a[4]*b[1] + a[5]*b[3] + b[5],
	}
}

// numbers reads an array of n numbers.
func numbers(ctx *model.Context, obj types.Object, n int) ([]float64, bool) {
	arr, err := ctx.DereferenceArray(obj)
	if err != nil || len(arr) != n {
		return nil, false
	}
	out := make([]float64, n)
	for i, o := range arr {
		v, err := ctx.DereferenceNumber(o)
		if err != nil {
			return nil, false
		}
		out[i] = v
	}
	return out, true
}
