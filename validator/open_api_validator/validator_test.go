package openapivalidator

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	api "request_validator/http/v1"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/stretchr/testify/require"
)

const correctRequest = `
{
	"id": "32d3e8f1-2f81-49c0-acb6-6dccd84f3dab",
	"firstName": "Jon",
	"lastName": "Snow"
}`

const missingMandatoryFieldRequest = `
{
	"firstName": "Jon",
	"lastName": "Snow"
}
`

const invalidFormatFieldRequest = `
{
	"id": "sadwefsds",
	"firstName": "Jon",
	"lastName": "Snow"
}`

func TestValidator(t *testing.T) {
	// create the validator
	ctx := context.Background()
	swaggerDoc, err := api.GetSwagger()
	require.NoError(t, err, "swagger recovery should not error")
	validator := MustCreateValidator(ctx, swaggerDoc)

	tests := []struct {
		name     string
		req      string
		url      string
		wantFunc func(t *testing.T, err error)
	}{
		{
			name: "given a request that whose ID is not of a UUID type, when we try to validate it, an error should be returned",
			req:  invalidFormatFieldRequest,
			url:  "http://api.example.com/v1/users/create",
			wantFunc: func(t *testing.T, err error) {
				require.Error(t, err, "validator should error")
				var schemaErr *openapi3.SchemaError
				require.True(t, errors.As(err, &schemaErr), "error should be of type SchemaError")
			},
		},
		{
			name: "given a valid request that does not come from one of the specified url server, when we try to validate it, an error should be returned",
			req:  correctRequest,
			url:  "XXXXXXXXXXX",
			wantFunc: func(t *testing.T, err error) {
				require.Error(t, err, "validator should error")
				require.True(t, errors.Is(err, routers.ErrPathNotFound), "error should be of type ErrPathNotFound")
			},
		},
		{
			name: "given a request that does not have a required field specified, when we try to validate it, an error should be returned",
			req:  missingMandatoryFieldRequest,
			url:  "http://api.example.com/v1/users/create",
			wantFunc: func(t *testing.T, err error) {
				var schemaErr *openapi3.SchemaError
				require.True(t, errors.As(err, &schemaErr), "error should be of type SchemaError")
				require.Error(t, err, "validator should error")
			},
		},
		{
			name:     "given a valid request from a registered server url, when we try to validate it, no error should be returned",
			req:      correctRequest,
			url:      "http://api.example.com/v1/users/create",
			wantFunc: func(t *testing.T, err error) { require.NoError(t, err, "validator should not error") },
		},
		{
			name: "given a valid request from a registered server url but with an invalid email formatted field, when we try to validate it, an error should be returned",
			req: `
			{
				"id": "32d3e8f1-2f81-49c0-acb6-6dccd84f3dab",
				"firstName": "Jon",
				"lastName": "Snow",
				"email: "this_is_a_test"
			}
			`,
			url: "http://api.example.com/v1/users/create",
			wantFunc: func(t *testing.T, err error) {
				// Since the additional validation of the email field is not done by the openapi3 library, the error will be of type ParseError because it can't parse
				// the request email field to the specific format we have defined in the validator
				var parseErr *openapi3filter.ParseError
				require.True(t, errors.As(err, &parseErr), "error should be of type ParseError")
				require.Error(t, err, "validator should error")
			},
		},
		{
			name: "given a valid request from a registered server url, when we try to validate it, no error should be returned",
			req: `
			{
				"id": "32d3e8f1-2f81-49c0-acb6-6dccd84f3dab",
				"firstName": "Jon",
				"lastName": "Snow",
				"email": "jon_snow@winterfell.com"
			}
			`,
			url:      "http://api.example.com/v1/users/create",
			wantFunc: func(t *testing.T, err error) { require.NoError(t, err, "validator should not error") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, tt.url, bytes.NewReader([]byte(tt.req)))
			httpRequest.Header.Add("Content-Type", "application/json")
			require.NoError(t, err, "http request creation should not error")

			// act
			err = validator.ValidateRequest(ctx, httpRequest)

			// assert
			tt.wantFunc(t, err)
		})
	}
}

func BenchmarkValidator(b *testing.B) {
	b.Run("benchmark with correct request", func(b *testing.B) {
		// arrange
		ctx := context.Background()

		// create the validator
		swaggerDoc, err := api.GetSwagger()
		require.NoError(b, err, "swagger recovery should not error")
		validator := MustCreateValidator(ctx, swaggerDoc)

		for i := 0; i < b.N; i++ {
			httpRequest, _ := http.NewRequestWithContext(ctx, http.MethodPost, "http://api.example.com/v1/users/create", bytes.NewReader([]byte(correctRequest)))
			httpRequest.Header.Add("Content-Type", "application/json")

			validator.ValidateRequest(ctx, httpRequest)
		}
	})

	b.Run("benchmark with invalid format request", func(b *testing.B) {
		// arrange
		ctx := context.Background()

		// create the validator
		swaggerDoc, err := api.GetSwagger()
		require.NoError(b, err, "swagger recovery should not error")
		validator := MustCreateValidator(ctx, swaggerDoc)

		for i := 0; i < b.N; i++ {
			httpRequest, _ := http.NewRequestWithContext(ctx, http.MethodPost, "http://api.example.com/v1/users/create", bytes.NewReader([]byte(invalidFormatFieldRequest)))
			httpRequest.Header.Add("Content-Type", "application/json")

			validator.ValidateRequest(ctx, httpRequest)
		}
	})

	b.Run("benchmark with missing field request", func(b *testing.B) {
		// arrange
		ctx := context.Background()

		// create the validator
		swaggerDoc, err := api.GetSwagger()
		require.NoError(b, err, "swagger recovery should not error")
		validator := MustCreateValidator(ctx, swaggerDoc)

		for i := 0; i < b.N; i++ {
			httpRequest, _ := http.NewRequestWithContext(ctx, http.MethodPost, "http://api.example.com/v1/users/create", bytes.NewReader([]byte(missingMandatoryFieldRequest)))
			httpRequest.Header.Add("Content-Type", "application/json")

			validator.ValidateRequest(ctx, httpRequest)
		}
	})
}
