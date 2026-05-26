// Copyright 2026 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tcell

import (
	"github.com/gdamore/tcell/v3/color"
)

// Style represents a complete text style, including both foreground color,
// background color, and additional attributes such as "bold" or "underline".
//
// Note that not all terminals can display all colors or attributes, and
// many might have specific incompatibilities between specific attributes
// and color combinations.
//
// To use Style, just declare a variable of its type.
type Style struct {
	fg      color.Color
	bg      color.Color
	ulColor color.Color
	attrs   AttrMask
	ulStyle UnderlineStyle
	url     *urlInfo
}

type urlInfo struct {
	url string
	id  string
}

// stripOSCControls removes control bytes that can terminate OSC payloads early.
func stripOSCControls(s string) string { _ = "STUB: not implemented"; return "" }

// stripOSCControlsIfNeeded returns the original string when it contains no
// control bytes and only allocates when stripping is required.
func stripOSCControlsIfNeeded(s string) string { _ = "STUB: not implemented"; return "" }

// StyleDefault represents a default style, based upon the context.
// It is the zero value.
var StyleDefault Style

// styleInvalid is just an arbitrary invalid style used internally.
var styleInvalid = Style{attrs: AttrInvalid}

// Foreground returns a new style based on s, with the foreground color set
// as requested.  ColorDefault can be used to select the global default.
func (s Style) Foreground(c color.Color) Style { _ = "STUB: not implemented"; return *new(Style) }

// Background returns a new style based on s, with the background color set
// as requested.  ColorDefault can be used to select the global default.
func (s Style) Background(c color.Color) Style { _ = "STUB: not implemented"; return *new(Style) }

func (s Style) setAttrs(attrs AttrMask, on bool) Style {
	_ = "STUB: not implemented"
	return *new(Style)
}

// Normal returns the style with all attributes disabled.
// Colors are preserved, as are hyperlinks.  (Underline color
// will also be preserved, but no underline is currently shown.
// Apart from color, the underline style is reset as well.)
func (s Style) Normal() Style { _ = "STUB: not implemented"; return *new(Style) }

// Bold returns a new style based on s, with the bold attribute set
// as requested.
func (s Style) Bold(on bool) Style { _ = "STUB: not implemented"; return *new(Style) }

// Blink returns a new style based on s, with the blink attribute set
// as requested.
func (s Style) Blink(on bool) Style { _ = "STUB: not implemented"; return *new(Style) }

// Dim returns a new style based on s, with the dim attribute set
// as requested.
func (s Style) Dim(on bool) Style { _ = "STUB: not implemented"; return *new(Style) }

// Italic returns a new style based on s, with the italic attribute set
// as requested.
func (s Style) Italic(on bool) Style { _ = "STUB: not implemented"; return *new(Style) }

// Reverse returns a new style based on s, with the reverse attribute set
// as requested.  (Reverse usually changes the foreground and background
// colors.)
func (s Style) Reverse(on bool) Style { _ = "STUB: not implemented"; return *new(Style) }

// StrikeThrough sets strike-through mode.
func (s Style) StrikeThrough(on bool) Style { _ = "STUB: not implemented"; return *new(Style) }

// Underline style.  Modern terminals have the option of rendering the
// underline using different styles, and even different colors.
type UnderlineStyle uint8

const (
	UnderlineStyleNone = UnderlineStyle(iota)
	UnderlineStyleSolid
	UnderlineStyleDouble
	UnderlineStyleCurly
	UnderlineStyleDotted
	UnderlineStyleDashed
)

// Underline returns a new style based on s, with the underline attribute set
// as requested.  The parameters can be:
//
// bool: on / off - enables just a simple underline
// UnderlineStyle: sets a specific style (should not coexist with the bool)
// Color: the color to use
func (s Style) Underline(params ...any) Style { _ = "STUB: not implemented"; return *new(Style) }

// GetForeground returns the foreground (text) color.
func (s Style) GetForeground() color.Color {
	_ = "STUB: not implemented"

	// GetBackground returns the background color.
	return *new(color.Color)
}

func (s Style) GetBackground() color.Color {
	_ = "STUB: not implemented"

	// GetUnderlineStyle returns the underline style for the style.
	return *new(color.Color)
}

func (s Style) GetUnderlineStyle() UnderlineStyle {
	_ = "STUB: not implemented"

	// GetUnderlineColor returns the underline color for the style.
	return *new(UnderlineStyle)
}

func (s Style) GetUnderlineColor() color.Color {
	_ = "STUB: not implemented"

	// Attributes returns a new style based on s, with its attributes set as
	// specified.
	//
	// Deprecated: Use direct functions instead.
	return *new(color.Color)
}

func (s Style) Attributes(attrs AttrMask) Style { _ = "STUB: not implemented"; return *new(Style) }

// GetAttributes gets the attributes for a style.
// Deprecated: Use individual properties instead.
func (s Style) GetAttributes() AttrMask {
	_ = "STUB: not implemented"

	// Url returns a style with the Url set.  If the provided Url is not empty,
	// and the terminal supports it, text will typically be marked up as a clickable
	// link to that Url.  If the Url is empty, then this mode is turned off.
	return *new(AttrMask)
}

func (s Style) Url(url string) Style { _ = "STUB: not implemented"; return *new(Style) }

// UrlId returns a style with the UrlId set. If the provided UrlId is not empty,
// any marked up Url with this style will be given the UrlId also. If the
// terminal supports it, any text with the same UrlId will be grouped as if it
// were one Url, even if it spans multiple lines.
func (s Style) UrlId(id string) Style { _ = "STUB: not implemented"; return *new(Style) }

// GetUrl returns the URL (id and actual URL) associated with the style.
// This is a hyper link that will be used for cells marked up with this style.
func (s Style) GetUrl() (id string, url string) { _ = "STUB: not implemented"; return "", "" }

// HasBold returns true if the style indicates bold text.
// Note that on some terminals bold text is simply brighter.
func (s Style) HasBold() bool { _ = "STUB: not implemented"; return false }

// HasBlink returns true if the style indicates blinking text.
func (s Style) HasBlink() bool { _ = "STUB: not implemented"; return false }

// HasReverse returns true if the style indicates reverse video text.
func (s Style) HasReverse() bool { _ = "STUB: not implemented"; return false }

// HasItalic returns true if the style indicates italicized text.
func (s Style) HasItalic() bool { _ = "STUB: not implemented"; return false }

// HasDim returns true if the style indicates dim or faint text.
func (s Style) HasDim() bool { _ = "STUB: not implemented"; return false }

// HasStrikeThrough returns true if the style indicates crossed-out text.
func (s Style) HasStrikeThrough() bool { _ = "STUB: not implemented"; return false }

// HasUnderline returns true if any underline style is set.
// Note that more detail is available via the GetUnderlineStyle
// and GetUnderlineColor methods.
func (s Style) HasUnderline() bool { _ = "STUB: not implemented"; return false }
