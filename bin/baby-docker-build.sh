#!/bin/sh

case $1 in
'')
  docker rmi baby-gateway
  docker build -t baby-gateway -f cmd/gateway/Dockerfile .
  docker rmi baby-user-account
  docker build -t baby-user-account -f cmd/userAccount/Dockerfile .
  docker rmi baby-account-dao
  docker build -t baby-account-dao -f cmd/accountDao/Dockerfile .
  docker rmi baby-space
  docker build -t baby-space -f cmd/space/Dockerfile .
  docker rmi baby-space-dao
  docker build -t baby-space-dao -f cmd/spaceDao/Dockerfile .
  docker rmi baby-im
  docker build -t baby-im -f cmd/im/Dockerfile .
  docker rmi baby-im-dao
  docker build -t baby-im-dao -f cmd/imDao/Dockerfile .
  docker rmi baby-connect
  docker build -t baby-connect -f cmd/connect/Dockerfile .
  docker rmi baby-file
  docker build -t baby-file -f cmd/file/Dockerfile .
;;
"gateway")
  docker rmi baby-gateway
  docker build -t baby-gateway -f cmd/gateway/Dockerfile .
;;
"user-account")
  docker rmi baby-user-account
  docker build -t baby-user-account -f cmd/userAccount/Dockerfile .
;;
"account-dao")
  docker rmi baby-account-dao
  docker build -t baby-account-dao -f cmd/accountDao/Dockerfile .
;;
"space")
  docker rmi baby-space
  docker build -t baby-space -f cmd/space/Dockerfile .
;;
"space-dao")
  docker rmi baby-space-dao
  docker build -t baby-space-dao -f cmd/spaceDao/Dockerfile .
;;
"im")
  docker rmi baby-im
  docker build -t baby-im -f cmd/im/Dockerfile .
;;
"im-dao")
  docker rmi baby-im-dao
  docker build -t baby-im-dao -f cmd/imDao/Dockerfile .
;;
"connect")
  docker rmi baby-connect
  docker build -t baby-connect -f cmd/connect/Dockerfile .
;;
"file")
  docker rmi baby-file
  docker build -t baby-file -f cmd/file/Dockerfile .
;;
esac
