#!/bin/bash
set -eux
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

CLOUD_URL=${CLOUD_URL:-http://localhost:5000}
CLOUD_IMAGE_TAG=${CLOUD_IMAGE_TAG:-main}
CLOUD_IMAGE=${CLOUD_IMAGE:-ghcr.io/calyptia/cloud:$CLOUD_IMAGE_TAG}
TOKEN_DIR=${TOKEN_DIR:-$SCRIPT_DIR/resources}
TOKENFILE="$TOKEN_DIR/token"

if ! docker pull "$CLOUD_IMAGE"; then
    echo "Unable to pull cloud image: $CLOUD_IMAGE"
    exit 1
fi

# TODO: remove once we have all-in-one container: https://github.com/calyptia/cloud/issues/308
CLOUD_DIR=${CLOUD_DIR:-$SCRIPT_DIR/resources/cloud}

if [[ ! -d "$CLOUD_DIR" ]]; then
    echo "No CLOUD_DIR at '$CLOUD_DIR', please clone and select branch to use first."
    exit 1
fi

if [[ ! -f "$CLOUD_DIR/docker-compose.yml" ]]; then
    echo "No compose stack defined at $CLOUD_DIR/docker-compose.yml"
    exit 1
fi

docker-compose up -d --project-dir="$CLOUD_DIR"
# END OF TODO

# TODO: only due to permissions issues: https://github.com/calyptia/cloud/issues/309
mkdir -p "$TOKEN_DIR"
rm -f "$TOKENFILE"
touch "$TOKENFILE"
chmod a+r "$TOKENFILE"
# END OF TODO

docker run -d --network=host \
    --name cloud \
    -e DEBUG=true \
    -e DEFAULT_TOKEN_FILE=/token/token -v "$TOKEN_DIR":/token:Z \
    "$CLOUD_IMAGE"

echo "Waiting for Cloud container to be ready"
until [[ $(curl -sL --write-out '%{http_code}' "$CLOUD_URL/healthz" -o /dev/null) -eq 200 ]]
do
    echo -n "."
    sleep 1
done
echo
echo "Container responding"

echo "Waiting for token to be created"
until [[ -f "$TOKENFILE" ]]
do
    echo -n "."
    sleep 1
done
echo
echo "Found token file"
