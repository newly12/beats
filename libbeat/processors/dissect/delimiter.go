// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package dissect

import (
	"fmt"
	"strings"
)

//delimiter represents a text section after or before a key, it keeps track of the needle and allows
// to retrieve the position where it starts from a haystack.
type delimiter interface {
	// IndexOf receives the haystack and a offset position and will return the absolute position where
	// the needle is found.
	IndexOf(haystack string, offset int) int

	// Len returns the length of the needle used to calculate boundaries.
	Len() int

	// String displays debugging information.
	String() string

	// Delimiter returns the actual delimiter string.
	Delimiter() string

	// IsGreedy return true if the next key should be greedy (end of string) or when explicitly
	// configured.
	IsGreedy() bool

	// MarkGreedy marks this delimiter as greedy.
	MarkGreedy()

	// Next returns the next delimiter in the chain.
	Next() delimiter

	//SetNext sets the next delimiter or nil if current delimiter is the last.
	SetNext(d delimiter)
}

type baseDelimiter struct {
	needle string
	length int
	greedy bool
	next   delimiter
}

func (b *baseDelimiter) IndexOf(haystack string, offset int) int {
	return offset + b.length
}

func (b *baseDelimiter) Len() int {
	return b.length
}

func (b *baseDelimiter) Delimiter() string {
	return b.needle
}

func (b *baseDelimiter) IsGreedy() bool {
	return b.greedy
}

func (b *baseDelimiter) MarkGreedy() {
	b.greedy = true
}

func (b *baseDelimiter) Next() delimiter {
	return b.next
}

func (b *baseDelimiter) SetNext(d delimiter) {
	b.next = d
}

// zeroByte represents a zero string delimiter its usually start of the line.
type zeroByte struct {
	baseDelimiter
}

func (z *zeroByte) String() string {
	return "delimiter: zerobyte"
}

// fixedLengthByte represents a delimiter with fixed-length bytes.
type fixedLengthByte struct {
	baseDelimiter
}

func (f *fixedLengthByte) String() string {
	return fmt.Sprintf("delimiter: fixedlengthbyte (len: %d)", f.length)
}

// multiByte represents a delimiter with at least one byte.
type multiByte struct {
	baseDelimiter
}

func (m *multiByte) IndexOf(haystack string, offset int) int {
	// todo if fixed length, m.needle should not be used here
	i := strings.Index(haystack[offset:], m.needle)
	if i != -1 {
		return i + offset
	}
	return -1
}

func (m *multiByte) Len() int {
	return len(m.needle)
}

func (m *multiByte) String() string {
	return fmt.Sprintf("delimiter: multibyte (match: '%s', len: %d)", string(m.needle), m.Len())
}

func newDelimiter(needle string) delimiter {
	if len(needle) == 0 {
		return &zeroByte{}
	}
	return &multiByte{baseDelimiter{needle: needle}}
}
