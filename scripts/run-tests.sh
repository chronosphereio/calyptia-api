#!/bin/bash
set -eux
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

CLOUD_URL=${CLOUD_URL:-http://localhost:5000}
SPEC_DIR=${SPEC_DIR:-$SCRIPT_DIR/../spec}
WORKER_COUNT=${WORKER_COUNT:-8}
TOKEN_DIR=${TOKEN_DIR:-$SCRIPT_DIR/resources}
TOKENFILE="$TOKEN_DIR/token"

if [[ -z "${TOKEN:-}" ]]; then
    if [[ ! -f "$TOKENFILE" ]]; then
        echo "No TOKEN or invalid TOKENFILE defined"
        exit 1
    fi
    TOKEN=$(cat "$TOKENFILE")
fi

docker run --rm --network=host -v "$SPEC_DIR/":/spec:ro \
    schemathesis/schemathesis:stable \
        run \
        --header "Authorization:" \
        --header "X-Project-Token:$TOKEN" \
        --stateful=links \
        --workers "$WORKER_COUNT" \
        /spec/open-api.yml \
        --base-url="$CLOUD_URL/"
