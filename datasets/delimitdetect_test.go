package datasets

import (
	"reflect"
	"testing"
)

func TestDetect(t *testing.T) {
	cases := []struct {
		in     string
		want   string
		fields []string
	}{
		{"one::two::three::four", "::", []string{"one", "two", "three", "four"}},
		{"a,b,c,d,e", ",", []string{"a", "b", "c", "d", "e"}},
		{"one, two, three", ", ", []string{"one", "two", "three"}},
		{"one two three four", " ", []string{"one", "two", "three", "four"}},
		{"a\tb\tc", "\t", []string{"a", "b", "c"}},
		{"a | b | c", " | ", []string{"a", "b", "c"}},
		{"a:::b:::c", ":::", []string{"a", "b", "c"}},
		{"::a::b::", "::", []string{"", "a", "b", ""}},
		{"one::two", "::", []string{"one", "two"}}, // single-occurrence fallback
		{"aXbXcXd", "X", []string{"a", "b", "c", "d"}},
		{"2024-01-02", "-", []string{"2024", "01", "02"}},
		// Nested delimiters are genuinely ambiguous: "=" occurs three times,
		// "&" only twice, so by "most common pattern" the answer is "=".
		{"k=1&k2=2&k3=3", "=", []string{"k", "1&k2", "2&k3", "3"}},
		{"the cat and the dog and the bird", " ",
			[]string{"the", "cat", "and", "the", "dog", "and", "the", "bird"}},
	}
	for _, c := range cases {
		got, ok := detectWithOptions(c.in, delimitOptions{})
		if !ok {
			t.Errorf("detectColumnDelimter(%q): no delimiter found, want %q", c.in, c.want)
			continue
		}
		if got.Delimiter != c.want {
			t.Errorf("detectColumnDelimter(%q) = %q, want %q", c.in, got.Delimiter, c.want)
		}
		if !reflect.DeepEqual(got.Fields, c.fields) {
			t.Errorf("detectColumnDelimter(%q) fields = %q, want %q", c.in, got.Fields, c.fields)
		}
	}
}

// Double-quoted sections are field content, never delimiter material.
func TestQuoted(t *testing.T) {
	cases := []struct {
		in     string
		want   string
		fields []string
		ok     bool
	}{
		{`a,"b,c",d`, ",", []string{"a", `"b,c"`, "d"}, true},
		{`one, "two, dos", three`, ", ", []string{"one", `"two, dos"`, "three"}, true},
		{`name;"last;first";age`, ";", []string{"name", `"last;first"`, "age"}, true},
		{`"a"::"b"::"c"`, "::", []string{`"a"`, `"b"`, `"c"`}, true},
		{`"a"::"b"`, "::", []string{`"a"`, `"b"`}, true}, // fallback, still quote-aware
		{`"one, two, three"`, "", nil, false},            // all one quoted field
		{`a,"b,c`, ",", []string{"a", `"b`, "c"}, true},  // unterminated quote is a plain rune
	}
	for _, c := range cases {
		got, ok := detectWithOptions(c.in, delimitOptions{})
		if ok != c.ok {
			t.Errorf("detectColumnDelimter(%q) ok = %v (%q), want %v", c.in, ok, got.Delimiter, c.ok)
			continue
		}
		if !ok {
			continue
		}
		if got.Delimiter != c.want {
			t.Errorf("detectColumnDelimter(%q) = %q, want %q", c.in, got.Delimiter, c.want)
		}
		if !reflect.DeepEqual(got.Fields, c.fields) {
			t.Errorf("detectColumnDelimter(%q) fields = %q, want %q", c.in, got.Fields, c.fields)
		}
	}
}

func TestNoDelimiter(t *testing.T) {
	for _, in := range []string{"", "x", "aaaa", "hello"} {
		if d := detectColumnDelimter(in, "\n"); d != "" {
			t.Errorf("detectColumnDelimter(%q) = %q, want no delimiter", in, d)
		}
	}
}

func TestPunctuationOnly(t *testing.T) {
	// By default a word-rune delimiter is allowed when nothing better repeats.
	if got, _ := detectWithOptions("aXbXcXd", delimitOptions{}); got.Delimiter != "X" {
		t.Errorf("default = %q, want %q", got.Delimiter, "X")
	}
	// With PunctuationOnly there is no candidate left at all.
	if got, ok := detectWithOptions("aXbXcXd", delimitOptions{PunctuationOnly: true}); ok {
		t.Errorf("PunctuationOnly = %q, want no delimiter", got.Delimiter)
	}
}

func TestUnicode(t *testing.T) {
	got, ok := detectWithOptions("één…twee…drie", delimitOptions{})
	if !ok || got.Delimiter != "…" {
		t.Errorf("got %q (%v), want \"…\"", got.Delimiter, ok)
	}
}

func BenchmarkDetect(b *testing.B) {
	s := ""
	for i := 0; i < 200; i++ {
		s += "field" + string(rune('a'+i%26)) + "::"
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detectColumnDelimter(s, "\n")
	}
}
