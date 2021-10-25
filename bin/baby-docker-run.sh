#!/bin/bash

case $1 in
"up")
  docker-compose -f deploy/docker-compose.yaml up -d
;;
"down")
  docker-compose -f deploy/docker-compose.yaml down
;;
esac