// Copyright 2025 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vt

// KeyEvent is a key event.
type KeyEvent struct {
	Down   bool     // true if event is for key down event
	Repeat int      // if > 1, a repeat count
	Key    Key      // Key symbol.
	Base   BaseKey  // base key code (physical key, e.g 'a'), may be zero if same as code
	VK     WinVK    // Windows virtual key. 0 for none, or if not known.
	SC     ScanCode // Windows scan code. 0 for none, or if not known.
	Mod    Modifier // modifiers
	Utf    string   // if non-empty, the unicode content for this
}

type Modifier int

const (
	ModNone   = Modifier(0)
	ModLShift = Modifier(1 << iota)
	ModRShift
	ModLCtrl
	ModRCtrl
	ModLAlt
	ModRAlt
	ModLMeta
	ModRMeta
	ModLHyper
	ModRHyper
	ModCapsLock
	ModNumLock
)

func (m Modifier) IsShift() bool    { _ = "STUB: not implemented"; return false }
func (m Modifier) IsCtrl() bool     { _ = "STUB: not implemented"; return false }
func (m Modifier) IsAlt() bool      { _ = "STUB: not implemented"; return false }
func (m Modifier) IsMeta() bool     { _ = "STUB: not implemented"; return false }
func (m Modifier) IsHyper() bool    { _ = "STUB: not implemented"; return false }
func (m Modifier) IsNumLock() bool  { _ = "STUB: not implemented"; return false }
func (m Modifier) IsCapsLock() bool { _ = "STUB: not implemented"; return false }
func (m Modifier) IsCapitals() bool { _ = "STUB: not implemented"; return false }
func (m Modifier) IsAltGr() bool    { _ = "STUB: not implemented"; return false }

// Button is the mouse button pressed or released.
type Button int

const (
	NoButton = Button(0)         // No buttons are pressed.
	Button1  = Button(1 << iota) // Usually left most button.
	Button2                      // Usually right most button.
	Button3                      // Usually middle button.
	Button4
	Button5
	Button6
	Button7
	Button8
	WheelUp    // Wheel motion up/away from user.
	WheelDown  // Wheel motion down/towards user.
	WheelLeft  // Wheel motion to left.
	WheelRight // Wheel motion to right.
)

// MouseEvent reports a single mouse event.  Only a single button
// may be reported for a given event.  The application will have
// to keep state.  As buttons are never pressed exactly simultaneously,
// the backend will send chords as a series of presses followed by a series
// of releases.
type MouseEvent struct {
	Position Coord    // Location of pointer.
	Button   Button   // Buttons pressed.
	Down     bool     // True on press, false on release.
	Motion   bool     // True if mouse moved at least once cell.
	Mod      Modifier // Modifiers (for modified click).
}

// encodeButton just encodes the XTerm style button details into a byte
func (ev MouseEvent) encodeButton() byte { _ = "STUB: not implemented"; return 0 }

// intentionally reversed with button 3
