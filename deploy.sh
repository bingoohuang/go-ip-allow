#!/usr/bin/env bash

# ./deploy.sh ci@ci01
targetHost=$1
deployName=go-ip-allow
./gobin.sh
rm -fr $deployName.linux.bin $deployName.linux.bin.bz2
env GOOS=linux GOARCH=amd64 go build -o $deployName.linux.bin
bzip2 $deployName.linux.bin
rsync -avz --human-readable --progress -e "ssh -p 22" ./$deployName.linux.bin.bz2 $targetHost:./firewall
ssh -tt $targetHost "bash -s" << eeooff
cd ./firewall
ps -ef|grep $deployName|grep -v grep|awk '{print \$2}'|xargs -r kill -9
rm $deployName.linux.bin
bzip2 -d $deployName.linux.bin.bz2
nohup ./$deployName.linux.bin 2>&1 >> go-ip-allow.log &
exit
eeooff