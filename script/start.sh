#!/bin/bash

set -ex

SH_DIR="$(cd "$(dirname "$0")"; pwd -P)"
ROOT_DIR=$(dirname $SH_DIR)
BUILD_DIR=${ROOT_DIR}

cd $BUILD_DIR

case "$1" in
  runall)
    docker-compose pull
    docker-compose up -d
    ;;
  runsimplify)
     docker-compose -f docker-compose-simplification.yml pull
     docker-compose -f docker-compose-simplification.yml up -d
     ;;
  down)
    docker-compose down -v
    ;;
  *)
    echo "unknow action."
    exit 1
esac
