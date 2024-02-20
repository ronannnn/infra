package user

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPasswordHashAndComparison(t *testing.T) {
	password := "123456a"
	hashedPassword, err := HashPassword(password)
	println(hashedPassword)
	require.NoError(t, err)
	require.Truef(t, CheckPassword(hashedPassword, password), "password comparison error")
}
