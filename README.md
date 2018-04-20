# Rainbow

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://godoc.org/github.com/arsham/rainbow?status.svg)](http://godoc.org/github.com/arsham/rainbow)
[![Build Status](https://travis-ci.org/arsham/rainbow.svg?branch=master)](https://travis-ci.org/arsham/rainbow)
[![Coverage Status](https://codecov.io/gh/arsham/rainbow/branch/master/graph/badge.svg)](https://codecov.io/gh/arsham/rainbow)
[![Go Report Card](https://goreportcard.com/badge/github.com/arsham/rainbow)](https://goreportcard.com/report/github.com/arsham/rainbow)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/6cc048fae4ba4129b05226308a0bd7e9)](https://www.codacy.com/app/arsham/rainbow?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=arsham/rainbow&amp;utm_campaign=Badge_Grade)

Tasty rainbows for your terminal like these:

![Screenshot](/docs/rainbow.png?raw=true "Rainbow")

This app was inspired by lolcats, but written in Go.

### Table of Contents

1. [Installation](#installation)
2. [Usage](#usage)
3. [As library](#as-library)
4. [See Also](#see-also)
5. [License](#license)

## Installation

Get the library:
```bash
$ go get github.com/arsham/rainbow
```

## Usage

You can pipe the text into the app in many ways. Choose one that is suitable for
you:
```bash
# File contents:
$ rainbow < filename.txt

# Echo a string:
$ echo "Any quotes" | rainbow

# Here string:
$ rainbow <<END
Consectetur aliqua do quis sed
proident enim fugiat occaecat nisi
in deserunt culpa aliquip do excepteur.
END

# Output of a program:
$ ls -l | rainbow
```

## As library
`Light` struct implements io.Reader and io.Writer:
```go
import "github.com/arsham/rainbow/rainbow"
// ...
rb := rainbow.Light{
    Reader: someReader, // to read from
    Writer: os.Stdout, // to write to
}
rb.Paint() // will rainbow everything it reads from reader to writer.
```
If you want the rainbow to be random, you can seed it this way:

```go
rb := rainbow.Light{
    Reader: buf,
    Writer: os.Stdout,
    Seed:   int(rand.Int31n(256)),
}
```

## See Also
See also [Figurine][figurine]. It prints beautiful ASCII arts in FIGlet.

## License
Use of this source code is governed by the Apache 2.0 license. License that can
be found in the [LICENSE](./LICENSE) file.

Please note that this was initially forked from [glolcat][glolcat], but I
decided to rewrite it as the library is not maintained anymore.

[figurine]: https://github.com/arsham/figurine
[glolcat]: https://github.com/cezarsa/glolcat
