#!/usr/bin/env bash

set -ex

APIHOST=${APIHOST:-api.poc.parvaeres.io}
MAOHOST=${MAOHOST:-160.85.252.48:8080}

docker run --rm --name parvaeres-ui \
    --tty \
    --interactive \
    --network k3d-parvaeres \
    --volume "$(pwd)/public:/go/src/github.com/parvaeres/go8s/public" \
    --volume "$(pwd)/app/views:/go/src/github.com/parvaeres/go8s/app/views" \
    --env APIHOST="${APIHOST}" \
    --env MAOHOST="${MAOHOST}" \
    --env APIVERSION=v1 \
    -p 9000:9000 \
    parvaeres/go8s-ui:latest
