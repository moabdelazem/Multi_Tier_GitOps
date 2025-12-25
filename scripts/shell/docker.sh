#!/bin/env bash

set -e

COMMIT_HASH=$(git rev-parse --short HEAD)
IMAGE_NAME="multitier-api"
DOCKERHUB_USERNAME="${DOCKERHUB_USERNAME:-moabdelazem}"

echo "##############################"
echo "Building API Docker Image And Pushing it to DockerHub"
echo "##############################"

# Build
docker build -t "${DOCKERHUB_USERNAME}/${IMAGE_NAME}:${COMMIT_HASH}" ./api

# Tag as latest
docker tag "${DOCKERHUB_USERNAME}/${IMAGE_NAME}:${COMMIT_HASH}" "${DOCKERHUB_USERNAME}/${IMAGE_NAME}:latest"

# Push both tags
docker push "${DOCKERHUB_USERNAME}/${IMAGE_NAME}:${COMMIT_HASH}"
docker push "${DOCKERHUB_USERNAME}/${IMAGE_NAME}:latest"

echo "##############################"
echo "Done! Pushed:"
echo "  - ${DOCKERHUB_USERNAME}/${IMAGE_NAME}:${COMMIT_HASH}"
echo "  - ${DOCKERHUB_USERNAME}/${IMAGE_NAME}:latest"
echo "##############################"
