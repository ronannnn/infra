package validator_test

import (
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

func TestValidatorCheck(t *testing.T) {
	translator, err := i18n.New(i18n.Cfg{BundleDir: "./testdata/"})
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

	errFields, _ := srv.Check(i18n.LanguageChinese, user)
	require.Equal(t, 2, len(errFields))
	ageErrField := errFields[0]
	require.Equal(t, "age", ageErrField.ErrorField)
	require.Equal(t, "Age必须小于或等于130", ageErrField.ErrorMsg)
	colorErrField := errFields[1]
	require.Equal(t, "favoriteColor", colorErrField.ErrorField)
	require.Equal(t, "FavoriteColor必须是一个有效的颜色", colorErrField.ErrorMsg)

}
