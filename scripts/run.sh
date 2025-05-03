#!/bin/sh

cd "$(dirname "$(realpath "$0")")/.." || exit 1

npm run server
