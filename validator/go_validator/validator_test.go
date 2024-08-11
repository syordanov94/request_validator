package govalidator

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidator(t *testing.T) {
	// create the validator
	ctx := context.Background()
	reqValidator := NewValidator()

	tests := []struct {
		name     string
		req      string
		wantFunc func(t *testing.T, err error)
	}{
		// {
		// 	name: "given a request that whose ID is not of a UUID type, when we try to validate it, an error should be returned",
		// 	req: `
		// 	{
		// 		"id": "sadwefsds",
		// 		"firstName": "Jon",
		// 		"lastName": "Snow"
		// 	}
		// 	`,
		// 	wantFunc: func(t *testing.T, err error) { require.Error(t, err, "validator should error") },
		// },
		// {
		// 	name: "given a request that does not have a required field specified, when we try to validate it, an error should be returned",
		// 	req: `
		// 	{
		// 		"firstName": "Jon",
		// 		"lastName": "Snow"
		// 	}
		// 	`,
		// 	wantFunc: func(t *testing.T, err error) { require.Error(t, err, "validator should error") },
		// },
		{
			name: "given a valid request, when we try to validate it, no error should be returned",
			req: `
			{
				"id": "32d3e8f1-2f81-49c0-acb6-6dccd84f3dab",
				"firstName": "Jon",
				"lastName": "Snow"
			}`,
			wantFunc: func(t *testing.T, err error) { require.NoError(t, err, "validator should not error") },
		},
		// {
		// 	name: "given a request with an invalid email formatted field, when we try to validate it, an error should be returned",
		// 	req: `
		// 	{
		// 		"id": "32d3e8f1-2f81-49c0-acb6-6dccd84f3dab",
		// 		"firstName": "Jon",
		// 		"lastName": "Snow",
		// 		"email: this_is_a_test
		// 	}`,
		// 	wantFunc: func(t *testing.T, err error) { require.Error(t, err, "validator should error") },
		// },
		// {
		// 	name: "given a valid request with a valid email format, when we try to validate it, no error should be returned",
		// 	req: `
		// 	{
		// 		"id": "32d3e8f1-2f81-49c0-acb6-6dccd84f3dab",
		// 		"firstName": "Jon",
		// 		"lastName": "Snow",
		// 		"email: jon_snow@winterfell.com
		// 	}`,
		// 	wantFunc: func(t *testing.T, err error) { require.NoError(t, err, "validator should not error") },
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, "", bytes.NewReader([]byte(tt.req)))
			httpRequest.Header.Add("Content-Type", "application/json")
			require.NoError(t, err, "http request creation should not error")

			// act
			_, err = reqValidator.ValidateRequest(ctx, httpRequest)

			// assert
			tt.wantFunc(t, err)
		})
	}
}
