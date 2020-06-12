# Contributing

By participating to this project, you agree to abide our [code of conduct](/CODE_OF_CONDUCT.md).

## Setup your machine

`gm` is written in [Go](https://golang.org/).

Prerequisites:

- [Go 1.14+](https://golang.org/doc/install)

Clone `gm` anywhere:

```sh
$ git clone git@github.com:kordamp/gm.git
```

Install the build and lint dependencies:

```sh
$ go install github.com/fzipp/gocyclo
$ go get ./...
```

A good way of making sure everything is all right is running the test suite:

```sh
$ go test -v ./...
```

## Test your change

You can create a branch for your changes and try to build from the source as you go:

```sh
$ go build gm.go
```

## Submit a pull request

Push your branch to your `gm` fork and open a pull request against the `development` branch.

