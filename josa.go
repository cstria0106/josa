// josa는 한국어의 조사가 있는 동적인 문자열을 쉽게 만들 수 있도록 도와주는 라이브러리입니다.
package josa

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"
)

var functionMap = map[string]func(string) string{
	"은": EunNeun, "는": EunNeun, "은는": EunNeun, "는은": EunNeun, "은(는)": EunNeun, "는(은)": EunNeun, "은/는": EunNeun, "는/은": EunNeun,
	"을": EulReul, "를": EulReul, "을를": EulReul, "를을": EulReul, "을(를)": EulReul, "를(을)": EulReul, "을/를": EulReul, "를/을": EulReul,
	"이": IGa, "가": IGa, "이가": IGa, "가이": IGa, "이(가)": IGa, "가(이)": IGa, "이/가": IGa, "가/이": IGa,
	"과": GwaWa, "와": GwaWa, "과와": GwaWa, "와과": GwaWa, "과(와)": GwaWa, "와(과)": GwaWa, "과/와": GwaWa, "와/과": GwaWa,
	"으로": EuroRo, "로": EuroRo, "으로로": EuroRo, "로으로": EuroRo, "으로(로)": EuroRo, "로(으로)": EuroRo, "(으)로": EuroRo, "으로/로": EuroRo, "로/으로": EuroRo,
	"아": AYa, "야": AYa, "아야": AYa, "야아": AYa, "아(야)": AYa, "야(아)": AYa, "아/야": AYa, "야/아": AYa,
	"이여": IyeoYeo, "여": IyeoYeo, "이여여": IyeoYeo, "여이": IyeoYeo, "이여(여)": IyeoYeo, "여(이여)": IyeoYeo, "이여/여": IyeoYeo, "여/이여": IyeoYeo,
	"이랑": IrangRang, "랑": IrangRang, "이랑랑": IrangRang, "랑이랑": IrangRang, "이랑(랑)": IrangRang, "랑(이랑)": IrangRang, "이랑/랑": IrangRang, "랑/이랑": IrangRang,
}

var josaList []string
var josaListMutex sync.Mutex

// 조사 목록을 함수 맵의 키로 부터 얻는다.
func evaluateJosaList() {
	josaListMutex.Lock()
	defer josaListMutex.Unlock()

	if josaList != nil {
		return
	}

	josaList = make([]string, len(functionMap))

	i := 0
	for key := range functionMap {
		josaList[i] = key
		i++
	}
}

// HasJongseong 함수는 한글로 이루어진 word 문자열의 마지막 글자가 종성을 갖는지, 갖지 않는지 검사한다.
//
// 예:
//		j := HasJongseong("사과")
//		// j := false
//
//		j := HasJongseong("구름")
//		// j := true
func HasJongseong(word string) bool {
	var c, _ = utf8.DecodeLastRuneInString(word)
	return (c-0xac00)%28 > 0
}

// Josa 함수는 한글로 이루어진 word 문자열에 붙을 수 있는 올바른 조사와 true 값을 반환한다.
// 만약 인자 s로써 주어지는 조사 형식 문자열이 올바르지 않다면 반환값으로 빈 문자열과 false 를 반환한다.
//
// 예:
//		j := Josa("사과", "은는")
//		// j := "는"
//
//		j := Josa("구름", "와과")
//		// j := "과"
//
// 조사 형식 문자열로 다음과 같은 종류가 있다.
//
// 	"은", "는", "은는", "는은", "은(는)", "는(은)", "은/는", "는/은",
// 	"을", "를", "을를", "를을", "을(를)", "를(을)", "을/를", "를/을",
// 	"이", "가", "이가", "가이", "이(가)", "가(이)", "이/가", "가/이",
// 	"과", "와", "과와", "와과", "과(와)", "와(과)", "과/와", "와/과",
// 	"으로", "로", "으로로", "로으로", "으로(로)", "로(으로)", "(으)로", "으로/로", "로/으로",
// 	"아", "야", "아야", "야아", "아(야)", "야(아)", "아/야", "야/아",
// 	"이여", "여", "이여여", "여이", "이여(여)", "여(이여)", "이여/여", "여/이여",
// 	"이랑", "랑", "이랑랑", "랑이랑", "이랑(랑)", "랑(이랑)", "이랑/랑", "랑/이랑"
func Josa(word string, s string) (string, bool) {
	s = strings.TrimSpace(s)

	f, ok := functionMap[s]

	if ok {
		s = f(word)
	}

	return s, ok
}

// Concat 함수는 한글로 이루어진 word 문자열에 올바른 조사를 붙여 반환한다.
// 만약 인자 s로써 주어지는 조사 형식 문자열이 올바르지 않다면 두 문자열을 그대로 붙여서 반환한다.
// 이용할 수 있는 조사 형식 문자열은 Josa 함수와 동일하다.
//
// 예:
//		j := Concat("사과", "은는")
//		// j := "사과는"
//
//		j := Concat("구름", "와과")
//		// j := "구름과"
func Concat(word string, s string) string {
	josa, ok := Josa(word, s)

	if ok {
		s = josa
	}

	return word + s
}

// Format 함수는 format 문자열을 적절한 조사가 붙도록 포맷하여 반환한다.
// "{{"와 "}}" 사이에 조사 형식 문자열을 나타내어 사용한다.
// 예를 들어 어떤 단어에 대해 주격조사 이/가를 붙여 나타내기 위해서는 "{{이/가}}", "{{이가}}", "{{이(가)}}" 등을 사용할 수 있다.
// 이용할 수 있는 조사 형식 문자열은 Josa 함수와 동일하다.
//
// 예:
// 		s := Format("일찍 일어나는 새{{이/가}} 벌레{{을/를}} 잡는다.")
//		// s := "일찍 일어나는 새가 벌레를 잡는다."
//
//		s := Format("나{{은}} 생각한다. 고로 나{{은(는)}} 존재한다.")
//		// s := "나는 생각한다. 고로 나는 존재한다."
func Format(format string) string {
	evaluateJosaList()

	for _, josa := range josaList {
		r := regexp.MustCompile(fmt.Sprintf("([가-힣ㄱ-ㅣ]+)\\{\\{%s\\}\\}", regexp.QuoteMeta(josa)))

		matches := r.FindAllStringSubmatchIndex(format, -1)
		for _, match := range matches {
			word := format[match[2]:match[3]]
			concatenated := Concat(word, josa)

			formatted := format[:match[0]] + concatenated + format[match[1]:]
			format = formatted
		}
	}

	return format
}

// EunNeun 함수는 한글로 이루어진 word 문자열에 "은" 혹은 "는" 중 사용될 수 있는 올바른 조사를 반환한다.
func EunNeun(word string) string {
	if HasJongseong(word) {
		return "은"
	}

	return "는"
}

// EulReul 함수는 한글로 이루어진 word 문자열에 "을" 혹은 "를" 중 사용될 수 있는 올바른 조사를 반환한다.
func EulReul(word string) string {
	if HasJongseong(word) {
		return "을"
	}

	return "를"
}

// IGa 함수는 한글로 이루어진 word 문자열에 "이" 혹은 "가" 중 사용될 수 있는 올바른 조사를 반환한다.
func IGa(word string) string {
	if HasJongseong(word) {
		return "이"
	}

	return "가"
}

// GwaWa 함수는 한글로 이루어진 word 문자열에 "과" 혹은 "와" 중 사용될 수 있는 올바른 조사를 반환한다.
func GwaWa(word string) string {
	if HasJongseong(word) {
		return "과"
	}

	return "와"
}

// EuroRo 함수는 한글로 이루어진 word 문자열에 "으로" 혹은 "로" 중 사용될 수 있는 올바른 조사를 반환한다.
func EuroRo(word string) string {
	if HasJongseong(word) {
		return "으로"
	}

	return "로"
}

// AYa 함수는 한글로 이루어진 word 문자열에 "아" 혹은 "야" 중 사용될 수 있는 올바른 조사를 반환한다.
func AYa(word string) string {
	if HasJongseong(word) {
		return "아"
	}

	return "야"
}

// IyeoYeo 함수는 한글로 이루어진 word 문자열에 "이여" 혹은 "여" 중 사용될 수 있는 올바른 조사를 반환한다.
func IyeoYeo(word string) string {
	if HasJongseong(word) {
		return "이여"
	}

	return "여"
}

// IrangRang 함수는 한글로 이루어진 word 문자열에 "이랑" 혹은 "랑" 중 사용될 수 있는 올바른 조사를 반환한다.z
func IrangRang(word string) string {
	if HasJongseong(word) {
		return "이랑"
	}

	return "랑"
}
