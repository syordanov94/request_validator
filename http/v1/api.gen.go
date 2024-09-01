// Package http_v1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.2.0 DO NOT EDIT.
package http_v1

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// CreateUserReq defines model for CreateUserReq.
type CreateUserReq struct {
	Email *openapi_types.Email `json:"email,omitempty"`

	// FirstName The user's first name
	FirstName string `json:"firstName"`

	// Id The user's unique identifier in UUID format
	Id openapi_types.UUID `json:"id"`

	// LastName The user's last names
	LastName string `json:"lastName"`
}

// PostUsersCreateJSONRequestBody defines body for PostUsersCreate for application/json ContentType.
type PostUsersCreateJSONRequestBody = CreateUserReq

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/5RUTU/cMBT8K9ZrpYKUTbILQpBTCz10pUJRC+oBcfDaLxtD/IH9Qlmh/PfKTijho6i9",
	"eb1jz8zzTO5BWO2sQUMBqnsIokHN0/LIIyc8D+i/403ccN469KQw/Y2aqzYt7rh2LUIFK06am49rSw3X",
	"ubAaMqit15ygGvEZ0MZFbCCvzBr6DGrlA51wjfEyiUF45UhZAxWcNci6gP5DYAnFTIRlE8pD3wl87VYl",
	"37yuM+qmQ6YkGlK1Qs+UYefny89sFDwl2VnIHdyv57NFvT+f7R6IcsbFam+2J4WQ+7v1juSrqdWuU/I1",
	"TS3/B6MRlHyGJxp+8o15xWifgcebTnmUUF1A4n0c6ITy8s9Ru7pCQdDHs8rU9qWab2nBW6a7llSrDDLr",
	"WVBm3eIs/Zzg4+AujqzW1hxzf3251RC5qihE2tLcX+fWr4sGW1dsx3u+nB1/zaMTRcnZj2SRfTpdQga3",
	"6MMgoszn+UGcmnVouFNQwSIv8xIycJyalMEiziwUIiU1RdQGesMP3hEaifK5/kf5U4Ex7TxClhIqOLWB",
	"YhnC0AsYBo+BDq3cRE5hDaFJ9Ny5Vol0trgKUcNDs+LqvccaKnhXPFavGHtXPC1d3w8PHJw1Yajdoixf",
	"Gjyx7AHEVGAGMZq86gIxapAtypIF4tQFJqzElJrQac39Bqqx54FxZvBXSmE+hCOgj88B1cVfBzpApuPM",
	"GObrnB1zZdiW81Z2Im5vj1DIoPMtVDCmhDuVjyGPH4zidg599v90S0PoE4L4Wpn1A7K2nhEGil15xjwi",
	"Z88UQH/Z/w4AAP//LJP9sBUFAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}