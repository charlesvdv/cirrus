#!/usr/bin/env bash

set -e

tests_dir="$(dirname ${0})"
compose_file="${tests_dir}/docker-compose.yml"

function cleanup {
    docker-compose -f ${compose_file} down
}
trap cleanup EXIT

docker-compose -f ${compose_file} rm --force
sudo docker-compose -f ${compose_file} up --build --detach
docker-compose -f ${compose_file} logs --follow > ${tests_dir}/debug.log &

python3 ${tests_dir}/main.py --port 8080 --host localhost