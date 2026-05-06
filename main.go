// Copyright 2020 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the Apache 2.0 license
// License that can be found in the LICENSE file.

// Package main is the entrypoint to the rainbow binary.
package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"strings"

	"github.com/arsham/rainbow/v2/rainbow"
)

func main() {
	var r io.Reader
	switch len(os.Args) {
	case 1:
		r = os.Stdin
	default:
		r = bytes.NewBufferString(strings.Join(os.Args[1:], " ") + "\n")
	}
	l := &rainbow.Light{
		Writer: os.Stdout,
		Seed:   rand.Int64N(256),
	}

	if _, err := io.Copy(l, r); err != nil {
		fmt.Println("Error painting your input:", err)
	}
}
