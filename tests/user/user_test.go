package user

import (
	"gorm.io/gorm/utils/tests"
	"main/utils"
	"reflect"
	"testing"
)

func TestRequestValidator(t *testing.T) {
	data := map[string]interface{}{
		"email":  "rianariansaputra@gmail.com",
		"email1": "rianariansaputragmail.com",
	}
	email, err := utils.RequestValidator("email", data, utils.RequestDataValidator{
		Type: reflect.String, Max: 10, Min: 5, Null: false, Email: true,
	})
	tests.AssertEqual(t, email, "rianariansaputra@gmail.com")
	tests.AssertEqual(t, err, nil)

	email1, err := utils.RequestValidator("email1", data, utils.RequestDataValidator{
		Type: reflect.String, Max: 10, Min: 5, Null: false, Email: true,
	})
	tests.AssertEqual(t, email1, nil)
}
