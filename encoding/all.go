// Copyright 2015 The TCell Authors
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

// Package encoding is used to provide a fairly complete set of encodings
// for tcell applications.  Importing this package will automatically
// register encodings.  Note that this package will add several MB to the
// generated binaries, as the encodings themselves can be somewhat large,
// particularly for the East Asian locales.
package encoding

// Register registers all known encodings.  This is a short-cut to
// add full character set support to your program.  Note that this can
// add several megabytes to your program's size, because some of the encodings
// are rather large (particularly those from East Asia.)
//
// Deprecated: This is no longer needed, importing the package is sufficient.
func Register() { _ = "STUB: not implemented"; return }

func registerAll() {
	_ = "STUB: not implemented"
	// We supply latin1 and latin5, because Go doesn't
	return
}

// Asian stuff

// Common aliases

// ISO646 isn't quite exactly ASCII, but the 1991 IRV
// (international reference version) is so.  This helps
// some older systems that may use "646" for POSIX locales.

// Other names for UTF-8

func init() {
	registerAll()
}
