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
