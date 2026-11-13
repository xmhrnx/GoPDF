//go:build !darwin

// Package platform holds the few choices that differ between operating
// systems, split by build tag.
package platform

import "github.com/xmhrnx/giode"

// TitleBarVariant draws Windows-style caption buttons on the right.
var TitleBarVariant = giode.TitleBarWindows
