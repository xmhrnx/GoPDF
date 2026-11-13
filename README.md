# GoPDF

A fast, native PDF editor built with [Giode](https://github.com/xmhrnx/giode) (Go on Gio, GPU
rendered, no web view). It runs in a custom frameless window with its own title bar, menus
and window buttons.

## Features

- **View:** continuous scrolling with lazily rendered pages, fit-width and zoom from 25 to
  400 %, page thumbnails, outline (bookmarks) and full-text search.
- **Annotate:** text boxes, pen, marker, area highlight, rectangle, ellipse, line, arrow and
  whiteout, with colors, stroke widths, fill and font sizes. Select markup to move it, restyle
  it, edit its text or delete it. On save the markup is drawn into the page content, so every
  PDF viewer shows it.
- **Organize:** rotate, duplicate, reorder, insert blank pages and delete pages. Markup
  moves with its page.
- **Combine:** merge PDFs, extract page ranges (`1-3, 5, 8-`) into a new file, and split a
  document into parts of N pages.
- **Stamp:** text watermarks (size, angle, opacity, color) and page numbers.
- **Protect:** open password-protected files, and encrypt with AES-256 on save.
- **Document:** edit title, author, subject and keywords. Export a page as PNG or all of
  the text as `.txt`.
- **Everything else:** undo and redo (100 steps), light and dark themes,
  drag-and-drop opening, recent files, a warning before unsaved changes are lost, and
  keyboard shortcuts for everything (Help › Keyboard Shortcuts).

## Build and run

You need Go 1.26+, a C compiler (MuPDF is linked through cgo) and the Giode CLI:

```bash
go install github.com/xmhrnx/giode-cli/cmd/giode@latest
```

```bash
giode doctor
```

```bash
giode run
```

```bash
giode build linux
```

You can pass a file on the command line (`gopdf report.pdf`). Building for Windows or
macOS from Linux needs a cross C compiler; see `giode.toml` and `giode doctor`.

```bash
go test ./...
```

## Layout

```
main.go              wiring only: fonts, icons, theme, window config
assets/              app icon, Lucide icons (ISC), Inter font (OFL)
internal/            the app's logic, with no UI imports and its own tests
  geom/              points, rectangles, rotation
  annot/             the markup model: kinds, styles, hit testing, rotation
  pdf/               open/decrypt, page operations, merge, watermark, encryption,
                     flattening markup into content streams, MuPDF rendering, search
  editor/            the session: snapshots, undo/redo, page ops that keep markup in step
  history/           a generic undo stack
  pagecache/         background page rendering with a bounded LRU cache
  recent/            the recent files list
ui/
  theme/             palettes, tokens, text styles, fonts, icons
  components/        reusable widgets (tool button, swatch, dialog frame, ...)
  platform/          title bar style per OS (build tags)
  studio/            the app: one view per file, plus the actions they run
```

Each view or component is in its own file. Documents are immutable byte slices: every edit
produces new bytes on a background goroutine and is committed on the UI goroutine, which
keeps undo cheap and the interface responsive.

## How it stays light and fast

- Only the pages in view are laid out and rendered (`LazyList`), on a background goroutine.
- Rendered pages are kept in a small LRU cache (10 pages and 240 thumbnails). Evicted images
  are released from the GPU cache, and while a sharper render is on its way the closest
  cached resolution is shown.
- MuPDF's resource cache is capped at 48 MB and the Go garbage collector runs a little more
  often than its default.
- The binary is stripped. A frame takes about 1–2 ms to build.

## Credits

Rendering: [MuPDF](https://mupdf.com) through [go-fitz](https://github.com/gen2brain/go-fitz)
(AGPL). Editing: [pdfcpu](https://github.com/pdfcpu/pdfcpu) (Apache 2.0). Icons:
[Lucide](https://lucide.dev) (ISC). Typeface: [Inter](https://rsms.me/inter/) (OFL).
