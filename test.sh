#!/usr/bin/env bash

# from https://github.com/codecov/example-go

set -e
echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
	# Run tests with coverage
    go test -v -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi

	# lint the package
	go vet $d
done
