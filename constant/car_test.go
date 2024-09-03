package constant_test

import (
	"regexp"
	"testing"

	"github.com/ronannnn/infra/constant"
)

func TestCarRegexp(t *testing.T) {
	carList := []string{
		"浙B55ZW12",
		"浙B55ZW1",
		"浙B12345",
		"浙B55ZW挂",
	}
	reg := regexp.MustCompile(constant.CnCarRegexp)
	for _, car := range carList {
		if !reg.MatchString(car) {
			t.Errorf("car %s should match regexp", car)
		}
	}
}
