#!/bin/sh

case $1 in
'')
  docker build -t baby-gateway -f cmd/gateway/Dockerfile .
  docker build -t baby-user-account -f cmd/userAccount/Dockerfile .
  docker build -t baby-account-dao -f cmd/accountDao/Dockerfile .
  docker build -t baby-space -f cmd/space/Dockerfile .
  docker build -t baby-space-dao -f cmd/spaceDao/Dockerfile .
  docker build -t baby-im -f cmd/im/Dockerfile .
  docker build -t baby-im-dao -f cmd/imDao/Dockerfile .
  docker build -t baby-connect -f cmd/connect/Dockerfile .
  docker build -t baby-file -f cmd/file/Dockerfile .
;;
"gateway")
  docker build -t baby-gateway -f cmd/gateway/Dockerfile .
;;
"user-account")
  docker build -t baby-user-account -f cmd/userAccount/Dockerfile .
;;
"account-dao")
  docker build -t baby-account-dao -f cmd/accountDao/Dockerfile .
;;
"space")
  docker build -t baby-space -f cmd/space/Dockerfile .
;;
"space-dao")
  docker build -t baby-space-dao -f cmd/spaceDao/Dockerfile .
;;
"im")
  docker build -t baby-im -f cmd/im/Dockerfile .
;;
"im-dao")
  docker build -t baby-im-dao -f cmd/imDao/Dockerfile .
;;
"connect")
  docker build -t baby-connect -f cmd/connect/Dockerfile .
;;
"file")
  docker build -t baby-file -f cmd/file/Dockerfile .
;;
esac
