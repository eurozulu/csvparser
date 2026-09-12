package csvfile

import "unicode"

// delimitOptions tunes the detector. The zero value is usable and gives the defaults
// described on each field.
type delimitOptions struct {
	// MaxLen is the longest delimiter considered, in runes.
	// 0 means "auto": min(16, len(s)/2).
	MaxLen int

	// MinCount is the minimum number of non-overlapping occurrences a
	// candidate needs. 0 means 2, so a pattern has to actually repeat before
	// it is called a delimiter.
	MinCount int

	// PunctuationOnly rejects any candidate containing a letter or a digit.
	// Off by default: a punctuation-only delimiter is already preferred (see
	// DetectWithOptions), but with this off "aXbXcXd" can still report "X".
	PunctuationOnly bool
}

// Result describes a detected delimiter.
type delimitResult struct {
	Delimiter string   // the delimiter itself
	Count     int      // non-overlapping occurrences used for the split
	Fields    []string // the string split on those occurrences
	Score     int      // Count * len([]rune(Delimiter)) -- runes covered
}

// DetectColumnDelimter returns the most likely delimiter in s, or ok == false when s shows no
// repeated separating pattern.
func DetectColumnDelimter(s string) (string, bool) {
	r, ok := detectWithOptions(s, delimitOptions{})
	return r.Delimiter, ok
}

// Split splits s on its detected delimiter. When no delimiter is found the
// whole string is returned as a single field.
//func Split(s string) []string {
//	if r, ok := DetectWithOptions(s, delimitOptions{}); ok {
//		return r.Fields
//	}
//	return []string{s}
//}

// detectWithOptions is DetectColumnDelimter with explicit tuning.
//
// The algorithm:
//
//  0. Mark every rune that sits inside a pair of double quotes, quotes
//     included. Those runes take no part in any of the steps below.
//  1. Record the start positions of every substring of length 1..MaxLen that
//     contains no marked rune.
//  2. For each candidate, walk its positions left to right and keep the
//     non-overlapping ones -- that is how many times it could actually be used
//     as a separator.
//  3. Reject a candidate unless every *interior* field it produces is
//     non-empty. This is what makes "one::two::three" answer "::" and not ":":
//     splitting on ":" would leave an empty field between each pair of colons.
//  4. Rank the survivors. A candidate free of letters and digits always beats
//     one containing them, because a pattern that mixes punctuation with field
//     text ("d, " in "red, green, red, blue") is a coincidence, not a separator.
//     Within that, score by the number of runes covered (count * length), so
//     " | " beats " " on "a | b | c" while a genuinely more frequent pattern
//     still wins. Ties break on higher count, then longer delimiter, then
//     earliest occurrence, which keeps the result deterministic.
//  5. If no punctuation-only pattern repeats, fall back to the longest single
//     run of punctuation with text on both sides, so "one::two" answers "::"
//     rather than the letter "o" that happens to occur twice. Only when that
//     fails too is a repeated word-rune pattern accepted ("aXbXcXd" -> "X").
//
// Cost is O(n * MaxLen) time and space: each of the n*MaxLen substrings is
// recorded once, and validating a candidate is linear in its own occurrences.
func detectWithOptions(s string, opt delimitOptions) (delimitResult, bool) {
	runes := []rune(s)
	n := len(runes)

	maxLen := opt.MaxLen
	if maxLen <= 0 {
		maxLen = n / 2
		if maxLen > 16 {
			maxLen = 16
		}
	}
	minCount := opt.MinCount
	if minCount <= 0 {
		minCount = 2
	}
	if n < 2 || maxLen < 1 {
		return delimitResult{}, false
	}

	// 0. Everything between a pair of double quotes is field content.
	quoted := quotedMask(runes)

	// 1. Collect start positions of every candidate substring.
	positions := make(map[string][]int, n*maxLen)
	for i := 0; i < n; i++ {
		if quoted[i] {
			continue
		}
		for l := 1; l <= maxLen && i+l <= n; l++ {
			if quoted[i+l-1] {
				break // any longer candidate from i reaches into the quotes too
			}
			if opt.PunctuationOnly && isWord(runes[i+l-1]) {
				break // any longer candidate from i contains this rune too
			}
			cand := string(runes[i : i+l])
			positions[cand] = append(positions[cand], i)
		}
	}

	// 2-4. Validate and score.
	var best delimitResult
	var bestFirst int
	found := false
	for cand, starts := range positions {
		if len(starts) < minCount {
			continue
		}
		l := len([]rune(cand))
		hits := nonOverlapping(starts, l)
		if len(hits) < minCount {
			continue
		}
		if !separates(hits, l) {
			continue
		}
		r := delimitResult{
			Delimiter: cand,
			Count:     len(hits),
			Fields:    fields(runes, hits, l),
			Score:     len(hits) * l,
		}
		if !found || better(r, best, hits[0], bestFirst) {
			best, bestFirst, found = r, hits[0], true
		}
	}
	// 5. A repeating punctuation pattern is the strongest evidence there is.
	if found && isPunct(best.Delimiter) {
		return best, true
	}
	// Otherwise a lone punctuation run beats a word pattern that merely recurs.
	if r, ok := punctuationFallback(runes, quoted); ok {
		return r, true
	}
	return best, found
}

// nonOverlapping keeps the greedy left-to-right subset of starts whose
// occurrences of length l do not overlap.
func nonOverlapping(starts []int, l int) []int {
	hits := make([]int, 0, len(starts))
	end := 0
	for _, p := range starts {
		if p >= end {
			hits = append(hits, p)
			end = p + l
		}
	}
	return hits
}

// separates reports whether the occurrences carve the string into fields whose
// interiors are all non-empty. Leading and trailing fields may be empty, so
// "::a::b::" still yields "::".
func separates(hits []int, l int) bool {
	for i := 1; i < len(hits); i++ {
		if hits[i] <= hits[i-1]+l {
			return false // two delimiters back to back: empty interior field
		}
	}
	return true
}

func fields(runes []rune, hits []int, l int) []string {
	out := make([]string, 0, len(hits)+1)
	prev := 0
	for _, p := range hits {
		out = append(out, string(runes[prev:p]))
		prev = p + l
	}
	return append(out, string(runes[prev:]))
}

// better ranks a against b: punctuation-only first, then coverage, then
// occurrence count, then length, then the earlier first occurrence.
func better(a, b delimitResult, aFirst, bFirst int) bool {
	if pa, pb := isPunct(a.Delimiter), isPunct(b.Delimiter); pa != pb {
		return pa
	}
	if a.Score != b.Score {
		return a.Score > b.Score
	}
	if a.Count != b.Count {
		return a.Count > b.Count
	}
	if la, lb := len([]rune(a.Delimiter)), len([]rune(b.Delimiter)); la != lb {
		return la > lb
	}
	return aFirst < bFirst
}

// punctuationFallback handles strings with a single separator, e.g. "one::two".
// Quoted runes count as text, so they bound a run without ever joining it.
func punctuationFallback(runes []rune, quoted []bool) (delimitResult, bool) {
	n := len(runes)
	text := func(i int) bool { return quoted[i] || isWord(runes[i]) }
	var best delimitResult
	for i := 0; i < n; {
		if text(i) {
			i++
			continue
		}
		j := i
		for j < n && !text(j) {
			j++
		}
		if i > 0 && j < n && j-i > best.Score { // text on both sides
			best = delimitResult{
				Delimiter: string(runes[i:j]),
				Count:     1,
				Fields:    []string{string(runes[:i]), string(runes[j:])},
				Score:     j - i,
			}
		}
		i = j
	}
	return best, best.Score > 0
}

// quotedMask marks every rune that lies inside a pair of double quotes, the
// quotes themselves included. A quote with no partner later in the string is
// left unmarked and treated as an ordinary rune, so a stray quote cannot swallow
// the rest of the line.
func quotedMask(runes []rune) []bool {
	mask := make([]bool, len(runes))
	for i := 0; i < len(runes); i++ {
		if runes[i] != '"' {
			continue
		}
		j := i + 1
		for j < len(runes) && runes[j] != '"' {
			j++
		}
		if j == len(runes) {
			break // unterminated: this quote is just a character
		}
		for k := i; k <= j; k++ {
			mask[k] = true
		}
		i = j
	}
	return mask
}

func isWord(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) }

// isPunct reports whether s is made purely of separator-ish runes.
func isPunct(s string) bool {
	for _, r := range s {
		if isWord(r) {
			return false
		}
	}
	return true
}
