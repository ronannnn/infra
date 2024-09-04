package validator_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/ronannnn/infra/i18n"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/validator"
	"github.com/stretchr/testify/require"
)

type User struct {
	models.Base
	FirstName     *string             `json:"firstName" validate:"required,not_blank"`
	LastName      *string             `json:"lastName" validate:"required"`
	Age           *uint8              `json:"age" validate:"gte=0,lte=130"`
	Email         *string             `json:"email" validate:"required,email"`
	Gender        *string             `json:"gender" validate:"oneof=male female prefer_not_to"`
	FavoriteColor *string             `json:"favoriteColor" validate:"iscolor"` // alias for 'hexcolor|rgb|rgba|hsl|hsla'
	GrossWt       *models.DecimalSafe `json:"grossWt" validate:"required,d_gt=1"`
	NetWt         *models.DecimalSafe `json:"netWt" validate:"required,d_lt=1"`
	TotalWt       *models.DecimalSafe `json:"totalWt" validate:"required,d_decimal_len_lte=2"`
	CarNo         *string             `json:"carNo" validate:"required,cn_car"`
	Cards         []*Card             `json:"cards" validate:"required,dive,required"`
}

type Card struct {
	Sign
	Number *string `json:"number" validate:"required"`
	Size   *uint8  `json:"size" validate:"gte=1,lte=10"`
}

type Sign struct {
	Username *string             `json:"username" validate:"required,min=1"`
	Password *string             `json:"password" validate:"required"`
	Wt       *models.DecimalSafe `json:"wt" validate:"required,d_gt=0"`
}

func TestValidatorCheck(t *testing.T) {
	translator, err := i18n.New(&i18n.Cfg{BundleDir: "./testdata/"})
	require.NoError(t, err)
	srv := validator.New(translator)

	var user User
	jsonData, err := os.ReadFile("./testdata/validator/user1.json")
	require.NoError(t, err)
	err = json.Unmarshal(jsonData, &user)
	require.NoError(t, err)

	errFields, _ := srv.Check(context.Background(), i18n.LanguageChinese, user)
	require.Equal(t, 7, len(errFields))
	firstNameErrField := errFields[0]
	require.Equal(t, "user.firstName", firstNameErrField.ErrorWithNamespace)
	require.Equal(t, "firstName", firstNameErrField.ErrorField)
	require.Equal(t, "名不能为空", firstNameErrField.ErrorMsg)
	ageErrField := errFields[1]
	require.Equal(t, "user.age", ageErrField.ErrorWithNamespace)
	require.Equal(t, "age", ageErrField.ErrorField)
	require.Equal(t, "年龄必须小于或等于130", ageErrField.ErrorMsg)
	colorErrField := errFields[2]
	require.Equal(t, "user.favoriteColor", colorErrField.ErrorWithNamespace)
	require.Equal(t, "favoriteColor", colorErrField.ErrorField)
	require.Equal(t, "最喜欢的颜色必须是一个有效的颜色", colorErrField.ErrorMsg)
	grossWtErrField := errFields[3]
	require.Equal(t, "user.grossWt", grossWtErrField.ErrorWithNamespace)
	require.Equal(t, "grossWt", grossWtErrField.ErrorField)
	require.Equal(t, "毛重必须大于1", grossWtErrField.ErrorMsg)
	netWtErrField := errFields[4]
	require.Equal(t, "user.netWt", netWtErrField.ErrorWithNamespace)
	require.Equal(t, "netWt", netWtErrField.ErrorField)
	require.Equal(t, "净重必须小于1", netWtErrField.ErrorMsg)
	totalWtErrField := errFields[5]
	require.Equal(t, "user.totalWt", totalWtErrField.ErrorWithNamespace)
	require.Equal(t, "totalWt", totalWtErrField.ErrorField)
	require.Equal(t, "总重小数点位数必须小于或等于2", totalWtErrField.ErrorMsg)
	cardSizeField := errFields[6]
	require.Equal(t, "user.cards.size", cardSizeField.ErrorWithNamespace)
	require.Equal(t, "size", cardSizeField.ErrorField)
	require.Equal(t, "大小必须大于或等于1", cardSizeField.ErrorMsg)
}

func TestValidatorCheckPartial(t *testing.T) {
	translator, err := i18n.New(&i18n.Cfg{BundleDir: "./testdata/"})
	require.NoError(t, err)
	srv := validator.New(translator)

	var user User
	jsonData, err := os.ReadFile("./testdata/validator/user2.json")
	require.NoError(t, err)
	err = json.Unmarshal(jsonData, &user)
	require.NoError(t, err)

	errFields, _ := srv.CheckPartial(context.Background(), i18n.LanguageChinese, user)
	require.Equal(t, 3, len(errFields))
	ageErrField := errFields[0]
	require.Equal(t, "user.age", ageErrField.ErrorWithNamespace)
	require.Equal(t, "age", ageErrField.ErrorField)
	require.Equal(t, "年龄必须小于或等于130", ageErrField.ErrorMsg)
	wtErrField := errFields[1]
	require.Equal(t, "user.cards.sign.wt", wtErrField.ErrorWithNamespace)
	require.Equal(t, "wt", wtErrField.ErrorField)
	require.Equal(t, "重量必须大于0", wtErrField.ErrorMsg)
	cardSizeField := errFields[2]
	require.Equal(t, "user.cards.size", cardSizeField.ErrorWithNamespace)
	require.Equal(t, "size", cardSizeField.ErrorField)
	require.Equal(t, "大小必须大于或等于1", cardSizeField.ErrorMsg)
}
