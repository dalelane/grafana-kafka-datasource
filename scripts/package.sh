#!/bin/sh

set -e

cd "$(dirname "$(realpath "$0")")/.." || exit 1

VERSION=0.0.1

rm -f dalelane-kafka-datasource-$VERSION.zip

mv dist dalelane-kafka-datasource

zip dalelane-kafka-datasource-$VERSION.zip dalelane-kafka-datasource -r

mv dalelane-kafka-datasource dist

docker run --platform linux/amd64  --pull=always -v ./dalelane-kafka-datasource-$VERSION.zip:/archive.zip grafana/plugin-validator-cli /archive.zip
