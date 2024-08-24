package utils

import (
	"strings"
	"unicode"
)

func LowercaseFirstLetter(word string) string {
	if len(word) == 0 {
		return word
	}
	r := []rune(word)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// LowercaseFirstLetterAndJoin 将字符串中的每个单词的首字母转换为小写，并使用指定的分隔符连接
func LowercaseFirstLetterAndJoin(str string, sep string) string {
	parts := strings.Split(str, sep)

	// 转换每个单词的首字母为小写
	for i, part := range parts {
		parts[i] = LowercaseFirstLetter(part)
	}

	return strings.Join(parts, sep)
}
