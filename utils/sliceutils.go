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

func Equals[S ~[]E, E comparable](s1, s2 S) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}
	return true
}
