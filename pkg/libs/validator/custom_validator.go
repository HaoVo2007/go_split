package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

// CustomValidator chứa validator instance và translator
type CustomValidator struct {
	Validator  *validator.Validate
	Translator ut.Translator
}

// NewCustomValidator creates custom validator with English messages
func NewCustomValidator() *CustomValidator {
	// Tạo validator
	v := validator.New()

	// Register tag name function to use JSON tag
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Tạo translator
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ := uni.GetTranslator("en")

	// Đăng ký translations mặc định
	enTranslations.RegisterDefaultTranslations(v, trans)

	// Register custom translations
	v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	v.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	v.RegisterTranslation("min", trans, func(ut ut.Translator) error {
		return ut.Add("min", "{0} must be at least {1} characters", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("min", fe.Field(), fe.Param())
		return t
	})

	v.RegisterTranslation("max", trans, func(ut ut.Translator) error {
		return ut.Add("max", "{0} must not exceed {1} characters", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("max", fe.Field(), fe.Param())
		return t
	})

	v.RegisterTranslation("oneof", trans, func(ut ut.Translator) error {
		return ut.Add("oneof", "{0} must be one of {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("oneof", fe.Field(), fe.Param())
		return t
	})

	v.RegisterTranslation("eqfield", trans, func(ut ut.Translator) error {
		return ut.Add("eqfield", "{0} must be equal to {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("eqfield", fe.Field(), fe.Param())
		return t
	})

	// Register custom validation for password strength
	v.RegisterValidation("strong_password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		// Check for at least 1 uppercase, 1 lowercase, 1 number
		hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
		hasNumber := strings.ContainsAny(password, "0123456789")

		return hasUpper && hasLower && hasNumber
	})

	v.RegisterTranslation("strong_password", trans, func(ut ut.Translator) error {
		return ut.Add("strong_password", "{0} must contain at least 1 uppercase, 1 lowercase, and 1 number", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("strong_password", fe.Field())
		return t
	})

	return &CustomValidator{
		Validator:  v,
		Translator: trans,
	}
}

// defaultValidator implement gin binding.Validator interface
type defaultValidator struct {
	Validate *validator.Validate
}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	return v.Validate.Struct(obj)
}

func (v *defaultValidator) Engine() interface{} {
	return v.Validate
}

// ValidateAndTranslate performs validation and returns translated error messages
func (cv *CustomValidator) ValidateAndTranslate(obj interface{}) map[string]string {
	err := cv.Validator.Struct(obj)
	if err == nil {
		return nil
	}

	// Check if it's validation errors
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string]string{"error": "unknown validation error"}
	}

	errors := make(map[string]string)

	for _, err := range validationErrors {
		// Check for custom message from struct tag
		field, ok := reflect.TypeOf(obj).Elem().FieldByName(err.Field())
		if ok {
			customMsg := field.Tag.Get("msg")
			if customMsg != "" {
				errors[strings.ToLower(err.Field())] = customMsg
				continue
			}
		}

		// Use translated message
		errors[strings.ToLower(err.Field())] = err.Translate(cv.Translator)
	}

	return errors
}

// ValidateAndTranslateString returns error messages as a string
func (cv *CustomValidator) ValidateAndTranslateString(obj interface{}) string {
	err := cv.Validator.Struct(obj)
	if err == nil {
		return ""
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		// Check for custom message from struct tag
		field, _ := reflect.TypeOf(obj).Elem().FieldByName(err.Field())
		customMsg := field.Tag.Get("msg")

		if customMsg != "" {
			messages = append(messages, customMsg)
		} else {
			messages = append(messages, err.Translate(cv.Translator))
		}
	}

	return strings.Join(messages, "; ")
}

// ValidateStruct performs basic validation only
func (cv *CustomValidator) ValidateStruct(obj interface{}) error {
	return cv.Validator.Struct(obj)
}

// Example: Using custom validation
type StrongPasswordRequest struct {
	Password string `json:"password" binding:"required,strong_password,min=8" msg:"password must be strong: at least 8 characters, 1 uppercase, 1 lowercase, 1 number"`
}
