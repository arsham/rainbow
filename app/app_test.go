// Copyright 2016 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the Apache 2.0 license
// License that can be found in the LICENSE file.

package app_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"github.com/arsham/rainbow/app"
)

var re = regexp.MustCompile(`\x1B\[[0-9;]*[JKmsu]`)

func setup(t *testing.T) func() {
	oldStdin := os.Stdin
	oldStdout := os.Stdout
	oldArgs := os.Args

	fin, err := ioutil.TempFile("", "testMain")
	if err != nil {
		t.Fatal(err)
	}
	fout, err := ioutil.TempFile("", "testMain")
	if err != nil {
		fin.Close()
		t.Fatal(err)
	}
	os.Stdin = fin
	os.Stdout = fout
	return func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
		os.Args = oldArgs
		fin.Close()
		fout.Close()
		os.Remove(fin.Name())
		os.Remove(fout.Name())
	}
}

func TestMainArg(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()
	input := "RjAuFZXfyIkaE5Ox LitwdlBD6E0GL4Y6p"
	os.Args = []string{"rainbow", input}
	app.Main()
	os.Stdout.Seek(0, 0)
	buf := new(bytes.Buffer)
	buf.ReadFrom(os.Stdout)

	out := buf.Bytes()
	got := re.ReplaceAll(out, []byte(""))
	if !bytes.Equal(got, []byte(input+"\n")) {
		t.Errorf("want %v = %v", got, []byte(input+"\n"))
	}
}

func TestMainPipe(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()
	input := "RjAuFZXfyIkaE5Ox LitwdlBD6E0GL4Y6p"
	os.Args = []string{"rainbow"}
	os.Stdin.WriteString(input)
	os.Stdin.Seek(0, 0)
	app.Main()
	os.Stdout.Seek(0, 0)
	buf := new(bytes.Buffer)
	buf.ReadFrom(os.Stdout)

	out := buf.Bytes()
	got := re.ReplaceAll(out, []byte(""))
	if !bytes.Equal(got, []byte(input)) {
		t.Errorf("want (%s) = (%s)", string(got), input)
	}
}
