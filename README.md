# Open API Request Validation Comparison

This project aims to present 2 different implementations of a validator that can assure us if a user request that a service receives follows the service's **Openapi** defined specification. 

The first implementation will be a direct validator that uses the OpenAPI schema specs (**yaml** file) and checks if the http request complies with the schema. This validator checks:
- Origin Server
- Path params
- Request body

For more information regarding this validator, you can check the specific pkg page [here](https://github.com/getkin/kin-openapi?tab=readme-ov-file#validating-http-requestsresponses)


The other implementation uses the **Go-Playground** validator that compares the unmarshalled request body against the Go structures generated from the OpenAPI specs. This  implementation DOES NOT validate the origin server neither does it validate the path params, it ONLY validates the request body. 

For more information on this validator you can check the specific pkg page [here](https://github.com/go-playground/validator). Also, you can check how you can generate the validation rules from the **OpenAPI** spec [here](https://github.com/oapi-codegen/oapi-codegen/blob/main/examples/extensions/xoapicodegenextratags/api.yaml).

## Prerequisites

- Golang 1.20 or higher installed
- _Recommended but not mandatory VS Code or a similar IDE_

## How to install and Run the project

This project does not have **main.go** because its objective is to compare both request validator implementations. We have added a benchmarks to be able to accomplish this goal. The benchmarks are the same for both implementations in order to be able to give the "fairest" comparison. 

- First you will have to clone the project from this github repository:

```bash
git clone https://github.com/syordanov94/request_validator.git
```

- Go to the specific implementation folder:

    - For OpenAPI Validator

        ```bash
        cd validator/open_api_validator/
        ```

    - For Go Validator:
    
        ```bash
        cd validator/go_validator/
        ```

- Run all the benchmarks for each one of the implementations with the following command (which will iterate over the benchmarks 1000 times)

        ```bash
        go test -bench ./... -benchmem -benchtime=1000x
        ```

## How to perform changes to schema yaml spec

The **Go-Playground Validator** requires us to generate the resulting Go files from the **OpenAPI** yaml specs. To do this, we need to navigate to the schema version. Currently we have 2 schema versions:

- V1. The version located in *http/v1/api.yaml* is used by the **OpenAPI** validator.
- V2. The version located in *http/v2/api.yaml* is used by the **Go** validator.

To generate the Go code, once in the specific folder, we can run the following command:

    ```bash
    go generate ./...
    ```


### Possible error during go generate for open-api

When you run:

```bash
go generate ./...
```

you might get the following error:

```bash
sh: oapi-codegen: command not found
gen.go:4: running "sh": exit status 127
```

or

```bash
example/gen.go:4: running "oapi-codegen": exec: "oapi-codegen": executable file not found in $PATH
```

To fix this issue just run the following command:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

This solution will be required multiple times since it does not permanently add the *GOPATH/bin** path to the PATH. To add it permanently you will need to add it in the **~/.bashrc** or any other shell used.

## Benchmark Results

- Open API Validator:

```bash
goos: darwin
goarch: arm64
pkg: request_validator/validator/open_api_validator
BenchmarkValidator/OpenAPI_Validator_benchmark_with_correct_request-12              1000              6625 ns/op            3763 B/op         45 allocs/op
BenchmarkValidator/OpenAPI_Validator_benchmark_with_correct_request-12              1000              5456 ns/op            3763 B/op         45 allocs/op
BenchmarkValidator/OpenAPI_Validator_benchmark_with_correct_request-12              1000              4987 ns/op            3763 B/op         45 allocs/op
BenchmarkValidator/OpenAPI_Validator_benchmark_with_correct_request-12              1000              4870 ns/op            3763 B/op         45 allocs/op
BenchmarkValidator/OpenAPI_Validator_benchmark_with_correct_request-12              1000              4494 ns/op            3763 B/op         45 allocs/op

BenchmarkValidator/OpenAPI_Validator_benchmark_with_invalid_format_request-12       1000              5950 ns/op            6875 B/op         74 allocs/op
BenchmarkValidator/OpenAPI_Validator_benchmark_with_invalid_format_request-12       1000              5988 ns/op            6876 B/op         74 allocs/op
BenchmarkValidator/OpenAPI_Validator_benchmark_with_invalid_format_request-12       1000              5518 ns/op            6837 B/op         74 allocs/op
BenchmarkValidator/OpenAPI_Validator_benchmark_with_invalid_format_request-12       1000              5480 ns/op            6877 B/op         74 allocs/op
BenchmarkValidator/OpenAPI_Validator_benchmark_with_invalid_format_request-12       1000              5688 ns/op            6877 B/op         74 allocs/op

BenchmarkValidator/OpenAPI_Validator_benchmark_with_missing_field_request-12        1000             35223 ns/op           46807 B/op        248 allocs/op
BenchmarkValidator/OpenAPI_Validator_benchmark_with_missing_field_request-12        1000             32339 ns/op           46818 B/op        248 allocs/op
BenchmarkValidator/OpenAPI_Validator_benchmark_with_missing_field_request-12        1000             32281 ns/op           46815 B/op        248 allocs/op
BenchmarkValidator/OpenAPI_Validator_benchmark_with_missing_field_request-12        1000             31690 ns/op           46819 B/op        248 allocs/op
BenchmarkValidator/OpenAPI_Validator_benchmark_with_missing_field_request-12        1000             31982 ns/op           46781 B/op        248 allocs/op
PASS
ok      request_validator/validator/open_api_validator  0.473s
```

- Go Validator:

```bash
goos: darwin
goarch: arm64
pkg: request_validator/validator/go_validator
BenchmarkValidator/Go_validator_benchmark_with_correct_request-12                   1000              3707 ns/op            2165 B/op         20 allocs/op
BenchmarkValidator/Go_validator_benchmark_with_correct_request-12                   1000              3282 ns/op            2165 B/op         20 allocs/op
BenchmarkValidator/Go_validator_benchmark_with_correct_request-12                   1000              2948 ns/op            2165 B/op         20 allocs/op
BenchmarkValidator/Go_validator_benchmark_with_correct_request-12                   1000              2974 ns/op            2165 B/op         20 allocs/op
BenchmarkValidator/Go_validator_benchmark_with_correct_request-12                   1000              2522 ns/op            2166 B/op         20 allocs/op

BenchmarkValidator/Go_validator_benchmark_with_invalid_format_request-12            1000              2243 ns/op            2331 B/op         25 allocs/op
BenchmarkValidator/Go_validator_benchmark_with_invalid_format_request-12            1000              2149 ns/op            2331 B/op         25 allocs/op
BenchmarkValidator/Go_validator_benchmark_with_invalid_format_request-12            1000              2060 ns/op            2331 B/op         25 allocs/op
BenchmarkValidator/Go_validator_benchmark_with_invalid_format_request-12            1000              2039 ns/op            2331 B/op         25 allocs/op
BenchmarkValidator/Go_validator_benchmark_with_invalid_format_request-12            1000              1825 ns/op            2331 B/op         25 allocs/op

BenchmarkValidator/Go_validator_benchmark_with_missing_field_request-12             1000              1612 ns/op            2291 B/op         24 allocs/op
BenchmarkValidator/Go_validator_benchmark_with_missing_field_request-12             1000              1760 ns/op            2291 B/op         24 allocs/op
BenchmarkValidator/Go_validator_benchmark_with_missing_field_request-12             1000              1775 ns/op            2291 B/op         24 allocs/op
BenchmarkValidator/Go_validator_benchmark_with_missing_field_request-12             1000              1629 ns/op            2291 B/op         24 allocs/op
BenchmarkValidator/Go_validator_benchmark_with_missing_field_request-12             1000              1334 ns/op            2291 B/op         24 allocs/op
PASS
ok      request_validator/validator/go_validator        0.333s
```