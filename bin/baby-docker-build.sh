#!/bin/sh

case $1 in
'')
  docker rmi baby-gateway
  docker build -t baby-gateway -f cmd/gateway/Dockerfile .
  docker rmi baby-manage
  docker build -t baby-manage -f cmd/manage/Dockerfile .
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
  docker rmi baby-shop
  docker build -t baby-shop -f cmd/shop/Dockerfile .
  docker rmi baby-shop-dao
  docker build -t baby-shop-dao -f cmd/shopDao/Dockerfile .
  docker rmi baby-sms-dao
  docker build -t baby-sms-dao -f cmd/smsDao/Dockerfile .
  docker rmi baby-game
  docker build -t baby-game -f cmd/game/Dockerfile .
  docker rmi baby-game-dao
  docker build -t baby-game-dao -f cmd/gameDao/Dockerfile .
  docker rmi baby-live
  docker build -t baby-live -f cmd/live/Dockerfile .
  docker rmi baby-live-dao
  docker build -t baby-live-dao -f cmd/liveDao/Dockerfile .
  docker rmi baby-coturn
  docker build -t baby-coturn -f cmd/coturn/Dockerfile
;;
"gateway")
  docker rmi baby-gateway
  docker build -t baby-gateway -f cmd/gateway/Dockerfile .
;;
"backend")
  docker rmi baby-manage
  docker build -t baby-manage -f cmd/manage/Dockerfile .
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
"shop")
  docker rmi baby-shop
  docker build -t baby-shop -f cmd/shop/Dockerfile .
;;
"shop-dao")
  docker rmi baby-shop-dao
  docker build -t baby-shop-dao -f cmd/shopDao/Dockerfile .
;;
"sms-dao")
  docker rmi baby-sms-dao
  docker build -t baby-sms-dao -f cmd/smsDao/Dockerfile .
;;
"game")
  docker rmi baby-game
  docker build -t baby-game -f cmd/game/Dockerfile .
;;
"game-dao")
  docker rmi baby-game-dao
  docker build -t baby-game-dao -f cmd/gameDao/Dockerfile .
;;
"live")
  docker rmi baby-live
  docker build -t baby-live -f cmd/live/Dockerfile .
;;
"live-dao")
  docker rmi baby-live-dao
  docker build -t baby-live-dao -f cmd/liveDao/Dockerfile .
;;
"coturn")
  docker rmi baby-coturn
  docker build -t baby-coturn -f cmd/coturn/Dockerfile .
;;
esac
