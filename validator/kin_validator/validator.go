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

// The Open API Validator will contain a router which is used to recover the route and path params that the http request relates to
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

func (v *Validator) ValidateRequest(ctx context.Context, httpReq *http.Request) error {

	route, pathParams, err := v.router.FindRoute(httpReq)
	if err != nil {
		return fmt.Errorf("error recovering request route: %w", err)
	}

	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
		// This set up ignores authentication validation. We do this because our main objective here is to only validate the request body. Security validation is already
		// performed by each service in it's middleware.
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