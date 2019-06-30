#!/bin/sh

# docker volume create --label prometheus

if [ "$1" != "stop" ]; then
docker run --rm \
    -d \
    --net host \
    -v prometheus:/prometheus \
    -v `pwd`/prometheus.yml:/etc/prometheus/prometheus.yml \
    --name prometheus \
    prom/prometheus
else
    docker stop prometheus
fi
