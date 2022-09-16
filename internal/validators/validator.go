package validators

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	trans     ut.Translator
	validator *validator.Validate
}

type IValidator interface {
	Validate(i interface{}) error
}

func NewValidator() *Validator {
	e := en.New()
	universalTranslator := ut.New(e, e)
	trans, _ := universalTranslator.GetTranslator("en")
	validate := validator.New()

	// Register custom validation here
	validate.RegisterValidation("password", ValidatePassword)
	validate.RegisterValidation("email", ValidateEmail)

	en_translations.RegisterDefaultTranslations(validate, trans)
	RegisterValidatePasswordTranslation(validate, trans)

	return &Validator{
		trans:     trans,
		validator: validate,
	}
}

func (v *Validator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	if err == nil {
		return nil
	}
	errs := err.(validator.ValidationErrors)

	fmt.Println(errors.New(errs[0].Translate(v.trans)))

	return errors.New(errs[0].Translate(v.trans))
}

// Custom validation functions
// validate password
func ValidatePassword(field validator.FieldLevel) bool {
	password := field.Field().String()
	secure := true
	tests := []string{".{8,}", ".*[a-z]", ".*[A-Z]", ".*[0-9]"}
	for _, test := range tests {
		matched, _ := regexp.MatchString(test, password)
		if !matched {
			secure = false
			break
		}
	}
	return secure
}

// validate e-mail
func ValidateEmail(field validator.FieldLevel) bool {
	email := field.Field().String()
	_, err := mail.ParseAddress(email)
	return err == nil
}

// custom validation message translation for the password
func RegisterValidatePasswordTranslation(validate *validator.Validate, trans ut.Translator) {
	validate.RegisterTranslation("password", trans, func(ut ut.Translator) error {
		return ut.Add("password-validator", "{0} must contain at least one uppercase letter, lower case letter, number, and special character.", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password-validator", fe.Field())

		return t
	})
}
