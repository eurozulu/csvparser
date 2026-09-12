package utils

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
)

func DoubleQuote(s string) string {
	if IsDoubleQuoted(s) {
		return s
	}
	sb := &strings.Builder{}
	sb.WriteRune('"')
	sb.WriteString(EscapeQuotes(s))
	sb.WriteRune('"')
	return sb.String()
}

func UndoubleQuote(s string) string {
	if !IsDoubleQuoted(s) {
		return s
	}
	return UnescapeQuotes(strings.Trim(s, "\""))
}

func IsDoubleQuoted(s string) bool {
	return strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")
}

func EscapeQuotes(s string) string {
	sb := &strings.Builder{}
	for len(s) > 0 {
		i := strings.Index(s, "\"")
		if i == -1 {
			sb.WriteString(s)
			break
		}
		sb.WriteString(s[:i])
		if i > 0 && s[i-1] != '\\' {
			sb.WriteRune('\\')
		}
		sb.WriteString(s[i : i+1])
		s = s[i+1:]
	}
	return sb.String()
}

func UnescapeQuotes(s string) string {
	return strings.ReplaceAll(s, "\\\"", "\"")
}

func IsIdentifier(s string, include ...rune) bool {
	r := []rune(s)
	if len(r) == 0 {
		return false
	}
	if !unicode.IsLetter(r[0]) && !slices.Contains(include, r[0]) {
		return false
	}
	for _, c := range r[1:] {
		if !unicode.IsLetter(c) &&
			!unicode.IsNumber(c) &&
			!slices.Contains(include, c) {
			return false
		}
	}
	return true
}

func RemoveEmptyStrings[S []E, E ~string](s S) S {
	var sz S
	for _, ss := range s {
		if ss == "" {
			continue
		}
		sz = append(sz, ss)
	}
	return sz
}

func StringerStrings[S []E, E fmt.Stringer](s S) []string {
	var sz []string
	for _, ss := range s {
		sz = append(sz, ss.String())
	}
	return sz
}
