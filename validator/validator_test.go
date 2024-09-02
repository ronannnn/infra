package validator_test

import (
	"context"
	"testing"

	"github.com/ronannnn/infra/i18n"
	"github.com/ronannnn/infra/validator"
	"github.com/stretchr/testify/require"
)

type User struct {
	FirstName     *string `json:"firstName" validate:"required"`
	LastName      *string `json:"lastName" validate:"required"`
	Age           *uint8  `json:"age" validate:"gte=0,lte=130"`
	Email         *string `json:"email" validate:"required,email"`
	Gender        *string `json:"gender" validate:"oneof=male female prefer_not_to"`
	FavoriteColor *string `json:"favoriteColor" validate:"iscolor"` // alias for 'hexcolor|rgb|rgba|hsl|hsla'
	Cards         []*Card `json:"cards" validate:"required,dive,required"`
}

type Card struct {
	Number *string `json:"number" validate:"required"`
	Size   *uint8  `json:"size" validate:"gte=1,lte=10"`
}

func TestValidatorCheck(t *testing.T) {
	translator, err := i18n.New(&i18n.Cfg{BundleDir: "./testdata/"})
	require.NoError(t, err)
	srv := validator.New(translator)

	firstName := "Badger"
	lastName := "Smith"
	age := uint8(135)
	email := "Badger.Smith@gmail.com"
	gender := "male"
	favoriteColor := "#000-"

	cardNumber := "1234567890"
	cardSize := uint8(0)

	user := &User{
		FirstName:     &firstName,
		LastName:      &lastName,
		Age:           &age,
		Email:         &email,
		Gender:        &gender,
		FavoriteColor: &favoriteColor,
		Cards: []*Card{
			{Number: &cardNumber, Size: &cardSize},
		},
	}

	errFields, _ := srv.Check(context.Background(), i18n.LanguageChinese, user)
	require.Equal(t, 3, len(errFields))
	ageErrField := errFields[0]
	require.Equal(t, "user.age", ageErrField.ErrorWithNamespace)
	require.Equal(t, "age", ageErrField.ErrorField)
	require.Equal(t, "年龄必须小于或等于130", ageErrField.ErrorMsg)
	colorErrField := errFields[1]
	require.Equal(t, "user.favoriteColor", colorErrField.ErrorWithNamespace)
	require.Equal(t, "favoriteColor", colorErrField.ErrorField)
	require.Equal(t, "最喜欢的颜色必须是一个有效的颜色", colorErrField.ErrorMsg)
	cardSizeField := errFields[2]
	require.Equal(t, "user.cards.size", cardSizeField.ErrorWithNamespace)
	require.Equal(t, "size", cardSizeField.ErrorField)
	require.Equal(t, "大小必须大于或等于1", cardSizeField.ErrorMsg)
}

func TestValidatorCheckPartial(t *testing.T) {
	translator, err := i18n.New(&i18n.Cfg{BundleDir: "./testdata/"})
	require.NoError(t, err)
	srv := validator.New(translator)

	age := uint8(135)
	email := "Badger.Smith@gmail.com"

	cardNumber := "1234567890"
	cardSize := uint8(0)

	user := &User{
		Age:   &age,
		Email: &email,
		Cards: []*Card{
			{Number: &cardNumber, Size: &cardSize},
		},
	}

	errFields, _ := srv.CheckPartial(context.Background(), i18n.LanguageChinese, user)
	require.Equal(t, 2, len(errFields))
	ageErrField := errFields[0]
	require.Equal(t, "user.age", ageErrField.ErrorWithNamespace)
	require.Equal(t, "age", ageErrField.ErrorField)
	require.Equal(t, "年龄必须小于或等于130", ageErrField.ErrorMsg)
	cardSizeField := errFields[1]
	require.Equal(t, "user.cards.size", cardSizeField.ErrorWithNamespace)
	require.Equal(t, "size", cardSizeField.ErrorField)
	require.Equal(t, "大小必须大于或等于1", cardSizeField.ErrorMsg)
}
