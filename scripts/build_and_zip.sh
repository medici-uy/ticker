#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
[[ "${DEBUG:-false}" == "true" ]] && set -o xtrace

binary_name="bootstrap"

go build -tags lambda.norpc -o "${binary_name}" main.go
zip "${ZIP_FILE}" bootstrap
ls
