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

//go:build (!js && !wasm) || (js && wasm)
// +build !js,!wasm js,wasm

package tcell

import (
	"bytes"
	"sync"
	"time"

	"github.com/gdamore/tcell/v3/color"
	"golang.org/x/text/transform"
)

// NewTerminfoScreen returns a Screen that uses the stock TTY interface
// and POSIX terminal control, combined with a terminfo description taken from
// the $TERM environment variable.  It returns an error if the terminal
// is not supported for any reason.
//
// For terminals that do not support dynamic resize events, the $LINES
// $COLUMNS environment variables can be set to the actual window size,
// otherwise defaults taken from the terminal database are used.
func NewTerminfoScreen(opts ...TerminfoScreenOption) (Screen, error) {
	_ = "STUB: not implemented"
	return *new(Screen), nil
}

type TerminfoScreenOption interface {
	apply(*tScreen)
}

// OptColors forces the number of colors, overriding the value
// of the color count that would be detected by the environment.
// If the value is 0, then color is forced off.  Other reasonable values
// are 8, 16, 88, 256, or 1<<24.  The latter case intrinsically enables
// 24-bit color as well.
type OptColors int

func (o OptColors) apply(t *tScreen) { _ = "STUB: not implemented"; return }

// OptTerm overrides the detection of $TERM.
type OptTerm string

func (o OptTerm) apply(t *tScreen) {
	_ = "STUB: not implemented"

	// OptAltScreen controls whether the alternate screen buffer is used.
	// The default is true. The TCELL_ALTSCREEN=disable environment override
	// is still honored.
	return
}

type OptAltScreen bool

func (o OptAltScreen) apply(t *tScreen) { _ = "STUB: not implemented"; return }

// OptSanitizeContent enables stripping control characters from content passed
// to Put and PutStr. This is safer, but a little slower than leaving content
// unsanitized.
type OptSanitizeContent bool

func (o OptSanitizeContent) apply(t *tScreen) { _ = "STUB: not implemented"; return }

// OptAdvancedKeys enables richer key reporting where supported.  In this mode
// key events may include release state, repeat counts, and physical keys, and
// ASCII control letters are reported as KeyRune with ModCtrl instead of
// KeyCtrlA through KeyCtrlZ.  Shift-Tab is reported as KeyTab with ModShift,
// rather than KeyBacktab.
type OptAdvancedKeys bool

func (o OptAdvancedKeys) apply(t *tScreen) { _ = "STUB: not implemented"; return }

// OptKeyboardProtocol forces the keyboard reporting protocol instead of using
// startup negotiation. The zero value forces legacy keyboard reporting.
type OptKeyboardProtocol KeyProtocol

func (o OptKeyboardProtocol) apply(t *tScreen) { _ = "STUB: not implemented"; return }

// OptNegotiation controls whether terminal capabilities are negotiated during
// startup. The default is true.
type OptNegotiation bool

func (o OptNegotiation) apply(t *tScreen) { _ = "STUB: not implemented"; return }

// OptControlStringLimit sets the maximum inbound control-string payload size
// accepted from the terminal before the parser drops the sequence. This limits
// OSC and XDA strings, including OSC 52 clipboard strings; OSC 52 is the
// protocol used for writing clipboard data through the terminal. The default is
// 64 KiB; a value of 0 disables the limit.
type OptControlStringLimit int

func (o OptControlStringLimit) apply(t *tScreen) { _ = "STUB: not implemented"; return }

// Some terminal escapes that are basically universal.
// We would really like to be able to use private mode queries for some of
// these but generally we've found that support for queries is not always present,
// even when the private modes can be controlled. It appears that *all* terminals
// will happily swallow the escapes that they do not recognize, with the small annoyance
// in "st" where it prints error messages to its stderr (which is usually not visible
// to the user unless they started it from another terminal session).  But apart from
// the complaint to stderr from "st", everything else is fine.
const (
	enableAutoMargin  = "\x1b[?7h" // dec private mode 7
	setCursorPosition = "\x1b[%[1]d;%[2]dH"
	sgr0              = "\x1b[m" // attrOff
	bold              = "\x1b[1m"
	dim               = "\x1b[2m"
	italic            = "\x1b[3m"
	underline         = "\x1b[4m"
	blink             = "\x1b[5m"
	reverse           = "\x1b[7m"
	strikeThrough     = "\x1b[9m"
	clear             = "\x1b[H\x1b[J"
	doubleUnder       = "\x1b[4:2m"
	curlyUnder        = "\x1b[4:3m"
	dottedUnder       = "\x1b[4:4m"
	dashedUnder       = "\x1b[4:5m"
	underColor        = "\x1b[58:5:%dm"
	underRGB          = "\x1b[58:2::%d:%d:%dm"
	underFg           = "\x1b[59m"
	enableAltChars    = "\x1b(B\x1b)0"                      // set G0 as US-ASCII, G1 as DEC line drawing
	startAltChars     = "\x0e"                              // aka Shift-Out
	endAltChars       = "\x0f"                              // aka Shift-In
	setFg8            = "\x1b[3%dm"                         // for colors less than 8
	setFg256          = "\x1b[38;5;%dm"                     // for colors less than 256
	setFgRgb          = "\x1b[38;2;%d;%d;%dm"               // for RGB
	setBg8            = "\x1b[4%dm"                         // color colors less than 8
	setBg256          = "\x1b[48;5;%dm"                     // for colors less than 256
	setBgRgb          = "\x1b[48;2;%d;%d;%dm"               // for RGB
	setFgBgRgb        = "\x1b[38;2;%d;%d;%d;48;2;%d;%d;%dm" // for RGB, in one shot
	enterCA           = "\x1b[?1049h"                       // alternate screen
	exitCA            = "\x1b[?1049l"                       // alternate screen
	enterKeypad       = "\x1b[?1h\x1b="                     // Note mode 1 might not be supported everywhere
	exitKeypad        = "\x1b[?1l\x1b>"                     // Also mode 1
	requestWindowSize = "\x1b[18t"                          // For modern terminals
	requestPrimaryDA  = "\x1b[c"                            // Request primary device attributes
	requestExtAttr    = "\x1b[>q"                           // Request extended attribute (emulator name and version)
	setClipboard      = "\x1b]52;c;%s\x1b\\"                // Clipboard content is base64
	notifyDesktop9    = "\x1b]9;%[2]s\x1b\\"                // Args are title, body (but OSC 9 only has body)
	notifyDesktop777  = "\x1b]777;notify;%s;%s\x1b\\"       // Most commonly supported
	queryKittyKbd     = "\x1b[?u"                           // Query for Kitty keyboard support
	enableKittyKbd    = "\x1b[=1u"                          // Technically this pushes
	enableKittyKbdAdv = "\x1b[=15u"                         // disambiguation, events, alternate keys, all keys
	disableKittyKbd   = "\x1b[=0u"                          // Technically this means pop previous mode
	queryXTermKbd     = "\x1b[?4m"                          // Query for XTerm modify other keys support
	enableXTermKbd    = "\x1b[>4;2m"                        // Enable modify other keys protocol
	disableXTermKbd   = "\x1b[>4;0m"                        // Disable modify other keys protocol
)

// NewTerminfoScreenFromTty returns a Screen using a custom Tty implementation.
// If the passed in tty is nil, then a reasonable default (typically /dev/tty)
// is presumed, at least on UNIX hosts. (Windows hosts will typically fail this
// call altogether.)
func NewTerminfoScreenFromTty(tty Tty, opts ...TerminfoScreenOption) (Screen, error) {
	_ = "STUB: not implemented"
	return *new(Screen), nil
}

// tScreen represents a screen backed by a terminfo implementation.
type tScreen struct {
	tty                Tty
	h                  int
	w                  int
	fini               bool
	cells              CellBuffer
	buffering          bool // true if we are collecting writes to buf instead of sending directly to out
	buf                bytes.Buffer
	curstyle           Style
	style              Style
	resizeQ            chan bool
	quit               chan struct{}
	keyQ               chan []byte
	cx                 int
	cy                 int
	cls                bool // clear screen
	cursorx            int
	cursory            int
	acs                map[rune]string
	charset            string
	encoder            transform.Transformer
	decoder            transform.Transformer
	fallback           map[rune]string
	ncolor             int
	colors             map[color.Color]color.Color
	palette            []color.Color
	truecolor          bool
	noColor            bool
	legacy             bool
	hasClipboard       bool // true if OSC 52 reported via DA1
	finiOnce           sync.Once
	enterUrl           string
	exitUrl            string
	setWinSize         string
	cursorStyles       map[CursorStyle]string
	cursorStyle        CursorStyle
	cursorColor        color.Color
	cursorRGB          string
	cursorFg           string
	stopQ              chan struct{}
	eventQ             chan Event
	initQ              chan Event
	initted            bool
	running            bool
	startTime          time.Time
	wg                 sync.WaitGroup
	mouseFlags         MouseFlags
	pasteEnabled       bool
	focusEnabled       bool
	setTitle           string
	saveTitle          string
	restoreTitle       string
	title              string
	setClipboard       string
	notifyDesktop      string
	termName           string
	termVers           string
	term               string // value from $TERM
	altScreen          bool
	inlineResize       bool
	haveMouse          bool
	haveMouseSgr       bool
	haveKittyKbd       bool
	haveWin32Kbd       bool
	haveXTermKbd       bool
	forcedKbd          KeyProtocol
	forceKbd           bool
	negotiate          bool
	mouseDisabled      bool
	advancedKeys       bool
	controlStringLimit int
	input              *inputParser
	sync.Mutex
}

func (t *tScreen) useAltScreen() bool { _ = "STUB: not implemented"; return false }

func validKeyboardProtocol(p KeyProtocol) bool { _ = "STUB: not implemented"; return false }

func parseKeyboardProtocol(s string) (KeyProtocol, bool) {
	_ = "STUB: not implemented"
	return *new(KeyProtocol), false
}

func (t *tScreen) forceKeyboardProtocol(p KeyProtocol) bool {
	_ = "STUB: not implemented"
	return false
}

func (t *tScreen) applyKeyboardProtocolOverride() { _ = "STUB: not implemented"; return }

func (t *tScreen) applyEnvironmentOverrides() { _ = "STUB: not implemented"; return }

func (t *tScreen) Init() error { _ = "STUB: not implemented"; return nil }

// environment overrides

// On Windows, enable 24-bit color by default (all terminals there are 24-bit capable)

// base 8-bit palette

// monochrome variants

// legacy DEC VT 100/220 etc. family.  (technically the VT525 can do ANSI, but they should set to ansi)

// best guess - this covers all the modern variants like ghostty,

// A user who wants to have his themes honored can set this environment variable.

// these terminals are "legacy" and not expected to support most OSC functions

// clip to reasonable limits

// identity map for our builtin colors

func (t *tScreen) processInitQ() {
	_ = "STUB: not implemented"
	// NB: called with lock held
	return
}

// terminal specific overrides

// Some terminals can use OSC 9.  Unfortunately we can only discover
// them using this means.  It appears that pretty much all of them
// except iTerm2 also support more standard OSC 777, and it seems like
// only Kitty has its OSC 99 thing, but it also does OSC 777 well.

func (t *tScreen) filterEvents() chan Event { _ = "STUB: not implemented"; return nil }

func (t *tScreen) prepareExtendedOSC() { _ = "STUB: not implemented"; return }

// OSC 8 is for enter/exit URL.

// CSI .. t is for window operations.

// this also tries to request that UTF-8 is allowed in the title

// OSC 52 is for saving to the clipboard.
// this string takes a base64 string and sends it to the clipboard.
// it will also be able to retrieve the clipboard using "?" as the
// sent string, when we support that.

// OSC 777 is the desktop notification supported by a variety of
// newer terminals.  (There was also OSC 9 and OSC 99, but they
// are not as widely deployed, and OSC 9 is not unique.)

func (t *tScreen) prepareCursorStyles() { _ = "STUB: not implemented"; return }

func (t *tScreen) Fini() {
	_ = "STUB: not implemented"
	// Ensure that enough time passes for terminals to  finish sending
	// their initial response (gnome-terminal sends terminal dimensions
	// asynchronously later than the response to primary DA for some reason.)
	return
}

func (t *tScreen) finish() { _ = "STUB: not implemented"; return }

func (t *tScreen) SetStyle(style Style) { _ = "STUB: not implemented"; return }

func (t *tScreen) encodeStr(s string) []byte { _ = "STUB: not implemented"; return nil }

// Combining characters are elided

// resolvePalette looks up a color to obtain the palette entry for it.
func (t *tScreen) resolvePalette(c Color) Color { _ = "STUB: not implemented"; return *new(Color) }

// sendFgBg sends the foreground and background.  It is assumed that sgr0
// was already emitted prior to calling this (so colors are already in default).
func (t *tScreen) sendFgBg(fg Color, bg Color, attr AttrMask) AttrMask {
	_ = "STUB: not implemented"
	return *

	// foreground vs background, we calculate luminance
	// and possibly do a reverse video
	new(AttrMask)
}

// emitAttrs dumps prints the attributes, aside from underline that is special
// The assumption is that sgr0 was already printed ahead of this.
func (t *tScreen) emitAttrs(attrs AttrMask) { _ = "STUB: not implemented"; return }

// emitUl dumps prints the underline, which may be colored.
// The assumption is that sgr0 was already printed ahead of this.
func (t *tScreen) emitUnderline(us UnderlineStyle, uc Color) { _ = "STUB: not implemented"; return }

// NB: under color should have been reset by sgr0

// to ensure everyone gets at least a basic underline

// emitUrl either emits a url (OSC 8), or if the string is empty
// then the OSC 8 to exit the URL.  It should only be called if we
// either have a new URL, or need to exit an old one, as it always emits
// the OSC 8 sequence (if OSC 8 is supported).
func (t *tScreen) emitUrl(u urlInfo) { _ = "STUB: not implemented"; return }

// urlNeedsEmission reports whether a hyperlink transition has any wire effect.
// Url ids can be staged before the Url itself, and id-only transitions have no
// OSC 8 representation of their own.
func urlNeedsEmission(oldUrl, newUrl urlInfo) bool { _ = "STUB: not implemented"; return false }

func (t *tScreen) drawCell(x, y int) int { _ = "STUB: not implemented"; return 0 }

// URL string can be long, so don't send it unless we really need to.

// now emit runes - taking care to not overrun width with a
// wide character, and to ensure that we emit exactly one regular
// character followed up by any residual combing characters

// No FullWidth character support

// too wide to fit; emit a single space instead

// Clobber over any content in the next cell.
// This fixes a problem with some terminals where overwriting two
// adjacent single cells with a wide rune would leave an image
// of the second cell.  This is a workaround for buggy terminals.

func (t *tScreen) ShowCursor(x, y int) { _ = "STUB: not implemented"; return }

func (t *tScreen) SetCursor(cs CursorStyle, cc Color) { _ = "STUB: not implemented"; return }

func (t *tScreen) HideCursor() { _ = "STUB: not implemented"; return }

func (t *tScreen) showCursor() { _ = "STUB: not implemented"; return }

func (t *tScreen) Write(b []byte) (int, error) { _ = "STUB: not implemented"; return 0, nil }

func (t *tScreen) Print(s string) { _ = "STUB: not implemented"; return }

func (t *tScreen) Printf(f string, args ...any) { _ = "STUB: not implemented"; return }

func (t *tScreen) Show() { _ = "STUB: not implemented"; return }

func (t *tScreen) clearScreen() { _ = "STUB: not implemented"; return }

func (t *tScreen) startBuffering() { _ = "STUB: not implemented"; return }

func (t *tScreen) endBuffering() { _ = "STUB: not implemented"; return }

func (t *tScreen) hideCursor() {
	_ = "STUB: not implemented"
	// just in case we cannot hide it, move it to the end
	return
}

// then hide it

func (t *tScreen) draw() {
	_ = "STUB: not implemented"
	// clobber cursor position, because we're going to change it all
	return
}

// make no style assumptions

// hide the cursor while we move stuff around

// this is necessary so that if we ever
// go back to drawing that cell, we
// actually will *draw* it.

// restore the cursor

func (t *tScreen) EnableMouse(flags ...MouseFlags) { _ = "STUB: not implemented"; return }

func (t *tScreen) enableMouse(f MouseFlags) {
	_ = "STUB: not implemented"
	// Rather than using terminfo to find mouse escape sequences, we rely on the fact that
	// pretty much *every* terminal that supports mouse tracking follows the
	// XTerm standards (the modern ones).  It is expected that all terminals understand
	// the same DEC private modes.  Note that the SGR mode is required for the mouse sequences
	// to be understood.
	return
}

// We rely on dec private mode queries for this.
// If your terminal doesn't support these, then ask them to fix it.
// Note that as of macOS 26, macOS Terminal does not support them,
// so we enable the mouse unconditionally unless we get a report
// that says we have mouse, but not SGR mouse.  This is suboptimal, but
// a concession forced by the sorry state of terminal emulators.

// start by disabling all tracking.

func (t *tScreen) DisableMouse() { _ = "STUB: not implemented"; return }

func (t *tScreen) EnablePaste() { _ = "STUB: not implemented"; return }

func (t *tScreen) DisablePaste() { _ = "STUB: not implemented"; return }

func (t *tScreen) enablePasting(on bool) { _ = "STUB: not implemented"; return }

func (t *tScreen) EnableFocus() { _ = "STUB: not implemented"; return }

func (t *tScreen) DisableFocus() { _ = "STUB: not implemented"; return }

func (t *tScreen) enableFocusReporting() { _ = "STUB: not implemented"; return }

func (t *tScreen) disableFocusReporting() { _ = "STUB: not implemented"; return }

func (t *tScreen) Size() (int, int) { _ = "STUB: not implemented"; return 0, 0 }

func (t *tScreen) resize() { _ = "STUB: not implemented"; return }

func (t *tScreen) Colors() int {
	_ = "STUB: not implemented"
	// this doesn't change, no need for lock
	return 0
}

// vtACSNames is a map of bytes defined by terminfo that are used in
// the terminals Alternate Character Set to represent other glyphs.
// For example, the upper left corner of the box drawing set can be
// displayed by printing "l" while in the alternate character set.
// It's not quite that simple, since the "l" is the terminfo name,
// and it may be necessary to use a different character based on
// the terminal implementation (or the terminal may lack support for
// this altogether).  These values are from the DEC VT100, and all
// modern terminal emulators support this as charset 0.
var vtACSNames = map[byte]rune{
	'`': RuneDiamond,
	'a': RuneCkBoard,
	'f': RuneDegree,
	'g': RunePlMinus,
	'h': RuneBoard,
	'i': RuneLantern,
	'j': RuneLRCorner,
	'k': RuneURCorner,
	'l': RuneULCorner,
	'm': RuneLLCorner,
	'n': RunePlus,
	'o': RuneS1,
	'p': RuneS3,
	'q': RuneHLine,
	'r': RuneS7,
	's': RuneS9,
	't': RuneLTee,
	'u': RuneRTee,
	'v': RuneBTee,
	'w': RuneTTee,
	'x': RuneVLine,
	'y': RuneLEqual,
	'z': RuneGEqual,
	'{': RunePi,
	'|': RuneNEqual,
	'}': RuneSterling,
	'~': RuneBullet,
}

// buildAcsMap builds a map of characters that we translate from Unicode to
// alternate character encodings.  To do this, we use the standard VT100 ACS
// maps.  This is only done if the terminal lacks support for Unicode; we
// always prefer to emit Unicode glyphs when we are able.
func (t *tScreen) buildAcsMap() { _ = "STUB: not implemented"; return }

func (t *tScreen) scanInput(buf *bytes.Buffer) {
	_ = "STUB: not implemented"
	// The end of the buffer isn't necessarily the end of the input, because
	// large inputs are chunked. Set atEOF to false so the UTF-8 validating decoder
	// returns ErrShortSrc instead of ErrInvalidUTF8 for incomplete multi-byte codepoints.
	return
}

func (t *tScreen) mainLoop(stopQ chan struct{}) { _ = "STUB: not implemented"; return }

func (t *tScreen) inputLoop(stopQ chan struct{}) { _ = "STUB: not implemented"; return }

func (t *tScreen) Sync() { _ = "STUB: not implemented"; return }

func (t *tScreen) CharacterSet() string { _ = "STUB: not implemented"; return "" }

func (t *tScreen) RegisterRuneFallback(orig rune, fallback string) {
	_ = "STUB: not implemented"
	return
}

func (t *tScreen) UnregisterRuneFallback(orig rune) { _ = "STUB: not implemented"; return }

func (t *tScreen) SetSize(w, h int) { _ = "STUB: not implemented"; return }

func (t *tScreen) Resize(int, int, int, int) { _ = "STUB: not implemented"; return }

func (t *tScreen) Suspend() error { _ = "STUB: not implemented"; return nil }

func (t *tScreen) Resume() error { _ = "STUB: not implemented"; return nil }

func (t *tScreen) Tty() (Tty, bool) { _ = "STUB: not implemented"; return *new(Tty), false }

func (t *tScreen) applyKnownTerminalProfile(goos, termProgram string) bool {
	_ = "STUB: not implemented"
	return false
}

// macOS Terminal.app cannot handle the startup queries, but it does
// support modern mouse reporting.

// The WezTerm keyboard protocol to use is in theory driven by its
// own configuration, but we have found this unreliable because it
// does not mask unsupported capabilities.  Furthermore, on Windows
// builds the kitty protocol implementation is broken, while on other
// builds win32-input-mode is broken.  This is a best effort to make
// WezTerm work reasonably; our stronger advice is to choose another
// terminal program altogether.  This workaround will probably not
// apply to ssh sessions, as TERM_PROGRAM is not normally propagated.

func useVTWindowSizeQuery(goos string) bool { _ = "STUB: not implemented"; return false }

func useXTermKeyboardQuery(goos string) bool { _ = "STUB: not implemented"; return false }

// engage is used to place the terminal in raw mode and establish screen size, etc.
// Think of this is as tcell "engaging" the clutch, as it's going to be driving the
// terminal interface.
func (t *tScreen) engage() error { _ = "STUB: not implemented"; return nil }

// engageLocked is engage's implementation when t's lock is already held.
func (t *tScreen) engageLocked() error { _ = "STUB: not implemented"; return nil }

// macOS Terminal.app is brain damaged
// https://garrett.damore.org/2025/12/macos-terminal-still-missing-mark-apple.html
// Eventually they'll hopefully fix this.  As the environment variable
// does not convey by default via ssh, remote sessions might see spurious characters
// emitted during startup.  See the blog post for alternatives.

// XTerm's modifyOtherKeys mode is mainly useful for XTerm
// itself, and we do not use it on Windows.

// NB: MUST BE LAST

// Technically this may not be right, but every terminal we know about
// (even Wyse 60) uses this to enter the alternate screen buffer, and
// possibly save and restore the window title and/or icon.
// (In theory there could be terminals that don't support X,Y cursor
// positions without a setup command, but we don't support them.)

// disengage is used to release the terminal back to support from the caller.
// Think of this as tcell disengaging the clutch, so that another application
// can take over the terminal interface.  This restores the TTY mode that was
// present when the application was first started.
func (t *tScreen) disengage() { _ = "STUB: not implemented"; return }

// disengageStart begins a disengage operation while t's lock is already held.
// It returns true when disengageFinish must be called after releasing the lock.
func (t *tScreen) disengageStart() bool { _ = "STUB: not implemented"; return false }

// disengageFinish completes a disengage operation after disengageStart has
// released the running loops.
func (t *tScreen) disengageFinish() {
	_ = "STUB: not implemented"
	// wait for everything to shut down
	return
}

// shutdown the screen and disable special modes (e.g. mouse and bracketed paste)

// Hack for Windows.

// t.Print(t.disableCsiU)

// Beep emits a beep to the terminal.
func (t *tScreen) Beep() error { _ = "STUB: not implemented"; return nil }

// finalize is used to at application shutdown, and restores the terminal
// to it's initial state.  It should not be called more than once.
func (t *tScreen) finalize() { _ = "STUB: not implemented"; return }

func (t *tScreen) StopQ() <-chan struct{} { _ = "STUB: not implemented"; return nil }

func (t *tScreen) EventQ() chan Event { _ = "STUB: not implemented"; return nil }

func (t *tScreen) GetCells() *CellBuffer { _ = "STUB: not implemented"; return nil }

func (t *tScreen) SetTitle(title string) { _ = "STUB: not implemented"; return }

func (t *tScreen) SetClipboard(data []byte) {
	_ = "STUB: not implemented"
	// Post binary data to the system clipboard.  It might be UTF-8, it might not be.
	return
}

func (t *tScreen) GetClipboard() { _ = "STUB: not implemented"; return }

func (t *tScreen) HasClipboard() bool { _ = "STUB: not implemented"; return false }

func (t *tScreen) ShowNotification(title string, body string) { _ = "STUB: not implemented"; return }

func (t *tScreen) Terminal() (string, string) { _ = "STUB: not implemented"; return "", "" }

func (t *tScreen) KeyboardProtocol() KeyProtocol {
	_ = "STUB: not implemented"
	return *new(KeyProtocol)
}
