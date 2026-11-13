package pdf

import (
	"context"
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// Encrypt protects data with AES-256. The user password opens the file;
// the owner password (the user password when empty) unlocks editing.
func Encrypt(data []byte, userPW, ownerPW string) ([]byte, error) {
	if ownerPW == "" {
		ownerPW = userPW
	}
	return transform(data, func(c context.Context, rs io.ReadSeeker, w io.Writer, conf *model.Configuration) error {
		conf.UserPW, conf.OwnerPW = userPW, ownerPW
		conf.EncryptUsingAES = true
		conf.EncryptKeyLength = 256
		return api.Encrypt(c, rs, w, conf)
	})
}
