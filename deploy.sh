#!/usr/bin/env bash

ALIYUN_IP=${ALIYUN_IP}

cd cmd

GOOS=linux GOARCH=amd64 go build -o aliyun .
scp aliyun root@${ALIYUN_IP}:~/aliyun_semi/