#!/bin/bash

set -e
unset CDPATH
cd "$(dirname "${BASH_SOURCE[0]}")/.." # cd to repo root dir

for dir in web shared browser packages/sourcegraph-extension-api packages/@sourcegraph/extension-api-types lsif/server lsif/extension; do
    (set -x; cd "$dir" && "$@")
done
