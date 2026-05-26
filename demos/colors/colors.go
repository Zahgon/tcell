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

// colors just displays a single centered rectangle that should pulse
// through available colors.  It uses the RGB color cube, bumping at
// predefined larger intervals (values of about 8) in order that the
// changes happen quickly enough to be appreciable.
//
// Press ESC to exit the program.
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

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
var interval = time.Millisecond * 50

func makeBox(s tcell.Screen) { _ = "STUB: not implemented"; return }

func flipCoin() bool { _ = "STUB: not implemented"; return false }

func main() {

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
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
				case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyCtrlQ:
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

	cnt := 0
	dur := time.Duration(0)
loop:
	for {
		select {
		case <-quit:
			break loop
		case <-time.After(interval):
		}
		start := time.Now()
		makeBox(s)
		dur += time.Since(start)
		cnt++
		if cnt%(256/int(inc)) == 0 {
			if flipCoin() {
				redi = -redi
			}
			if flipCoin() {
				grni = -grni
			}
			if flipCoin() {
				blui = -blui
			}
		}
	}

	s.Fini()
	fmt.Printf("Finished %d boxes in %s\n", cnt, dur)
	fmt.Printf("Average is %0.3f ms / box\n", (float64(dur)/float64(cnt))/1000000.0)
}
