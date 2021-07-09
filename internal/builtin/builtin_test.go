// See doc.go

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package log implements a simple logging package. It defines a type, Builtin,
// with methods for formatting output. It also has a predefined 'standard'
// Builtin accessible through helper functions Print[f|ln], Fatal[f|ln], and
// Panic[f|ln], which are easier to use than creating a Builtin manually.
// That logger writes to standard error and prints the date and time
// of each logged message.
// The Fatal functions call os.Exit(1) after writing the log message.
// The Panic functions call panic after writing the log message.
package builtin

// These tests are too simple.

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

const (
	Rdate         = `[0-9][0-9][0-9][0-9]/[0-9][0-9]/[0-9][0-9]`
	Rtime         = `[0-9][0-9]:[0-9][0-9]:[0-9][0-9]`
	Rmicroseconds = `\.[0-9][0-9][0-9][0-9][0-9][0-9]`
	Rline         = `(60|62):` // must update if the calls to l.Printf / l.Print below move
	Rlongfile     = `.*/[A-Za-z0-9_\-]+\.go:`
	Rshortfile    = `[A-Za-z0-9_\-]+\.go:` + Rline
)

type tester struct {
	flag    int
	prefix  string
	pattern string // regexp that log output must match; we add ^ and expected_text$ always
}

var tests = []tester{
	// individual pieces:
	{0, "", ""},
	{0, "XXX", "XXX"},
	{Ldate, "", Rdate + " "},
	{Ltime, "", Rtime + " "},
	{Ltime | Lmsgprefix, "XXX", Rtime + " XXX"},
	{Ltime | Lmicroseconds, "", Rtime + Rmicroseconds + " "},
	{Lmicroseconds, "", Rtime + Rmicroseconds + " "}, // microsec implies time
}

// Test using Println("hello", 23, "world") or using Printf("hello %d world", 23)
func testPrint(t *testing.T, flag int, prefix string, pattern string, useFormat bool) {
	buf := new(bytes.Buffer)
	SetOutput(buf)
	SetFlags(flag)
	SetPrefix(prefix)
	if useFormat {
		Printf("hello %d world", 23)
	} else {
		Println("hello", 23, "world")
	}
	line := buf.String()
	pattern = "^" + pattern + "hello 23 world"
	matched, err := regexp.MatchString(pattern, line)
	if err != nil {
		t.Fatal("pattern did not compile:", err)
	}
	if !matched {
		t.Errorf("log output should match is %q", line)
	}
	SetOutput(os.Stderr)
}

func TestDefault(t *testing.T) {
	if got := Default(); got != std {
		t.Errorf("Default [%p] should be std [%p]", got, std)
	}
}

func TestAll(t *testing.T) {
	for _, testcase := range tests {
		testPrint(t, testcase.flag, testcase.prefix, testcase.pattern, false)
		testPrint(t, testcase.flag, testcase.prefix, testcase.pattern, true)
	}
}

func TestOutput(t *testing.T) {
	const testString = "test"
	var b bytes.Buffer
	l := NewBuiltin(&b, "", 0)
	l.Println(testString)
	if expect := testString + "\n"; b.String() != expect {
		t.Errorf("log output should match %q is %q", expect, b.String())
	}
}

func TestOutputRace(t *testing.T) {
	var b bytes.Buffer
	l := NewBuiltin(&b, "", 0)
	for i := 0; i < 100; i++ {
		go func() {
			l.SetFlags(0)
		}()
		l.OutputBuiltin(0, "")
	}
}

func TestFlagAndPrefixSetting(t *testing.T) {
	var b bytes.Buffer
	l := NewBuiltin(&b, "Test:", LstdFlags)
	f := l.Flags()
	if f != LstdFlags {
		t.Errorf("Flags 1: expected %x got %x", LstdFlags, f)
	}
	l.SetFlags(f | Lmicroseconds)
	f = l.Flags()
	if f != LstdFlags|Lmicroseconds {
		t.Errorf("Flags 2: expected %x got %x", LstdFlags|Lmicroseconds, f)
	}
	p := l.Prefix()
	if p != "Test:" {
		t.Errorf(`Prefix: expected "Test:" got %q`, p)
	}
	l.SetPrefix("Reality:")
	p = l.Prefix()
	if p != "Reality:" {
		t.Errorf(`Prefix: expected "Reality:" got %q`, p)
	}
	// Verify a log message looks right, with our prefix and microseconds present.
	l.Print("hello")
	pattern := "^Reality:" + Rdate + " " + Rtime + Rmicroseconds + " hello"
	matched, err := regexp.Match(pattern, b.Bytes())
	if err != nil {
		t.Fatalf("pattern %q did not compile: %s", pattern, err)
	}
	if !matched {
		t.Error("message did not match pattern")
	}
}

func TestUTCFlag(t *testing.T) {
	var b bytes.Buffer
	l := NewBuiltin(&b, "Test:", LstdFlags)
	l.SetFlags(Ldate | Ltime | LUTC)
	// Verify a log message looks right in the right time zone. Quantize to the second only.
	now := time.Now().UTC()
	l.Print("hello")
	want := fmt.Sprintf("Test:%d/%.2d/%.2d %.2d:%.2d:%.2d hello\n",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	got := b.String()
	if got == want {
		return
	}
	want = fmt.Sprintf("Test:%d/%.2d/%.2d %.2d:%.2d:%.2d hello",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	if got == want {
		return
	}
	t.Errorf("got %q; want %q", got, want)
}

func TestEmptyPrintCreatesLine(t *testing.T) {
	var b bytes.Buffer
	l := NewBuiltin(&b, "Header:", LstdFlags)
	l.Println("")
	output := b.String()
	if n := strings.Count(output, "Header"); n != 1 {
		t.Errorf("expected 2 headers, got %d", n)
	}
	if n := strings.Count(output, "\n"); n != 1 {
		t.Errorf("expected 2 lines, got %d", n)
	}
}

func BenchmarkItoa(b *testing.B) {
	dst := make([]byte, 0, 64)
	for i := 0; i < b.N; i++ {
		dst = dst[0:0]
		itoa(&dst, 2015, 4)   // year
		itoa(&dst, 1, 2)      // month
		itoa(&dst, 30, 2)     // day
		itoa(&dst, 12, 2)     // hour
		itoa(&dst, 56, 2)     // minute
		itoa(&dst, 0, 2)      // second
		itoa(&dst, 987654, 6) // microsecond
	}
}

func BenchmarkPrintln(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := NewBuiltin(&buf, "", LstdFlags)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Println(testString)
	}
}

func BenchmarkPrintlnNoFlags(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := NewBuiltin(&buf, "", 0)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Println(testString)
	}
}
