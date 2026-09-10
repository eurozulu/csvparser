package main

import "strings"

type Arguments []string

func (a Arguments) ContainsFlag(name string) bool {
	return a.IndexOfFlag(name) != -1
}

func (a Arguments) IndexOfFlag(name string) int {
	for i := 0; i < len(a); i++ {
		if !strings.HasPrefix(a[i], "-") {
			continue
		}
		f := strings.TrimLeft(a[i], "-")
		if f == name {
			return i
		}
	}
	return -1
}

func (a Arguments) Flag(name string) (string, bool) {
	i := a.IndexOfFlag(name)
	if i == -1 {
		return "", false
	}
	if i+1 >= len(a) || strings.HasPrefix(a[i+1], "-") {
		return "", true
	}
	return a[i+1], true
}

func (a Arguments) Parameters() []string {
	var params []string
	for i := 0; i < len(a); i++ {
		if !strings.HasPrefix(a[i], "-") {
			params = append(params, a[i])
			continue
		}
		if i+1 < len(a) && !strings.HasPrefix(a[i+1], "-") {
			i++
		}
	}
	return params
}
