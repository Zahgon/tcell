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

package vt

import (
	"sync"
	"time"

	"github.com/gdamore/tcell/v3/tty"
)

// mockTerm implements MockTerm.
type mockTerm struct {
	mb MockBackend
	em Emulator
	ks *KeyboardState
}

// Stop the terminal.
func (mt *mockTerm) Stop() error { _ = "STUB: not implemented"; return nil }

// Start the terminal.
func (mt *mockTerm) Start() error { _ = "STUB: not implemented"; return nil }

// Drain all output from the terminal, ensuring
// any queued commands are processed.
func (mt *mockTerm) Drain() error { _ = "STUB: not implemented"; return nil }

// Read data from the terminal. This is called by a terminal
// application (e.g. via tcell Tty.)  Read data will include
// key strokes, mouse events, and responses to terminal queries.
func (mt *mockTerm) Read(data []byte) (int, error) {
	_ = "STUB: not implemented"
	return 0,

		// Write data to the terminal, typically either commands or data
		// that should be displayed on the virtual screen.
		nil
}

func (mt *mockTerm) Write(b []byte) (n int, err error) {
	_ = "STUB: not implemented"
	return 0,

		// WindowSize obtains the dimensions of the window.
		nil
}

func (mt *mockTerm) WindowSize() (tty.WindowSize, error) {
	_ = "STUB: not implemented"
	return *

	// No pixel sizes for now
	new(tty.WindowSize), nil
}

// NotifyResize registers a channel to be signaled when a resize has occurred.
// In real terminal emulators this would be posted (non-blocking) by a signal handler.
func (mt *mockTerm) NotifyResize(resizeq chan<- bool) { _ = "STUB: not implemented"; return }

// Close closes the terminal, after which it should no longer be used. Stop is implied.
func (mt *mockTerm) Close() error {
	_ = "STUB: not implemented"

	// Pos returns the cursor position.
	return nil
}

func (mt *mockTerm) Pos() Coord { _ = "STUB: not implemented"; return *new(Coord) }

// GetCell returns the contents of the cell at the given coordinates, or a zero value
// if the coordinates are out of range.
func (mt *mockTerm) GetCell(pos Coord) Cell {
	_ = "STUB: not implemented"
	return *

	// Bells counts the number of times the bell has rung.
	new(Cell)
}

func (mt *mockTerm) Bells() int { _ = "STUB: not implemented"; return 0 }

// KeyEvent is used to inject a key event.  Call this to inject
// a synthetic, fully specified key event.  Most uses should just use
// the KeyPress, KeyRelease, or even simpler KeyTap APIs.
func (mt *mockTerm) KeyEvent(ev KeyEvent) { _ = "STUB: not implemented"; return }

// Inject a delay to simulate human typing.
// Necessary to disambiguate Escape from other sequences.

// KeyPress implements MockTerm.KeyPress.
func (mt *mockTerm) KeyPress(k Key) { _ = "STUB: not implemented"; return }

// KeyRelease implements MockTerm.KeyRelease.
func (mt *mockTerm) KeyRelease(k Key) { _ = "STUB: not implemented"; return }

// KeyTap implements MockTerm.KeyTap.
func (mt *mockTerm) KeyTap(keys ...Key) { _ = "STUB: not implemented"; return }

// SetRepeat sets the repeat interval for the keyboard.
// Set the interval to zero to disable repeat.
func (mt *mockTerm) SetRepeat(delay, interval time.Duration) { _ = "STUB: not implemented"; return }

// MouseEvent implements MockTerm.MouseEvent.
func (mt *mockTerm) MouseEvent(ev MouseEvent) { _ = "STUB: not implemented"; return }

// FocusEvent implements MockTerm.FocusEvent.
func (mt *mockTerm) FocusEvent(focused bool) { _ = "STUB: not implemented"; return }

// GetTitle returns the current window title.
func (mt *mockTerm) GetTitle() string { _ = "STUB: not implemented"; return "" }

// SetSize is used to change the terminal size.
func (mt *mockTerm) SetSize(size Coord) { _ = "STUB: not implemented"; return }

// Backend returns the backend for testing.
func (mt *mockTerm) Backend() MockBackend {
	_ = "STUB: not implemented"

	// SendRaw is used to inject raw bytes to the read stream of the app.
	// Use this for fuzz testing.
	return *new(MockBackend)
}

func (mt *mockTerm) SendRaw(data []byte) { _ = "STUB: not implemented"; return }

// SetLayout sets the keyboard layout.
func (mt *mockTerm) SetLayout(km *Layout) { _ = "STUB: not implemented"; return }

// MockTerm is a mock terminal (emulator).  It can be used to
// test the emulator itself, or to test applications (or tcell) that
// uses the terminal.  It also implements the Tty interface used
// by tcell itself.
type MockTerm interface {
	tty.Tty

	// Pos reports the current cursor position.
	Pos() Coord

	// GetCell returns the cell at the given coordinates.
	// The coordinates must be valid.
	GetCell(Coord) Cell

	// Bells returns the number of times the bell has been rung.
	Bells() int

	// Inject a keyboard event - this is a full event, and bypasses
	// the layout and keyboard state processor.
	KeyEvent(KeyEvent)

	// Inject a key press.
	KeyPress(Key)

	// Inject a key release.
	KeyRelease(Key)

	// SetRepeat configures keyboard repeating. Repeat keystrokes
	// will be assumed after the key has been held for at least delay,
	// with new keys added each interval.
	SetRepeat(delay, interval time.Duration)

	// Inject one or more key press and releases.
	// The keys are pressed in the order, and released in reverse order.
	// Thus modifiers should be listed first.  This should not be used
	// to simulate typing a sequence (e.g. a word), but if you wanted to
	// test say N-Key rollover you could do that here.
	KeyTap(...Key)

	// Inject a mouse event.
	MouseEvent(MouseEvent)

	// Inject a focus event.
	FocusEvent(bool)

	// GetTitle obtains the current window title.
	GetTitle() string

	// SetSize is used to resize the terminal.
	SetSize(Coord)

	// SendRaw is used to send raw data to the application.
	// This is mostly intended to facilitate fuzz testing the application.
	SendRaw([]byte)

	// Backend returns the backend (used for testing).
	Backend() MockBackend

	// SetLayout sets the keyboard layout to use.
	// If not specified, a US standard ANSI keyboard will be assumed.
	SetLayout(*Layout)
}

type noMockBlit struct {
	MockBackend
	Blit struct{} // prevents use as Blitter
}

// NewMockTerm gives a mock terminal emulator.
func NewMockTerm(opts ...MockOpt) MockTerm { _ = "STUB: not implemented"; return *new(MockTerm) }

// MockBackend provides additional mock-specific capabilities on top of Backend.
// This is meant to facilitate test cases
type MockBackend interface {
	Backend

	// GetCell returns the cell at the given position, or an empty cell if the
	// position is out of the bounds of the window.
	GetCell(Coord) Cell

	// Bells counts the number of bells rung.
	Bells() int

	// GetTitle gets the current window title.
	GetTitle() string

	// SetSize is used to resize the window.
	// Newly added cells are empty, and content in old cells that out of range is lost.
	SetSize(Coord)

	// GetCursor is used to obtain the current cursor style.
	GetCursor() CursorStyle

	// SetClipboard sets the clipboard contents (copy buffer).
	SetClipboard([]byte)

	// GetClipboard returns the clipboard (copy buffer).
	GetClipboard() []byte

	// IsAdvancedKeyboard returns true, as we always support the full keyboard protocol.
	IsAdvancedKeyboard() bool
}

// mockBackend is a mock of a backend device for use with the emulator.
// It implements the following interfaces:
// vt.Backend, vt.Beeper, vt.Colorer, vt.Titler, vt.Resizer, vt.Blitter
type mockBackend struct {
	cells        []Cell // Content of cells
	size         Coord
	pos          Coord
	colors       int
	style        Style
	defaultStyle Style
	notifyQ      chan<- bool
	resized      bool
	newSize      Coord
	modes        map[PrivateMode]ModeStatus
	bells        int
	errs         int
	title        string
	clipboard    []byte
	cursor       CursorStyle
	lock         sync.Mutex
}

func (mb *mockBackend) GetSize() Coord { _ = "STUB: not implemented"; return *new(Coord) }

func (mb *mockBackend) Beep() { _ = "STUB: not implemented"; return }

func (mb *mockBackend) SetMouse(MouseReporting) { _ = "STUB: not implemented"; return }

func (mb *mockBackend) GetPrivateMode(pm PrivateMode) ModeStatus {
	_ = "STUB: not implemented"
	return *new(ModeStatus)
}

// note default (zero) value is ModeNA

func (mb *mockBackend) SetPrivateMode(pm PrivateMode, status ModeStatus) error {
	_ = "STUB: not implemented"
	return nil
}

func (mb *mockBackend) Put(pos Coord, cell Cell) {
	_ = "STUB: not implemented" // grapheme string, width int, style Style) {
	return
}

// writing to a cell right after a wide
// character clears that wide character (but leaves style/attributes)

// wide characters delete the next cell

func (mb *mockBackend) isPositionValid(pos Coord) bool { _ = "STUB: not implemented"; return false }

// index calculates the index in the cells array.  If the coordinates are invalid,
// -1 will be returned.
func (mb *mockBackend) index(pos Coord) int { _ = "STUB: not implemented"; return 0 }

func (mb *mockBackend) GetCell(pos Coord) Cell { _ = "STUB: not implemented"; return *new(Cell) }

func (mb *mockBackend) Bells() int { _ = "STUB: not implemented"; return 0 }

func (mb *mockBackend) GetPosition() Coord { _ = "STUB: not implemented"; return *new(Coord) }

func (mb *mockBackend) SetPosition(pos Coord) { _ = "STUB: not implemented"; return }

func (mb *mockBackend) Colors() int { _ = "STUB: not implemented"; return 0 }

func (mb *mockBackend) SetStyle(style Style) { _ = "STUB: not implemented"; return }

// SetWindowTitle implements the Titler interface.
func (mb *mockBackend) SetWindowTitle(title string) { _ = "STUB: not implemented"; return }

// GetTitle allows test code to observe what was set with SetWindowTitle.
func (mb *mockBackend) GetTitle() string { _ = "STUB: not implemented"; return "" }

// NotifyResize registers a channel to be written to (non-blocking) if the
// backend changes size.
func (mb *mockBackend) NotifyResize(rq chan<- bool) { _ = "STUB: not implemented"; return }

// checkSize performs a possible terminal resize. Cells that are
// added are treated as empty, while cells that are removed are just lost.
// (Note that at least one other emulator erases content on a resize.  There is no
// standard for what to do here.) This is done inline when calculating the index.
// The caller is expected to hold mb.lock.
func (mb *mockBackend) checkSize() { _ = "STUB: not implemented"; return }

func (mb *mockBackend) RaiseResize() { _ = "STUB: not implemented"; return }

// SetSize is used to change the size of the virtual terminal.
func (mb *mockBackend) SetSize(size Coord) { _ = "STUB: not implemented"; return }

// Reset the terminal to startup defaults.
func (mb *mockBackend) Reset() { _ = "STUB: not implemented"; return }

func (mb *mockBackend) Blit(src, dst, dim Coord) { _ = "STUB: not implemented"; return }

// clip to visible source

// and clip to final destination

// gap represents decrement when shifting to the next row --
// skipping over the irrelevant cells. (The increment in the
// index when going from last cell of row to first cell of next row,
// or vice versa.)

// the following logic is carefully constructed to avoid expensive
// operations in the loops (only addition or subtraction)
// source appears later, so we can forward copy

// advance to next row

// source appears earlier, so we have to reverse copy

// Buffering is not supported by the mockBackend, and there is little point in it.
func (mb *mockBackend) Buffering(bool) {
	_ = "STUB: not implemented"

	// SetCursor is used to set how the cursor is displayed.
	return
}

func (mb *mockBackend) SetCursor(cs CursorStyle) { _ = "STUB: not implemented"; return }

// GetCursor returns the current cursor style.
func (mb *mockBackend) GetCursor() CursorStyle { _ = "STUB: not implemented"; return *new(CursorStyle) }

// SetClipboard sets the current clipboard contents.
func (mb *mockBackend) SetClipboard(data []byte) { _ = "STUB: not implemented"; return }

// GetClipboard gets the current clipboard contents.
func (mb *mockBackend) GetClipboard() []byte { _ = "STUB: not implemented"; return nil }

// IsAdvancedKeyboard returns true - we always implement
// the raw keyboard protocol.
func (mb *mockBackend) IsAdvancedKeyboard() bool {
	_ = "STUB: not implemented"

	// MockOpt is an interface by which options can change the behavior of the mocked terminal.
	// This is intended to permit easier testing.
	return false
}

type MockOpt interface{ SetMockOpt(mb *mockBackend) }

// MockOptSize changes the default terminal size, which is normally 80x24.
type MockOptSize Coord

func (o MockOptSize) SetMockOpt(mb *mockBackend) {
	_ = "STUB: not implemented"

	// MockOptColors changes the number of colors the terminal supports.
	return
}

type MockOptColors int

func (o MockOptColors) SetMockOpt(mb *mockBackend) {
	_ = "STUB: not implemented"

	// MockOptNoBlit suppresses the blitter interface.
	return
}

type MockOptNoBlit struct{}

func (MockOptNoBlit) SetMockOpt(mb *mockBackend) {
	_ = "STUB: not implemented"

	// MockOpt8BitControls enables raw 8-bit and UTF-8 encoded C1 controls in the
	// emulator. The default is to accept only 7-bit ESC-prefixed controls.
	return
}

type MockOpt8BitControls struct{}

func (MockOpt8BitControls) SetMockOpt(mb *mockBackend) {
	_ = "STUB: not implemented"

	// NewMockBackend returns a MockBackend modified by the given options.
	// The default is a fully featured 256-color backend with initial size 80x24.
	return
}

func NewMockBackend(options ...MockOpt) MockBackend {
	_ = "STUB: not implemented"
	return *new(MockBackend)
}
