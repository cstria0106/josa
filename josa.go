package josa

import "strings"

func Concat(word string, s string) string {
	eunNeun := strings.ContainsAny(s, "은는")
	eulReul := strings.ContainsAny(s, "을를")
	iGa := strings.ContainsAny(s, "이가")

	switch {
	case eunNeun:
		return word + EunNeun(word)
	case eulReul:
		return word + EulReul(word)
	case iGa:
		return word + IGa(word)
	}

	return word + s
}

func EunNeun(word string) string {

}

func EulReul(word string) string {

}

func IGa(word string) string {

}
