package login_test

import (
	"testing"

	"github.com/ronannnn/infra/service/login"
	"github.com/stretchr/testify/require"
)

func TestPasswordHashAndComparison(t *testing.T) {
	password := "123456a"
	hashedPassword, err := login.HashPassword(password)
	println(hashedPassword)
	require.NoError(t, err)
	require.Truef(t, login.CheckPassword(hashedPassword, password), "password comparison error")
}
