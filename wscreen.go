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

//go:build js && wasm
// +build js,wasm

package tcell

import (
	"sync"
	"syscall/js"

	"github.com/gdamore/tcell/v3/tty"
)

// initialize installs the browser-backed TTY used by tScreen on js/wasm.
func (t *tScreen) initialize() error { _ = "STUB: not implemented"; return nil }

func getCharset() string { _ = "STUB: not implemented"; return "" }

type browserTty struct {
	mu      sync.Mutex
	cond    *sync.Cond
	started bool
	drained bool
	closed  bool
	input   []byte
	resizeQ chan<- bool

	writeFunc  js.Value
	sizeFunc   js.Value
	closeFuncs []js.Func
}

func newBrowserTty() *browserTty { _ = "STUB: not implemented"; return nil }

func (t *browserTty) Start() error { _ = "STUB: not implemented"; return nil }

func (t *browserTty) Stop() error { _ = "STUB: not implemented"; return nil }

func (t *browserTty) Drain() error { _ = "STUB: not implemented"; return nil }

func (t *browserTty) NotifyResize(resizeQ chan<- bool) { _ = "STUB: not implemented"; return }

func (t *browserTty) WindowSize() (tty.WindowSize, error) {
	_ = "STUB: not implemented"
	return *new(tty.WindowSize), nil
}

func (t *browserTty) Read(b []byte) (int, error) { _ = "STUB: not implemented"; return 0, nil }

func (t *browserTty) Write(b []byte) (int, error) { _ = "STUB: not implemented"; return 0, nil }

func (t *browserTty) Close() error { _ = "STUB: not implemented"; return nil }

func (t *browserTty) enqueue(data []byte) { _ = "STUB: not implemented"; return }
