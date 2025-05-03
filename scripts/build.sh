#!/bin/sh

set -e

cd "$(dirname "$(realpath "$0")")/.." || exit 1

if [ "$(node -v)" \< "v20" ]; then
  echo "You need to have node version 20 or greater installed"
  exit 1
fi

if [ ! -d node_modules ]; then
  npm install
fi

if [ ! -f connect/kafka-connect-loosehangerjeans-source-0.5.0-jar-with-dependencies.jar ]; then
  curl -L -o connect/kafka-connect-loosehangerjeans-source-0.5.0-jar-with-dependencies.jar https://github.com/IBM/kafka-connect-loosehangerjeans-source/releases/download/0.5.0/kafka-connect-loosehangerjeans-source-0.5.0-jar-with-dependencies.jar
fi

npm run lint:fix

go fmt

mage

npm run build
