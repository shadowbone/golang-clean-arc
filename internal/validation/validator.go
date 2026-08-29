package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"golang-rest-api/internal/repository"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	// Pakai tag `json` sebagai nama field di pesan error,
	// supaya client melihat "email", bukan "Email"
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})
}

// Struct memvalidasi dan mengembalikan error yang cocok dengan
// errors.Is(err, repository.ErrValidation)
func Struct(s any) error {
	if err := validate.Struct(s); err != nil {
		var ve validator.ValidationErrors
		if !errors.As(err, &ve) {
			return err
		}

		msgs := make([]string, 0, len(ve))
		for _, fe := range ve {
			msgs = append(msgs, pesan(fe))
		}

		return repository.Validationf("%s", strings.Join(msgs, "; "))
	}
	return nil
}

func pesan(fe validator.FieldError) string {
	f := fe.Field()

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s wajib diisi", f)
	case "email":
		return fmt.Sprintf("%s harus berupa email yang valid", f)
	case "min":
		return fmt.Sprintf("%s minimal %s karakter", f, fe.Param())
	case "max":
		return fmt.Sprintf("%s maksimal %s karakter", f, fe.Param())
	case "gt":
		return fmt.Sprintf("%s harus lebih dari %s", f, fe.Param())
	case "gte":
		return fmt.Sprintf("%s minimal %s", f, fe.Param())
	case "lte":
		return fmt.Sprintf("%s maksimal %s", f, fe.Param())
	case "oneof":
		return fmt.Sprintf("%s harus salah satu dari: %s", f, fe.Param())
	default:
		return fmt.Sprintf("%s tidak valid (%s)", f, fe.Tag())
	}
}
