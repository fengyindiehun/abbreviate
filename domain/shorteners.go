package domain

import (
	"strings"
	"unicode"
)

type parseState int

const (
	starting parseState = iota
	inWord
	betweenWords
)

type Sequences []string

func NewSequences(original string) Sequences {
	seq := Sequences{}

	set := ""
	for _, ch := range []rune(original) {
		if !unicode.IsLetter(ch) {
			if set != "" {
				seq.AddBack(set)
			}
			seq.AddBack(string(ch))
			set = ""
		} else if unicode.IsUpper(ch) {
			if set != "" {
				seq.AddBack(set)
			}
			set = string(ch)
		} else if !unicode.IsLetter(ch) {
			if set != "" {
				seq.AddBack(set)
			}
			seq.AddBack(string(ch))
			set = ""
		} else {
			set += string(ch)
		}
	}
	if set != "" {
		seq.AddBack(set)
	}

	return seq
}

func (all Sequences) String() string {
	str := ""
	for _, s := range all {
		str += s
	}

	return str
}

func (all *Sequences) AddFront(str string) {
	initial := Sequences{str}
	*all = append(initial, *all...)
}

func (all *Sequences) AddBack(str string) {
	*all = append(*all, str)
}

func (all Sequences) Len() int {
	l := 0
	for _, s := range all {
		l += len(s)
	}

	return l
}

// Shortener represents an algorithm that can be used to shorten a string
// by substituting words for abbreviations
type Shortener func(matcher Matcher, original string, max int) string

// ShortenFromBack discovers words using camel case and non letter characters,
// starting from the back until the string has less than 'max' characters
// or it can't shorten any more.
func ShortenFromBack(matcher *Matcher, original string, max int) string {
	if len(original) < max {
		return original
	}

	shortened := NewSequences(original)
	for pos := len(shortened) - 1; pos >= 0 && shortened.Len() > max; pos-- {
		str := shortened[pos]
		abbr := matcher.Match(strings.ToLower(str))
		if isTitleCase(str) {
			abbr = makeTitle(abbr)
		}
		shortened[pos] = abbr
	}

	return shortened.String()
}

func isTitleCase(str string) bool {
	ch := first(str)
	return unicode.IsUpper(ch)
}

func makeTitle(str string) string {
	if str == "" {
		return ""
	}

	ch := first(str)
	ch = unicode.ToUpper(ch)
	result := string(ch)
	if len(str) > 1 {
		result += str[1:]
	}

	return result
}

func first(str string) rune {
	if str == "" {
		return rune(0)
	}

	return []rune(str)[0]
}

func lastChar(str string) (string, rune) {
	l := len(str)

	switch l {
	case 0:
		return "", rune(0)

	case 1:
		return "", []rune(str)[0]
	}

	return str[0 : l-1], []rune(str)[l-1:][0]
}
