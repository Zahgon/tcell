//go:build plan9
// +build plan9

// Copyright 2025 The TCell Authors
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

package tty

import (
	"os"
	"sync"
	"sync/atomic"
)

// p9Tty implements tcell.Tty using Plan 9's /dev/cons and /dev/consctl.
// Raw mode is enabled by writing "rawon" to /dev/consctl while the fd stays open.
// Resize notifications are read from /dev/wctl: the first read returns geometry,
// subsequent reads block until the window changes (rio(4)).
//
// References:
// - kbdfs(8): cons/consctl rawon|rawoff semantics
// - rio(4): wctl geometry and blocking-on-change behavior
// - vt(1): VT100 emulator typically used for TUI programs on Plan 9
//
// Limitations:
//   - We assume VT100-level capabilities (often no colors, no mouse).
//   - Window size is conservative: we return 80x24 unless overridden.
//     Set LINES/COLUMNS (or TCELL_LINES/TCELL_COLS) to refine.
//   - Mouse and bracketed paste are not wired; terminfo/xterm queries
//     are not attempted because vt(1) may not support them.
type p9Tty struct {
	cons    *os.File // /dev/cons (read+write)
	consctl *os.File // /dev/consctl (write "rawon"/"rawoff")
	wctl    *os.File // /dev/wctl (resize notifications)

	closed  atomic.Bool
	started bool

	onResize atomic.Value // resize channel
	wg       sync.WaitGroup
	stopCh   chan struct{}
}

func NewDevTty() (Tty, error) {
	_ = "STUB: not implemented" // tcell signature
	return *new(Tty), nil
}

func NewStdIoTty() (Tty, error) {
	_ = "STUB: not implemented" // also required by tcell
	return *
	// On Plan 9 there is no POSIX tty discipline on stdin/stdout;
	// use /dev/cons explicitly for robustness.
	new(Tty), nil
}

func NewDevTtyFromDev(_ string) (Tty, error) {
	_ = "STUB: not implemented" // required by tcell
	// Plan 9 does not have multiple "ttys" in the POSIX sense;
	// always bind to /dev/cons and /dev/consctl.
	return *new(Tty), nil
}

func newPlan9TTY() (Tty, error) { _ = "STUB: not implemented"; return *new(Tty), nil }

// /dev/wctl may not exist (console without rio); best-effort.

func (t *p9Tty) Start() error { _ = "STUB: not implemented"; return nil }

// Recreate stop channel if absent or closed (supports resume).

// Put console into raw mode; remains active while consctl is open.

// Reopen /dev/wctl on resume; best-effort (system console may lack it).

func (t *p9Tty) Drain() error {
	_ = "STUB: not implemented"
	// Per tcell docs, this may reasonably be a no-op on non-POSIX ttys.
	// Read deadlines are not available on plan9 os.File; we rely on Stop().
	return nil
}

func (t *p9Tty) Stop() error {
	_ = "STUB: not implemented"

	// Signal watcher to stop (if not already).
	return nil
}

// Exit raw mode first.

// Closing wctl unblocks watchResize; nil it so Start() can reopen later.

// Ensure watcher goroutine has exited before returning.

func (t *p9Tty) Close() error { _ = "STUB: not implemented"; return nil }

func (t *p9Tty) Read(p []byte) (int, error) { _ = "STUB: not implemented"; return 0, nil }

func (t *p9Tty) Write(p []byte) (int, error) { _ = "STUB: not implemented"; return 0, nil }

func (t *p9Tty) NotifyResize(resizeQ chan<- bool) { _ = "STUB: not implemented"; return }

func (t *p9Tty) WindowSize() (WindowSize, error) {
	_ = "STUB: not implemented"
	// Strategy:
	// 1) honor explicit overrides (TCELL_LINES/TCELL_COLS, LINES/COLUMNS),
	// 2) otherwise return conservative 80x24.
	// Reading /dev/wctl gives pixel geometry, but char cell metrics are
	// not generally available to non-draw clients; vt(1) is fixed-cell.
	return *new(WindowSize), nil
}

// watchResize blocks on /dev/wctl reads; each read returns when the window
// changes size/position/state, per rio(4). We ignore the parsed geometry and
// just notify tcell to re-query WindowSize().
func (t *p9Tty) watchResize() { _ = "STUB: not implemented"; return }

// Each read delivers something like:
// "   minx        miny        maxx        maxy   visible current\n"
// We don't need to parse here; just signal.

// transient errors: continue

func envInt(name string) int { _ = "STUB: not implemented"; return 0 }

// helper: safe check if a channel is closed
func isClosed(ch <-chan struct{}) bool { _ = "STUB: not implemented"; return false }
