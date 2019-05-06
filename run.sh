#!/usr/bin/env bash
set -ex

docker-compose up -d --build
docker-compose exec app sh -c "curl http://localhost:8080/number"
docker-compose exec app sh -c "wrk -c5 -t1 -d30s --latency http://localhost:8080/number"
docker-compose run --rm kafkacat -b kafka:9092 -t test -e | tail -n1
docker-compose exec app sh -c "curl http://localhost:8080/number"