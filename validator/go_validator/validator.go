package govalidator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"

	api "request_validator/http/v2"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() Validator {
	ret := Validator{validate: validator.New()}
	return ret
}

func (v *Validator) ValidateRequest(ctx context.Context, r *http.Request) (*api.CreateUserReq, error) {

	// --- (1) ----
	// Try to decode the request body into the struct.
	var req api.CreateUserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal request body: %w", err)
	}

	// --- (2) ----
	// Validate the unmarshalled struct
	err = v.validate.StructCtx(ctx, req)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, validationErrors
	}

	return &req, nil
}
