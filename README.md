<div align="center">

# 📘&nbsp;&nbsp;go-parameters

**Parameter multi-tool that parses json, msg pack, or multipart form data into a parameter object.**

<br/>

<a href="https://github.com/mrz1836/go-parameters/releases"><img src="https://img.shields.io/github/release-pre/mrz1836/go-parameters?include_prereleases&style=flat-square&logo=github&color=black" alt="Release"></a>
<a href="https://golang.org/"><img src="https://img.shields.io/github/go-mod/go-version/mrz1836/go-parameters?style=flat-square&logo=go&color=00ADD8" alt="Go Version"></a>
<a href="https://github.com/mrz1836/go-parameters/blob/master/LICENSE"><img src="https://img.shields.io/github/license/mrz1836/go-parameters?style=flat-square&color=blue" alt="License"></a>

<br/>

<table align="center" border="0">
  <tr>
    <td align="right">
       <code>CI / CD</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://github.com/mrz1836/go-parameters/actions"><img src="https://img.shields.io/github/actions/workflow/status/mrz1836/go-parameters/fortress.yml?branch=master&label=build&logo=github&style=flat-square" alt="Build"></a>
       <a href="https://github.com/mrz1836/go-parameters/actions"><img src="https://img.shields.io/github/last-commit/mrz1836/go-parameters?style=flat-square&logo=git&logoColor=white&label=last%20update" alt="Last Commit"></a>
    </td>
    <td align="right">
       &nbsp;&nbsp;&nbsp;&nbsp; <code>Quality</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://goreportcard.com/report/github.com/mrz1836/go-parameters"><img src="https://goreportcard.com/badge/github.com/mrz1836/go-parameters?style=flat-square" alt="Go Report"></a>
       <a href="https://codecov.io/gh/mrz1836/go-parameters"><img src="https://codecov.io/gh/mrz1836/go-parameters/branch/master/graph/badge.svg?style=flat-square" alt="Coverage"></a>
    </td>
  </tr>

  <tr>
    <td align="right">
       <code>Security</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://scorecard.dev/viewer/?uri=github.com/mrz1836/go-parameters"><img src="https://api.scorecard.dev/projects/github.com/mrz1836/go-parameters/badge?style=flat-square" alt="Scorecard"></a>
       <a href=".github/SECURITY.md"><img src="https://img.shields.io/badge/policy-active-success?style=flat-square&logo=security&logoColor=white" alt="Security"></a>
    </td>
    <td align="right">
       &nbsp;&nbsp;&nbsp;&nbsp; <code>Community</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://github.com/mrz1836/go-parameters/graphs/contributors"><img src="https://img.shields.io/github/contributors/mrz1836/go-parameters?style=flat-square&color=orange" alt="Contributors"></a>
       <a href="https://mrz1818.com/"><img src="https://img.shields.io/badge/donate-bitcoin-ff9900?style=flat-square&logo=bitcoin" alt="Bitcoin"></a>
    </td>
  </tr>
</table>

</div>

<br/>
<br/>

<div align="center">

### <code>Project Navigation</code>

</div>

<table align="center">
  <tr>
    <td align="center" width="33%">
       🚀&nbsp;<a href="#installation"><code>Installation</code></a>
    </td>
    <td align="center" width="33%">
       🧪&nbsp;<a href="#examples--tests"><code>Examples&nbsp;&&nbsp;Tests</code></a>
    </td>
    <td align="center" width="33%">
       📚&nbsp;<a href="#documentation"><code>Documentation</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
       🤝&nbsp;<a href="#contributing"><code>Contributing</code></a>
    </td>
    <td align="center">
      🛠️&nbsp;<a href="#code-standards"><code>Code&nbsp;Standards</code></a>
    </td>
    <td align="center">
      ⚡&nbsp;<a href="#benchmarks"><code>Benchmarks</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
      🤖&nbsp;<a href="#-ai-usage--assistant-guidelines"><code>AI&nbsp;Usage</code></a>
    </td>
    <td align="center">
       ⚖️&nbsp;<a href="#license"><code>License</code></a>
    </td>
    <td align="center">
       👥&nbsp;<a href="#maintainers"><code>Maintainers</code></a>
    </td>
  </tr>
</table>
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
- This package uses the fastest router: Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter)
- Works with `json`, `msgpack`, and `multi-part` forms
- Handles all standard types for `GetParams`
- Handler methods like `MakeParsedReq()` for `httprouter` use
- `Imbue` and `Permit` helper methods
- `GetParams()` parses parameters only once

<details>
<summary><strong><code>Development Setup (Getting Started)</code></strong></summary>
<br/>

Install [MAGE-X](https://github.com/mrz1836/mage-x) build tool for development:

```bash
# Install MAGE-X for development and building
go install github.com/mrz1836/mage-x/cmd/magex@latest
magex update:install
```
</details>

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

This project uses [goreleaser](https://github.com/goreleaser/goreleaser) for streamlined binary and library deployment to GitHub. To get started, install it via:

```bash
brew install goreleaser
```

The release process is defined in the [.goreleaser.yml](.goreleaser.yml) configuration file.

Then create and push a new Git tag using:

```bash
magex version:bump bump=patch push=true branch=master
```

This process ensures consistent, repeatable releases with properly versioned artifacts and citation metadata.

</details>

<details>
<summary><strong><code>Build Commands</code></strong></summary>
<br/>

View all build commands

```bash script
magex help
```

</details>

<details>
<summary><strong>GitHub Workflows</strong></summary>
<br/>

All workflows are driven by modular configuration in [`.github/env/`](.github/env/README.md) — no YAML editing required.

**[View all workflows and the control center →](.github/docs/workflows.md)**

</details>

<details>
<summary><strong><code>Updating Dependencies</code></strong></summary>
<br/>

To update all dependencies (Go modules, linters, and related tools), run:

```bash
magex deps:update
```

This command ensures all dependencies are brought up to date in a single step, including Go modules and any managed tools. It is the recommended way to keep your development environment and CI in sync with the latest versions.

</details>

<br/>

## Examples & Tests
All unit tests and [examples](examples) run via [GitHub Actions](https://github.com/mrz1836/go-template/actions) and use [Go version 1.24.x](https://go.dev/doc/go1.24). View the [configuration file](.github/workflows/fortress.yml).

Run all tests (fast):

```bash script
magex test
```

Run all tests with race detector (slower):
```bash script
magex test:race
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
magex bench
```
<br/>

## Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br>

## 🤖 AI Usage & Assistant Guidelines
Read the [AI Usage & Assistant Guidelines](.github/tech-conventions/ai-compliance.md) for details on how AI is used in this project and how to interact with AI assistants.

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


[![Stars](https://img.shields.io/github/stars/mrz1836/go-parameters?label=Please%20like%20us&style=social)](https://github.com/mrz1836/go-parameters/stargazers)

<br/>

## License

[![License](https://img.shields.io/github/license/mrz1836/go-parameters.svg?style=flat)](LICENSE)
