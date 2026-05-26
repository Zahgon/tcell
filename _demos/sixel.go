//go:build ignore
// +build ignore

// Copyright 2025 The TCell Authors
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
//
// sixel displays a sixel and demonstrates how to use direct drawing
package main

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/gdamore/tcell/v3/encoding"

	"github.com/mattn/go-sixel"
)

type imageData struct {
	width  int // width in pixels
	height int // height in pixels
	data   *bytes.Buffer
}

func displayHelloWorld(s tcell.Screen) { _ = "STUB: not implemented"; return }

func displaySixel(s tcell.Screen, img *imageData, lock bool) { _ = "STUB: not implemented"; return }

// Get the dimensions of a single cell

// Calculate the image dimensions in cells. We round up to prevent
// drawing on a partially filled cell

// Center the image horizontally

// Lock the region where we will draw the sixel, this prevents tcell
// from drawing over this area

// Move the cursor to our draw position

// Draw the sixel data

func loadImage(path string) (image.Image, error) {
	_ = "STUB: not implemented"
	return *new(image.Image), nil
}

func main() {
	encoding.Register()

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(color.Black).
		Foreground(color.White)
	s.SetStyle(defStyle)

	raw, err := loadImage("./logos/tcell.png")
	if err != nil {
		s.Fini()
		log.Println("couldn't load image. try running from the root directory")
		log.Fatalf("        go run ./_demos/sixel.go")
	}

	img := &imageData{
		width:  raw.Bounds().Dx(),
		height: raw.Bounds().Dy(),
		data:   bytes.NewBuffer(nil),
	}
	enc := sixel.NewEncoder(img.data)
	if err := enc.Encode(raw); err != nil {
		s.Fini()
		log.Fatal(err)
	}

	lock := true
	displayHelloWorld(s)
	displaySixel(s, img, lock)

	for {
		ev := <-s.EventQ()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			displayHelloWorld(s)
			displaySixel(s, img, lock)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
			if ev.Key() == tcell.KeyEnter {
				lock = !lock
				displaySixel(s, img, lock)
			}
		}
	}
}
