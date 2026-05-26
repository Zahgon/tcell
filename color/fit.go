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

package color

// Find attempts to find a given color, or the best match possible for it,
// from the palette given.  This is an expensive operation, so results should
// be cached by the caller.
func Find(c Color, palette []Color) Color { _ = "STUB: not implemented"; return *new(Color) }

// CIE94 is more accurate, but really really expensive.

// NB: nd < dist is false if is NaN.
// We have never seen a case where the CIE76 algorithm returns NaN.
