// Copyright 2016 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the Apache 2.0 license
// License that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
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
	input := gofakeit.Sentence(20)
	os.Args = []string{"rainbow", input}
	main()
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
	input := gofakeit.Sentence(20)
	os.Args = []string{"rainbow"}
	os.Stdin.WriteString(input)
	os.Stdin.Seek(0, 0)
	main()
	os.Stdout.Seek(0, 0)
	buf := new(bytes.Buffer)
	buf.ReadFrom(os.Stdout)

	out := buf.Bytes()
	got := re.ReplaceAll(out, []byte(""))
	assert.Equal(t, []byte(input), got)
}
