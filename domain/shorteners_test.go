package domain

import (
	"reflect"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestShortenFromBack(t *testing.T) {
	matcher := NewMatcherFromString(`a=aaa
b=bbb
c=ccc
d=ddd
stg=strategy
ltd=limited`)
	tests := []struct {
		name     string
		original string
		max      int
		want     string
	}{
		{name: "Length longer than origin with '-'", original: "aaa-bbb-ccc", max: 99, want: "aaa-bbb-ccc"},
		{name: "Length is 0 with '-'", original: "aaa-bbb-ccc", max: 0, want: "a-b-c"},
		{name: "Partial abbreviation with '-'", original: "aaa-bbb-ccc", max: 10, want: "aaa-bbb-c"},
		{name: "Length longer than origin with camel case", original: "AaaBbbCcc", max: 99, want: "AaaBbbCcc"},
		{name: "Length is 0 with camel case", original: "AaaBbbCcc", max: 0, want: "ABC"},
		{name: "Length is 0 with camel case, matching case", original: "aaaBbbCcc", max: 0, want: "aBC"},
		{name: "Partial abbreviation with camel case", original: "AaaBbbCcc", max: 8, want: "AaaBbbC"},
		{name: "Doesn't match wrong casing", original: "AaaBBbCcc", max: 0, want: "ABBbC"},
		{name: "Mixed camel case and non word seperators", original: "AaaBbb-ccc", max: 0, want: "AB-c"},
		{name: "Mixed camel case and non word seperators with same borders", original: "Aaa-Bbb-Ccc", max: 0, want: "A-B-C"},
		{name: "Real example, full short", original: "strategy-limited", max: 0, want: "stg-ltd"},
		{name: "Real example, shorter than total", original: "strategy-limited", max: 13, want: "strategy-ltd"},
		{name: "Real example, max same as shorted", original: "strategy-limited", max: 12, want: "strategy-ltd"},
		{name: "Real example, max on seperator", original: "strategy-limited", max: 9, want: "stg-ltd"},
		{name: "Real example, max shorter than first word", original: "strategy-limited", max: 6, want: "stg-ltd"},
		{name: "Real example, no short", original: "strategy-limited", max: 99, want: "strategy-limited"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShortenFromBack(matcher, tt.original, tt.max); got != tt.want {
				t.Errorf("ShortenFromBack('%s', %d) = '%v', want '%v'", tt.original, tt.max, got, tt.want)
			}
		})
	}
}

func Test_lastChar(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		remains string
		last    rune
	}{
		{name: "1", str: "string", remains: "strin", last: 'g'},
		{name: "2", str: "s", remains: "", last: 's'},
		{name: "3", str: "", remains: "", last: rune(0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := lastChar(tt.str)
			if got != tt.remains {
				t.Errorf("lastChar() got = '%v', want '%v'", got, tt.remains)
			}
			if got1 != tt.last {
				t.Errorf("lastChar() got1 = %v, want %v", got1, tt.last)
			}
		})
	}
}

func TestSequences_AddFront(t *testing.T) {
	seqs := Sequences{}

	seqs.AddFront("a")
	seqs.AddFront("b")
	seqs.AddFront("cd")

	assert.Equal(t, 3, len(seqs))
	assert.Equal(t, "cd", seqs[0])
	assert.Equal(t, "b", seqs[1])
	assert.Equal(t, "a", seqs[2])
	assert.Equal(t, "cdba", seqs.String())
	assert.Equal(t, 4, seqs.Len())
}

func TestSequences_AddBack(t *testing.T) {
	seqs := Sequences{}

	seqs.AddBack("a")
	seqs.AddBack("b")
	seqs.AddBack("cd")

	assert.Equal(t, 3, len(seqs))
	assert.Equal(t, "a", seqs[0])
	assert.Equal(t, "b", seqs[1])
	assert.Equal(t, "cd", seqs[2])
	assert.Equal(t, "abcd", seqs.String())
}

func TestNewSequences(t *testing.T) {
	tests := []struct {
		name string
		orig string
		want Sequences
	}{
		{name: "1", orig: "abc", want: Sequences{"abc"}},
		{name: "2", orig: "a-b-c", want: Sequences{"a", "-", "b", "-", "c"}},
		{name: "3", orig: "ABC", want: Sequences{"A", "B", "C"}},
		{name: "4", orig: "a--b--c", want: Sequences{"a", "-", "-", "b", "-", "-", "c"}},
		{name: "5", orig: "aa-bb", want: Sequences{"aa", "-", "bb"}},
		{name: "6", orig: "aaBbCc", want: Sequences{"aa", "Bb", "Cc"}},
		{name: "7", orig: "aa-Bb-cc", want: Sequences{"aa", "-", "Bb", "-", "cc"}},
		{name: "8", orig: "AaaBBbCcc", want: Sequences{"Aaa", "B", "Bb", "Ccc"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSequences(tt.orig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSequences() = %v (%d), want %v (%d)",
					[]string(got), len(got),
					[]string(tt.want), len(tt.want),
				)
			}
		})
	}
}

func Test_isTitleCase(t *testing.T) {
	assert.Equal(t, true, isTitleCase("Abc"))
	assert.Equal(t, true, isTitleCase("A"))
	assert.Equal(t, true, isTitleCase("ABC"))
	assert.Equal(t, false, isTitleCase("abc"))
	assert.Equal(t, false, isTitleCase("aBC"))
	assert.Equal(t, false, isTitleCase("a"))
}

func Test_first(t *testing.T) {
	assert.Equal(t, rune(0), first(""))
	assert.Equal(t, rune('a'), first("a"))
	assert.Equal(t, rune('c'), first("cba"))
	assert.Equal(t, rune('B'), first("Bac"))
	assert.Equal(t, false, unicode.IsTitle(first("")))
}
