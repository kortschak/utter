utter
=====

[![Build Status](https://travis-ci.org/kortschak/utter.png?branch=master)]
(https://travis-ci.org/kortschak/utter) [![Coverage Status]
(https://coveralls.io/repos/kortschak/utter/badge.png?branch=master)]
(https://coveralls.io/r/kortschak/utter?branch=master)

utter is a fork of the outstanding [go-spew tool](https://github.com/davecgh/go-spew).
Where go-spew is an aid for debugging, providing annotation of dumped datastructures,
utter is a tool for taking snapshots of data structures to include in tests or other
code. An utter dump will not construct cyclic structure literals and a number of
pseudo-code representations of pointer-based structures will require subsequent
processing.

utter implements a deep pretty printer for Go data structures to aid in
debugging.  A comprehensive suite of tests with 100% test coverage is provided
to ensure proper functionality.  See `test_coverage.txt` for the gocov coverage
report.  utter is licensed under the liberal ISC license, so it may be used in
open source or commercial projects.

If you're interested in reading about how this package came to life and some
of the challenges involved in providing a deep pretty printer, there is a blog
post about it
[here](https://blog.cyphertite.com/go-spew-a-journey-into-dumping-go-data-structures/).

## Documentation

[![GoDoc](https://godoc.org/github.com/kortschak/utter?status.png)]
(http://godoc.org/github.com/kortschak/utter)

Full `go doc` style documentation for the project can be viewed online without
installing this package by using the excellent GoDoc site here:
http://godoc.org/github.com/kortschak/utter

You can also view the documentation locally once the package is installed with
the `godoc` tool by running `godoc -http=":6060"` and pointing your browser to
http://localhost:6060/pkg/github.com/kortschak/utter

## Installation

```bash
$ go get -u github.com/kortschak/utter
```

## Quick Start

To dump a variable with full newlines, indentation, type, and pointer
information use Dump, Fdump, or Sdump:

```Go
utter.Dump(myVar1, myVar2, ...)
utter.Fdump(someWriter, myVar1, myVar2, ...)
str := utter.Sdump(myVar1, myVar2, ...)
```

Alternatively, if you would prefer to use format strings with a compacted inline
printing style, use the convenience wrappers Printf, Fprintf, etc with %v (most
compact), %+v (adds pointer addresses), %#v (adds types), or %#+v (adds types
and pointer addresses):

```Go
utter.Printf("myVar1: %v -- myVar2: %+v", myVar1, myVar2)
utter.Printf("myVar3: %#v -- myVar4: %#+v", myVar3, myVar4)
utter.Fprintf(someWriter, "myVar1: %v -- myVar2: %+v", myVar1, myVar2)
utter.Fprintf(someWriter, "myVar3: %#v -- myVar4: %#+v", myVar3, myVar4)
```

## Sample Dump Output

```
(main.Foo) {
 unexportedField: (*main.Bar)(0xf84002e210)({
  flag: (main.Flag) flagTwo,
  data: (uintptr) <nil>
 }),
 ExportedField: (map[interface {}]interface {}) {
  (string) "one": (bool) true
 }
}
([]uint8) {
 00000000  11 12 13 14 15 16 17 18  19 1a 1b 1c 1d 1e 1f 20  |............... |
 00000010  21 22 23 24 25 26 27 28  29 2a 2b 2c 2d 2e 2f 30  |!"#$%&'()*+,-./0|
 00000020  31 32                                             |12|
}
```

## Sample Formatter Output

Double pointer to a uint8:
```
	  %v: <**>5
	 %+v: <**>(0xf8400420d0->0xf8400420c8)5
	 %#v: (**uint8)5
	%#+v: (**uint8)(0xf8400420d0->0xf8400420c8)5
```

Pointer to circular struct with a uint8 field and a pointer to itself:
```
	  %v: <*>{1 <*><shown>}
	 %+v: <*>(0xf84003e260){ui8:1 c:<*>(0xf84003e260)<shown>}
	 %#v: (*main.circular){ui8:(uint8)1 c:(*main.circular)<shown>}
	%#+v: (*main.circular)(0xf84003e260){ui8:(uint8)1 c:(*main.circular)(0xf84003e260)<shown>}
```

## Configuration Options

Configuration of utter is handled by fields in the ConfigState type. For
convenience, all of the top-level functions use a global state available via the
utter.Config global.

It is also possible to create a ConfigState instance that provides methods
equivalent to the top-level functions. This allows concurrent configuration
options. See the ConfigState documentation for more details.

```
* Indent
	String to use for each indentation level for Dump functions.
	It is a single space by default.  A popular alternative is "\t".

* SortKeys
	Specifies map keys should be sorted before being printed. Use
	this to have a more deterministic, diffable output.  Note that
	only native types (bool, int, uint, floats, uintptr and string)
	are supported with other types sorted according to the
	reflect.Value.String() output which guarantees display stability.
	Natural map order is used by default.
```

## License

utter is licensed under the liberal ISC License.
