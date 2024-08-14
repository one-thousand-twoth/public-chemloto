package validator

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var ruTranslation = map[string]string{
	"Name": `"Имя"`,
	"Code": `"Код"`,
}
var regexAlphaNumSpace = regexp.MustCompile(`^[^\s][a-zA-Zа-яА-Я0-9- ]+[^\s]$`)

var (
	Ins *validator.Validate
	// uni   *ut.UniversalTranslator
	// Trans ut.Translator
)

func init() {
	Ins = validator.New()
	Ins.RegisterValidation("safeinput", ValidateAlphaNumSpace)
	// ru := ru.New()
	// uni = ut.New(ru, ru)
	Ins.RegisterTagNameFunc(func(field reflect.StructField) string {
		return ruTranslation[field.Name]
	})

	// Trans, _ = uni.GetTranslator("ru")

	// ru_translations.RegisterDefaultTranslations(Ins, Trans)
}

func ValidateAlphaNumSpace(fl validator.FieldLevel) bool {
	return regexAlphaNumSpace.MatchString(fl.Field().String())
}

func ValidationError(errs validator.ValidationErrors) []string {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("Поле %s является обязательным", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		case "safeinput":
			errMsgs = append(errMsgs, fmt.Sprintf("Поле %s должно содержать только буквы и цифры", err.Field()))
		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("Поле %s должно быть больше %s символов", err.Field(), err.Param()))
		case "gt":
			errMsgs = append(errMsgs, fmt.Sprintf("Поле %s должно быть больше %s", err.Field(), err.Param()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid: %s", err.Field(), err.Error()))
		}
	}

	return errMsgs
}
