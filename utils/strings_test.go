package utils_test

import (
	"testing"

	"github.com/ronannnn/infra/utils"
	"github.com/stretchr/testify/require"
)

func TestLowercaseFirstLetter(t *testing.T) {
	require.Equal(t, "hello.World", utils.LowercaseFirstLetter("Hello.World"))
	require.Equal(t, "helloWorld", utils.LowercaseFirstLetter("HelloWorld"))
}

func TestLowercaseFirstLetterAndJoin(t *testing.T) {
	require.Equal(t, "hello.world", utils.LowercaseFirstLetterAndJoin("Hello.World", "."))
	require.Equal(t, "helloWorld", utils.LowercaseFirstLetterAndJoin("HelloWorld", "."))
}

func TestFindAndRemoveIndexLikeStrings(t *testing.T) {
	strWithoutIndexes, numbers := utils.FindAndRemoveIndexLikeStrings("hello[1].world[2]")
	require.Equal(t, "hello.world", strWithoutIndexes)
	require.Equal(t, []int{1, 2}, numbers)
}

func TestRemoveIndexLikeStrings(t *testing.T) {
	strWithoutIndexes := utils.RemoveIndexLikeStrings("hello[1].world[2]")
	require.Equal(t, "hello.world", strWithoutIndexes)
}
