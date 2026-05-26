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
	"bytes"
	"io"
	"sync"
	"unicode/utf8"

	"github.com/clipperhouse/uax29/v2/graphemes"
	"github.com/gdamore/tcell/v3/color"
)

// Emulator is a terminal emulator API. It implements the state machinery
// (escape parsing and so forth) associated with being a terminal emulator.
// The backend handles rendering the content, and some low level details.
//
// NOTE: This is not a committed interface yet, its entirely a work in progress.
type Emulator interface {
	// SetId sets our identity.
	SetId(name string, version string)

	// SendRaw sends raw data to the consumer.  This bypasses the normal encoding,
	// so it should be used with caution.
	SendRaw([]byte)

	// KeyEvent injects a keyboard event into the emulator, which will ultimately
	// result in data being sent via SendRaw.
	KeyEvent(ev KeyEvent)

	// ResizeEvent is called by a backend when the terminal has resized
	// This will send in-band resize notifications if the client has requested them.
	ResizeEvent(Coord)

	// MouseEvent is called by a backend to report mouse activity.
	MouseEvent(ev MouseEvent)

	// FocusEvent is called by a backend to report that focus is gained (true) or lost (false).
	FocusEvent(bool)

	// Drain waits until any queued but not processed input has finished processing.
	// It also wakes the reader.
	Drain() error

	// Start starts processing.
	Start() error

	// Stop stops processing.
	Stop() error

	// Reader reads data from the emulator.  These are bytes that would be transmitted
	// to a remote party.
	io.Reader

	// Writer writes data to the emulator.  These are commands that the emulator should process.
	io.Writer
}

// Style represents the styling of a cell.
// This is an interface to prevent direct modification.
type Style interface {
	Fg() color.Color              // Fg returns the foreground color.
	Bg() color.Color              // Bg returns the background color.
	Uc() color.Color              // Uc returns the underline color.k
	Attr() Attr                   // Attr returns the associated attributes.
	Url() (string, string)        // Url returns the URL and associated id if one was set.
	WithFg(color.Color) Style     // WithFg creates a new style with the foreground
	WithBg(color.Color) Style     // WithBg creates a new style with the background.
	WithUc(color.Color) Style     // WithUc creates a new style with the underline color
	WithAttr(Attr) Style          // WithAttr creates a new style with the attributes.
	WithUrl(string, string) Style // WithLink creates a new style with the URL and id.
	Equal(Style) bool             // Equal returns true if the styles are the same.
}

// styleStruct implements Style.  Note that it is possible to make this even more
// compact, but we don't think further optimization here on size will justify the
// complexity and runtime performance hit to do so.  We're also already only storing
// a class reference to this per cell.
type styleStruct struct {
	fg   color.Color
	bg   color.Color
	uc   color.Color // underline color
	attr Attr
	url  string // URL
	id   string // Id for link
}

var BaseStyle = &styleStruct{}

var asciiRuneStrings = func() [utf8.RuneSelf]string {
	var table [utf8.RuneSelf]string
	for i := 0; i < utf8.RuneSelf; i++ {
		table[i] = string(rune(i))
	}
	return table
}()

const (
	runeStringCacheSize      = 32
	clusterStringCacheSize   = 32
	clusterStringCacheMaxLen = 128
)

type runeStringCache struct {
	entries [runeStringCacheSize]runeStringCacheEntry
	n       int
}

type runeStringCacheEntry struct {
	r rune
	s string
}

func (c *runeStringCache) stringFor(r rune) string { _ = "STUB: not implemented"; return "" }

type clusterStringCache struct {
	entries [clusterStringCacheSize]clusterStringCacheEntry
	n       int
}

type clusterStringCacheEntry struct {
	n int
	b [clusterStringCacheMaxLen]byte
	s string
}

func (c *clusterStringCache) stringFor(cluster []byte) string { _ = "STUB: not implemented"; return "" }

func (ss *styleStruct) Fg() color.Color              { _ = "STUB: not implemented"; return *new(color.Color) }
func (ss *styleStruct) Bg() color.Color              { _ = "STUB: not implemented"; return *new(color.Color) }
func (ss *styleStruct) Uc() color.Color              { _ = "STUB: not implemented"; return *new(color.Color) }
func (ss *styleStruct) Attr() Attr                   { _ = "STUB: not implemented"; return *new(Attr) }
func (ss *styleStruct) Url() (string, string)        { _ = "STUB: not implemented"; return "", "" }
func (ss *styleStruct) WithFg(fg color.Color) Style  { _ = "STUB: not implemented"; return *new(Style) }
func (ss *styleStruct) WithBg(bg color.Color) Style  { _ = "STUB: not implemented"; return *new(Style) }
func (ss *styleStruct) WithUc(uc color.Color) Style  { _ = "STUB: not implemented"; return *new(Style) }
func (ss *styleStruct) WithAttr(a Attr) Style        { _ = "STUB: not implemented"; return *new(Style) }
func (ss *styleStruct) WithUrl(url, id string) Style { _ = "STUB: not implemented"; return *new(Style) }
func (ss *styleStruct) Equal(other Style) bool       { _ = "STUB: not implemented"; return false }

// We have chosen not to support alternative implementations for this compare.
// We could delegate to the other style, but that could lead to a loop if they
// do the same.

// Cell is a representation of a display cell. Most consumers will not need this.
// Note, this is not the simplest possible representation, and a 256x256 cell
// display is going to need about 3MB to store it all, but it's simple, and adequate
// to retain pretty much all of what we need for Unicode.  We could save some memory
// by using explicit struct pointers and by eliminating grapheme cluster support,
// but modern users expect these features.
type Cell struct {
	C string // Content, it will be a grapheme cluster
	S Style  // Style, a pointer is used efficiency
	W int    // Display width (0, 1, or 2)
}

// EmulatorOpt configures an Emulator.
type EmulatorOpt interface {
	setEmulatorOpt(*emulator)
}

// EmulatorOpt8BitControls enables parsing of C1 controls, such as CSI and OSC,
// when they are presented as raw 8-bit bytes or UTF-8 encoded C1 controls.
// The default is to only accept the 7-bit ESC-prefixed forms.
type EmulatorOpt8BitControls struct{}

func (EmulatorOpt8BitControls) setEmulatorOpt(em *emulator) { _ = "STUB: not implemented"; return }

// NewEmulator creates an emulator instance on top of the given backend.
// The input is relative to the emulator, so it receives data from the host,
// whereas the emulator sends data to the application through the output.
func NewEmulator(be Backend, opts ...EmulatorOpt) Emulator {
	_ = "STUB: not implemented"
	return *new(Emulator)
}

// we never support VT52 mode (note ON means ANSI mode)

// add mouse modes - we also add focus reporting mode

// emulator is an implementation of a terminal emulator built on top of
// a Backend.  It implements the common escape sequence handling and high
// level functionality that a real terminal emulator, or a mock, would need.
type emulator struct {
	stopQ          chan bool
	writeQ         chan any // queues data from application to emulator
	readQ          chan any // queues data from emulator to application
	be             Backend
	inBuf          *bytes.Buffer // buffer queued for input
	inb            func(byte)    // input byte function (faster than state switch)
	style          Style
	defaultStyle   Style
	utfLen         int
	pos            Coord
	buffering      uint         // reference count - number of (re-entrant) buffering calls
	autoWrap       bool         // next character will wrap (auto margin, deferred until char emitted)
	c1Allowed      bool         // allow C1 controls in raw 8-bit and UTF-8 encodings
	c1Enabled      bool         // C1 controls are currently enabled
	c1Prefix       bool         // string parser has seen the first byte of a UTF-8 encoded C1 control
	appKeyPad      bool         // use application key pad keys?
	name           string       // name of this emulator (used for extended attributes)
	vers           string       // version string of this emulator (used for extended attributes)
	saved          savedCursor  // data saved by save cursor (DECSC)
	sendLock       sync.Mutex   // ensures that send data cannot be intermixed
	modeLock       sync.RWMutex // protects localModes/ansiModes and related derived state
	tabStops       []Col        // tab stops, ordered. if nil every 8th position is used
	lastIndex      int          // index of last cell written + 1 (for grapheme clustering) (zero means none)
	graphemeBuf    []byte       // scratch buffer for grapheme clustering checks
	graphemeIter   graphemes.Iterator[[]byte]
	runeStrings    runeStringCache
	clusterStrings clusterStringCache
	cells          []Cell         // content of cells, we have to maintain our own copy (backend might or might not)
	mouseReports   MouseReporting // whether we have enabled mouse reports
	size           Coord          // physical window size
	topMargin      Row            // top margin, scrollable region includes this row
	botMargin      Row            // bottom margin, scrollable region includes this row
	ltMargin       Col            // left margin, scrollable region to the right
	rtMargin       Col            // right margin, scrollable region to the left
	cursor         CursorStyle    // current cursor style (visibility, blink, shape)

	localModes map[PrivateMode]ModeStatus // some modes we handle locally
	ansiModes  map[AnsiMode]ModeStatus    // some modes we handle locally
}

// savedCursor is the content we save when saving the cursor,
// which is more than just the cursor location itself.
type savedCursor struct {
	pos      Coord
	style    Style
	autoWrap bool
	// We should probably store OSC 8 data here, eventually.
	// TODO: Character sets
	// TODO: Origin mode (DEC Mode 6)
}

func (em *emulator) saveCursor() { _ = "STUB: not implemented"; return }

func (em *emulator) restoreCursor() { _ = "STUB: not implemented"; return }

func (em *emulator) bufferingStart() { _ = "STUB: not implemented"; return }

func (em *emulator) bufferingEnd() { _ = "STUB: not implemented"; return }

// inbInit processes bytes received in the "default" state. Most often these are just
// text characters to display on screen, but if ESC is seen then additional processing will result.
func (em *emulator) inbInit(b byte) {
	_ = "STUB: not implemented"

	// hot path - just doing ASCII directly.
	return
}

// plain ascii

// For C1 controls, the raw 8-bit form is the same as ESC followed by
// (b - 0x40). This is disabled by default because modern protocols
// generally treat these forms as insecure.

// TODO: To support non-UTF-8 locales, include a check here for > 0x7F.  Those locales
// might preclude 8-bit control sequences - 8859 character sets are fine, but e.g. KOI8,
// and ShiftJIS use values in those ranges.

// ESC (escape)

// BEL (bell)

// BS (backspace)

// horizontal tab

// LF (line feed), VF, FF

// CR (carriage return)

// TODO: SO

// TODO: SI

//TODO Cancel (reset parser)

// TODO: consider separating Unicode from other 8-bit character sets

// inbEsc processes the next byte after an escape character is seen.
func (em *emulator) inbEsc(b byte) {
	_ = "STUB: not implemented"

	// By default, reset to init state. Other states will be set explicitly as needed.
	return
}

// 0x20 - 0x2F -- usually followed by just one terminating character, but could include others

// privacy message (PM)

// application program command (APC)

// down one line (IND)

// next line (NEL)

// set tab stop (HTS) - VT52 is go home, but we do not support VT52

// up one line (RI)

// single shift two (SS2) (TODO)
// single shift three (SS3) (TODO)

// device control string (DCS) (TODO)
// start of string (SOS)

// DECID, obsolete form to get primary DA

// RIS, soft reset

// back index (DECBI, VT420, not widely supported)

// save cursor (DECSC, VT100)

// restore cursor (DECRC, VT100)

// forward index (DECFI, VT420, not widely supported)

// ESC-V and ESC-W are for guarded area (TODO)

// inbNF processes bytes that are part of an "nF" sequence (see ECMA-48).
func (em *emulator) inbNF(b byte) { _ = "STUB: not implemented"; return }

// not a valid sequence

// DECALN - fill screen with 'E'

// TODO: Reset DECOM (when we implement origin mode)

// most implementations leave the cursor at home for this

// case "%@": // TODO: select 8859-1
// case "%G": // TODO: select UTF-8
// case "(A": // TODO: select G0 as UK
// case "(B": // TODO: select G0 as US
// case "(C", "(5": // TODO: select G0 as Finnish
// case "(H", "(7": // TODO: select G0 as Swedish
// case "(K": // TODO: select G0 as German
// case "(Q", "(9": // TODO: select G0 as French Canadian
// case "(R", "(f": // TODO: select G0 as French
// case "(Y": // TODO: select G0 as Italian
// S7C1T - send/use 7-bit C1 controls

// S8C1T - send/use 8-bit C1 controls

// inbCSI handles bytes that are part of a CSI based sequence.
func (em *emulator) inbCSI(b byte) { _ = "STUB: not implemented"; return }

// parameter bytes

// intermediate bytes

// error state

// inbOSC handles bytes that are part of on OSC sequences (operating system command).
func (em *emulator) inbOSC(b byte) { _ = "STUB: not implemented"; return }

// inbStr handles PM, SOS, and any other string we want to consume and discard.
func (em *emulator) inbStr(b byte) { _ = "STUB: not implemented"; return }

func (em *emulator) inbStringC1(b byte, done func()) bool { _ = "STUB: not implemented"; return false }

// inbUTF handles continuation bytes for UTF-8 sequences.
func (em *emulator) inbUTF(b byte) {
	_ = "STUB: not implemented"

	// good continuation byte
	return
}

func (em *emulator) beep() { _ = "STUB: not implemented"; return }

// numericParams splits the string consisting of numeric parameters into integers.
// It ensures a minimum number are present (needed for some safety cases).
// Empty strings default to zero.
func numericParams(str string, minimumLen int) ([]int, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

// parseSgrColor grabs either 2 arguments, or 4 arguments for palette or rgb values
// used with SGR 38 and 48. The arguments must be numbers, and are returned as such.
func (em *emulator) parseSgrColor(args []string, words []string) (color.Color, []string, error) {
	_ = "STUB: not implemented"
	return *new(color.Color), nil, nil
}

// RGB (direct)

// palette index

// RGB color

// palette index

func (em *emulator) pickColor(c color.Color, def color.Color) (color.Color, bool) {
	_ = "STUB: not implemented"
	return *new(color.Color), false
}

// processSgr processes SGR commands (things that change how characters are displayed).
func (em *emulator) processSgr(str string) { _ = "STUB: not implemented"; return }

// technically parameters for 38 or 48 should be separated by colons, but due to historical
// accident it is more common to see semicolon separation.  Underline styles are also separated
// by a colon, if present.

// we do this instead of a range so we can lop off
// multiple words for SGR38 and 48.

// just swallow it for now

// ignore, its for invisible

// Doubly underlined, per ECMA

// simple foreground colors

// simple background colors

// processCursorUp implements CUU.
func (em *emulator) processCursorUp(str string) { _ = "STUB: not implemented"; return }

// processCursorDown implements CUD.
func (em *emulator) processCursorDown(str string) { _ = "STUB: not implemented"; return }

// processCursorForward implements CUF.
func (em *emulator) processCursorForward(str string) { _ = "STUB: not implemented"; return }

// processCursorBackward implements CUB.
func (em *emulator) processCursorBackward(str string) { _ = "STUB: not implemented"; return }

// processCursorNextLine implements CNL.
func (em *emulator) processCursorNextLine(str string) { _ = "STUB: not implemented"; return }

// processCursorPreviousLine implements CPL.
func (em *emulator) processCursorPreviousLine(str string) { _ = "STUB: not implemented"; return }

// processCursorColumn implements CHA.
func (em *emulator) processCursorColumn(str string) { _ = "STUB: not implemented"; return }

// TODO: possibly clip to margins (origin mode)

// processCursorPosition implements CUP, and also HVP.
func (em *emulator) processCursorPosition(str string) { _ = "STUB: not implemented"; return }

// processCursorTab implements CHT.
func (em *emulator) processCursorTab(str string) { _ = "STUB: not implemented"; return }

// Note: tab does not clear this field.

// processCursorBackTab implements CBT.
func (em *emulator) processCursorBackTab(str string) { _ = "STUB: not implemented"; return }

// processEraseDisplay implements ED.
func (em *emulator) processEraseDisplay(str string) { _ = "STUB: not implemented"; return }

// erase below

// erase above

// erase all

// others not supported (3 is erase saved lines)

// processEraseLine implements EL.
func (em *emulator) processEraseLine(str string) { _ = "STUB: not implemented"; return }

// processEraseCharacter implements ECH.
// This ignores the margin.
func (em *emulator) processEraseCharacter(str string) { _ = "STUB: not implemented"; return }

// TODO: delete wide character if we are splitting it at the start

// processScrollUp implements SU (VT420.)
func (em *emulator) processScrollUp(str string) { _ = "STUB: not implemented"; return }

// TODO: consider faster jump scroll.
// This should be something tunable as well.

// processScrollDown implements SD (VT420.)
func (em *emulator) processScrollDown(str string) { _ = "STUB: not implemented"; return }

// TODO: consider faster jump scroll.
// This should be something tunable as well.

// processWindowOps handles CSI ... t window operations.
func (em *emulator) processWindowOps(str string) { _ = "STUB: not implemented"; return }

// Resize window: CSI 8 ; rows ; cols t

// Report text area size: CSI 8 ; rows ; cols t

// processVerticalMargins implements DECSTBM (set top and bottom margins, VT220.)
func (em *emulator) processVerticalMargins(str string) { _ = "STUB: not implemented"; return }

// no change if values are out of range

// processHorizontalMargins implements DECSLRM (set left and right margins, VT400.)
// It only works if Private Mode 69 (Left and Right margins)
func (em *emulator) processHorizontalMargins(str string) { _ = "STUB: not implemented"; return }

// For compat with SCO and ANSI.SYS.

// no change if values are out of range

// processIndex moves down, unless already on the bottom margin, in which case it scrolls Up.
func (em *emulator) processIndex() { _ = "STUB: not implemented"; return }

// processCarriageReturn handle CR.
func (em *emulator) processCarriageReturn() { _ = "STUB: not implemented"; return }

// processLineFeed is like IND, but if ANSI mode 20 is set, then a CR is appended as well.
func (em *emulator) processLineFeed() { _ = "STUB: not implemented"; return }

// processReverseIndex moves up, unless already on the top margin, in which case it scrolls down.
func (em *emulator) processReverseIndex() { _ = "STUB: not implemented"; return }

// processCursorRow implements VPA (set vertical position absolute)
func (em *emulator) processCursorRow(str string) { _ = "STUB: not implemented"; return }

// processCursorRowAdvance implements VPR (vertical position relative).
func (em *emulator) processCursorRowAdvance(str string) { _ = "STUB: not implemented"; return }

// processInsertLine implements IL.
func (em *emulator) processInsertLine(str string) {
	_ = "STUB: not implemented"
	// insert line only takes effect within the scrolling region
	return
}

// process as a scroll down at our position.

// processDeleteLine implements DL.
func (em *emulator) processDeleteLine(str string) {
	_ = "STUB: not implemented"
	// delete line only takes effect within the scrolling region
	return
}

// process as a scroll up at our position.

// processDeleteCharacter implements DCH.
func (em *emulator) processDeleteCharacter(str string) {
	_ = "STUB: not implemented"
	// only takes effect within the scrolling region
	return
}

// this is essentially a one line scroll left

// if we are breaking a wide rune, delete it (but preserve style)

// processInsertCharacter implements ICH.
func (em *emulator) processInsertCharacter(str string) {
	_ = "STUB: not implemented"

	// only takes effect within the scrolling region -- HOWEVER,
	// ICH still resets auto-wrap in this case, unlike DCH.
	return
}

// this is essentially a one line scroll right

// if we are breaking a wide rune, delete it

// NB: We don't use eraseCell, because we need to preserve attributes.

// if we clipped off the end of a wide character, then delete it.

// processSetMode implements SM (set ANSI mode).
func (em *emulator) processSetMode(str string) { _ = "STUB: not implemented"; return }

// processResetMode implements RM (reset ANSI mode).
func (em *emulator) processResetMode(str string) { _ = "STUB: not implemented"; return }

// processRequestMode implements DECRQM for ANSI modes.
// Only a single numeric parameter (mode number) can be supplied (VT300+)
func (em *emulator) processRequestMode(str string) { _ = "STUB: not implemented"; return }

// processSetPrivateMode implements DECSET (set private mode).
func (em *emulator) processSetPrivateMode(str string) { _ = "STUB: not implemented"; return }

// processResetPrivateMode implements DECRST (reset private mode).
func (em *emulator) processResetPrivateMode(str string) { _ = "STUB: not implemented"; return }

// processRequestPrivateMode implements DECRQM for private modes.
// Only a single numeric parameter (mode number) can be supplied (VT300+)
func (em *emulator) processRequestPrivateMode(str string) { _ = "STUB: not implemented"; return }

// processTabReset implements DECST8C (set tab stops to every 8 chars)
func (em *emulator) processTabReset(str string) { _ = "STUB: not implemented"; return }

// processTabClear implements TBC (clear horizontal tab).
func (em *emulator) processTabClear(str string) { _ = "STUB: not implemented"; return }

// clear stop at current column

// clear all columns
// this is distinct from nil

// processPrimaryAttributes implements send DA.
func (em *emulator) processPrimaryAttributes(str string) { _ = "STUB: not implemented"; return }

// processExtendedAttributes implements XTVERSION (send terminal name and version).
func (em *emulator) processExtendedAttributes(str string) { _ = "STUB: not implemented"; return }

// processCursorStyle implements DECSCUSR (set cursor style).
func (em *emulator) processCursorStyle(str string) {
	_ = "STUB: not implemented"
	// get previous visibility state, as we don't change it with this call.
	return
}

// processCsi processes CSI sequences.
func (em *emulator) processCsi(final byte) {
	_ = "STUB: not implemented"

	// CSI sequences are supported in several different possible ways:
	// parameters may have a prefix character that is not numeric, typically
	// indicating a whole different mode of operation than the final byte.
	// There may also be intermediate bytes, but we only look for one, because
	// the use cases we have this are that only a single intermediate byte is
	// sometimes used to affect function.  (E.g. $ in some cases.)
	return
}

// processClipboard handles OSC 52 commands.
func (em *emulator) processClipboard(str string) { _ = "STUB: not implemented"; return }

// first parameter is the target.  We only have a single
// target, and alias all possibilities to the same.

// request for clipboard content

// processHyperLink handles OSC 8 commands.
func (em *emulator) processHyperLink(str string) {
	_ = "STUB: not implemented"
	// format is params;URI params are colon separated key value pairs.
	// if the URI is absent, then the link is terminated.
	return
}

// No URI

// processOSC processes an operating system command.
func (em *emulator) processOSC() {
	_ = "STUB: not implemented"

	// Every OSC we support has a number, semicolon, then string.
	return
}

// Set window title

// TODO: possibly validate the UTF-8 content?

func (em *emulator) getPosition() Coord { _ = "STUB: not implemented"; return *new(Coord) }

func (em *emulator) setPosition(pos Coord) { _ = "STUB: not implemented"; return }

func (em *emulator) deviceReport(s string) { _ = "STUB: not implemented"; return }

// ignore

func (em *emulator) moveUpN(count Row) { _ = "STUB: not implemented"; return }

func (em *emulator) moveDownN(count Row) { _ = "STUB: not implemented"; return }

func (em *emulator) moveLeftN(count Col) { _ = "STUB: not implemented"; return }

func (em *emulator) moveRightN(count Col) { _ = "STUB: not implemented"; return }

// moveDown moves down, to the limit of either bottom margin, or the bottom of the screen if outside the margin.
func (em *emulator) moveDown() { _ = "STUB: not implemented"; return }

// moveUp moves up, to the limit of either top margin, or zero if outside the margin.
func (em *emulator) moveUp() { _ = "STUB: not implemented"; return }

func (em *emulator) moveLeft() { _ = "STUB: not implemented"; return }

func (em *emulator) moveRight() { _ = "STUB: not implemented"; return }

// nextLine is like CNL with 1, but it optionally also scrolls.
func (em *emulator) nextLine() { _ = "STUB: not implemented"; return }

// blit performs a data move operation.  It does ignores margins.
func (em *emulator) blit(src, dst, dim Coord) { _ = "STUB: not implemented"; return }

// save the source and destination for the backend blit

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

// Now we possibly blit underneath.  We'll use the underlying
// implementation's blit operation if it has one, else we'll
// just rewrite the cells in linear order.

// The backend implements what should be a fast blit.

// This does math for each cell, so you're looking at a lot of multiplications
// for a bit display -- 100x100 means 10,000x2 multiplications.  This could be
// optimized, but this is really a fallback as every backend really *should* have
// an efficient blit operation. (The ones that don't probably don't keep their own
// state, such as a wrappers on top of TTYs.)  In these cases the cost of writing
// the content is probably substantially dominant anyway.

func (em *emulator) scrollUp() { _ = "STUB: not implemented"; return }

// TODO: deal with wide characters broken across the margin

func (em *emulator) scrollDown() { _ = "STUB: not implemented"; return }

// nextTab advances to the next tab stop, or the end of
// the line if there is no further tab.
func (em *emulator) nextTab() { _ = "STUB: not implemented"; return }

// already at end

// just advance to the next one

func (em *emulator) prevTab() { _ = "STUB: not implemented"; return }

// initTabStops initializes the tab stops assuming every 8th column
// is a tab stop.  This should only be called if the user is intentionally
// changing the tab stops, because it will no longer support expanding
// tab stops on resizing.
func (em *emulator) initTabStops() { _ = "STUB: not implemented"; return }

// no tab stop at offset 0 since that would be pointless

// setTabStop sets a tab stop at the given location.
// This calls  initTabStops - please see the description of that function for ramifications.
func (em *emulator) setTabStop(ts Col) { _ = "STUB: not implemented"; return }

// clrTabStop clears the tab stop at the given column.  This calls
// initTabStops - please see the description of that function for ramifications.
func (em *emulator) clrTabStop(ts Col) { _ = "STUB: not implemented"; return }

// index obtains the index in the cells slice for the given coordinates,
// which must be within the bounds of the display size.
func (em *emulator) index(c Coord) int { _ = "STUB: not implemented"; return 0 }

// putRune puts out a single rune.  This might be a subsequent part of a grapheme cluster, in
// which case it will be emitted together with the preceding base character.
func (em *emulator) putRune(r rune) { _ = "STUB: not implemented"; return }

// ASCII-to-ASCII pairs cannot extend a grapheme cluster, except CRLF.

// fall through to the normal single-rune path

// maybe we need to update the last index

// we are adding to a cluster

// we may have to move position if this switches to wide, so recalculate expected end

// erase the next cell before putting down a character

// we leave the em.lastIndex for now, we might keep extending this cluster

// Advance the cursor. This will stop at the margin.
// Note that if auto margin is enabled, we will have set
// autoWrap above if we were at the margin already.

func (em *emulator) runeString(r rune) string { _ = "STUB: not implemented"; return "" }

func (em *emulator) clusterString(cluster []byte) string { _ = "STUB: not implemented"; return "" }

func shouldCheckGrapheme(prev byte, r rune) bool { _ = "STUB: not implemented"; return false }

func isRegionalIndicator(r rune) bool { _ = "STUB: not implemented"; return false }

// eraseCell erases a single cell at the given offset.
// It clears attributes, but leaves the colors intact.
func (em *emulator) eraseCell(c Coord) { _ = "STUB: not implemented"; return }

// eraseBelow erases from (and including) the current cursor position to the end of the window.
func (em *emulator) eraseBelow() { _ = "STUB: not implemented"; return }

// eraseAbove erases from the origin to (and including) the current cursor position.
func (em *emulator) eraseAbove() { _ = "STUB: not implemented"; return }

// eraseAll erases the entire screen. It uses the color, but resets all other attributes.
func (em *emulator) eraseAll() { _ = "STUB: not implemented"; return }

// eraseToLineEnd erases to the end of the line, including the cursor position.
func (em *emulator) eraseToLineEnd() { _ = "STUB: not implemented"; return }

// eraseToLineStart erases to the start of the line, including the cursor position.
func (em *emulator) eraseToLineStart() { _ = "STUB: not implemented"; return }

// eraseLine erases the entire line.
func (em *emulator) eraseLine() { _ = "STUB: not implemented"; return }

// softReset performs a soft reset.
func (em *emulator) softReset() {
	_ = "STUB: not implemented"
	// TODO:
	// Select default character sets
	return
}

// start by resetting all modes

// NB: No effect for non-changeable modes

// NB: No effect for non-changeable modes

// and set any that should reset on (auto-margin)

// set default cursor - matches VT defaults

func (em *emulator) ansiModeKeys() []AnsiMode { _ = "STUB: not implemented"; return nil }

func (em *emulator) privateModeKeys() []PrivateMode { _ = "STUB: not implemented"; return nil }

// sendDA ends the primary device attributes.
func (em *emulator) sendDA() { _ = "STUB: not implemented"; return }

// 9 for NRC?
// 15 for graphics?

// setAnsiMode sets the ANSI mode.
func (em *emulator) setAnsiMode(mode AnsiMode, ms ModeStatus) { _ = "STUB: not implemented"; return }

func (em *emulator) getAnsiMode(mode AnsiMode) ModeStatus {
	_ = "STUB: not implemented"
	return *new(ModeStatus)
}

// getPrivateMode returns the value of a DEC private mode.
func (em *emulator) getPrivateMode(pm PrivateMode) ModeStatus {
	_ = "STUB: not implemented"
	return *new(ModeStatus)
}

func (em *emulator) updateMouseReporting() { _ = "STUB: not implemented"; return }

// setPrivateMode sets the DEC private mode.
func (em *emulator) setPrivateMode(pm PrivateMode, ms ModeStatus) {
	_ = "STUB: not implemented"
	return
}

func (em *emulator) mouseReportingLocked() MouseReporting {
	_ = "STUB: not implemented"
	return *new(MouseReporting)
}

// SendRaw allows raw data to be sent to the application.
// This is done in a thread-safe way, so that content is not intermingled.
func (em *emulator) SendRaw(b []byte) { _ = "STUB: not implemented"; return }

// Do not attempt to send *anything* if we are stopped.

// Try to write to the readQ, but if we cannot, then wait until
// either we can, or the stopQ is fired.  This ensures that we avoid
// breaking up content if at all possible.

// KeyEvent injects a keyboard event into the emulator
func (em *emulator) KeyEvent(ev KeyEvent) { _ = "STUB: not implemented"; return }

// eliminate "control" keys (which keyboard maps provide) from consideration.
// (We handle control keys explicitly.)

// TODO: more add support for kitty, and maybe modify other keys

// ResizeEvent is called by the backend when a resize occurs.  A real backend with a child
// process (essentially a "real emulator") should probably also fire SIGWINCH if appropriate.
// That would be the job of something other than this code.
func (em *emulator) ResizeEvent(size Coord) { _ = "STUB: not implemented"; return }

func (em *emulator) applyResize(size Coord) {
	_ = "STUB: not implemented"
	// resize clobbers our content, until it is redrawn
	return
}

// resizing resets the margins

// NB: we never support "ModeOnLocked"
// NB: for now we do not support pixel sizes

// Send a SIGWINCH or similar.

var legacyKeys = map[Key]struct {
	K  string // unmodified key
	A  string // unmodified in application cursor mode (smkx)
	S  string // with shift (if empty use regular modifier)
	C  string // with control (if empty use regular modifier)
	CS string // with ctrl-shift
}{
	KeyF1:        {K: "\x1bOP"}, // SS3 P
	KeyF2:        {K: "\x1bOQ"}, // SS3 Q
	KeyF3:        {K: "\x1bOR"}, // SS3 R
	KeyF4:        {K: "\x1bOS"}, // SS3 S
	KeyF5:        {K: "\x1b[15~"},
	KeyF6:        {K: "\x1b[17~"},
	KeyF7:        {K: "\x1b[18~"},
	KeyF8:        {K: "\x1b[19~"},
	KeyF9:        {K: "\x1b[20~"},
	KeyF10:       {K: "\x1b[21~"},
	KeyF11:       {K: "\x1b[23~"},
	KeyF12:       {K: "\x1b[24~"},
	KeyF13:       {K: "\x1b[25~"},
	KeyF14:       {K: "\x1b[26~"},
	KeyF15:       {K: "\x1b[28~"},
	KeyF16:       {K: "\x1b[29~"},
	KeyF17:       {K: "\x1b[31~"},
	KeyF18:       {K: "\x1b[32~"},
	KeyF19:       {K: "\x1b[33~"},
	KeyF20:       {K: "\x1b[34~"},
	KeyUp:        {K: "\x1b[A", A: "\x1bOA"},
	KeyDown:      {K: "\x1b[B", A: "\x1bOB"},
	KeyRight:     {K: "\x1b[C", A: "\x1bOC"},
	KeyLeft:      {K: "\x1b[D", A: "\x1bOD"},
	KeyHome:      {K: "\x1b[H", A: "\x1bOH"},
	KeyEnd:       {K: "\x1b[F", A: "\x1bOF"},
	KeyPgUp:      {K: "\x1b[5~"},
	KeyPgDn:      {K: "\x1b[6~"},
	KeyDelete:    {K: "\x1b[3~"},
	KeyInsert:    {K: "\x1b[2~"},
	KeyMenu:      {K: "\x1b[29~"}, // also F16
	KeyTab:       {K: "\t", S: "\x1b[Z", CS: "\x1b[Z"},
	KeyBackspace: {K: "\x7f", S: "\x7f", C: "\x08", CS: "\x08"},
	KeySpace:     {K: " ", S: " ", C: "\x00", CS: "\x00"},
	KeyEnter:     {K: "\r", S: "\r", CS: "\r"}, // NB: consider using kitty encoding here
	KeyPadEnter:  {K: "\r", S: "\r", CS: "\r"}, // NB: consider using kitty encoding here
	KeyEsc:       {K: "\x1b", S: "\x1b", C: "\x1b"},
}

var legacyControls = map[Key]string{
	// These ones are weird legacy control sequences that we mostly
	// do not care about.  We don't include shifted variants.
	Key2:      "\x00",
	Key3:      "\x1b",
	Key4:      "\x1c",
	Key5:      "\x1d",
	Key6:      "\x1e",
	Key7:      "\x1f",
	Key8:      "\x7f",
	KeyLBrace: "\x1b",
	KeySlash:  "\x1c",
	KeyRBrace: "\x1d",
}

// legacyPadKeys are keys that are on the keypad, when not in numeric keypad mode.
// Note that num lock overrides this.
var legacyPadKeys = map[Key]struct {
	app string
	num string
}{
	KeyPadEnter: {"\x1bOM", "\r"},
	KeyPadMul:   {"\x1bOj", "*"},
	KeyPadAdd:   {"\x1bOk", "+"},
	KeyPadSub:   {"\x1bOm", "-"},
	KeyPadDiv:   {"\x1bOo", "/"},
	KeyPadDec:   {"\x1b[3~", "."}, // Del
	KeyPad0:     {"\x1b[2~", "0"}, // Ins
	KeyPad1:     {"\x1bOF", "1"},  // End
	KeyPad2:     {"\x1b[B", "2"},  // Down
	KeyPad3:     {"\x1b[6~", "3"}, // PgDn
	KeyPad4:     {"\x1b[D", "4"},  // Left
	KeyPad5:     {"\x1b[E", "5"},  // Clear/Begin
	KeyPad6:     {"\x1b[C", "6"},  // Right
	KeyPad7:     {"\x1bOH", "7"},  // Home
	KeyPad8:     {"\x1b[A", "8"},  // Up
	KeyPad9:     {"\x1b[5~", "9"}, // PgUp
	KeyPadEqual: {"\x1bOX", "="},
}

// repeatRaw is called to provide key repeat.  We limit key repeating to just 40,
// and we ensure that at least one is included.  We only repeat if key repeat is enabled.
func (em *emulator) repeatRaw(ev KeyEvent, data []byte) { _ = "STUB: not implemented"; return }

// noRepeatRaw is used to send a key that should never repeat.
// It will only send if the repeat count is zero.
func (em *emulator) noRepeatRaw(ev KeyEvent, data []byte) { _ = "STUB: not implemented"; return }

// keyLegacy handles a keyboard event when in legacy vt220 style mode.
func (em *emulator) keyLegacy(ev KeyEvent) {
	_ = "STUB: not implemented"
	// legacy protocol does not support key release
	return
}

// legacy protocol does not support these

// Shift-Ctrl keys are never sent in the legacy protocol.  We do have to ensure
// that if we are sending other Utf (for example with AltGr), then we still might
// send it, but this is only an issue for non-ASCII runes. Also, this filter only
// applies for "regular" keys (i.e. not function keys, cursor keys, etc.)

// keypad sequences

// For control keys (e.g. control-J) we never emit a rune directly -- but we might later
// add after decoding the key accordingly.

// ASCII might get alt

// otherwise send the UTF as-is

// some weird number control sequences - legacy compatibility
// We do not repeat these.

// AnsiMode 20 sends newline, but only in legacy mode.

// IsCtrl & IsShift

// No specific modifiers present, lets add them. There are two cases,
// one for SS3 based keys and another for CSI based keys.  SS3 based
// keys are converted to CSI - 1 ; mod ; final
// Note: legacy encoding does not use modifiers for alt or super - alt will be
// determined by sending an escape prefix.

// no repeating ALT sequences
// alt sends leading escape

// no repeating CTRL sequences

// but other sequences (should just be shifted or unmodified)
// are fine.  (E.g. we want to allow repeats of cursor keys)

// fallback control key handling

/* ctrl-A */

var win32NoRepeat = map[Key]bool{
	KeyLShift:   true,
	KeyRShift:   true,
	KeyLCtrl:    true,
	KeyRCtrl:    true,
	KeyLAlt:     true,
	KeyRAlt:     true,
	KeyLMeta:    true,
	KeyRMeta:    true,
	KeyCapsLock: true,
	KeyNumLock:  true,
	KeyEnter:    true,
	KeyScrLock:  true,
	KeyPause:    true,
	KeyPrtScr:   true,
}

// keyWin32IM generates the sequence for a key event when in Win32 input mode.
// Win32 input mode is ESC [ Vk ; Sc ; Uc ; Kd ; Cs ; Rc _
// Note that we specifically do NOT doubly encode non-keyboard events -- those
// are already unambiguously handled within the protocol.  (Windows Terminal behaves
// the same way, but most 3rd party terminals do doubly encode.)
func (em *emulator) keyWin32IM(ev KeyEvent) {
	_ = "STUB: not implemented"
	// Some keys that never repeat
	return
}

// Modifiers

// NB: 0x40 is for scroll lock, we don't support it for now

// enhanced

func (em *emulator) MouseEvent(ev MouseEvent) { _ = "STUB: not implemented"; return }

// suppress motion events if the user didn't request

// if entire event was just motion, bail

// Old style reporting (via 1000h).
// Limitations of legacy VT200 reporting are that the coordinates must be between
// 1 and 223 inclusive, and that once any release occurs all buttons are assumed
// to be released.  (Please use SGR mode if at all possible.)
// Further, this mode is not CSI compliant as the encoded values that arrive ahead of
// the final character may be within the range of technically legal CSI final bytes.

// legacy X10 reporting only

// NB: we intentionally reverse buttons 2 & 3 (for xterm compatibility)

func (em *emulator) FocusEvent(focused bool) { _ = "STUB: not implemented"; return }

// SetId sets the terminal name and version.
func (em *emulator) SetId(name string, version string) { _ = "STUB: not implemented"; return }

// Start the terminal emulator.
func (em *emulator) Start() error { _ = "STUB: not implemented"; return nil }

// already running

// Stop the terminal emulator.  This also wakes any blocked
// Read or Write calls, which will return an error.
func (em *emulator) Stop() error { _ = "STUB: not implemented"; return nil }

// Drain pending output to the terminal emulator.
func (em *emulator) Drain() error { _ = "STUB: not implemented"; return nil }

// make sure to wake the reader

// Write data to the emulator (commands).
func (em *emulator) Write(data []byte) (n int, err error) { _ = "STUB: not implemented"; return 0, nil }

// we add the drainQ for synchronization, so that we only
// return after the the emulator has processed this.

// Read data (key events, etc.) from the emulator.
func (em *emulator) Read(data []byte) (n int, err error) { _ = "STUB: not implemented"; return 0, nil }

// The data arriving in the channel may be a byte, or it might be a bool
// trying to force a wake up.  Note that the bool may be intermingled with other
// bytes, so we check it. Also data may have arrived since the bool was posted,
// so make sure we don't terminate until we have collected all the relevant data
// that we can (up to the limit of what was requested.)

func (em *emulator) run(stopQ <-chan bool) { _ = "STUB: not implemented"; return }

// resize notification
