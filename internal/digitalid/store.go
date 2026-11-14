package digitalid

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"time"

	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

// Errors reported to the user.
var (
	ErrPassword  = errors.New("the password is wrong")
	ErrShortPass = errors.New("use a password of at least 6 characters")
	ErrNotRSA    = errors.New("this digital ID's key type isn't supported; use an RSA ID")
)

// minPassword is the shortest password a new ID accepts.
const minPassword = 6

// Store keeps IDs as PKCS#12 files in a folder, with an index of their
// public details.
type Store struct {
	dir string
}

// Signer is an unlocked ID: what signing needs.
type Signer struct {
	Info  Info
	Key   crypto.Signer
	Cert  *x509.Certificate
	Chain []*x509.Certificate // issuers, nearest first; empty for self-signed
}

// NewStore returns the store in the configuration folder of app.
func NewStore(app string) (*Store, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	return &Store{dir: filepath.Join(dir, app, "digital-ids")}, nil
}

// List returns the stored IDs.
func (s *Store) List() ([]Info, error) {
	b, err := os.ReadFile(s.index())
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var list []Info
	return list, json.Unmarshal(b, &list)
}

// Create makes a self-signed ID for id, protected by password.
func (s *Store) Create(id Identity, password string) (Info, error) {
	if len(password) < minPassword {
		return Info{}, ErrShortPass
	}
	u, err := newSelfSigned(id, time.Now())
	if err != nil {
		return Info{}, err
	}
	p12, err := pkcs12.Modern.Encode(u.key, u.cert, nil, password)
	if err != nil {
		return Info{}, err
	}
	return s.add(p12, u.cert)
}

// Import stores a PKCS#12 (.p12, .pfx) ID, after checking password
// opens it.
func (s *Store) Import(p12 []byte, password string) (Info, error) {
	u, err := decode(p12, password)
	if err != nil {
		return Info{}, err
	}
	return s.add(p12, u.Cert)
}

// Unlock opens the ID with fingerprint for signing.
func (s *Store) Unlock(fingerprint, password string) (*Signer, error) {
	p12, err := os.ReadFile(s.file(fingerprint))
	if err != nil {
		return nil, err
	}
	return decode(p12, password)
}

// Delete removes the ID with fingerprint.
func (s *Store) Delete(fingerprint string) error {
	list, err := s.List()
	if err != nil {
		return err
	}
	list = slices.DeleteFunc(list, func(i Info) bool { return i.Fingerprint == fingerprint })
	if err := os.Remove(s.file(fingerprint)); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return s.writeIndex(list)
}

func (s *Store) add(p12 []byte, cert *x509.Certificate) (Info, error) {
	info := infoOf(cert)
	if err := os.MkdirAll(s.dir, 0o700); err != nil {
		return Info{}, err
	}
	if err := os.WriteFile(s.file(info.Fingerprint), p12, 0o600); err != nil {
		return Info{}, err
	}
	list, err := s.List()
	if err != nil {
		return Info{}, err
	}
	list = slices.DeleteFunc(list, func(i Info) bool { return i.Fingerprint == info.Fingerprint })
	return info, s.writeIndex(append(list, info))
}

func decode(p12 []byte, password string) (*Signer, error) {
	key, cert, chain, err := pkcs12.DecodeChain(p12, password)
	if errors.Is(err, pkcs12.ErrIncorrectPassword) {
		return nil, ErrPassword
	}
	if err != nil {
		return nil, fmt.Errorf("can't read the digital ID: %w", err)
	}
	signer, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrNotRSA
	}
	return &Signer{Info: infoOf(cert), Key: signer, Cert: cert, Chain: chain}, nil
}

func (s *Store) writeIndex(list []Info) error {
	b, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.index(), b, 0o600)
}

func (s *Store) index() string { return filepath.Join(s.dir, "index.json") }

func (s *Store) file(fingerprint string) string {
	return filepath.Join(s.dir, filepath.Base(fingerprint)+".p12")
}

// NewStoreAt returns a store in dir.
func NewStoreAt(dir string) *Store { return &Store{dir: dir} }
