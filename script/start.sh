#!/bin/bash

set -ex

SH_DIR="$(cd "$(dirname "$0")"; pwd -P)"
ROOT_DIR=$(dirname $SH_DIR)
BUILD_DIR=${ROOT_DIR}

cd $BUILD_DIR

case "$1" in
  runall)
    docker-compose --profile general pull
    docker-compose --profile general up -d
    ;;
  runsimplify)
     docker-compose --profile simplification pull
     docker-compose --profile simplification up -d
     ;;
  down)
    docker-compose down -v
    ;;
  *)
    echo "unknow action."
    exit 1
esac
