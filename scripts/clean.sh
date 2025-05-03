#!/bin/sh

cd "$(dirname "$(realpath "$0")")/.." || exit 1

rm -rf dist
rm dalelane-kafka-datasource-*.zip
docker compose rm -s -f kafka connect grafana renderer
docker network rm dalelane-kafka-datasource_default
docker rmi dalelane-kafka-datasource-grafana
