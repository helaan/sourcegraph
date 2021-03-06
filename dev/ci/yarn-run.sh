#!/usr/bin/env bash

set -e

echo "--- yarn"
yarn --frozen-lockfile --network-timeout 60000
yarn --cwd lsif/server --frozen-lockfile --network-timeout 60000
yarn --cwd lsif/extension --frozen-lockfile --network-timeout 60000

for cmd in "$@"
do
    echo "--- $cmd"
    yarn -s run $cmd
done
