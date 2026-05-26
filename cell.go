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

package tcell

type cell struct {
	currStr   string
	lastStr   string
	currStyle Style
	lastStyle Style
	width     int
	lock      bool
}

func (c *cell) setDirty(dirty bool) {
	_ = "STUB: not implemented"

	// Empty cells use currStr == "" until they are first drawn, at which
	// point SetDirty(false) normalizes them to a space.  Using "" as the
	// dirty marker for an untouched empty cell would therefore leave
	// lastStr == currStr and fail to force a redraw.
	return
}

// CellBuffer represents a two-dimensional array of character cells.
// This is primarily intended for use by Screen implementers; it
// contains much of the common code they need.  To create one, just
// declare a variable of its type; no explicit initialization is necessary.
//
// CellBuffer is not thread safe.
type CellBuffer struct {
	w               int
	h               int
	cells           []cell
	sanitizeContent bool
}

// Put a single styled grapheme using the given string and style
// at the same location.  Note that only the first grapheme in the string
// will be displayed, using only the 1 or 2 (depending on width) cells
// located at x, y. It returns the rest of the string, and the width used.
func (cb *CellBuffer) Put(x int, y int, str string, style Style) (string, int) {
	_ = "STUB: not implemented"
	return "", 0
}

func (cb *CellBuffer) put(x int, y int, str string, style Style) (string, int) {
	_ = "STUB: not implemented"
	return "", 0
}

// Identical re-Put (a full-screen redraw): the grapheme split is
// unchanged, so reuse the measured width instead of segmenting.

// Wide characters: we want to mark the "wide" cells
// dirty as well as the base cell, to make sure we consider
// both cells as dirty together.  We only need to do this
// if we're changing content

// Prevent unnecessary bounds checks for first cell, since we already
// received that one.

// Get the contents of a character cell (or two adjacent cells), including the
// the style and the display width in cells.  (The width can be either 1, normally,
// or 2 for East Asian full-width characters.  If the width is 0, then the cell is
// is empty.)
func (cb *CellBuffer) Get(x, y int) (string, Style, int) {
	_ = "STUB: not implemented"
	return "", *new(Style), 0
}

// Size returns the (width, height) in cells of the buffer.
func (cb *CellBuffer) Size() (int, int) {
	_ = "STUB: not implemented"

	// Invalidate marks all characters within the buffer as dirty.
	return 0, 0
}

func (cb *CellBuffer) Invalidate() { _ = "STUB: not implemented"; return }

// Dirty checks if a character at the given location needs to be
// refreshed on the physical display.  This returns true if the cell
// content is different since the last time it was marked clean.
func (cb *CellBuffer) Dirty(x, y int) bool { _ = "STUB: not implemented"; return false }

// SetDirty is normally used to indicate that a cell has
// been displayed (in which case dirty is false), or to manually
// force a cell to be marked dirty.
func (cb *CellBuffer) SetDirty(x, y int, dirty bool) { _ = "STUB: not implemented"; return }

// LockCell locks a cell from being drawn, effectively marking it "clean" until
// the lock is removed. This can be used to prevent tcell from drawing a given
// cell, even if the underlying content has changed. For example, when drawing a
// sixel graphic directly to a TTY screen an implementer must lock the region
// underneath the graphic to prevent tcell from drawing on top of the graphic.
func (cb *CellBuffer) LockCell(x, y int) { _ = "STUB: not implemented"; return }

// UnlockCell removes a lock from the cell and marks it as dirty
func (cb *CellBuffer) UnlockCell(x, y int) { _ = "STUB: not implemented"; return }

// Resize is used to resize the cells array, with different dimensions,
// while preserving the original contents.  The cells will be invalidated
// so that they can be redrawn.
func (cb *CellBuffer) Resize(w, h int) { _ = "STUB: not implemented"; return }

// Fill fills the entire cell buffer array with the specified character
// and style.  Normally choose ' ' to clear the screen.  This API doesn't
// support combining characters, or characters with a width larger than one.
// If either the foreground or background are ColorNone, then the respective
// color is unchanged.
func (cb *CellBuffer) Fill(r rune, style Style) { _ = "STUB: not implemented"; return }
