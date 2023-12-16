# alert-system
> A go microservice for managing alerts and runs alongside Bitcoin SV nodes utilizing RPC

[![Release](https://img.shields.io/github/release-pre/bitcoin-sv/alert-system.svg?logo=github&style=flat&v=2)](https://github.com/bitcoin-sv/alert-system/releases)
[![Build](https://github.com/bitcoin-sv/alert-system/workflows/run-go-tests/badge.svg?branch=master&v=1)](https://github.com/bitcoin-sv/alert-system/actions)
[![Report](https://goreportcard.com/badge/github.com/bitcoin-sv/alert-system?style=flat&v=2)](https://goreportcard.com/report/github.com/bitcoin-sv/alert-system)
[![Go](https://img.shields.io/badge/Go-1.21.xx-blue.svg?v=1)](https://golang.org/)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat&v=2)](https://github.com/RichardLitt/standard-readme)
[![Makefile Included](https://img.shields.io/badge/Makefile-Supported%20-brightgreen?=flat&logo=probot&v=2)](Makefile)
<br> <!-- [![Go](https://img.shields.io/github/go-mod/go-version/bitcoin-sv/alert-system?v=2)](https://golang.org/) -->

<br/>

## Table of Contents
- [Installation](#installation)
- [Documentation](#documentation)
- [Examples & Tests](#examples--tests)
- [Benchmarks](#benchmarks)
- [Code Standards](#code-standards)
- [Contributing](#contributing)
- [License](#license)

<br/>

## Installation

**alert-system** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).

To run the application, clone this repository locally and run:
```shell script
export ALERT_SYSTEM_ENVIRONMENT=local && go run cmd/main.go
```

To run this application with a custom configuration file, run:
```shell script
export ALERT_SYSTEM_CONFIG_FILEPATH=path/to/file/config.json && go run cmd/main.go
```

Configuration files can be found in the [config](app/config/envs) directory.

<br/>

## Container Environment
**Note:** to use a custom settings file, it needs to be mounted and the appropriate environment variables set. Running it as below will run an ephemeral database but the container should sync up from the peers on the network on startup.
### podman
```
$ podman run -u root -e P2P_PORT=9908 -e P2P_IP=0.0.0.0  --expose 9908 docker.io/bsvb/alert-system:0.0.2
```

## Documentation
View the generated [documentation](https://pkg.go.dev/github.com/bitcoin-sv/alert-system)

[![GoDoc](https://godoc.org/github.com/bitcoin-sv/alert-system?status.svg&style=flat&v=2)](https://pkg.go.dev/github.com/bitcoin-sv/alert-system)

<br/>

<details>
<summary><strong><code>Makefile Commands</code></strong></summary>
<br/>

View all `makefile` commands
```shell script
make help
```

List of all current commands:
```text
all                   Runs multiple commands
clean                 Remove previous builds and any cached data
clean-mods            Remove all the Go mod cache
coverage              Shows the test coverage
diff                  Show the git diff
generate              Runs the go generate command in the base of the repo
godocs                Sync the latest tag with GoDocs
help                  Show this help message
install               Install the application
install-go            Install the application (Using Native Go)
install-releaser      Install the GoReleaser application
lint                  Run the golangci-lint application (install if not found)
release               Full production release (creates release in GitHub)
release               Runs common.release then runs godocs
release-snap          Test the full release (build binaries)
release-test          Full production test release (everything except deploy)
replace-version       Replaces the version in HTML/JS (pre-deploy)
tag                   Generate a new tag and push (tag version=0.0.0)
tag-remove            Remove a tag if found (tag-remove version=0.0.0)
tag-update            Update an existing tag to current commit (tag-update version=0.0.0)
test                  Runs lint and ALL tests
test-ci               Runs all tests via CI (exports coverage)
test-ci-no-race       Runs all tests via CI (no race) (exports coverage)
test-ci-short         Runs unit tests via CI (exports coverage)
test-no-lint          Runs just tests
test-short            Runs vet, lint and tests (excludes integration tests)
test-unit             Runs tests and outputs coverage
uninstall             Uninstall the application (and remove files)
update-linter         Update the golangci-lint package (macOS only)
vet                   Run the Go vet application
```
</details>

<br/>

## Examples & Tests
All unit tests and examples run via [GitHub Actions](https://github.com/bitcoin-sv/alert-system/actions) and
uses [Go version 1.21.x](https://golang.org/doc/go1.21). View the [configuration file](.github/workflows/run-tests.yml).

<br/>

Run all tests (including integration tests)
```shell script
make test
```

<br/>

Run tests (excluding integration tests)
```shell script
make test-short
```

<br/>

## Benchmarks
Run the Go benchmarks:
```shell script
make bench
```

<br/>

## Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

## Contributing
View the [contributing guidelines](.github/CONTRIBUTING.md) and follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

<br/>

## License

[![License](https://img.shields.io/badge/license-OpenBSV-green.svg?style=flat&v=2)](LICENSE)