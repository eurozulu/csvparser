package utils

import "strings"

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
