#!/usr/bin/env bash

# Only run tests with code coverage for Go 1.10. For other versions of Go, run
# tests without code coverage.
if [[ "${TRAVIS_GO_VERSION}" =~ ^1\.10 ]]; then
    go test -v -coverprofile="coverage.txt" ./...
else
    go test -v ./...
fi
