#!/bin/sh

if [ "$1" != "stop" ]; then
docker run \
  -d \
  --rm \
  --net host \
  --name=grafana \
  -e "GF_SERVER_ROOT_URL=http://portia.schweikert.ch:3000" \
  -e "GF_SECURITY_ADMIN_PASSWORD=lou6iser" \
  grafana/grafana
else
docker stop grafana
fi
