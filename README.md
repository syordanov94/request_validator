# Open API Request Validation Comparison

This project aims to present 2 different implementations of a validator that can assure us if a user request that a service receives follows the service's **Openapi** defined specification.

## Prerequisites

- Golang 1.20 or higher installed
- _Recommended but not mandatory VS Code or a similar IDE_

## How to install and Run the project

### How to perform changes to schema yaml spec

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