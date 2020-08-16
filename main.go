// Copyright 2020 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the Apache 2.0 license
// License that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/arsham/rainbow/rainbow"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var r io.Reader
	switch len(os.Args) {
	case 1:
		r = os.Stdin
	default:
		r = bytes.NewBufferString(strings.Join(os.Args[1:], " ") + "\n")
	}
	l := &rainbow.Light{
		Writer: os.Stdout,
		Seed:   rand.Int63n(256),
	}

	if _, err := io.Copy(l, r); err != nil {
		log.Fatal(err)
	}
}
