package digitalid

import (
	"errors"
	"testing"
)

func TestCreateUnlockImportDelete(t *testing.T) {
	s := &Store{dir: t.TempDir()}
	info, err := s.Create(Identity{Name: "Jane Doe", Email: "jane@example.com", Organization: "Acme"}, "secret-pass")
	if err != nil {
		t.Fatal(err)
	}
	if !info.SelfSigned || info.Identity.Email != "jane@example.com" || info.Identity.Organization != "Acme" {
		t.Fatalf("info = %+v", info)
	}
	if _, err := s.Unlock(info.Fingerprint, "wrong"); !errors.Is(err, ErrPassword) {
		t.Fatalf("wrong password: %v", err)
	}
	signer, err := s.Unlock(info.Fingerprint, "secret-pass")
	if err != nil {
		t.Fatal(err)
	}
	if signer.Cert.Subject.CommonName != "Jane Doe" {
		t.Fatalf("cert subject = %v", signer.Cert.Subject)
	}
	list, _ := s.List()
	if len(list) != 1 {
		t.Fatalf("list = %v", list)
	}
	if err := s.Delete(info.Fingerprint); err != nil {
		t.Fatal(err)
	}
	if list, _ := s.List(); len(list) != 0 {
		t.Fatalf("after delete = %v", list)
	}
}

func TestCreateRejectsWeakInput(t *testing.T) {
	s := &Store{dir: t.TempDir()}
	if _, err := s.Create(Identity{Name: "Jane"}, "123"); !errors.Is(err, ErrShortPass) {
		t.Fatalf("short password: %v", err)
	}
	if _, err := s.Create(Identity{Name: "  "}, "long-enough"); !errors.Is(err, ErrNoName) {
		t.Fatalf("no name: %v", err)
	}
}
