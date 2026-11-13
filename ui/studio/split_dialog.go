package studio

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/input"

	"gopdf/internal/editor"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// splitDialog cuts the document into files of N pages in a folder.
type splitDialog struct {
	s       *Studio
	open    bool
	problem string
	every   *input.Input
	split   *giode.ButtonElement
	cancel  *giode.ButtonElement
	dialog  *dialog.Dialog
}

func newSplitDialog(s *Studio) *splitDialog {
	d := &splitDialog{s: s}
	d.every = components.NewTextInput("1").Keyboard(giode.KeyboardNumeric).Submit(func(string) { d.submit() })
	d.split = components.NewPrimaryButton("Choose Folder…", d.submit)
	d.cancel = components.NewSecondaryButton("Cancel", func() { d.open = false })
	d.dialog = components.NewDialog(&d.open, d.content)
	return d
}

func (d *splitDialog) show() {
	d.every.SetText("1")
	d.problem, d.open = "", true
	d.every.Focus()
}

func (d *splitDialog) submit() {
	n, err := strconv.Atoi(strings.TrimSpace(d.every.Text()))
	if err != nil || n < 1 {
		d.problem = "Enter a whole number of pages, 1 or more."
		return
	}
	d.open = false
	snap, base := d.s.doc.Snapshot(), baseName(d.s.doc.Name)
	var written int
	d.s.runBackground("Splitting", func() error {
		dir, err := d.s.win.PickFolder(context.Background(), giode.FolderOptions{Title: "Where should the parts go?"})
		if err != nil {
			return err
		}
		parts, err := editor.Split(snap, n)
		if err != nil {
			return err
		}
		for i, part := range parts {
			name := filepath.Join(dir, fmt.Sprintf("%s-part-%d.pdf", base, i+1))
			if err := os.WriteFile(name, part, 0o644); err != nil {
				return err
			}
			written++
		}
		return nil
	}, func() { d.s.toasts.Show(giode.ToastSuccess, fmt.Sprintf("Saved %d files", written)) })
}

func (d *splitDialog) content() giode.Element {
	body := []giode.Element{components.Field("Pages per file", d.every)}
	if d.problem != "" {
		body = append(body, giode.Text(d.problem, giode.Styles{FontSize: 12, Color: theme.Danger}))
	}
	return components.DialogFrame("scissors", "Split document",
		fmt.Sprintf("Saves the %d pages as separate PDFs in a folder you choose.", d.s.pageCount()),
		giode.Stack(giode.Styles{Width: -1, Gap: 8}, body...), d.cancel, d.split)
}
