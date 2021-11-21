#!/bin/bash

case $1 in
"up")
  docker-compose -f deploy/docker-compose.yaml up -d coturn redis mysql etcd nsqd nsqlookupd nsqadmin
;;
"down")
  docker rm -f baby-redis baby-mysql etcd nsqd nsqlookupd nsqadmin
;;
esac
