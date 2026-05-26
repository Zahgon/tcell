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

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris zos

package tty

import (
	"os"

	"golang.org/x/term"
)

// Use -1 to differentiate from stdin (fd=0).
const uninitializedTtyFd = -1

// devTty is an implementation of the Tty API based upon /dev/tty.
type devTty struct {
	fd      int
	f       *os.File
	saved   *term.State
	sig     chan os.Signal
	dev     string
	started bool
}

func (tty *devTty) Read(b []byte) (int, error) { _ = "STUB: not implemented"; return 0, nil }

func (tty *devTty) Write(b []byte) (int, error) { _ = "STUB: not implemented"; return 0, nil }

func (tty *devTty) Close() error { _ = "STUB: not implemented"; return nil }

func (tty *devTty) Start() error { _ = "STUB: not implemented"; return nil }

// We open another copy of /dev/tty.  This is a workaround for unusual behavior
// observed in macOS, apparently caused when a subshell (for example) closes our
// own tty device (when it exits for example).  Getting a fresh new one seems to
// resolve the problem.  (We believe this is a bug in the macOS tty driver that
// fails to account for dup() references to the same file before applying close()
// related behaviors to the tty.)  (Note that when using stdin/stdout instead of
// /dev/tty this problem is not observed.)

// also sets vMin and vTime

func (tty *devTty) Drain() error { _ = "STUB: not implemented"; return nil }

func (tty *devTty) Stop() error {
	_ = "STUB: not implemented"
	// unconditionally set this, because we cannot recover
	// if we fail anyway, so this gives the best hope of
	// picking up the pieces in such a circumstance
	return nil
}

// close our tty device -- we'll get another one if we Start again later.

func (tty *devTty) WindowSize() (WindowSize, error) {
	_ = "STUB: not implemented"
	return *new(WindowSize), nil
}

// If WindowSize is called when the tty isn't yet running, the fd for /dev/tty won't be initialized,
// so open the file just long enough to retrieve the window size.

// default

// default

func (tty *devTty) NotifyResize(resizeQ chan<- bool) { _ = "STUB: not implemented"; return }

// queue full, so nvm.

// NewDevTty opens a /dev/tty based Tty.
func NewDevTty() (Tty, error) { _ = "STUB: not implemented"; return *new(Tty), nil }

// NewDevTtyFromDev opens a tty device given a path.  This can be useful to bind to other nodes.
func NewDevTtyFromDev(dev string) (Tty, error) { _ = "STUB: not implemented"; return *new(Tty), nil }

// Only open the file long enough to check that the device
// represents a TTY.  We will reopen it in start.  We do collect
// the terminal state so we can restore it later though.
