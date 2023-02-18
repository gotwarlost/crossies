#!/bin/bash

set -euxo pipefail

GOOS=linux CGO_ENABLED=0 go build -o ./site/cgi-bin/api.fcgi ./cmd/crossie-fcgi
ssh "${CROSSIE_MACHINE}" killall api.fcgi || true
rsync -rvtl ./site/ "${CROSSIE_MACHINE}:${CROSSIE_HOME}/"
ssh "${CROSSIE_MACHINE}" "${CROSSIE_HOME}/cgi-bin/cache-buster.sh"
