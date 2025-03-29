# Description

Custom linter to check err before struct usage

# Run in main

For local development run main in go-check-err-before-use/cmd/go-check-err-before-use with args: `example_files/positive.go example_files/negative.go`

Positive cases must pass without linter comments.

Negative cases must be commented by linter.

# Run as linter

Install local linter:

`go install ./cmd/...`

Then call it in some project:

`go-check-err-before-use ./...`

It will show related errors.

# Test

TestNegative fails now, that's ok.
