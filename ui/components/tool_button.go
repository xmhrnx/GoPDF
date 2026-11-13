// Package components holds small reusable pieces of the GoPDF
// interface: buttons, swatches, sections and dialog frames.
package components

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/tooltip"
	"github.com/xmhrnx/giode/styles/properties"

	"gopdf/ui/theme"
)

// Square sizes of tool buttons, in dp.
const (
	toolButtonSize = 32
	compactSize    = 22
)

// ToolButton is a square icon button with a tooltip and an optional
// active (toggled) state. Create it once; call View every frame.
type ToolButton struct {
	icon    string
	size    int
	press   *giode.PressableElement
	tip     *tooltip.Tooltip
	active  func() bool
	enabled func() bool
}

// NewToolButton returns a button showing icon, explaining itself with tip.
func NewToolButton(icon, tip string, onClick func()) *ToolButton {
	b := &ToolButton{icon: icon, size: toolButtonSize, active: never, enabled: always}
	b.press = giode.Pressable(b.content).Label(tip).OnClick(func(giode.ClickEvent) { onClick() })
	b.tip = giode.Tooltip(b.press, tip)
	return b
}

func never() bool  { return false }
func always() bool { return true }

// Active makes the button show as toggled on while fn reports true.
func (b *ToolButton) Active(fn func() bool) *ToolButton {
	b.active = fn
	return b
}

// Enabled disables the button while fn reports false.
func (b *ToolButton) Enabled(fn func() bool) *ToolButton {
	b.enabled = fn
	return b
}

// Compact shrinks the button for dense bars (a status bar).
func (b *ToolButton) Compact() *ToolButton {
	b.size = compactSize
	return b
}

// View returns the button for this frame.
func (b *ToolButton) View() giode.Element {
	b.press.Disabled(!b.enabled())
	return b.tip
}

func (b *ToolButton) content(s giode.PressState) giode.Element {
	bg, fg := "transparent", theme.TextColor
	switch {
	case s.Disabled:
		fg = theme.Faint
	case b.active():
		bg, fg = theme.SelectedStrong, theme.Primary
	case s.Pressed:
		bg = theme.Pressed
	case s.Hovered:
		bg = theme.Hover
	}
	icon := theme.IconSize
	if b.size == compactSize {
		icon = theme.SmallIcon
	}
	return giode.Box(giode.Styles{
		Width: properties.Width(b.size), Height: properties.Height(b.size), BorderRadius: theme.Radius,
		Background: bg, Align: giode.AlignCenter, Justify: giode.JustifyCenter,
	}, theme.Icon(b.icon, icon, fg))
}

// SetIcon replaces the icon (a theme toggle showing sun or moon).
func (b *ToolButton) SetIcon(icon string) { b.icon = icon }
