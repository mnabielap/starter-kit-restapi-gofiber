package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type XValidator struct {
	validator *validator.Validate
}

var validate = validator.New()

func New() *XValidator {
	return &XValidator{validator: validate}
}

func (v *XValidator) Validate(data interface{}) map[string]string {
	err := v.validator.Struct(data)
	if err == nil {
		return nil
	}
	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errors[strings.ToLower(err.Field())] = fmt.Sprintf("failed on tag: %s", err.Tag())
	}
	return errors
}

func ParseAndValidate(c *fiber.Ctx, out interface{}) error {
	// Support both Body, Query and Params
	if err := c.BodyParser(out); err != nil {
		// Try query parser if body fails or is empty
		if err := c.QueryParser(out); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
		}
	}
	
	// Explicit check for Params if needed (usually handled manually)

	v := New()
	if errs := v.Validate(out); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  errs,
		})
	}
	return nil
}