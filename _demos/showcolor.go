//go:build ignore
// +build ignore

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

// colors just displays a single centered rectangle that should pulse
// through available colors.  It uses the RGB color cube, bumping at
// predefined larger intervals (values of about 8) in order that the
// changes happen quickly enough to be appreciable.
//
// Press ESC to exit the program.
package main

import (
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

var red = int32(rand.Int() % 256)
var grn = int32(rand.Int() % 256)
var blu = int32(rand.Int() % 256)
var inc = int32(8) // rate of color change
var redi = int32(inc)
var grni = int32(inc)
var blui = int32(inc)

func makebox(s tcell.Screen, name string, fill color.Color) { _ = "STUB: not implemented"; return }

func fatal(v string, args ...any) { _ = "STUB: not implemented"; return }

func main() {

	if len(os.Args) != 2 {
		fatal("usage: %s color", os.Args[0])
	}

	var color_ color.Color

	if strings.HasPrefix(os.Args[1], "#") {
		name := os.Args[1][1:]
		if len(name) == 3 {
			hex, err := strconv.ParseUint(name, 16, 12)
			if err != nil {
				fatal("invalid color specification")
			}
			// expand 12 bit color to 24 bit
			r := int32((hex & 0xf00) >> 8)
			g := int32((hex & 0x0f0) >> 4)
			b := int32(hex & 0x00f)
			r = r<<4 + r
			g = g<<4 + g
			b = b<<4 + b
			color_ = color.NewRGBColor(r, g, b)
		} else if len(name) == 6 {
			hex, err := strconv.ParseUint(name, 16, 32)
			if err != nil {
				fatal("invalid color specification")
			}
			color_ = color.NewHexColor(int32(hex))
		} else {
			fatal("invalid color specification")
		}
	} else if val, err := strconv.Atoi(os.Args[1]); err == nil {
		color_ = color.Black + color.Color(val)
	} else if val, ok := color.Names[strings.ToLower(os.Args[1])]; ok {
		color_ = val
	} else {
		fatal("invalid color specification")
	}

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fatal("new screen: %v", e.Error())
	}
	if e = s.Init(); e != nil {
		fatal("new screen: %v", e.Error())
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(color.Black).
		Background(color.White))
	s.Clear()

	quit := make(chan struct{})
	go func() {
		for {
			ev := <-s.EventQ()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	makebox(s, os.Args[1], color_)
	<-quit

	s.Fini()
}
