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

//go:build windows
// +build windows

package tty

import (
	"sync"
	"syscall"
)

var (
	k32 = syscall.NewLazyDLL("kernel32.dll")
)

var (
	procReadConsoleInput              = k32.NewProc("ReadConsoleInputW")
	procGetNumberOfConsoleInputEvents = k32.NewProc("GetNumberOfConsoleInputEvents")
	procFlushConsoleInputBuffer       = k32.NewProc("FlushConsoleInputBuffer")
	procWaitForMultipleObjects        = k32.NewProc("WaitForMultipleObjects")
	procSetConsoleMode                = k32.NewProc("SetConsoleMode")
	procGetConsoleMode                = k32.NewProc("GetConsoleMode")
	procGetConsoleScreenBufferInfo    = k32.NewProc("GetConsoleScreenBufferInfo")
	procCreateEvent                   = k32.NewProc("CreateEventW")
	procSetEvent                      = k32.NewProc("SetEvent")
)

const (
	keyEvent    uint16 = 1
	mouseEvent  uint16 = 2
	resizeEvent uint16 = 4
	menuEvent   uint16 = 8 // don't use
	focusEvent  uint16 = 16
)

const (
	w32Infinite    = ^uintptr(0)
	w32WaitObject0 = uintptr(0)
)

const (
	// Input modes
	modeExtendFlg = uint32(0x0080)
	modeMouseEn   = uint32(0x0010)
	modeResizeEn  = uint32(0x0008)
	modeVtInput   = uint32(0x0200)
	// modeCooked    = uint32(0x0001)

	// Output modes
	modeCookedOut = uint32(0x0001)
	modeVtOutput  = uint32(0x0004)
	modeNoAutoNL  = uint32(0x0008)
	modeUnderline = uint32(0x0010) // ENABLE_LVB_GRID_WORLDWIDE, needed for underlines
	// modeWrapEOL   = uint32(0x0002)
)

type coord struct {
	x int16
	y int16
}

type rect struct {
	left   int16
	top    int16
	right  int16
	bottom int16
}

type consoleInfo struct {
	size  coord
	pos   coord
	attrs uint16
	win   rect
	maxsz coord
}

type inputRecord struct {
	typ  uint16
	_    uint16
	data [16]byte
}

type winTty struct {
	buf        chan byte
	out        syscall.Handle
	in         syscall.Handle
	cancelFlag syscall.Handle
	running    bool
	stopQ      chan struct{}
	resizeQ    chan<- bool
	cols       uint16
	rows       uint16
	pair       []uint16 // for surrogate pairs (UTF-16)
	oimode     uint32   // original input mode
	oomode     uint32   // original output mode
	oscreen    consoleInfo
	wg         sync.WaitGroup
	surrogate  rune
	sync.Mutex
}

func (w *winTty) Read(b []byte) (int, error) {
	_ = "STUB: not implemented"
	// first character read blocks
	return 0, nil
}

// stopping, so make sure we eat everything, which might require
// very short sleeps to ensure all buffered data is consumed.

// second character read is non-blocking

func (w *winTty) Write(b []byte) (int, error) { _ = "STUB: not implemented"; return 0, nil }

func (w *winTty) Close() error { _ = "STUB: not implemented"; return nil }

func (w *winTty) Drain() error { _ = "STUB: not implemented"; return nil }

func (w *winTty) getConsoleInput() error {
	_ = "STUB: not implemented"
	// cancelFlag comes first as WaitForMultipleObjects returns the lowest index
	// in the event that both events are signaled.
	return nil
}

// As arrays are contiguous in memory, a pointer to the first object is the
// same as a pointer to the array itself.

// WaitForMultipleObjects returns WAIT_OBJECT_0 + the index.

// w.cancelFlag

// w.in

// we normally only expect to see ascii, but paste data may come in as UTF-16.

// We normally expect only to see ASCII (win32-input-mode),
// but apparently pasted data can arrive in UTF-16 here.

func (w *winTty) scanInput() { _ = "STUB: not implemented"; return }

func (w *winTty) Start() error { _ = "STUB: not implemented"; return nil }

func (w *winTty) Stop() error { _ = "STUB: not implemented"; return nil }

func (tty *winTty) NotifyResize(resizeQ chan<- bool) { _ = "STUB: not implemented"; return }

func (w *winTty) WindowSize() (WindowSize, error) {
	_ = "STUB: not implemented"
	return *new(WindowSize), nil
}

func NewDevTty() (Tty, error) { _ = "STUB: not implemented"; return *new(Tty), nil }

func NewDevTtyFromDev(dev string) (Tty, error) { _ = "STUB: not implemented"; return *new(Tty), nil }

func NewStdIoTty() (Tty, error) { _ = "STUB: not implemented"; return *new(Tty), nil }
