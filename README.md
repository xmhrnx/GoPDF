# GoPDF

A fast, native PDF editor built with [Giode](https://github.com/xmhrnx/giode) (Go on Gio, GPU
rendered, no web view). It runs in a custom frameless window with its own title bar, menus
and window buttons.

Website: <https://xmhrnx.github.io/GoPDF/> (source in `site/`).

## Features

- **View:** continuous scrolling with lazily rendered pages, fit-width and zoom from 25 to
  400 %, page thumbnails, outline (bookmarks) and full-text search.
- **Annotate:** text boxes, pen, marker, area highlight, rectangle, ellipse, line, arrow and
  whiteout, with colors, stroke widths, fill and font sizes. Select markup to move it (with
  the mouse, or the arrow keys: 1 pt, 10 pt with Shift), restyle it, edit its text, copy,
  cut, paste, duplicate or delete it. On save the markup is drawn into the page content, so
  every PDF viewer shows it.
- **Sign:** draw a signature, or type your name and pick one of six script typefaces;
  keep it and place it in any document with the Signature tool. Typed signatures are saved
  as vector outlines, so they look the same in every PDF viewer without embedding a font.
  Several signatures can be saved; pick the one to use in the inspector.
- **Digital signatures:** turn a placed signature into a digital signature with a digital
  ID: create a self-signed one in GoPDF (name, email, organization, password) or import a
  `.p12`/`.pfx` from a certificate authority. On save the signature becomes a signature
  field sealed with the ID (PAdES, `ETSI.CAdES.detached`), recording the signer, time,
  reason and location; adding a signature to a signed file keeps the earlier ones valid.
  Opening a signed PDF shows a banner with the result, and clicking a signature shows who
  signed and whether the document changed since. Checks run offline.
- **Remembers:** the last color, stroke width, font size and signature size used for each
  kind of markup, between runs.
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

You need Go 1.27+, a C compiler (MuPDF is linked through cgo) and the Giode CLI:

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

You can pass a file on the command line (`gopdf report.pdf`).

### Windows

Install a MinGW-w64 gcc and make sure it is on `PATH`; MSYS2's MINGW64 and UCRT64
toolchains both work. Then run `giode build windows`.

The MuPDF library go-fitz ships for Windows was compiled against UCRT and calls
`__intrinsic_setjmpex`, which msvcrt-based toolchains don't define;
`internal/pdf/setjmp_windows.go` supplies it there by forwarding to msvcrt's `_setjmpex`.

If the link fails with `.rsrc merge failure: corrupt .rsrc section`, delete every
`*.syso` file in the project folder: the giode CLI writes the icon resource there, and an
extra copy linked in alongside it breaks the link.

Cross-compiling for Windows or macOS from Linux needs a matching cross C compiler; see
`giode.toml` and `giode doctor`.

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
  prefs/             the style remembered for each kind of markup
  signature/         saved signatures: normalizing strokes, typed names as outlines
  digitalid/         digital IDs: self-signed creation, PKCS#12 import, locked storage
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

Digital signatures: [pdfsign](https://github.com/digitorus/pdfsign) (BSD 2-Clause) and
[go-pkcs12](https://github.com/SSLMate/go-pkcs12) (BSD 3-Clause).
Rendering: [MuPDF](https://mupdf.com) through [go-fitz](https://github.com/gen2brain/go-fitz)
(AGPL). Editing: [pdfcpu](https://github.com/pdfcpu/pdfcpu) (Apache 2.0). Icons:
[Lucide](https://lucide.dev) (ISC). Typefaces: [Inter](https://rsms.me/inter/), and Great
Vibes, Allura, Alex Brush, Sacramento, Mrs Saint Delafield and Cedarville Cursive from
[Google Fonts](https://fonts.google.com) (all SIL OFL).
