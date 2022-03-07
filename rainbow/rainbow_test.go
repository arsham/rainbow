// Copyright 2016 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the Apache 2.0 license
// License that can be found in the LICENSE file.

package rainbow

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func readFile(t *testing.T, name string) []byte {
	f, err := os.Open("testdata/" + name)
	require.NoError(t, err)
	b, err := io.ReadAll(f)
	require.NoError(t, err)
	return b
}

func TestLightPaint(t *testing.T) {
	t.Parallel()
	plain := string(readFile(t, "plain.txt"))
	painted := string(readFile(t, "painted.txt"))
	//nolint:stylecheck // this is on purpose.
	tcs := []struct {
		sample  string
		painted string
	}{
		{
			"2d7gMRSgGLj9F0c tPjSmsdRsTej4x7BJiOp R9HUHEiyH0G1Ld XeL5fjQ1KkxI3",
			`[38;5;154m2[0m[38;5;154md[0m[38;5;154m7[0m[38;5;154mg[0m[38;5;154mM[0m[38;5;154mR[0m[38;5;154mS[0m[38;5;148mg[0m[38;5;184mG[0m[38;5;184mL[0m[38;5;184mj[0m[38;5;184m9[0m[38;5;184mF[0m[38;5;184m0[0m[38;5;184mc[0m[38;5;184m [0m[38;5;184mt[0m[38;5;184mP[0m[38;5;184mj[0m[38;5;178mS[0m[38;5;214mm[0m[38;5;214ms[0m[38;5;214md[0m[38;5;214mR[0m[38;5;214ms[0m[38;5;214mT[0m[38;5;214me[0m[38;5;214mj[0m[38;5;214m4[0m[38;5;208mx[0m[38;5;208m7[0m[38;5;208mB[0m[38;5;208mJ[0m[38;5;208mi[0m[38;5;208mO[0m[38;5;208mp[0m[38;5;208m [0m[38;5;208mR[0m[38;5;209m9[0m[38;5;203mH[0m[38;5;203mU[0m[38;5;203mH[0m[38;5;203mE[0m[38;5;203mi[0m[38;5;203my[0m[38;5;203mH[0m[38;5;203m0[0m[38;5;203mG[0m[38;5;203m1[0m[38;5;203mL[0m[38;5;204md[0m[38;5;198m [0m[38;5;198mX[0m[38;5;198me[0m[38;5;198mL[0m[38;5;198m5[0m[38;5;198mf[0m[38;5;198mj[0m[38;5;198mQ[0m[38;5;198m1[0m[38;5;199mK[0m[38;5;199mk[0m[38;5;199mx[0m[38;5;199mI[0m[38;5;199m3[0m`,
		},
		{
			"11✂1",
			`[38;5;154m1[0m[38;5;154m1[0m[38;5;154m✂[0m[38;5;154m1[0m`,
		},
		{
			"🏧-✂1",
			`[38;5;154m🏧[0m[38;5;154m-[0m[38;5;154m✂[0m[38;5;154m1[0m`,
		},
		{
			plain,
			painted,
		},
	}
	for _, tc := range tcs {
		r := strings.NewReader(tc.sample)
		w := &bytes.Buffer{}
		l := &Light{
			Reader: r,
			Writer: w,
			Seed:   1,
		}
		err := l.Paint()
		assert.NoError(t, err)
		assert.EqualValues(t, tc.painted, w.String())
	}
}

func BenchmarkLightPaint(b *testing.B) {
	bcs := []struct {
		lines   int
		letters int
	}{
		{1, 10},
		{1, 240},
		{10, 10},
		{10, 240},
		{100, 10},
		{100, 240},
		{500, 500},
	}
	for _, bc := range bcs {
		var (
			totalLen int
			name     = fmt.Sprintf("lines%d_let%d", bc.lines, bc.letters)
			line     = make([]byte, bc.letters)
			r        = &bytes.Buffer{}
			w        = &bytes.Buffer{}
		)
		rand.Read(line)
		for i := 0; i < bc.lines; i++ {
			r.Write(line)
			r.WriteString("\n")
			totalLen += len(line) + 1
		}
		b.ResetTimer()
		b.Run(name, func(b *testing.B) {
			b.Run("Serial", func(b *testing.B) {
				b.ResetTimer()
				l := &Light{
					Writer: w,
					Reader: r,
					Seed:   1,
				}
				for i := 0; i < b.N; i++ {
					l.Paint()
				}
			})
			b.Run("Parallel", func(b *testing.B) {
				b.ResetTimer()
				b.RunParallel(func(bp *testing.PB) {
					l := &Light{
						Writer: w,
						Reader: r,
						Seed:   1,
					}
					for bp.Next() {
						l.Paint()
					}
				})
			})
		})
	}
}

func TestPlotPos(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name  string
		x     float64
		red   int
		green int
		blue  int
	}{
		{"0", 0, 128, 237, 18},
		{"1", 1, 140, 231, 12},
		{"5", 5, 188, 194, 1},
		{"10", 10, 234, 133, 15},
		{"50", 50, 6, 220, 157},
		{"100", 100, 58, 70, 254},
		{"360", 360, 2, 176, 205},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, got1, got2 := plotPos(tc.x)
			assert.Equal(t, tc.red, int(got), "red value")
			assert.Equal(t, tc.green, int(got1), "green value")
			assert.Equal(t, tc.blue, int(got2), " blue value")
		})
	}
}

var got, got1, got2 float64

func BenchmarkPlotPos(b *testing.B) {
	b.Run("Serial", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			got, got1, got2 = plotPos(100)
		}
	})
	b.Run("Parallel", func(b *testing.B) {
		b.RunParallel(func(b *testing.PB) {
			for b.Next() {
				got, got1, got2 = plotPos(100)
			}
		})
	})
}

// this test is here for keeping the logic in sync when we refactor the codes.
func TestColour(t *testing.T) {
	t.Parallel()
	randColour := func() int32 {
		return rand.Int31n(256)
	}
	bc := func(value float64, factor int) int {
		return int(6*value/256) * factor
	}
	check := func(red, green, blue float64) int {
		return 16 + bc(red, 36) + bc(green, 6) + bc(blue, 1)
	}
	for i := 0; i < 1000; i++ {
		red, green, blue := float64(randColour()), float64(randColour()), float64(randColour())
		got := colour(red, green, blue)
		want := check(red, green, blue)
		assert.Equalf(t, want, got, "colour(%f, %f, %f)", red, green, blue)
	}
}

var intval int

func BenchmarkColour(b *testing.B) {
	b.Run("Serial", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			intval = colour(6, 100, 1000)
		}
	})
	b.Run("Parallel", func(b *testing.B) {
		b.RunParallel(func(b *testing.PB) {
			for b.Next() {
				intval = colour(6, 100, 1000)
			}
		})
	})
}

type writeError func([]byte) (int, error)

func (w *writeError) Write(p []byte) (int, error) { return (*w)(p) }

func TestLightWrite(t *testing.T) {
	errExam := errors.New("this error")
	wrErr := writeError(func([]byte) (int, error) { return 0, errExam })
	//nolint:stylecheck // this is on purpose.
	tcs := map[string]struct {
		data     []byte
		want     []byte
		checkErr bool
	}{
		"new line":           {[]byte("\n"), []byte("\n"), true},
		"tab":                {[]byte("\t"), []byte("\t"), true},
		"NL tab":             {[]byte("\n\t"), []byte("\n\t"), true},
		"tab NL":             {[]byte("\t\n"), []byte("\t\n"), true},
		`033[38;5;2m`:        {[]byte("\033[38;5;2m"), []byte(""), false},
		`033[38;5;2K`:        {[]byte("\033[38;5;2K"), []byte(""), false},
		`033[32K`:            {[]byte("\033[32K"), []byte(""), false},
		`033[3K`:             {[]byte("\033[3K"), []byte(""), false},
		`033[3KARSHAM bytes`: {[]byte("\033[3KARSHAM"), []byte{27, 91, 51, 56, 59, 53, 59, 49, 53, 52, 109, 65, 27, 91, 48, 109, 27, 91, 51, 56, 59, 53, 59, 49, 53, 52, 109, 82, 27, 91, 48, 109, 27, 91, 51, 56, 59, 53, 59, 49, 53, 52, 109, 83, 27, 91, 48, 109, 27, 91, 51, 56, 59, 53, 59, 49, 53, 52, 109, 72, 27, 91, 48, 109, 27, 91, 51, 56, 59, 53, 59, 49, 53, 52, 109, 65, 27, 91, 48, 109, 27, 91, 51, 56, 59, 53, 59, 49, 53, 52, 109, 77, 27, 91, 48, 109}, true},
		`033[3KARSHAM string`: {[]byte("\033[3KARSHAM"), []byte("[38;5;154mA[0m[38;5;154mR[0m[38;5;154mS[0m[38;5;154mH[0m[38;5;154mA[0m[38;5;154mM[0m"), true},
	}
	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			l := &Light{
				Writer: buf,
				Seed:   1,
			}
			_, err := l.Write(tc.data)
			assert.NoError(t, err)
			got := buf.Bytes()
			if !bytes.Equal(got, tc.want) {
				t.Errorf("got (%v), want (%v)", got, tc.want)
			}
			if !tc.checkErr {
				return
			}
			l.Writer = &wrErr
			_, err = l.Write(tc.data)
			assert.Error(t, err)
		})
	}
	r := &bytes.Buffer{}
	l := &Light{
		Reader: r,
	}
	n, err := l.Write([]byte("blah"))
	assert.Error(t, err)
	assert.Zero(t, n)
}

func TestLightWriteRace(t *testing.T) {
	var (
		wg    sync.WaitGroup
		count = 1000
		data  = bytes.Repeat([]byte("abc def\n"), 10)
		l     = &Light{
			Writer: io.Discard,
			Seed:   1,
		}
	)
	wg.Add(3)
	go func() {
		for i := 0; i < count; i++ {
			_, err := l.Write(data)
			require.NoError(t, err)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < count; i++ {
			_, err := l.Write(data)
			require.NoError(t, err)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < count; i++ {
			_, err := l.Write(data)
			require.NoError(t, err)
		}
		wg.Done()
	}()
	wg.Wait()
}

func BenchmarkLightWrite(b *testing.B) {
	bcs := []struct {
		line   int
		length int
		data   string
	}{
		{1, 1, "\n"},
		{5, 10, strings.Repeat("aaaaa\n", 10)},
		{15, 50, strings.Repeat("aaaaabbbbbccccc\n", 15)},
		{15, 100, strings.Repeat(strings.Repeat("abcde", 20)+"\n", 15)},
		{50, 50, strings.Repeat(strings.Repeat("abcde", 10)+"\n", 50)},
		{100, 120, strings.Repeat(strings.Repeat("a", 120)+"\n", 100)},
	}
	b.ResetTimer()
	b.Run("Serial", func(b *testing.B) {
		for _, bc := range bcs {
			bc := bc
			l := &Light{
				Writer: io.Discard,
				Seed:   1,
			}
			name := fmt.Sprintf("line%d_len%d", bc.line, bc.length)
			b.Run(name, func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					l.Write([]byte(bc.data))
				}
			})
		}
	})
	b.Run("Parallel", func(b *testing.B) {
		for _, bc := range bcs {
			bc := bc
			l := &Light{
				Writer: io.Discard,
				Seed:   1,
			}
			name := fmt.Sprintf("line%d_len%d", bc.line, bc.length)
			b.Run(name, func(b *testing.B) {
				b.ResetTimer()
				b.RunParallel(func(b *testing.PB) {
					for b.Next() {
						l.Write([]byte(bc.data))
					}
				})
			})
		}
	})
}

// making sure we are not altering any texts.
func TestLightWriteRevert(t *testing.T) {
	re := regexp.MustCompile(`\x1B\[[0-9;]*[JKmsu]`)
	tcs := []string{
		"u9VGCQ1E4KCr8bO8 3ULdtlHL3WsjulJU kqUneSFT6 tvyAfih1Qew 5wBKffL4Yc",
		"Y5LFQNulLC0GTKB W4buVQmQTMu6C7aFs uGL6 x2OgVRlUZHCq46kgk sjr 4HKIb",
		"MamsTagRix6bEYwBGR b9FK7b1L 5x1YtTo8nFLrz0dZ rIdZdY0 b0bC05T42bHfV",
		"Ubmba7pnOCoUw9 xGgDjSIPU 7vOwUoiPHeCoxT XtywrjciBYZR cBQySayKWb rx",
		"EPz0kb3pQUuVv LL0t4t8mNaRklyuZTPi wHI2H35IReZbbdb9akXw gMCuJ PRrVK",
		"9qR4HkZ86 enPzoAyIWfz3bFg 8LSokcvV47J0 XfHM ngASr MfkM53zAkwLZmY  ",
		"G7R9mjoYz9FU306 ATfgJ2C6AokIgtf BmL5uLWSJFTDs VO85P JXUjV 5n4OEuvl",
		"32Tj lwI9YaHDjsbqZUBIi 5XiJ7tOh 4eCqaaHT1i WpsJh3JA7s HKHk5R9yPKjL",
		"188qrmdf0GLU9 E2g vp5iX4g2CJ ueKSPvwY369daXZU 0bhg7IGhzeegWeGk5Fj ",
		"YXXUN751 qkTsoR95Udu EtFAN NgBAEe97uzp wpu2VKcX0W20P084d V1LLH28Rk",
		"nSZ6 Qww4GOxAsduJkEPIuNln GGbS rNq YEfhc4jaTAzoHC eT1ugr0mMYANX5JX",
		"wFZkZVWqV3ag7pYtlyN mKtPoVvvZMgU3p6E1 u6zSuUFHuk8nZsvQQ5 qeq 14YPo",
		"uXDLM l2J sznWecxn Ayjv9Ii3akvRD ArTwAkrA THyyrqT6LSAGnfJSSMx4Mes ",
		"m3Lo ufeV1XmoFxqs 4LE21GCiI5UT51 pj9YV B2a43pJxjCg b5CqaspplX5N2dq",
		"9BwGLsCTNGXDfLXNl ZcPNImZhzmDp8S 3ZG177TSj tjOSRMxBZ rhjP0zwJU1K o",
		"4MySJAF rybjdvOUuZhf VqPKvVuw5jaNsldI 8oYI8ZTL2s mZ7X 4awf4PPeLIHh",
		"LTLj3ayz83gM N5So9T32GVipPB2B Ccy 8UutGx N7u6DeZ bUN1hsBDoSC4Z0 sd",
		"gybKEl1 73H5qkvjR UbTUZl jAeMVFeDCtAGVbGhCq Fi0ZQctnmhtk0edD5 gkU3",
		"YX4dsZWyblp6PJIqTR vTgLmrYZd3 MlnwDVOtS1wKZDpoqlxY vu0ivZaYrqWhzAV",
		"Vb1Nks0QpBl1CTkJF QiO7nMLehiI0u S2XwbQXs9Znz Mbe0BQ13JTAOmSjmh WFH",
		"GNIzszrfT5 8WBHpE00a5j7Srfnx e8Qrrhomy8tw7XIa kQFG7ZazYio x5z PZIQ",
	}
	for _, tc := range tcs {
		w := &bytes.Buffer{}
		r := bytes.NewReader([]byte(tc))
		l := Light{
			Reader: r,
			Writer: w,
		}
		err := l.Paint()
		require.NoError(t, err)

		got := re.ReplaceAll(w.Bytes(), []byte(""))
		assert.EqualValues(t, tc, got)
	}
}
