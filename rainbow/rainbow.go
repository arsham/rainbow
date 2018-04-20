// Copyright 2016 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the Apache 2.0 license
// License that can be found in the LICENSE file.

// Package rainbow prints texts in beautiful rainbows in terminal. Usage is very
// simple:
//
//   import "github.com/arsham/rainbow/rainbow"
//   // ...
//   rb := rainbow.Light{
//       Reader: someReader, // to read from
//       Writer: os.Stdout, // to write to
//   }
//   rb.Paint() // will rainbow everything it reads from reader to writer.
//
// If you want the rainbow to be random, you can seed it this way:
//   rb := rainbow.Light{
//       Reader: buf,
//       Writer: os.Stdout,
//       Seed:   int(rand.Int31n(256)),
//   }
//
package rainbow

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"regexp"
)

var (
	colorMatch = regexp.MustCompile("^\033" + `\[(\d+)(;\d+)?(;\d+)?[m|K]`)
	tabs       = []byte("        ")
)

const (
	freq   = 0.1
	spread = 3
)

// Light reads from the reader and paints to the writer. You should seed it
// everytime otherwise you get the same results.
type Light struct {
	io.Reader
	io.Writer
	Seed int
}

// Paint returns an error if it could not copy the data.
func (rb *Light) Paint() error {
	if rb.Seed == 0 {
		rb.Seed = int(rand.Int31n(256))
	}
	_, err := io.Copy(rb, rb.Reader)
	return err
}

func (rb *Light) Write(data []byte) (int, error) {
	var offset int
	for i := 0; i < len(data); i++ {
		c := data[i]
		switch c {
		case '\n':
			offset = 0
			rb.Seed++
			if _, err := rb.Writer.Write([]byte{'\n'}); err != nil {
				return 0, err
			}
		case '\t':
			offset += len(tabs)
			if _, err := rb.Writer.Write(tabs); err != nil {
				return 0, err
			}
		default:
			pos := colorMatch.FindIndex(data[i:])
			if pos != nil {
				i += pos[1] - 1
				continue
			}
			r, g, b := plotPos(float64(rb.Seed) + (float64(offset) / spread))
			if _, err := colourise(rb.Writer, c, r, g, b); err != nil {
				return 0, err
			}
			offset++
		}
	}
	return len(data), nil
}

func plotPos(x float64) (int, int, int) {
	red := math.Sin(freq*x)*127 + 128
	green := math.Sin(freq*x+2*math.Pi/3)*127 + 128
	blue := math.Sin(freq*x+4*math.Pi/3)*127 + 128
	return int(red), int(green), int(blue)
}

func colourise(w io.Writer, c byte, r, g, b int) (int, error) {
	return fmt.Fprintf(w, "\033[38;5;%dm%c\033[0m", colour(float64(r), float64(g), float64(b)), c)
}

func colour(red, green, blue float64) int {
	return 16 + baseColor(red, 36) + baseColor(green, 6) + baseColor(blue, 1)
}

func baseColor(value float64, factor int) int {
	return int(6*value/256) * factor
}
