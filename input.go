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

// This file describes a generic VT input processor.  It parses key sequences,
// (input bytes) and loads them into events.  It expects UTF-8 or UTF-16 as the input
// feed, along with ECMA-48 sequences.  The assumption here is that all potential
// key sequences are unambiguous between terminal variants (analysis of extant terminfo
// data appears to support this conjecture). This allows us to implement  this once,
// in the most efficient and terminal-agnostic way possible.
//
// There is unfortunately *one* conflict, with aixterm, for CSI-P - which is KeyDelete
// in aixterm, but F1 in others.

//go:build (!js && !wasm) || (js && wasm)
// +build !js,!wasm js,wasm

package tcell

import (
	"sync"
	"time"

	"github.com/gdamore/tcell/v3/vt"
)

type inputState int

const (
	istInit = inputState(iota)
	istUtf  // utf8 state
	istEsc  // escape
	istCsi  // control sequence introducer
	istOsc  // operating system command
	istDcs  // device control string
	istSos  // start of string (unused)
	istPm   // privacy message (unused)
	istApc  // application program command
	istSt   // string terminator
	istSs2  // single shift 2
	istSs3  // single shift 3
	istLnx  // linux F-key (not ECMA-48 compliant - bogus CSI)
	istXda  // extended device attributes (ESC P Ps ST)
)

// defaultControlStringLimit caps inbound OSC/XDA control-string payloads
// before they can grow without bound while waiting for a string terminator.
const defaultControlStringLimit = 64 * 1024

func newInputParser(eq chan<- Event) *inputParser { _ = "STUB: not implemented"; return nil }

type inputParser struct {
	buf              []rune       // bytes to process (ingest data)
	utfBuf           []byte       // accrued UTF8 bytes
	strBuf           []byte       // accrued string data (for ST, OSC, etc.)
	csiParams        []byte       // accrued parameter bytes for CSI (and SS3)
	csiInterm        []byte       // accrued intermediate bytes for CSI
	escChar          byte         // last byte for escape
	escaped          bool         // true if next key should be modified by ESC
	btnsDown         ButtonMask   // mouse buttons down (excludes wheel buttons)
	state            inputState   // tracks processor state
	strState         inputState   // saved str state (needed for ST)
	l                sync.Mutex   // protects local state
	evch             chan<- Event // where events are routed
	rows             int          // used for clipping mouse coordinates
	cols             int          // used for clipping mouse coordinates
	pixelMouse       bool         // mouse reports in pixels (CSI ?1016h); skip cell clipping
	keyTime          time.Time    // time of last key press / byte ingested
	nested           *inputParser // for buggy win32-input-mode implementations
	surrogate        rune         // high surrogate pair seen (for Win32 input mode)
	advanced         bool         // use advanced key reporting semantics
	controlStringMax int          // maximum inbound OSC/XDA payload size; 0 means unlimited
	discardString    bool         // drop the rest of an over-limit OSC/XDA sequence
}

func keyFromInt(n int) (Key, bool) { _ = "STUB: not implemented"; return *new(Key), false }

func keyFromRune(r rune) (Key, bool) { _ = "STUB: not implemented"; return *new(Key), false }

func asciiByteFromInt(n int) (byte, bool) { _ = "STUB: not implemented"; return 0, false }

// Waiting returns true if the processor is waiting for
// some more input (i.e. we are not in in the initial state.)
// This can occur when we have ambiguous escape sequences, such
// as the lone escape.  If this is typed, we expect at least a minimal
// inter-key delay before the next stroke occurs, and the caller
// should check for waiting, and call Scan() or ScanUTF8() to
// finish the processing.  (Typically after a delay of around 100ms.)
func (ip *inputParser) Waiting() bool { _ = "STUB: not implemented"; return false }

// SetPixelMouse toggles whether SGR mouse reports are interpreted as
// pixel coordinates (CSI ?1016h) rather than character cells (CSI ?1006h).
// When enabled, mouse coordinates are not clipped to the screen size.
// The setting is also forwarded to the lazily-created nested parser used
// for win32-input-mode, if one exists, so both stay in sync.
func (ip *inputParser) SetPixelMouse(on bool) { _ = "STUB: not implemented"; return }

func (ip *inputParser) SetSize(w, h int) { _ = "STUB: not implemented"; return }

func (ip *inputParser) post(ev Event) { _ = "STUB: not implemented"; return }

func (ip *inputParser) newKey(k Key, str string, mod ModMask, pressed bool, physical Key, repeat int) *EventKey {
	_ = "STUB: not implemented"
	return nil
}

func (ip *inputParser) postKey(k Key, str string, mod ModMask) { _ = "STUB: not implemented"; return }

func (ip *inputParser) postKeyEx(k Key, str string, mod ModMask, pressed bool, physical Key, repeat int) {
	_ = "STUB: not implemented"
	return
}

func (ip *inputParser) postControlKey(r rune, mod ModMask) { _ = "STUB: not implemented"; return }

type csiParamMode struct {
	M rune // Mode
	P int  // Parameter (first)
}

type keyMap struct {
	Key  Key
	Mod  ModMask
	Rune rune
}

var csiAllKeys = map[csiParamMode]keyMap{
	{M: 'A'}:         {Key: KeyUp},
	{M: 'B'}:         {Key: KeyDown},
	{M: 'C'}:         {Key: KeyRight},
	{M: 'D'}:         {Key: KeyLeft},
	{M: 'E'}:         {Key: KeyClear},
	{M: 'F'}:         {Key: KeyEnd},
	{M: 'H'}:         {Key: KeyHome},
	{M: 'L'}:         {Key: KeyInsert},
	{M: 'P'}:         {Key: KeyF1}, // except for aixterm, where this is Delete
	{M: 'Q'}:         {Key: KeyF2},
	{M: 'S'}:         {Key: KeyF4},
	{M: 'Z'}:         {Key: KeyBacktab},
	{M: 'a'}:         {Key: KeyUp, Mod: ModShift},
	{M: 'b'}:         {Key: KeyDown, Mod: ModShift},
	{M: 'c'}:         {Key: KeyRight, Mod: ModShift},
	{M: 'd'}:         {Key: KeyLeft, Mod: ModShift},
	{M: 'q', P: 1}:   {Key: KeyF1}, // all these 'q' are for aixterm
	{M: 'q', P: 2}:   {Key: KeyF2},
	{M: 'q', P: 3}:   {Key: KeyF3},
	{M: 'q', P: 4}:   {Key: KeyF4},
	{M: 'q', P: 5}:   {Key: KeyF5},
	{M: 'q', P: 6}:   {Key: KeyF6},
	{M: 'q', P: 7}:   {Key: KeyF7},
	{M: 'q', P: 8}:   {Key: KeyF8},
	{M: 'q', P: 9}:   {Key: KeyF9},
	{M: 'q', P: 10}:  {Key: KeyF10},
	{M: 'q', P: 11}:  {Key: KeyF11},
	{M: 'q', P: 12}:  {Key: KeyF12},
	{M: 'q', P: 13}:  {Key: KeyF13},
	{M: 'q', P: 14}:  {Key: KeyF14},
	{M: 'q', P: 15}:  {Key: KeyF15},
	{M: 'q', P: 16}:  {Key: KeyF16},
	{M: 'q', P: 17}:  {Key: KeyF17},
	{M: 'q', P: 18}:  {Key: KeyF18},
	{M: 'q', P: 19}:  {Key: KeyF19},
	{M: 'q', P: 20}:  {Key: KeyF20},
	{M: 'q', P: 21}:  {Key: KeyF21},
	{M: 'q', P: 22}:  {Key: KeyF22},
	{M: 'q', P: 23}:  {Key: KeyF23},
	{M: 'q', P: 24}:  {Key: KeyF24},
	{M: 'q', P: 25}:  {Key: KeyF25},
	{M: 'q', P: 26}:  {Key: KeyF26},
	{M: 'q', P: 27}:  {Key: KeyF27},
	{M: 'q', P: 28}:  {Key: KeyF28},
	{M: 'q', P: 29}:  {Key: KeyF29},
	{M: 'q', P: 30}:  {Key: KeyF30},
	{M: 'q', P: 31}:  {Key: KeyF31},
	{M: 'q', P: 32}:  {Key: KeyF32},
	{M: 'q', P: 33}:  {Key: KeyF33},
	{M: 'q', P: 34}:  {Key: KeyF34},
	{M: 'q', P: 35}:  {Key: KeyF35},
	{M: 'q', P: 36}:  {Key: KeyF36},
	{M: 'q', P: 144}: {Key: KeyClear},
	{M: 'q', P: 146}: {Key: KeyEnd},
	{M: 'q', P: 150}: {Key: KeyPgUp},
	{M: 'q', P: 154}: {Key: KeyPgDn},
	{M: 'z', P: 214}: {Key: KeyHome},
	{M: 'z', P: 216}: {Key: KeyPgUp},
	{M: 'z', P: 220}: {Key: KeyEnd},
	{M: 'z', P: 222}: {Key: KeyPgDn},
	{M: 'z', P: 224}: {Key: KeyF1},
	{M: 'z', P: 225}: {Key: KeyF2},
	{M: 'z', P: 226}: {Key: KeyF3},
	{M: 'z', P: 227}: {Key: KeyF4},
	{M: 'z', P: 228}: {Key: KeyF5},
	{M: 'z', P: 229}: {Key: KeyF6},
	{M: 'z', P: 230}: {Key: KeyF7},
	{M: 'z', P: 231}: {Key: KeyF8},
	{M: 'z', P: 232}: {Key: KeyF9},
	{M: 'z', P: 233}: {Key: KeyF10},
	{M: 'z', P: 234}: {Key: KeyF11},
	{M: 'z', P: 235}: {Key: KeyF12},
	{M: 'z', P: 247}: {Key: KeyInsert},
	{M: '^', P: 1}:   {Key: KeyHome, Mod: ModCtrl},
	{M: '^', P: 2}:   {Key: KeyInsert, Mod: ModCtrl},
	{M: '^', P: 3}:   {Key: KeyDelete, Mod: ModCtrl},
	{M: '^', P: 4}:   {Key: KeyEnd, Mod: ModCtrl},
	{M: '^', P: 5}:   {Key: KeyPgUp, Mod: ModCtrl},
	{M: '^', P: 6}:   {Key: KeyPgDn, Mod: ModCtrl},
	{M: '^', P: 7}:   {Key: KeyHome, Mod: ModCtrl},
	{M: '^', P: 8}:   {Key: KeyEnd, Mod: ModCtrl},
	{M: '^', P: 11}:  {Key: KeyF23},
	{M: '^', P: 12}:  {Key: KeyF24},
	{M: '^', P: 13}:  {Key: KeyF25},
	{M: '^', P: 14}:  {Key: KeyF26},
	{M: '^', P: 15}:  {Key: KeyF27},
	{M: '^', P: 17}:  {Key: KeyF28}, // 16 is a gap
	{M: '^', P: 18}:  {Key: KeyF29},
	{M: '^', P: 19}:  {Key: KeyF30},
	{M: '^', P: 20}:  {Key: KeyF31},
	{M: '^', P: 21}:  {Key: KeyF32},
	{M: '^', P: 23}:  {Key: KeyF33}, // 22 is a gap
	{M: '^', P: 24}:  {Key: KeyF34},
	{M: '^', P: 25}:  {Key: KeyF35},
	{M: '^', P: 26}:  {Key: KeyF36}, // 27 is a gap
	{M: '^', P: 28}:  {Key: KeyF37},
	{M: '^', P: 29}:  {Key: KeyF38}, // 30 is a gap
	{M: '^', P: 31}:  {Key: KeyF39},
	{M: '^', P: 32}:  {Key: KeyF40},
	{M: '^', P: 33}:  {Key: KeyF41},
	{M: '^', P: 34}:  {Key: KeyF42},
	{M: '@', P: 23}:  {Key: KeyF43},
	{M: '@', P: 24}:  {Key: KeyF44},
	{M: '@', P: 1}:   {Key: KeyHome, Mod: ModShift | ModCtrl},
	{M: '@', P: 2}:   {Key: KeyInsert, Mod: ModShift | ModCtrl},
	{M: '@', P: 3}:   {Key: KeyDelete, Mod: ModShift | ModCtrl},
	{M: '@', P: 4}:   {Key: KeyEnd, Mod: ModShift | ModCtrl},
	{M: '@', P: 5}:   {Key: KeyPgUp, Mod: ModShift | ModCtrl},
	{M: '@', P: 6}:   {Key: KeyPgDn, Mod: ModShift | ModCtrl},
	{M: '@', P: 7}:   {Key: KeyHome, Mod: ModShift | ModCtrl},
	{M: '@', P: 8}:   {Key: KeyEnd, Mod: ModShift | ModCtrl},
	{M: '$', P: 1}:   {Key: KeyHome, Mod: ModShift},
	{M: '$', P: 2}:   {Key: KeyInsert, Mod: ModShift},
	{M: '$', P: 3}:   {Key: KeyDelete, Mod: ModShift},
	{M: '$', P: 5}:   {Key: KeyPgUp, Mod: ModShift},
	{M: '$', P: 6}:   {Key: KeyPgDn, Mod: ModShift},
	{M: '$', P: 7}:   {Key: KeyHome, Mod: ModShift},
	{M: '$', P: 8}:   {Key: KeyEnd, Mod: ModShift},
	{M: '$', P: 23}:  {Key: KeyF21},
	{M: '$', P: 24}:  {Key: KeyF22},
	{M: '~', P: 1}:   {Key: KeyHome},
	{M: '~', P: 2}:   {Key: KeyInsert},
	{M: '~', P: 3}:   {Key: KeyDelete},
	{M: '~', P: 4}:   {Key: KeyEnd},
	{M: '~', P: 5}:   {Key: KeyPgUp},
	{M: '~', P: 6}:   {Key: KeyPgDn},
	{M: '~', P: 7}:   {Key: KeyHome},
	{M: '~', P: 8}:   {Key: KeyEnd},
	{M: '~', P: 11}:  {Key: KeyF1},
	{M: '~', P: 12}:  {Key: KeyF2},
	{M: '~', P: 13}:  {Key: KeyF3},
	{M: '~', P: 14}:  {Key: KeyF4},
	{M: '~', P: 15}:  {Key: KeyF5},
	{M: '~', P: 17}:  {Key: KeyF6},
	{M: '~', P: 18}:  {Key: KeyF7},
	{M: '~', P: 19}:  {Key: KeyF8},
	{M: '~', P: 20}:  {Key: KeyF9},
	{M: '~', P: 21}:  {Key: KeyF10},
	{M: '~', P: 23}:  {Key: KeyF11},
	{M: '~', P: 24}:  {Key: KeyF12},
	{M: '~', P: 25}:  {Key: KeyF13},
	{M: '~', P: 26}:  {Key: KeyF14},
	{M: '~', P: 28}:  {Key: KeyF15}, // aka KeyHelp
	{M: '~', P: 29}:  {Key: KeyF16},
	{M: '~', P: 31}:  {Key: KeyF17},
	{M: '~', P: 32}:  {Key: KeyF18},
	{M: '~', P: 33}:  {Key: KeyF19},
	{M: '~', P: 34}:  {Key: KeyF20},
	{M: '~', P: 200}: {Key: keyPasteStart},
	{M: '~', P: 201}: {Key: keyPasteEnd},
}

// keys reported using Kitty csi-u protocol
var csiUKeys = map[int]keyMap{
	27:    {Key: KeyESC},
	9:     {Key: KeyTAB},
	13:    {Key: KeyEnter},
	127:   {Key: KeyBS},
	57358: {Key: KeyCapsLock},
	57359: {Key: KeyScrollLock},
	57360: {Key: KeyNumLock},
	57361: {Key: KeyPrint},
	57362: {Key: KeyPause},
	57363: {Key: KeyMenu},
	57376: {Key: KeyF13},
	57377: {Key: KeyF14},
	57378: {Key: KeyF15},
	57379: {Key: KeyF16},
	57380: {Key: KeyF17},
	57381: {Key: KeyF18},
	57382: {Key: KeyF19},
	57383: {Key: KeyF20},
	57384: {Key: KeyF21},
	57385: {Key: KeyF22},
	57386: {Key: KeyF23},
	57387: {Key: KeyF24},
	57388: {Key: KeyF25},
	57389: {Key: KeyF26},
	57390: {Key: KeyF27},
	57391: {Key: KeyF28},
	57392: {Key: KeyF29},
	57393: {Key: KeyF30},
	57394: {Key: KeyF31},
	57395: {Key: KeyF32},
	57396: {Key: KeyF33},
	57397: {Key: KeyF34},
	57398: {Key: KeyF35},
	57399: {Key: KeyRune, Rune: '0'}, // KP 0
	57400: {Key: KeyRune, Rune: '1'}, // KP 1
	57401: {Key: KeyRune, Rune: '2'}, // KP 2
	57402: {Key: KeyRune, Rune: '3'}, // KP 3
	57403: {Key: KeyRune, Rune: '4'}, // KP 4
	57404: {Key: KeyRune, Rune: '5'}, // KP 5
	57405: {Key: KeyRune, Rune: '6'}, // KP 6
	57406: {Key: KeyRune, Rune: '7'}, // KP 7
	57407: {Key: KeyRune, Rune: '8'}, // KP 8
	57408: {Key: KeyRune, Rune: '9'}, // KP 9
	57409: {Key: KeyRune, Rune: '.'}, // KP_DECIMAL
	57410: {Key: KeyRune, Rune: '/'}, // KP_DIVIDE
	57411: {Key: KeyRune, Rune: '*'}, // KP_MULTIPLY
	57412: {Key: KeyRune, Rune: '-'}, // KP_SUBTRACT
	57413: {Key: KeyRune, Rune: '+'}, // KP_ADD
	57414: {Key: KeyEnter},           // KP_ENTER
	57415: {Key: KeyRune, Rune: '='}, // KP_EQUAL
	57416: {Key: KeyClear},           // KP_SEPARATOR
	57417: {Key: KeyLeft},            // KP_LEFT
	57418: {Key: KeyRight},           // KP_RIGHT
	57419: {Key: KeyUp},              // KP_UP
	57420: {Key: KeyDown},            // KP_DOWN
	57421: {Key: KeyPgUp},            // KP_PG_UP
	57422: {Key: KeyPgDn},            // KP_PG_DN
	57423: {Key: KeyHome},            // KP_HOME
	57424: {Key: KeyEnd},             // KP_END
	57425: {Key: KeyInsert},          // KP_INSERT
	57426: {Key: KeyDelete},          // KP_DELETE
	// 57427: {Key: KeyBegin},          // KP_BEGIN
	57441: {Key: KeyShift}, // LEFT_SHIFT
	57442: {Key: KeyCtrl},  // LEFT_CONTROL
	57443: {Key: KeyAlt},   // LEFT_ALT
	57444: {Key: KeyMeta},  // LEFT_SUPER
	57447: {Key: KeyShift}, // RIGHT_SHIFT
	57448: {Key: KeyCtrl},  // RIGHT_CONTROL
	57449: {Key: KeyAlt},   // RIGHT_ALT
	57450: {Key: KeyMeta},  // RIGHT_SUPER

	// TODO: Media keys
}

// windows virtual key codes per microsoft
var winKeys = map[int]Key{
	0x03: KeyCancel,    // vkCancel
	0x08: KeyBackspace, // vkBackspace
	0x09: KeyTab,       // vkTab
	0x0d: KeyEnter,     // vkReturn
	0x13: KeyPause,     // vkPause
	0x1b: KeyEscape,    // vkEscape
	0x21: KeyPgUp,      // vkPrior
	0x22: KeyPgDn,      // vkNext
	0x23: KeyEnd,       // vkEnd
	0x24: KeyHome,      // vkHome
	0x25: KeyLeft,      // vkLeft
	0x26: KeyUp,        // vkUp
	0x27: KeyRight,     // vkRight
	0x28: KeyDown,      // vkDown
	0x2a: KeyPrint,     // vkPrint
	0x2c: KeyPrint,     // vkPrtScr
	0x2d: KeyInsert,    // vkInsert
	0x2e: KeyDelete,    // vkDelete
	0x2f: KeyHelp,      // vkHelp
	0x70: KeyF1,        // vkF1
	0x71: KeyF2,        // vkF2
	0x72: KeyF3,        // vkF3
	0x73: KeyF4,        // vkF4
	0x74: KeyF5,        // vkF5
	0x75: KeyF6,        // vkF6
	0x76: KeyF7,        // vkF7
	0x77: KeyF8,        // vkF8
	0x78: KeyF9,        // vkF9
	0x79: KeyF10,       // vkF10
	0x7a: KeyF11,       // vkF11
	0x7b: KeyF12,       // vkF12
	0x7c: KeyF13,       // vkF13
	0x7d: KeyF14,       // vkF14
	0x7e: KeyF15,       // vkF15
	0x7f: KeyF16,       // vkF16
	0x80: KeyF17,       // vkF17
	0x81: KeyF18,       // vkF18
	0x82: KeyF19,       // vkF19
	0x83: KeyF20,       // vkF20
	0x84: KeyF21,       // vkF21
	0x85: KeyF22,       // vkF22
	0x86: KeyF23,       // vkF23
	0x87: KeyF24,       // vkF24
}

// keys by their SS3 - used in application mode usually (legacy VT-style)
var ss3Keys = map[rune]Key{
	'A': KeyUp,
	'B': KeyDown,
	'C': KeyRight,
	'D': KeyLeft,
	'E': KeyClear,
	'F': KeyEnd,
	'H': KeyHome,
	'P': KeyF1,
	'Q': KeyF2,
	'R': KeyF3,
	'S': KeyF4,
	't': KeyF5,
	'u': KeyF6,
	'v': KeyF7,
	'l': KeyF8,
	'w': KeyF9,
	'x': KeyF10,
}

// linux terminal uses these non ECMA keys prefixed by CSI-[
var linuxFKeys = map[rune]Key{
	'A': KeyF1,
	'B': KeyF2,
	'C': KeyF3,
	'D': KeyF4,
	'E': KeyF5,
}

func (ip *inputParser) scan() { _ = "STUB: not implemented"; return }

// 8-bit extended Unicode we just treat as such - this will swallow anything else queued up

// ISO 2022 control chars

// we fall through so it will be treated as the 7-bit equivalent

// escape.. pending

// Control keys - legacy handling

// no known uses

// string terminator reached, (orphaned?)

// Linux console only, does not conform to ECMA

// leading ESC to capture alt

// treat as alt-key ... legacy emulators only (no CSI-u or other)

// usual case for incoming keys
// NB: rxvt uses terminating '$' which is not a legal CSI terminator,
// for certain shifted key sequences.  We special case this, and it's ok
// because no other terminal seems to use this for CSI intermediates from
// the terminal to the host (queries in the other direction can use it.)
// However, this is only true if the first parameter does not have a "?",
// because it *does* collide with DEC private mode queries otherwise.

// Per ECMA-48 §5.3.1, ESC restarts the escape
// sequence machine from any intermediate state.

// parameter bytes

// rxvt non-standard

// intermediate bytes, rarely used

// final byte

// bad parse, just swallow it all

// No known uses for SS2

// typically application mode keys or older terminals

// some SS3 sequences (old VTE) encode modifiers here just like CSI

// Per ECMA-48 §5.3.1, ESC restarts the escape
// sequence machine from any intermediate state.

// If there are no parameters, then it's simple without modifiers.
// The options for parameters are "1;<modifiers>" , or ";modifiers" (empty
// first parameter defaults to 1), or just <modifiers>.  If a sequence has
// parameters that do not match one of these forms, we just discard it.

// simple SS3 case

// SS3 with modifier (old style).  Note old terminfo would declare these as high
// numbered function keys, but we encode as modified since that's how they are entered.

// these we just eat

// bell - some send this instead of ST

// not sure if used

// linux console does not follow ECMA

// if we take too long between bytes, reset the state machine.

func (ip *inputParser) appendStringBytes(bs ...byte) { _ = "STUB: not implemented"; return }

func (ip *inputParser) handleOsc(str string) { _ = "STUB: not implemented"; return }

func (ip *inputParser) handleXda(str string) { _ = "STUB: not implemented"; return }

// two approaches, one with version like (1.23) another with just spaces

func calcModifier(n int) ModMask { _ = "STUB: not implemented"; return *new(ModMask) }

// kitty calls this Super

// for now not separating from Super

// Not doing (kitty only):
// caps_lock 0b1000000   (64)
// num_lock  0b10000000  (128)

func calcWinModifier(n int, advanced bool) ModMask { _ = "STUB: not implemented"; return *new(ModMask) }

// Bits through 0x0100 match Win32 dwControlKeyState. 0x0040 and
// 0x0080 are ScrollLock and CapsLock, not Meta. The 0x0200 and
// 0x0400 bits are tcell extensions used by the WASM browser shim,
// which has Meta keys but no native Win32 bit assignment for them.

func winModifierKey(vk int) (Key, ModMask, bool) {
	_ = "STUB: not implemented"
	return *new(Key), *new(ModMask), false
}

func kittyModifierKey(code int) ModMask { _ = "STUB: not implemented"; return *new(ModMask) }

func (ip *inputParser) handleMouse(mode rune, params []int) {
	_ = "STUB: not implemented"

	// XTerm mouse events only report at most one button at a time,
	// which may include a wheel button.  Wheel motion events are
	// reported as single impulses, while other button events are reported
	// as separate press & release events.
	return
}

// Some terminals will report mouse coordinates outside the
// screen, especially with click-drag events.  Clip the coordinates
// to the screen in that case.  In pixel-reporting mode (CSI ?1016h)
// the values are already pixels rather than cells, so skip the clip
// and pass them through unchanged for the application to interpret.

// Mouse wheel has bit 6 set, no release events.  It should be noted
// that wheel events are sometimes misdelivered as mouse button events
// during a click-drag, so we debounce these, considering them to be
// button press events unless we see an intervening release event.
// This excludes motion (bit 5) and modifiers (bits 2, 3, 4) for now.

// Note we prefer to treat right as button 2

// And the middle button as button 3

// a release without a corresponding press, so clear it

// Ghostty may send out motion signals that indicate a button has
// been pressed, even when the button is not actually pressed.
// Do not create a synthetic button-down state from these packets.

// record this press

// and use the full set so can see chords

// mice wheel do not have release events

func (ip *inputParser) handleWinKey(P []int) {
	_ = "STUB: not implemented"
	// win32-input-mode
	//
	//	^[ [ Vk ; Sc ; Uc ; Kd ; Cs ; Rc _
	//
	// Vk: the value of wVirtualKeyCode - any number. If omitted, defaults to '0'.
	// Sc: the value of wVirtualScanCode - any number. If omitted, defaults to '0'.
	// Uc: the decimal value of UnicodeChar - for example, NUL is "0", LF is
	//
	//	"10", the character 'A' is "65". If omitted, defaults to '0'.
	//
	// Kd: the value of bKeyDown - either a '0' or '1'. If omitted, defaults to '0'.
	// Cs: the value of dwControlKeyState - any number. If omitted, defaults to '0'.
	// Rc: the value of wRepeatCount - any number. If omitted, defaults to '1'.
	//
	// Note that some 3rd party terminal emulators (not Terminal) suffer from a bug
	// where other events, such as mouse events, are doubly encoded, using Vk 0
	// for each character.  (So a CSI-M sequence is encoded as a series of CSI-_
	// sequences.)  We consider this a bug in those terminal emulators -- Windows 11
	// Terminal does not suffer this brain damage. (We've observed this with both Alacritty
	// and WezTerm.)
	return
}

// ensure sufficient length

// key up event ignore ignore

// these terminals never send ambiguous escapes

// only ASCII in win32-input-mode

// Already decoded.

// high surrogate pair

// low surrogate pair

// Lone modifier releases are ignored unless advanced mode is enabled.

// filter out lone shift for printable chars

// Filter out ctrl+alt (it means AltGr)

func (ip *inputParser) handlePrimaryDA(params []int) { _ = "STUB: not implemented"; return }

func (ip *inputParser) handlePrivateModeResponse(params []int) { _ = "STUB: not implemented"; return }

func (ip *inputParser) handleKittyMode(params []int) { _ = "STUB: not implemented"; return }

func (ip *inputParser) handleXTermMode(params []int) { _ = "STUB: not implemented"; return }

func (ip *inputParser) handleCsi(mode rune, params []byte, intermediate []byte) {
	_ = "STUB: not implemented"

	// reset state
	return
}

// extract numeric parameters

// mouse event, we only do SGR tracking

// we don't know what to do with these for now

// focus in

// focus out

// linux console F-key - CSI-[ modifies next key

// CSI-u kitty keyboard protocol, is unambiguous

// window size report

// window resize report

// apply modifiers if present

// aixterm hack - conflicts with kitty protocol

// this might have been an SS3 style key with modifiers applied

// if we got here we just swallow the unknown sequence

func (ip *inputParser) ScanUTF8(b []byte) { _ = "STUB: not implemented"; return }

// fast path, basic ascii, also includes ISO2022 8-bit controls

// discard the leading byte as bad,
// hopefully it will recover.

// Scan scans the existing input, but does not take new content.
// This is typically called after a delay when Waiting() is true.
func (ip *inputParser) Scan() { _ = "STUB: not implemented"; return }

// Private events between input and tscreen.

// eventPrimaryAttributes is for primary device attributes -- this should be
// the last event returned during initial handshaking
type eventPrimaryAttributes struct {
	EventTime
	Class         int  // Terminal class, 1 is vt100, vt101, 6 is vt102, > 60 for vt200 and up
	ReGIS         bool // Terminal supports ReGIS graphics (DA 3)
	Sixel         bool // Terminal supports Sixel graphics (DA 4)
	National      bool // Terminal supports national replacement character sets (DA 9)
	SerboCroation bool // Serbo-Croatian(DA 12)
	Color         bool // Terminal supports color (DA 22)
	Greek         bool // Greek (DA 23)
	Turkish       bool // Turkish (DA 24)
	Latin2        bool // ISO Latin-2 (DA 42)
	Clipboard     bool // OSC 52 support (DA 52)
}

// eventTermName is for extended attributes
type eventTermName struct {
	EventTime
	Name    string
	Version string
}

type eventPrivateMode struct {
	EventTime
	Mode   vt.PrivateMode // numeric mode e.g. 7 for auto-margin, 1006 for SGR mouse reports, etc
	Status vt.ModeStatus  // value of status
}

type KittyKbdMode uint16

const (
	KittyKbdModeOff       = KittyKbdMode(0)  // Disable Kitty keyboard mode
	KittyKbdModeBase      = KittyKbdMode(1)  // Enable disambiguated keys
	KittyKbdModeEvents    = KittyKbdMode(2)  // Report event types (e.g. key release)
	KittyKbdModeAlternate = KittyKbdMode(4)  // Report alternate keys
	KittyKbdModeAll       = KittyKbdMode(8)  // Report all keys using kitty keyboard protocol
	KittyKbdModeText      = KittyKbdMode(16) // Report associated text
)

type eventKittyKbdMode struct {
	EventTime
	Mode KittyKbdMode
}

type XtermKbdMode uint16

const (
	XtermKbdModeOff  = XtermKbdMode(0) // Disabled
	XtermKbdModeBase = XtermKbdMode(1) // Enabled except for ones with legacy behavior
	XtermKbdModeExt  = XtermKbdMode(2) // Enabled for all modified keys
	XtermKbdModeAll  = XtermKbdMode(3) // Send all keys (including unmodified)
)

type eventXTermKbdMode struct {
	EventTime
	Mode XtermKbdMode
}
