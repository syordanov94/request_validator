package govalidator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() Validator {
	ret := Validator{validate: validator.New()}
	return ret
}

func (v *Validator) ValidateRequest(ctx context.Context, r *http.Request, req interface{}) error {

	// --- (1) ----
	// Try to decode the request body into the struct.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return fmt.Errorf("unable to unmarshal request body: %w", err)
	}

	// --- (2) ----
	// Validate the unmarshalled struct
	err = v.validate.StructCtx(ctx, req)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return validationErrors
	}

	return nil
}
