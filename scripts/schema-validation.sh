#!/bin/bash
set -eux
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

CLOUD_URL=${CLOUD_URL:-http://localhost:5000}
SPEC_DIR=${SPEC_DIR:-$SCRIPT_DIR/../spec}
WORKER_COUNT=${WORKER_COUNT:-8}

if [[ -n "${TOKENFILE:-}" ]]; then
    if [[ ! -f "$TOKENFILE" ]]; then
        echo "No such file: $TOKENFILE"
    fi
    TOKEN=$(cat "$TOKENFILE")
fi

if [[ -n "${TOKEN:-}" ]]; then
    echo "No TOKEN or invalid TOKENFILE defined"
    exit 1
fi

echo "Waiting for Cloud container to be ready"
until [[ $(curl -sL --write-out '%{http_code}' "$CLOUD_URL/healthz" -o /dev/null) -eq 200 ]]
do
    echo -n "."
    sleep 1
done
echo
echo "Container responding"

docker run --rm --network=host -v "$SPEC_DIR/":/spec:ro \
    schemathesis/schemathesis:stable \
        run \
        --header "Authorization:" \
        --header "X-Project-Token:$TOKEN" \
        --stateful=links \
        --workers "$WORKER_COUNT" \
        /spec/open-api.yml \
        --base-url="$CLOUD_URL/"
