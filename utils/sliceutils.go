package utils

import "slices"

func ContainsAll[S ~[]E, E comparable](s S, elements S) bool {
	for _, v := range elements {
		if !slices.Contains(s, v) {
			return false
		}
	}
	return true
}
