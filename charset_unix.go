//go:build unix
// +build unix

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

func getCharset() string {
	_ = "STUB: not implemented"
	// Determine the character set.  This can help us later.
	// Per POSIX, we search for LC_ALL first, then LC_CTYPE, and
	// finally LANG.  First one set wins.
	return ""
}

// Default assumption, and on Linux we can see LC_ALL
// without a character set, which we assume implies UTF-8.
