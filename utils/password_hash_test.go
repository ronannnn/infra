package utils_test

import (
	"testing"

	"github.com/ronannnn/infra/utils"
	"github.com/stretchr/testify/require"
)

func TestPasswordHashAndComparison(t *testing.T) {
	password := "123456a"
	hashedPassword, err := utils.HashPassword(password)
	println(hashedPassword)
	require.NoError(t, err)
	require.Truef(t, utils.CheckPassword(hashedPassword, password), "password comparison error")
}
