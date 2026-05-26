// Copyright 2026 The TCell Authors
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

package tests

import (
	"testing"

	"github.com/gdamore/tcell/v3/color"
	"github.com/gdamore/tcell/v3/vt"
)

type MockTerm = vt.MockTerm
type Row = vt.Row
type Col = vt.Col
type Coord = vt.Coord
type Attr = vt.Attr

// WriteF writes the string, and ensures it is fully flushed
// before returning.
func WriteF(t *testing.T, term MockTerm, str string, args ...any) {
	_ = "STUB: not implemented"
	return
}

// ReadF reads content from the term and returns it as a string.
func ReadF(t *testing.T, term MockTerm) string { _ = "STUB: not implemented"; return "" }

// VerifyF validates the condition, printing the message on failure.
func VerifyF(t *testing.T, b bool, fmt string, args ...any) { _ = "STUB: not implemented"; return }

// AssertF validates the condition, and aborts the test if it fails.
func AssertF(t *testing.T, b bool, fmt string, args ...any) { _ = "STUB: not implemented"; return }

// MustClose closes or calls Fatalf.
func MustClose(t *testing.T, term MockTerm) { _ = "STUB: not implemented"; return }

// MustStart starts the terminal or calls Fatalf.
func MustStart(t *testing.T, term MockTerm) { _ = "STUB: not implemented"; return }

// CheckPos is verifies the current cursor position of terminal.
func CheckPos(t *testing.T, term MockTerm, x Col, y Row) { _ = "STUB: not implemented"; return }

// CheckContent verifies the content at a given cell of the terminal.
func CheckContent(t *testing.T, term MockTerm, x Col, y Row, s string) {
	_ = "STUB: not implemented"
	return
}

// CheckAttrs verifies the attributes of a given cell.
func CheckAttrs(t *testing.T, term MockTerm, x Col, y Row, a Attr) {
	_ = "STUB: not implemented"
	return
}

// CheckColors verifies the colors of a given cell.
func CheckColors(t *testing.T, term MockTerm, x Col, y Row, fg color.Color, bg color.Color) {
	_ = "STUB: not implemented"
	return
}

// CheckRead verifies that a read matches.
func CheckRead(t *testing.T, term MockTerm, want string) { _ = "STUB: not implemented"; return }
