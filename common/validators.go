package common

import (
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()

	Validate.RegisterValidation("validFoodType", ValidateFoodType)
	Validate.RegisterValidation("validMealTime", ValidateMealTime)
}

func ValidateFoodType(fl validator.FieldLevel) bool {
	validTypes := ValidFoodTypes
	for _, validType := range validTypes {
		if fl.Field().String() == validType {
			return true
		}
	}
	return false
}

func ValidateMealTime(fl validator.FieldLevel) bool {
	validMealTimes := ValidMealTimes
	for _, validMealTime := range validMealTimes {
		if fl.Field().String() == validMealTime {
			return true
		}
	}
	return false
}
