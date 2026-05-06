// Copyright 2016 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the Apache 2.0 license
// License that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"regexp"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var re = regexp.MustCompile(`\x1B\[[0-9;]*[JKmsu]`)

const binaryName = "rainbow"

func setup(t *testing.T) {
	t.Helper()
	oldStdin := os.Stdin
	oldStdout := os.Stdout
	oldArgs := os.Args

	fin, err := os.CreateTemp(t.TempDir(), "testMain")
	require.NoError(t, err)
	os.Stdin = fin
	t.Cleanup(func() {
		os.Stdin = oldStdin
		assert.NoError(t, fin.Close())
	})

	fout, err := os.CreateTemp(t.TempDir(), "testMain")
	require.NoError(t, err)
	os.Stdout = fout
	t.Cleanup(func() {
		os.Stdout = oldStdout
		os.Args = oldArgs
		assert.NoError(t, fout.Close())
	})
}

func TestMain(t *testing.T) {
	t.Run("WithArgs", testMainWithArgs)
	t.Run("WithPipe", testMainWithPipe)
	t.Run("CopyError", testMainCopyError)
}

func testMainWithArgs(t *testing.T) {
	setup(t)
	input := gofakeit.Sentence(20)
	os.Args = []string{binaryName, input}
	main()
	os.Stdout.Seek(0, 0)
	buf := &bytes.Buffer{}
	buf.ReadFrom(os.Stdout)

	out := buf.Bytes()
	got := re.ReplaceAll(out, []byte(""))
	assert.Equal(t, []byte(input+"\n"), got)
}

func testMainWithPipe(t *testing.T) {
	setup(t)
	input := gofakeit.Sentence(20)
	os.Args = []string{binaryName}
	os.Stdin.WriteString(input)
	os.Stdin.Seek(0, 0)
	main()
	os.Stdout.Seek(0, 0)
	buf := &bytes.Buffer{}
	buf.ReadFrom(os.Stdout)

	out := buf.Bytes()
	got := re.ReplaceAll(out, []byte(""))
	assert.Equal(t, []byte(input), got)
}

func testMainCopyError(t *testing.T) {
	setup(t)
	fin, err := os.CreateTemp("", "testMain")
	require.NoError(t, err)
	require.NoError(t, fin.Close())
	os.Stdin = fin

	os.Args = []string{binaryName}
	main()
	os.Stdout.Seek(0, 0)
	buf := &bytes.Buffer{}
	buf.ReadFrom(os.Stdout)
	got := buf.String()
	assert.Contains(t, got, "already closed")
}
