package govalidator

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/require"

	api "request_validator/http/v2"
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
}
`

func TestValidator(t *testing.T) {
	// create the validator
	ctx := context.Background()
	reqValidator := NewValidator()

	tests := []struct {
		name     string
		req      string
		wantFunc func(t *testing.T, err error, req *api.CreateUserReq)
	}{
		{
			name: "given a request that whose ID is not of a UUID type, when we try to validate it, an error should be returned",
			req:  invalidFormatFieldRequest,
			wantFunc: func(t *testing.T, err error, req *api.CreateUserReq) {
				require.Error(t, err, "validator should error")
				var validationErrors validator.ValidationErrors
				require.True(t, errors.As(err, &validationErrors), "error should be of type validator.ValidationErrors")
			},
		},
		{
			name: "given a request that does not have a required field specified, when we try to validate it, an error should be returned",
			req:  missingMandatoryFieldRequest,
			wantFunc: func(t *testing.T, err error, req *api.CreateUserReq) {
				require.Error(t, err, "validator should error")
				var validationErrors validator.ValidationErrors
				require.True(t, errors.As(err, &validationErrors), "error should be of type validator.ValidationErrors")
			},
		},
		{
			name: "given a valid request, when we try to validate it, no error should be returned",
			req:  correctRequest,
			wantFunc: func(t *testing.T, err error, req *api.CreateUserReq) {
				require.Equal(t, req.Id, "32d3e8f1-2f81-49c0-acb6-6dccd84f3dab")
				require.Equal(t, req.FirstName, "Jon")
				require.Equal(t, req.LastName, "Snow")
				require.NoError(t, err, "validator should not error")
			},
		},
		{
			name: "given a request with an invalid email formatted field, when we try to validate it, an error should be returned",
			req: `
			{
				"id": "32d3e8f1-2f81-49c0-acb6-6dccd84f3dab",
				"firstName": "Jon",
				"lastName": "Snow",
				"email": "this_is_a_test"
			}`,
			wantFunc: func(t *testing.T, err error, req *api.CreateUserReq) {
				require.Error(t, err, "validator should error")
				var validationErrors validator.ValidationErrors
				require.True(t, errors.As(err, &validationErrors), "error should be of type validator.ValidationErrors")
			},
		},
		{
			name: "given a valid request with a valid email format, when we try to validate it, no error should be returned",
			req: `
			{
				"id": "32d3e8f1-2f81-49c0-acb6-6dccd84f3dab",
				"firstName": "Jon",
				"lastName": "Snow",
				"email": "jon_snow@winterfell.com"
			}`,
			wantFunc: func(t *testing.T, err error, req *api.CreateUserReq) {
				require.NoError(t, err, "validator should not error")
				require.Equal(t, req.Id, "32d3e8f1-2f81-49c0-acb6-6dccd84f3dab")
				require.Equal(t, req.FirstName, "Jon")
				require.Equal(t, req.LastName, "Snow")
				require.Equal(t, *req.Email, "jon_snow@winterfell.com")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, "", bytes.NewReader([]byte(tt.req)))
			httpRequest.Header.Add("Content-Type", "application/json")
			require.NoError(t, err, "http request creation should not error")

			// act
			var req api.CreateUserReq
			err = reqValidator.ValidateRequest(ctx, httpRequest, &req)

			// assert
			tt.wantFunc(t, err, &req)
		})
	}
}

func BenchmarkValidator(b *testing.B) {
	b.Run("Go validator benchmark with correct request", func(b *testing.B) {
		// arrange
		ctx := context.Background()
		var req api.CreateUserReq

		// create the validator
		reqValidator := NewValidator()

		for i := 0; i < b.N; i++ {
			httpRequest, _ := http.NewRequestWithContext(ctx, http.MethodPost, "", bytes.NewReader([]byte(correctRequest)))
			httpRequest.Header.Add("Content-Type", "application/json")

			reqValidator.ValidateRequest(ctx, httpRequest, &req)
		}
	})

	b.Run("Go validator benchmark with invalid format request", func(b *testing.B) {
		// arrange
		ctx := context.Background()
		var req api.CreateUserReq

		// create the validator
		reqValidator := NewValidator()

		for i := 0; i < b.N; i++ {
			httpRequest, _ := http.NewRequestWithContext(ctx, http.MethodPost, "", bytes.NewReader([]byte(invalidFormatFieldRequest)))
			httpRequest.Header.Add("Content-Type", "application/json")

			reqValidator.ValidateRequest(ctx, httpRequest, &req)
		}
	})

	b.Run("Go validator benchmark with missing field request", func(b *testing.B) {
		// arrange
		ctx := context.Background()
		var req api.CreateUserReq

		// create the validator
		reqValidator := NewValidator()

		for i := 0; i < b.N; i++ {
			httpRequest, _ := http.NewRequestWithContext(ctx, http.MethodPost, "", bytes.NewReader([]byte(missingMandatoryFieldRequest)))
			httpRequest.Header.Add("Content-Type", "application/json")

			reqValidator.ValidateRequest(ctx, httpRequest, &req)
		}
	})

}
