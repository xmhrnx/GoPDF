package studio

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/xmhrnx/giode"
)

func isPDF(path string) bool { return strings.EqualFold(filepath.Ext(path), ".pdf") }

func isCanceled(err error) bool { return errors.Is(err, giode.ErrCanceled) }

// baseName returns the file name without its extension.
func baseName(name string) string { return strings.TrimSuffix(name, filepath.Ext(name)) }

// writeFileAtomic writes data next to path and renames it into place, so
// a failed save never leaves a half-written document.
func writeFileAtomic(path string, data []byte) error {
	mode := os.FileMode(0o644)
	if info, err := os.Stat(path); err == nil {
		mode = info.Mode().Perm()
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".gopdf-*")
	if err != nil {
		return err
	}
	if err := tmp.Chmod(mode); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// writeAndClose writes data to w and reports where it went, when the
// writer is a file.
func writeAndClose(w io.WriteCloser, data []byte) (path string, err error) {
	if f, ok := w.(interface{ Name() string }); ok {
		path = f.Name()
	}
	if _, err := w.Write(data); err != nil {
		w.Close()
		return "", err
	}
	return path, w.Close()
}
