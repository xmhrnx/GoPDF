// Package platform holds the few choices that differ between operating
// systems, split by build tag.
package platform

import "github.com/xmhrnx/giode"

// TitleBarVariant draws macOS traffic lights on the left.
var TitleBarVariant = giode.TitleBarMac
