# go-parameters
**go-parameters** is a parameter multi-tool that parses json, msg pack, or multi-part form data into a parameters object.

[![Go](https://img.shields.io/github/go-mod/go-version/mrz1836/go-parameters)](https://golang.org/)
[![Build Status](https://travis-ci.com/mrz1836/go-parameters.svg?branch=master)](https://travis-ci.com/mrz1836/go-parameters)
[![Report](https://goreportcard.com/badge/github.com/mrz1836/go-parameters?style=flat&p=1)](https://goreportcard.com/report/github.com/mrz1836/go-parameters)
[![codecov](https://codecov.io/gh/mrz1836/go-parameters/branch/master/graph/badge.svg)](https://codecov.io/gh/mrz1836/go-parameters)
[![Release](https://img.shields.io/github/release-pre/mrz1836/go-parameters.svg?style=flat)](https://github.com/mrz1836/go-parameters/releases)
[![GoDoc](https://godoc.org/github.com/mrz1836/go-parameters?status.svg&style=flat)](https://pkg.go.dev/github.com/mrz1836/go-parameters)

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

## Documentation
You can view the generated [documentation here](https://pkg.go.dev/github.com/mrz1836/go-parameters).

### Features
- Uses the fastest router: Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter)
- Works with `json`, `msgpack`, and `multi-part` forms
- Handles all standard types for `GetParams`
- Handler methods like `MakeParsedReq()` for `httprouter` use
- `Imbue` and `Permit` helper methods
- `GetParams()` parses parameters only once

<details>
<summary><strong><code>Package Dependencies</code></strong></summary>

- Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter) package.
- Gorilla's [mux](https://github.com/gorilla/mux) package.
- Ugorji's [go](https://github.com/ugorji/go) package.
</details>

<details>
<summary><strong><code>Library Deployment</code></strong></summary>

[goreleaser](https://github.com/goreleaser/goreleaser) for easy binary or library deployment to Github and can be installed via: `brew install goreleaser`.

The [.goreleaser.yml](.goreleaser.yml) file is used to configure [goreleaser](https://github.com/goreleaser/goreleaser).

Use `make release-snap` to create a snapshot version of the release, and finally `make release` to ship to production.
</details>

<details>
<summary><strong><code>Makefile Commands</code></strong></summary>

View all `makefile` commands
```bash
$ make help
```

List of all current commands:
```text
bench                          Run all benchmarks in the Go application
clean                          Remove previous builds and any test cache data
clean-mods                     Remove all the Go mod cache
coverage                       Shows the test coverage
godocs                         Sync the latest tag with GoDocs
help                           Show all make commands available
lint                           Run the Go lint application
release                        Full production release (creates release in Github)
release-test                   Full production test release (everything except deploy)
release-snap                   Test the full release (build binaries)
run-examples                   Runs all the examples
tag                            Generate a new tag and push (IE: make tag version=0.0.0)
tag-remove                     Remove a tag if found (IE: make tag-remove version=0.0.0)
tag-update                     Update an existing tag to current commit (IE: make tag-update version=0.0.0)
test                           Runs vet, lint and ALL tests
test-short                     Runs vet, lint and tests (excludes integration tests)
update                         Update all project dependencies
update-releaser                Update the goreleaser application
vet                            Run the Go vet application
```
</details>

## Examples & Tests
All unit tests and [examples](examples/examples.go) run via [Travis CI](https://travis-ci.org/mrz1836/go-parameters) and uses [Go version 1.14.x](https://golang.org/doc/go1.14). View the [deployment configuration file](.travis.yml).

Run all tests (including integration tests)
```bash
$ make test
```

Run tests (excluding integration tests)
```bash
$ make test-short
```

## Benchmarks
Run the Go benchmarks:
```bash
$ make bench
```

## Code Standards
Read more about this Go project's [code standards](CODE_STANDARDS.md).

## Usage
View the [examples](examples/examples.go)

Basic implementation:
```go
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

	name, _ := params.GetStringOk("name")

	_, _ = fmt.Fprintf(w, `{"hello":"%s"}`, name)
}

func main() {
    router := httprouter.New()
	router.GET("/hello/:name", parameters.GeneralJSONResponse(Hello))
	log.Fatal(http.ListenAndServe(":8080", router))
}
```

## Maintainers

| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) | [<img src="https://github.com/kayleg.png" height="50" alt="kayleg" />](https://github.com/kayleg) |
|:---:|:---:|
| [MrZ](https://github.com/mrz1836) | [kayleg](https://github.com/kayleg) |

## Contributing

This project uses Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter) package.

This project uses Gorilla's [mux](https://github.com/gorilla/mux) package.

This project uses Ugorji's [go](https://github.com/ugorji/go) package.

View the [contributing guidelines](CONTRIBUTING.md) and follow the [code of conduct](CODE_OF_CONDUCT.md).

Support the development of this project 🙏

[![Donate](https://img.shields.io/badge/donate-bitcoin-brightgreen.svg)](https://mrz1818.com/?tab=tips&af=go-parameters)

## License

![License](https://img.shields.io/github/license/mrz1836/go-parameters.svg?style=flat&p=1)
