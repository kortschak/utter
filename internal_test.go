/*
 * Copyright (c) 2013 Dave Collins <dave@davec.name>
 * Copyright (c) 2015 Dan Kortschak <dan.kortschak@adelaide.edu.au>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

/*
This test file is part of the utter package rather than than the utter_test
package because it needs access to internals to properly test certain cases
which are not possible via the public interface since they should never happen.
*/

package utter

import (
	"bytes"
	"reflect"
	"testing"
	"unsafe"
)

// dummyFmtState implements a fake fmt.State to use for testing invalid
// reflect.Value handling.  This is necessary because the fmt package catches
// invalid values before invoking the formatter on them.
type dummyFmtState struct {
	bytes.Buffer
}

func (dfs *dummyFmtState) Flag(f int) bool {
	return f == int('+')
}

func (dfs *dummyFmtState) Precision() (int, bool) {
	return 0, false
}

func (dfs *dummyFmtState) Width() (int, bool) {
	return 0, false
}

// TestInvalidReflectValue ensures the dump and formatter code handles an
// invalid reflect value properly.  This needs access to internal state since it
// should never happen in real code and therefore can't be tested via the public
// API.
func TestInvalidReflectValue(t *testing.T) {
	i := 1

	// Dump invalid reflect value.
	v := new(reflect.Value)
	buf := new(bytes.Buffer)
	d := dumpState{w: buf, cs: &Config}
	d.dump(*v, false, true, false, 0)
	s := buf.String()
	want := "<invalid>"
	if s != want {
		t.Errorf("InvalidReflectValue #%d\n got: %s want: %s", i, s, want)
	}
}

// changeKind uses unsafe to intentionally change the kind of a reflect.Value to
// the maximum kind value which does not exist.  This is needed to test the
// fallback code which punts to the standard fmt library for new types that
// might get added to the language.
func changeKind(v *reflect.Value, readOnly bool) {
	rvf := (*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + offsetFlag))
	*rvf = *rvf | ((1<<flagKindWidth - 1) << flagKindShift)
	if readOnly {
		*rvf |= flagRO
	} else {
		*rvf &= ^uintptr(flagRO)
	}
}

// TestAddedReflectValue tests functionaly of the dump and formatter code which
// falls back to the standard fmt library for new types that might get added to
// the language.
func TestAddedReflectValue(t *testing.T) {
	i := 1

	// Dump using a reflect.Value that is exported.
	v := reflect.ValueOf(int8(5))
	changeKind(&v, false)
	buf := new(bytes.Buffer)
	d := dumpState{w: buf, cs: &Config}
	d.dump(v, false, true, false, 0)
	s := buf.String()
	want := "int8(5)"
	if s != want {
		t.Errorf("TestAddedReflectValue #%d\n got: %s want: %s", i, s, want)
	}
	i++

	// Dump using a reflect.Value that is not exported.
	changeKind(&v, true)
	buf.Reset()
	d.dump(v, false, true, false, 0)
	s = buf.String()
	want = "int8(<int8 Value>)"
	if s != want {
		t.Errorf("TestAddedReflectValue #%d\n got: %s want: %s", i, s, want)
	}
}

// SortMapByKeyVals makes the internal sortMapByKeyVals function available
// to the test package.
func SortMapByKeyVals(keys, vals []reflect.Value) {
	sortMapByKeyVals(keys, vals)
}
