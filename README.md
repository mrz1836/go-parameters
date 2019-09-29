# go-parameters
**go-parameters** is a parameter multi-tool that parses json, msg pack, or multi-part form data into a parameters object.

| | | | | | | |
|-|-|-|-|-|-|-|
| ![License](https://img.shields.io/github/license/mrz1836/go-parameters.svg?style=flat&p=1) | [![Report](https://goreportcard.com/badge/github.com/mrz1836/go-parameters?style=flat&p=1)](https://goreportcard.com/report/github.com/mrz1836/go-parameters)  | [![Codacy Badge](https://api.codacy.com/project/badge/Grade/0b377a0d1dde4b6ba189545aa7ee2e17)](https://www.codacy.com/app/mrz1818/go-parameters?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mrz1836/go-parameters&amp;utm_campaign=Badge_Grade) |  [![Build Status](https://travis-ci.com/mrz1836/go-parameters.svg?branch=master)](https://travis-ci.com/mrz1836/go-parameters)   |  [![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat)](https://github.com/RichardLitt/standard-readme) | [![Release](https://img.shields.io/github/release-pre/mrz1836/go-parameters.svg?style=flat)](https://github.com/mrz1836/go-parameters/releases) | [![GoDoc](https://godoc.org/github.com/mrz1836/go-parameters?status.svg&style=flat)](https://godoc.org/github.com/mrz1836/go-parameters) |

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

**go-parameters** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy) and [dep](https://github.com/golang/dep).
```bash
$ go get -u github.com/mrz1836/go-parameters
```

Updating dependencies in **go-parameters**:
```bash
$ cd ../go-parameters
$ dep ensure -update -v
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
All unit tests and [examples](examples/examples.go) run via [Travis CI](https://travis-ci.com/mrz1836/go-parameters) and uses [Go version 1.13.x](https://golang.org/doc/go1.13). View the [deployment configuration file](.travis.yml).

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
todo: @mrz1836
View the [examples](examples/examples.go)

Basic implementation:
todo: @mrz1836
```golang
package main

import (

)

func main() {

}
```

## Maintainers

[@MrZ1836](https://github.com/mrz1836)

## Contributing

This project uses Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter) package.

This project uses Gorilla's [mux](https://github.com/gorilla/mux) package.

This project uses Ugorji's [go](https://github.com/ugorji/go) package.

View the [contributing guidelines](CONTRIBUTING.md) and follow the [code of conduct](CODE_OF_CONDUCT.md).

Support the development of this project 🙏

[![Donate](https://img.shields.io/badge/donate-bitcoin-brightgreen.svg)](https://mrz1818.com/?tab=tips&af=go-parameters)

## License

![License](https://img.shields.io/github/license/mrz1836/go-parameters.svg?style=flat&p=1)
