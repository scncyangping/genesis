package util

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

// TitleLevenshteinFilter 相似度
func KwLevenshteinFilter(kw, body string) (string, error) {
	// find max length title
	keys := strings.FieldsFunc(body, func(r rune) bool {
		return r == ' ' || r == '?' || r == '？' || r == '!' || r == '：' || r == '(' || r == ')' || r == '-' ||
			r == ',' || r == '，' || r == ':' || r == '！' || r == '。' || r == '[' || r == ']' || r == '_'
	})

	key := lo.MaxBy(keys, func(a string, b string) bool {
		return len([]rune(a)) > len([]rune(b))
	})
	// step1. if contains => true
	if strings.Contains(key, kw) && (len([]rune(key)) < 20) {
		return key, nil
	}

	d1, d2 := lo.Difference([]rune(kw), []rune(key))
	ol := len([]rune(kw))
	kr := len(d1)
	kl := len(d2)

	if kr > ol/2 || kl > ol {
		return "", errors.New("not matching")
	}
	return key, nil
}
