#!/bin/bash
set -eux
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

CLOUD_PORT=${CLOUD_PORT:-5000}
CLOUD_EXPOSED_PORT=${CLOUD_EXPOSED_PORT:-5000}
CLOUD_URL=${CLOUD_URL:-http://localhost:$CLOUD_EXPOSED_PORT}
CLOUD_IMAGE_TAG=${CLOUD_IMAGE_TAG:-main}
CLOUD_IMAGE=${CLOUD_IMAGE:-ghcr.io/calyptia/cloud/all-in-one:$CLOUD_IMAGE_TAG}
TOKEN_DIR=${TOKEN_DIR:-$SCRIPT_DIR/resources}
TEST_AUTH0_DOMAIN=${TEST_AUTH0_DOMAIN:?}
TEST_AUTH0_MANAGEMENT_CLIENT_ID=${TEST_AUTH0_MANAGEMENT_CLIENT_ID:?}
TEST_AUTH0_MANAGEMENT_CLIENT_SECRET=${TEST_AUTH0_MANAGEMENT_CLIENT_SECRET:?}

docker rm -f cloud
docker run -d \
    --name cloud \
    -e AUTH0_DOMAIN="$TEST_AUTH0_DOMAIN" \
    -e AUTH0_MANAGEMENT_CLIENT_ID="$TEST_AUTH0_MANAGEMENT_CLIENT_ID" \
    -e AUTH0_MANAGEMENT_CLIENT_SECRET="$TEST_AUTH0_MANAGEMENT_CLIENT_SECRET" \
    -e DEBUG=true \
    -e DEFAULT_TOKEN_FILE=/token/token -v "$TOKEN_DIR":/token:Z \
    -p "$CLOUD_EXPOSED_PORT:$CLOUD_PORT" \
    "$CLOUD_IMAGE"

echo "Waiting for Cloud container to be ready"
until [[ $(curl -sL --write-out '%{http_code}' "$CLOUD_URL/healthz" -o /dev/null) -eq 200 ]]
do
    echo -n "."
    sleep 1
done
echo
echo "Container responding"
