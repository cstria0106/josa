package josa

import (
	"fmt"
	"testing"
)

type josaTestSet struct {
	word         string
	hasJongseong bool
}

var testSets = []josaTestSet{
	{"사과", false},
	{"과자", false},
	{"코드", false},
	{"구름", true},
	{"곰", true},
	{"맛있는 것", true},
}

var josaLists = [][]string{
	{"은", "는", "은(는)", "는(은)"},
	{"을", "를", "을(를)", "를(을)"},
	{"이", "가", "이(가)", "가(이)"},
	{"과", "와", "과(와)", "와(과)"},
	{"으로", "로", "으로(로)", "로(으로)", "(으)로"},
}

var josaListAfterJongseong = []string{"은", "을", "이", "과", "으로"}
var josaListAfterNotJongseong = []string{"는", "를", "가", "와", "로"}

func TestJongseong(t *testing.T) {
	for _, set := range testSets {
		actual := HasJongseong(set.word)
		if set.hasJongseong && !actual {
			t.Errorf("word '%s' was expected to have jongseong. but actually not", set.word)
		} else if !set.hasJongseong && actual {
			t.Errorf("word '%s' was not expected to have jongseong. but actually it was", set.word)
		}
	}
}

func TestConcat(t *testing.T) {
	for _, set := range testSets {
		for i, josaList := range josaLists {
			for _, josa := range josaList {
				result := Concat(set.word, josa)

				var expected string

				if set.hasJongseong {
					expected = josaListAfterJongseong[i]
				} else {
					expected = josaListAfterNotJongseong[i]
				}

				if result != fmt.Sprintf("%s%s", set.word, expected) {
					t.Errorf("concatenated string for word '%s' and josa '%s' is expected to be '%s%s'. but actual value was '%s'", set.word, josa, set.word, expected, result)
				}
			}
		}
	}
}

func TestFormat(t *testing.T) {
	source := "구름{{와과}} 고양이{{와(과)}} 멍멍이{{는은}} 화성{{(으)로}} 기나긴 여행을 떠난다"
	formatted := Format(source)

	if formatted != "구름과 고양이와 멍멍이는 화성으로 기나긴 여행을 떠난다" {
		t.Errorf(formatted)
	}
}
