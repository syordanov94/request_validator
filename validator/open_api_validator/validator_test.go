package openapivalidator

import (
	"bytes"
	"context"
	"net/http"
	api "request_validator/http"
	"testing"

	"github.com/stretchr/testify/require"
)

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
			req: `
			{
				"id": "sadwefsds",
				"firstName": "Jon",
				"lastName": "Snow"
			}
			`,
			url:      "http://api.example.com/v1/users/create",
			wantFunc: func(t *testing.T, err error) { require.Error(t, err, "validator should error") },
		},
		{
			name: "given a valid request that does not come from one of the specified url server, when we try to validate it, an error should be returned",
			req: `
			{
				"id": "32d3e8f1-2f81-49c0-acb6-6dccd84f3dab",
				"firstName": "Jon",
				"lastName": "Snow"
			}
			`,
			url:      "XXXXXXXXXXX",
			wantFunc: func(t *testing.T, err error) { require.Error(t, err, "validator should error") },
		},
		{
			name: "given a request that does not have a required field specified, when we try to validate it, an error should be returned",
			req: `
			{
				"firstName": "Jon",
				"lastName": "Snow"
			}
			`,
			url:      "http://api.example.com/v1/users/create",
			wantFunc: func(t *testing.T, err error) { require.Error(t, err, "validator should error") },
		},
		{
			name: "given a valid request from a registered server url, when we try to validate it, no error should be returned",
			req: `
			{
				"id": "32d3e8f1-2f81-49c0-acb6-6dccd84f3dab",
				"firstName": "Jon",
				"lastName": "Snow"
			}
			`,
			url:      "http://api.example.com/v1/users/create",
			wantFunc: func(t *testing.T, err error) { require.NoError(t, err, "validator should not error") },
		},
		{
			name: "given a valid request from a registered server url but with an invalid email formatted field, when we try to validate it, an error should be returned",
			req: `
			{
				"id": "32d3e8f1-2f81-49c0-acb6-6dccd84f3dab",
				"firstName": "Jon",
				"lastName": "Snow"
				"email: this_is_a_test
			}
			`,
			url:      "http://api.example.com/v1/users/create",
			wantFunc: func(t *testing.T, err error) { require.Error(t, err, "validator should error") },
		},
		{
			name: "given a valid request from a registered server url, when we try to validate it, no error should be returned",
			req: `
			{
				"id": "32d3e8f1-2f81-49c0-acb6-6dccd84f3dab",
				"firstName": "Jon",
				"lastName": "Snow"
				"email: jon_snow@winterfell.com
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
