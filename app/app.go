// Copyright 2016 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the Apache 2.0 license
// License that can be found in the LICENSE file.

// Package app bootstraps the application.
package app

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

// Main reads the os.Args and uses everything after the first one as input.
func Main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var (
		l    *rainbow.Light
		r    io.Reader
		seed = int(rand.Int31n(256))
	)
	switch len(os.Args) {
	case 1:
		r = os.Stdin
	default:
		r = bytes.NewBufferString(strings.Join(os.Args[1:], " ") + "\n")
	}
	l = &rainbow.Light{
		Writer: os.Stdout,
		Reader: r,
		Seed:   seed,
	}
	if err := l.Paint(); err != nil {
		log.Fatal(err)
	}
}
