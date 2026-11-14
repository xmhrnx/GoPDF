package signature

import "slices"

// Without returns list minus the signature with id.
func Without(list []Signature, id string) []Signature {
	return slices.DeleteFunc(slices.Clone(list), func(sig Signature) bool { return sig.ID == id })
}
