package validate

import "github.com/go-playground/validator"

var Validator *validator.Validate

func init() {
	Validator = validator.New()
}
