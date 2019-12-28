# go-parameters
**go-parameters** is a parameter multi-tool that parses json, msg pack, or multi-part form data into a parameters object.

[![Build Status](https://travis-ci.org/mrz1836/go-parameters.svg?branch=master)](https://travis-ci.org/mrz1836/go-parameters)
[![Report](https://goreportcard.com/badge/github.com/mrz1836/go-parameters?style=flat&p=1)](https://goreportcard.com/report/github.com/mrz1836/go-parameters)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/4859aa01ee8b435d9cd94711589f5086)](https://www.codacy.com/app/mrz1818/go-parameters?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mrz1836/go-parameters&amp;utm_campaign=Badge_Grade)
[![Release](https://img.shields.io/github/release-pre/mrz1836/go-parameters.svg?style=flat)](https://github.com/mrz1836/go-parameters/releases)
[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat)](https://github.com/RichardLitt/standard-readme)
[![GoDoc](https://godoc.org/github.com/mrz1836/go-parameters?status.svg&style=flat)](https://godoc.org/github.com/mrz1836/go-parameters)

## Table of Contents
- [Installation](#installation)
- [Documentation](#documentation)
- [Examples & Tests](#examples--tests)
- [Benchmarks](#benchmarks)
- [Code Standards](#code-standards)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

## Installation

**go-parameters** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```bash
$ go get -u github.com/mrz1836/go-parameters
```

### Package Dependencies
- Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter) package.
- Gorilla's [mux](https://github.com/gorilla/mux) package.
- Ugorji's [go](https://github.com/ugorji/go) package.

## Documentation
You can view the generated [documentation here](https://godoc.org/github.com/mrz1836/go-parameters).

### Features
- Uses the fastest router: Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter)
- Works with `json`, `msgpack`, and `multi-part` forms
- Handles all standard types for `GetParams`
- Handler methods like `MakeParsedReq()` for `httprouter` use
- `Imbue` and `Permit` helper methods
- `GetParams()` parses parameters only once

## Examples & Tests
All unit tests and [examples](examples/examples.go) run via [Travis CI](https://travis-ci.org/mrz1836/go-parameters) and uses [Go version 1.13.x](https://golang.org/doc/go1.13). View the [deployment configuration file](.travis.yml).

Run all tests (including integration tests)
```bash
$ cd ../go-parameters
$ go test ./... -v
```

Run tests (excluding integration tests)
```bash
$ cd ../go-parameters
$ go test ./... -v -test.short
```

## Benchmarks
Run the Go benchmarks:
```bash
$ cd ../go-parameters
$ go test -bench . -benchmem
```

## Code Standards
Read more about this Go project's [code standards](CODE_STANDARDS.md).

## Usage
View the [examples](examples/examples.go)

Basic implementation:
```golang
package main

import (
    "fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mrz1836/go-parameters"
)

func Hello(w http.ResponseWriter, req *http.Request) {

	params := parameters.GetParams(req)

	name := params.GetStringOk("name")

	_, _ = fmt.Fprintf(w, `{"hello":"%s"}`, name)
}

func main() {
    router := httprouter.New()
	router.GET("/hello/:name", parameters.GeneralJSONResponse(Hello))
	log.Fatal(http.ListenAndServe(":8080", router))
}
```

## Maintainers

[@MrZ](https://github.com/mrz1836)

## Contributing

This project uses Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter) package.

This project uses Gorilla's [mux](https://github.com/gorilla/mux) package.

This project uses Ugorji's [go](https://github.com/ugorji/go) package.

View the [contributing guidelines](CONTRIBUTING.md) and follow the [code of conduct](CODE_OF_CONDUCT.md).

Support the development of this project 🙏

[![Donate](https://img.shields.io/badge/donate-bitcoin-brightgreen.svg)](https://mrz1818.com/?tab=tips&af=go-parameters)

## License

![License](https://img.shields.io/github/license/mrz1836/go-parameters.svg?style=flat&p=1)
