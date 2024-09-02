package utils

import (
	"fmt"
	"regexp"
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

// 删除字符串中的索引字符串，例如：hello[1].world[2] -> hello.world
func RemoveIndexLikeStrings(input string) string {
	re := regexp.MustCompile(`\[(\d+)\]`)
	return re.ReplaceAllString(input, "")
}

// FindAndRemoveIndexLikeStrings 查找并删除字符串中的索引字符串，例如：hello[1].world[2] -> hello.world
func FindAndRemoveIndexLikeStrings(input string) (strWithoutIndexes string, numbers []int) {
	re := regexp.MustCompile(`\[(\d+)\]`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		indexStr := match[1]
		index := parseToInt(indexStr)
		numbers = append(numbers, index)
	}

	// Remove the index substrings from the input string
	strWithoutIndexes = re.ReplaceAllString(input, "")

	return
}

func parseToInt(s string) int {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		fmt.Println("Error parsing string to int:", err)
	}
	return result
}

func JoinNonEmptyStrings(sep string, strList ...string) string {
	var nonEmptyStrings []string
	for _, str := range strList {
		if str != "" {
			nonEmptyStrings = append(nonEmptyStrings, str)
		}
	}
	return strings.Join(nonEmptyStrings, sep)
}
