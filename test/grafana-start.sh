#!/bin/sh

password=$(pwgen -1)
url=http://portia.schweikert.ch:3000

if [ "$1" != "stop" ]; then
    docker run \
      -d \
      --rm \
      --net host \
      --name=grafana \
      -e "GF_SERVER_ROOT_URL=$url" \
      -e "GF_SECURITY_ADMIN_PASSWORD=$password" \
      grafana/grafana

    echo "grafana started: $url"
    echo "admin password: $password"
    echo
else
    docker stop grafana
fi
