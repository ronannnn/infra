package validator_test

import (
	"context"
	"testing"

	"github.com/ronannnn/infra/i18n"
	"github.com/ronannnn/infra/validator"
	"github.com/stretchr/testify/require"
)

type User struct {
	FirstName     string `json:"firstName" validate:"required"`
	LastName      string `json:"lastName" validate:"required"`
	Age           uint8  `json:"age" validate:"gte=0,lte=130"`
	Email         string `json:"email" validate:"required,email"`
	Gender        string `json:"gender" validate:"oneof=male female prefer_not_to"`
	FavoriteColor string `json:"favoriteColor" validate:"iscolor"` // alias for 'hexcolor|rgb|rgba|hsl|hsla'
}

type Card struct {
	Number string `json:"number" validate:"required"`
	Size   uint8  `json:"size" validate:"gte=1,lte=10"`
}

func TestValidatorCheck(t *testing.T) {
	translator, err := i18n.New(&i18n.Cfg{BundleDir: "./testdata/"})
	require.NoError(t, err)
	srv := validator.New(translator)

	user := &User{
		FirstName:     "Badger",
		LastName:      "Smith",
		Age:           135,
		Gender:        "male",
		Email:         "Badger.Smith@gmail.com",
		FavoriteColor: "#000-",
	}

	errFields, _ := srv.Check(context.Background(), i18n.LanguageChinese, user)
	require.Equal(t, 2, len(errFields))
	ageErrField := errFields[0]
	require.Equal(t, "user.age", ageErrField.ErrorWithNamespace)
	require.Equal(t, "age", ageErrField.ErrorField)
	require.Equal(t, "年龄必须小于或等于130", ageErrField.ErrorMsg)
	colorErrField := errFields[1]
	require.Equal(t, "user.favoriteColor", colorErrField.ErrorWithNamespace)
	require.Equal(t, "favoriteColor", colorErrField.ErrorField)
	require.Equal(t, "最喜欢的颜色必须是一个有效的颜色", colorErrField.ErrorMsg)
}
