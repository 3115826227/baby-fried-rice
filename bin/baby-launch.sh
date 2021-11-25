#!/bin/bash

DIR=$PWD
CMD=../cmd

# Kill all baby-* stuff
function cleanup {
	pkill baby
}

cd $CMD/gateway
exec -a baby-gateway ./gateway &
cd $DIR

cd $CMD/manage
exec -a baby-manage ./manage &
cd $DIR

cd $CMD/userAccount
exec -a baby-user-account ./userAccount &
cd $DIR

cd $CMD/accountDao
exec -a baby-account-dao ./accountDao &
cd $DIR

cd $CMD/spaceDao
exec -a baby-space-dao ./spaceDao &
cd $DIR

cd $CMD/space
exec -a baby-space ./space &
cd $DIR

cd $CMD/imDao
exec -a baby-im-dao ./imDao &
cd $DIR

cd $CMD/im
exec -a baby-im ./im &
cd $DIR

cd $CMD/connect
exec -a baby-connect ./connect &
cd $DIR

cd $CMD/file
exec -a baby-file ./file &
cd $DIR

cd $CMD/shopDao
exec -a baby-shop-dao ./shopDao &
cd $DIR

cd $CMD/shop
exec -a baby-shop ./shop &
cd $DIR

cd $CMD/smsDao
exec -a baby-sms-dao ./smsDao &
cd $DIR

cd $CMD/gameDao
exec -a baby-game-dao ./gameDao &
cd $DIR

cd $CMD/game
exec -a baby-game ./game &
cd $DIR

cd $CMD/liveDao
exec -a baby-live-dao ./liveDao &
cd $DIR

cd $CMD/live
exec -a baby-live ./live &
cd $DIR

trap cleanup EXIT

while : ; do sleep 1 ; done