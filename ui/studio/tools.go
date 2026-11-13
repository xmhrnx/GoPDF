package studio

import "gopdf/internal/annot"

// Tool is what a press on a page does: select markup, or draw a kind.
type Tool int

// ToolSelect selects and moves markup; every other Tool draws the
// annot.Kind of the same value.
const ToolSelect Tool = -1

// Kind returns the annotation kind a drawing tool creates.
func (t Tool) Kind() annot.Kind { return annot.Kind(t) }

// toolDef describes a tool button.
type toolDef struct {
	tool Tool
	icon string
	name string
	key  string
}

// tools are the toolbar's tools, in order.
var tools = []toolDef{
	{ToolSelect, "mouse-pointer-2", "Select", "V"},
	{Tool(annot.Text), "type", "Text", "T"},
	{Tool(annot.Pen), "pen-line", "Pen", "P"},
	{Tool(annot.Marker), "highlighter", "Marker", "M"},
	{Tool(annot.Highlight), "scan", "Highlight area", "H"},
	{Tool(annot.Rect), "square", "Rectangle", "R"},
	{Tool(annot.Ellipse), "circle", "Ellipse", "E"},
	{Tool(annot.Line), "minus", "Line", "L"},
	{Tool(annot.Arrow), "move-up-right", "Arrow", "A"},
	{Tool(annot.Whiteout), "eraser", "Whiteout", "W"},
}
