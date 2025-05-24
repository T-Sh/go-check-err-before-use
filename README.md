# Description

Linter to check err before struct usage

# Install

`go install github.com/T-Sh/go-check-err-before-use@latest`

Linter requires Go 1.23 or later version.

# Use

To check all source code in current directory run:

`go-check-err-before-use ./...`

Supported flags in command line:

 - distance - the maximum acceptable distance between the assignment of an error and its checking

# Testing and development

For local development run main in go-check-err-before-use/cmd/go-check-err-before-use with args: `example_files/positive.go example_files/negative.go`

Positive cases must pass without linter comments.

Negative cases must be commented by linter.

TestNegative fails now, that's ok.
