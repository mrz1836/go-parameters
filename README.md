# go-parameters
> Parameter multi-tool that parses json, msg pack, or multipart form data into a parameter object.

[![Release](https://img.shields.io/github/release-pre/mrz1836/go-parameters.svg?logo=github&style=flat)](https://github.com/mrz1836/go-parameters/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/mrz1836/go-parameters/run-tests.yml?branch=master&logo=github&v=2)](https://github.com/mrz1836/go-parameters/actions)
[![Report](https://goreportcard.com/badge/github.com/mrz1836/go-parameters?style=flat&p=1)](https://goreportcard.com/report/github.com/mrz1836/go-parameters)
[![codecov](https://codecov.io/gh/mrz1836/go-parameters/branch/master/graph/badge.svg)](https://codecov.io/gh/mrz1836/go-parameters)
[![Go](https://img.shields.io/github/go-mod/go-version/mrz1836/go-parameters)](https://golang.org/)
[![Sponsor](https://img.shields.io/badge/sponsor-MrZ-181717.svg?logo=github&style=flat&v=3)](https://github.com/sponsors/mrz1836)
[![Donate](https://img.shields.io/badge/donate-bitcoin-ff9900.svg?logo=bitcoin&style=flat)](https://mrz1818.com/?tab=tips&utm_source=github&utm_medium=sponsor-link&utm_campaign=go-parameters&utm_term=go-parameters&utm_content=go-parameters)

<br/>

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

<br/>

## Installation

**go-parameters** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get -u github.com/mrz1836/go-parameters
```

<br/>

## Documentation
View the generated [documentation](https://pkg.go.dev/github.com/mrz1836/go-parameters)
 
[![GoDoc](https://godoc.org/github.com/mrz1836/go-parameters?status.svg&style=flat)](https://pkg.go.dev/github.com/mrz1836/go-parameters)

### Features
- Uses the fastest router: Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter)
- Works with `json`, `msgpack`, and `multi-part` forms
- Handles all standard types for `GetParams`
- Handler methods like `MakeParsedReq()` for `httprouter` use
- `Imbue` and `Permit` helper methods
- `GetParams()` parses parameters only once

<details>
<summary><strong><code>Package Dependencies</code></strong></summary>
<br/>

- Gorilla's [mux](https://github.com/gorilla/mux) package.
- Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter) package.
- Ugorji's [go codec](https://github.com/ugorji/go) package.
</details>

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

[goreleaser](https://github.com/goreleaser/goreleaser) for easy binary or library deployment to GitHub and can be installed via: `brew install goreleaser`.

The [.goreleaser.yml](.goreleaser.yml) file is used to configure [goreleaser](https://github.com/goreleaser/goreleaser).

Use `make release-snap` to create a snapshot version of the release, and finally `make release` to ship to production.
</details>

<details>
<summary><strong><code>Makefile Commands</code></strong></summary>
<br/>

View all `makefile` commands
```shell script
make help
```

List of all current commands:
```text
all                  Runs multiple commands
clean                Remove previous builds and any test cache data
clean-mods           Remove all the Go mod cache
coverage             Shows the test coverage
godocs               Sync the latest tag with GoDocs
help                 Show this help message
install              Install the application
install-go           Install the application (Using Native Go)
lint                 Run the golangci-lint application (install if not found)
release              Full production release (creates release in Github)
release              Runs common.release then runs godocs
release-snap         Test the full release (build binaries)
release-test         Full production test release (everything except deploy)
replace-version      Replaces the version in HTML/JS (pre-deploy)
run-examples         Runs all the examples
tag                  Generate a new tag and push (tag version=0.0.0)
tag-remove           Remove a tag if found (tag-remove version=0.0.0)
tag-update           Update an existing tag to current commit (tag-update version=0.0.0)
test                 Runs vet, lint and ALL tests
test-ci              Runs all tests via CI (exports coverage)
test-ci-no-race      Runs all tests via CI (no race) (exports coverage)
test-ci-short        Runs unit tests via CI (exports coverage)
test-short           Runs vet, lint and tests (excludes integration tests)
uninstall            Uninstall the application (and remove files)
update-linter        Update the golangci-lint package (macOS only)
vet                  Run the Go vet application
```
</details>

<br/>

## Examples & Tests
All unit tests and [examples](examples) run via [GitHub Actions](https://github.com/mrz1836/go-parameters/actions) and
uses [Go version 1.19.x](https://golang.org/doc/go1.19). View the [configuration file](.github/workflows/run-tests.yml).

Run all tests (including integration tests)
```shell script
make test
```

Run tests (excluding integration tests)
```shell script
make test-short
```

<br/>

## Benchmarks
The following benchmarks were conducted to measure the performance of various functions in the `github.com/mrz1836/go-parameters` package. All tests were run on a machine with the following specifications:

- **Operating System:** macOS (Darwin)
- **Architecture:** ARM64
- **CPU:** Apple M1 Max

<br/>

### Benchmark Results
View the latest [benchmark results](benchmarks/benchmarks-latest.md)

| Benchmark                   |  Iterations |    ns/op | B/op | allocs/op |
|-----------------------------|------------:|---------:|-----:|----------:|
| **UniqueUint64**            |  13,989,841 | 84.49 ns | 64 B |         1 |
| **GetParams_ParseJSONBody** | 209,700,817 | 5.721 ns |  0 B |         0 |
| **GetParams**               | 209,668,894 | 5.719 ns |  0 B |         0 |
| **Params_GetStringOk**      |  37,573,434 | 31.96 ns | 16 B |         1 |
| **Params_GetBoolOk**        |  36,349,316 | 33.15 ns | 16 B |         1 |
| **Params_GetBytesOk**       |  38,539,616 | 31.52 ns | 16 B |         1 |
| **Params_GetBool**          |  35,637,433 | 32.57 ns | 16 B |         1 |
| **Params_GetFloatOk**       |  22,064,586 | 54.69 ns | 16 B |         1 |
| **Params_GetIntOk**         |  28,992,304 | 41.78 ns | 16 B |         1 |
| **Params_GetInt64Ok**       |  28,752,844 | 41.77 ns | 16 B |         1 |
| **Params_GetIntSliceOk**    |  38,099,671 | 31.59 ns | 16 B |         1 |
| **Params_GetUint64Ok**      |  29,921,580 | 39.68 ns | 16 B |         1 |

<br/>

### Benchmark Details

- **`UniqueUint64`**: Measures the performance of generating unique `uint64` values.
- **`GetParams_ParseJSONBody`**: Benchmarks parsing a JSON body into parameters.
- **`GetParams`**: Tests retrieving parameters without parsing.
- **`Params_GetStringOk`**: Evaluates fetching a string parameter with success indication.
- **`Params_GetBoolOk`**: Assesses fetching a boolean parameter with success indication.
- **`Params_GetBytesOk`**: Measures retrieving a byte slice parameter with success indication.
- **`Params_GetBool`**: Benchmarks fetching a boolean parameter without success indication.
- **`Params_GetFloatOk`**: Tests fetching a float parameter with success indication.
- **`Params_GetIntOk`**: Evaluates fetching an integer parameter with success indication.
- **`Params_GetInt64Ok`**: Measures fetching a 64-bit integer parameter with success indication.
- **`Params_GetIntSliceOk`**: Benchmarks retrieving a slice of integers with success indication.
- **`Params_GetUint64Ok`**: Tests fetching an unsigned 64-bit integer parameter with success indication.

### Benchmark Notes

- **Iterations**: The number of times the benchmark function was executed.
- **ns/op**: Nanoseconds per operation, indicating the average time taken for each operation.
- **B/op**: Bytes allocated per operation, showing the memory usage.
- **allocs/op**: Allocations per operation, indicating how many memory allocations occurred per operation.

<br/>

### How to Run Benchmarks
Run the Go benchmarks:
```shell script
make bench
```
<br/>

## Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

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

<br/>

## Maintainers
| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) | [<img src="https://github.com/kayleg.png" height="50" alt="kayleg" />](https://github.com/kayleg) |
|:------------------------------------------------------------------------------------------------:|:-------------------------------------------------------------------------------------------------:|
|                                [MrZ](https://github.com/mrz1836)                                 |                                [kayleg](https://github.com/kayleg)                                |

<br/>

## Contributing
View the [contributing guidelines](.github/CONTRIBUTING.md) and please follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

### How can I help?
All kinds of contributions are welcome :raised_hands:! 
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:. 
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/mrz1836) :clap: 
or by making a [**bitcoin donation**](https://mrz1818.com/?tab=tips&utm_source=github&utm_medium=sponsor-link&utm_campaign=go-parameters&utm_term=go-parameters&utm_content=go-parameters) to ensure this journey continues indefinitely! :rocket:

<br/>

## License

[![License](https://img.shields.io/github/license/mrz1836/go-parameters.svg?style=flat&p=1)](LICENSE)
