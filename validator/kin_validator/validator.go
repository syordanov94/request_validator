package kinvalidator

import (
	"context"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

var (
	emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
)

type Validator struct {
	router routers.Router
}

func MustCreateValidator(ctx context.Context, doc *openapi3.T) *Validator {

	// Set specific validation format for UUID and Email format typed fields
	openapi3.DefineStringFormatValidator("uuid", openapi3.NewRegexpFormatValidator(openapi3.FormatOfStringForUUIDOfRFC4122))
	openapi3.DefineStringFormatValidator("email", openapi3.NewRegexpFormatValidator(emailRegex))

	err := doc.Validate(ctx)
	if err != nil {
		panic(fmt.Errorf("unable to validate open api specs: %w", err))
	}

	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		panic(fmt.Errorf("unable to create router: %w", err))
	}

	return &Validator{
		router: router,
	}
}

func (v *Validator) ValidateRequest(ctx context.Context, httpRq *http.Request) error {

	r, params, err := v.router.FindRoute(httpRq)
	if err != nil {
		return fmt.Errorf("error finding request route: %w", err)
	}

	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpRq,
		PathParams: params,
		Route:      r,
		Options: &openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		},
	}
	err = openapi3filter.ValidateRequest(ctx, requestValidationInput)
	if err != nil {
		return fmt.Errorf("error validating request: %w", err)
	}
	return nil
}
