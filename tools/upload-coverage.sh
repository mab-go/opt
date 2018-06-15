#!/usr/bin/env bash

# Only upload code coverage reports for Go 1.10.
if [[ "${TRAVIS_GO_VERSION}" =~ ^1\.10 ]]; then
    bash <(curl -s https://codecov.io/bash)
fi
