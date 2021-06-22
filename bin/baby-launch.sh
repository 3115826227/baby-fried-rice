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

cd $CMD/userAccount
exec -a baby-user-account ./userAccount &
cd $DIR

cd $CMD/adminAccount
exec -a baby-admin-account ./adminAccount &
cd $DIR

cd $CMD/rootAccount
exec -a baby-root-account ./rootAccount &
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

trap cleanup EXIT

while : ; do sleep 1 ; done