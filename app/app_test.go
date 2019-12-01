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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var re = regexp.MustCompile(`\x1B\[[0-9;]*[JKmsu]`)

func setup(t *testing.T) func() {
	t.Helper()
	oldStdin := os.Stdin
	oldStdout := os.Stdout
	oldArgs := os.Args

	fin, err := ioutil.TempFile("", "testMain")
	require.NoError(t, err)
	fout, err := ioutil.TempFile("", "testMain")
	if err != nil {
		assert.NoError(t, fin.Close())
		t.Fatal(err)
	}
	os.Stdin = fin
	os.Stdout = fout
	return func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
		os.Args = oldArgs
		assert.NoError(t, fin.Close())
		assert.NoError(t, fout.Close())
		assert.NoError(t, os.Remove(fin.Name()))
		assert.NoError(t, os.Remove(fout.Name()))
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
	assert.Equal(t, []byte(input+"\n"), got)
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
	assert.Equal(t, []byte(input), got)
}
