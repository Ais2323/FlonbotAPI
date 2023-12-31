package helper

import (
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	_MAX_WORD_LENGTH    = 50
	_MIN_CN_WORD_LENGTH = 2
	_MIN_EN_WORD_LENGTH = 1
	_EN_WORD_REGEXP     = "^[A-Za-z0-9_ \\s \\! \\' \\\" \\\\ \\+ \\* \\? \\^ \\$ \\[ \\] \\{ \\} \\( \\) \\| \\/]+$"
)

// 拆字
func SpliteWord(oriWord string) []string {
	isEnglish, err := regexp.MatchString(_EN_WORD_REGEXP, oriWord)
	result := []string{}
	if err != nil || isEnglish == false {
		result = splitChineseWord(oriWord)
	} else {
		result = splitEnglishWord(oriWord)
	}

	return result
}
func splitEnglishWord(oriWord string) []string {
	subWords := strings.Split(oriWord, " ")
	if len(subWords) > _MAX_WORD_LENGTH {
		// 隨機取50字
		maxEndIdx := len(subWords) - _MAX_WORD_LENGTH
		randSeed := rand.NewSource(time.Now().UnixNano())
		randData := rand.New(randSeed)
		targetIdx := randData.Intn(maxEndIdx)
		subWords = subWords[targetIdx : targetIdx+_MAX_WORD_LENGTH]
	}

	// 文字拆分
	result := []string{}
	if len(subWords) < _MIN_EN_WORD_LENGTH {
		return result
	}
	maxLength := len(subWords)
	// 目前長度 慢慢縮減字數
	nowLength := maxLength
	for nowLength >= _MIN_EN_WORD_LENGTH {
		for startIdx := 0; startIdx <= (maxLength - nowLength); startIdx++ {
			target := strings.Join(subWords[startIdx:startIdx+nowLength], " ")
			result = append(result, target)
		}
		nowLength--
	}

	return result
}
func splitChineseWord(oriWord string) []string {
	// 限制最大長度
	if len([]rune(oriWord)) > _MAX_WORD_LENGTH {
		// 隨機取50字
		maxEndIdx := len([]rune(oriWord)) - _MAX_WORD_LENGTH
		randSeed := rand.NewSource(time.Now().UnixNano())
		randData := rand.New(randSeed)
		targetIdx := randData.Intn(maxEndIdx)
		oriWord = string([]rune(oriWord)[targetIdx : targetIdx+_MAX_WORD_LENGTH])
	}
	// 文字拆分
	result := []string{}
	if len([]rune(oriWord)) < _MIN_EN_WORD_LENGTH {
		return result
	}
	// 分文字段落
	subString := strings.FieldsFunc(oriWord, splitSign)
	for _, word := range subString {
		maxLength := len([]rune(word))
		// 目前長度 慢慢縮減字數
		nowLength := maxLength
		for nowLength >= _MIN_CN_WORD_LENGTH {
			for startIdx := 0; startIdx <= (maxLength - nowLength); startIdx++ {
				target := string([]rune(word)[startIdx : startIdx+nowLength])
				result = append(result, target)
			}
			nowLength--
		}
	}
	// 按字串長度倒序排序
	sort.Slice(result, func(i, j int) bool {
		return len([]rune(result[i])) > len([]rune(result[j]))
	})
	return result
}

func splitSign(r rune) bool {
	signs := []rune{' ', ','}
	for _, sign := range signs {
		if sign == r {
			return true
		}
	}
	return false
}

// 跳脫字元
var html_escape_table = map[string]string{
	"&amp;":    "&",
	"&quot;":   "\"",
	"&apos;":   "'",
	"&gt;":     ">",
	"&lt;":     "<",
	"&circ;":   "^",
	"&tilde;":  "~",
	"&ensp;":   " ",
	"&emsp;":   "　",
	"&ndash;":  "–",
	"&lsquo;":  "‘",
	"&rsquo;":  "’",
	"&sbquo;":  ",",
	"&ldquo;":  "“",
	"&rdquo;":  "”",
	"&bdquo;":  "„",
	"&permil;": "‰",
	"&lsaquo;": "‹",
	"&rsaquo;": "›",
	"&euro;":   "€",
	"&copy;":   "©",
	"&reg;":    "®",
	"&deg;":    "°",
	"<br />":   "\n",
	"&colon;":  ":",
}

func HtmlEscapeSign(text string, unescape bool) string {
	for s_from, s_to := range html_escape_table {
		if unescape {
			text = strings.Replace(text, s_to, s_from, -1)
		} else {
			text = strings.Replace(text, s_from, s_to, -1)
		}

	}
	return text
}
