#!/usr/bin/env bash

ALIYUN_IP=${ALIYUN_IP}
ALIYUN_IP=47.99.114.8

cd cmd

GOOS=linux GOARCH=amd64 go build -o test .
scp test root@${ALIYUN_IP}:~/