#!/bin/bash

export GO_ENV=test

go test ./internal/...

export GO_ENV=run
