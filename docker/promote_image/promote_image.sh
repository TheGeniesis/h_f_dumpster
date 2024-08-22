#!/bin/bash

set -ex
for ((COMMIT_INDEX=1; COMMIT_INDEX <= 3; ++COMMIT_INDEX))
do
  COMMIT_ID=$(curl -X GET --header "PRIVATE-TOKEN: ${GITLAB_ACCESS_TOKEN}" "${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/repository/commits?ref_name=${CI_COMMIT_TAG}&page=${COMMIT_INDEX}&per_page=1" | jq -r '.[].id');

  IMAGE_EXISTS="false"
  for IMAGE in ${IMAGES_LIST[*]}
  do
    for INDEX in ${PREV_APP_ENV[*]}
    do
      # if this will evaluate to true for the first image then all other images are also available
      if docker manifest inspect "${PREV_CONTAINER_REGISTRY}/${IMAGE}/${INDEX}:${COMMIT_ID}" >/dev/null; then
        docker pull "${PREV_CONTAINER_REGISTRY}/${IMAGE}/${INDEX}:${COMMIT_ID}"
        docker tag "${PREV_CONTAINER_REGISTRY}/${IMAGE}/${INDEX}:${COMMIT_ID}" "${CONTAINER_REGISTRY}/${IMAGE}/${APP_ENV}:${NEW_TAG_NAME}"
        docker tag "${PREV_CONTAINER_REGISTRY}/${IMAGE}/${INDEX}:${COMMIT_ID}" "${CONTAINER_REGISTRY}/${IMAGE}/${APP_ENV}:latest"
        docker push "${CONTAINER_REGISTRY}/${IMAGE}/${APP_ENV}:${NEW_TAG_NAME}"
        docker push "${CONTAINER_REGISTRY}/${IMAGE}/${APP_ENV}:latest"
        IMAGE_EXISTS="true"
        break;
      fi
    done;
  done;
  if [[ "${IMAGE_EXISTS}" == "true" ]]; then
    break
  fi;
done;

if [[ "${IMAGE_EXISTS}" == "false" ]]; then
  echo "Image: ${PREV_CONTAINER_REGISTRY}/${IMAGE}/${INDEX}:${COMMIT_ID} doesn't exists, skipping promotion... "
  echo "SHOULD_BUILD_NEW_IMAGES=true" >> art.prom.env;
  exit 0;
fi

echo "SHOULD_BUILD_NEW_IMAGES=false" >> art.prom.env;
