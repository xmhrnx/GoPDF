package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"

	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// shortcutsDialog lists every keyboard shortcut.
type shortcutsDialog struct {
	open   bool
	close  *giode.ButtonElement
	dialog *dialog.Dialog
}

func newShortcutsDialog(*Studio) *shortcutsDialog {
	d := &shortcutsDialog{}
	d.close = components.NewPrimaryButton("Done", func() { d.open = false })
	d.dialog = components.NewDialog(&d.open, d.content)
	return d
}

func (d *shortcutsDialog) content() giode.Element {
	var cols []giode.Element
	for _, g := range shortcutGroups {
		rows := []giode.Element{giode.Text(g.title, theme.Overline)}
		for _, it := range g.items {
			rows = append(rows, giode.Row(giode.Styles{Width: -1, Justify: giode.JustifySpaceBetween, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Gap: 12},
				giode.Text(it.what, theme.Body), components.Kbd(it.spec)))
		}
		cols = append(cols, giode.Stack(giode.Styles{Width: 196, Gap: 7}, rows...))
	}
	body := giode.Row(giode.Styles{Width: -1, Gap: 24, Wrap: giode.WrapWrap}, cols...)
	return components.DialogFrame("keyboard", "Keyboard shortcuts", "Work faster without the mouse.", body, d.close)
}
