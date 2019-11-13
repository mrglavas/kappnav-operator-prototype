#!/bin/bash
VERSION=0.0.1
IMAGE=kappnav.io/kappnav-operator

CURRENT=`pwd`
PROJECT=`basename "$CURRENT"`

echo "Building ${IMAGE} ${VERSION}"
operator-sdk build --image-build-args "--pull --build-arg VERSION=${VERSION} --build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') --build-arg PROJECT_NAME=${PROJECT}" ${IMAGE}:${VERSION}