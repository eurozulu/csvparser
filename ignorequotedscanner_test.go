package csvparser

import (
	"strings"
	"testing"
)

func TestIgnoreQuotedSplit(t *testing.T) {
	s := "a;b;"
	ss := IgnoreQuotedSplit(s, ";")
	if len(ss) != 3 {
		t.Errorf("xIgnoreQuotedSplit returned %d rows, expected 2", len(ss))
	}

	ss = strings.Split(s, ";")
	if len(ss) != 3 {
		t.Errorf("yIgnoreQuotedSplit returned %d rows, expected 2", len(ss))
	}

}
