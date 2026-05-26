// Copyright 2021 The TCell Authors
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

// stdIoTty is an implementation of the Tty API based upon stdin/stdout.
type stdIoTty struct {
	fd      int
	in      *os.File
	out     *os.File
	saved   *term.State
	sig     chan os.Signal
	started bool
}

func (tty *stdIoTty) Read(b []byte) (int, error) { _ = "STUB: not implemented"; return 0, nil }

func (tty *stdIoTty) Write(b []byte) (int, error) { _ = "STUB: not implemented"; return 0, nil }

func (tty *stdIoTty) Close() error { _ = "STUB: not implemented"; return nil }

func (tty *stdIoTty) Start() error { _ = "STUB: not implemented"; return nil }

// also sets vMin and vTime

func (tty *stdIoTty) Drain() error { _ = "STUB: not implemented"; return nil }

func (tty *stdIoTty) Stop() error { _ = "STUB: not implemented"; return nil }

func (tty *stdIoTty) WindowSize() (WindowSize, error) {
	_ = "STUB: not implemented"
	return *new(WindowSize), nil
}

// default

// default

func (tty *stdIoTty) NotifyResize(resizeQ chan<- bool) { _ = "STUB: not implemented"; return }

// queue full, so nvm.

// NewStdioTty opens a tty using standard input/output.
func NewStdIoTty() (Tty, error) { _ = "STUB: not implemented"; return *new(Tty), nil }
