#!/usr/bin/env bash

cd cmd

GOOS=linux GOARCH=amd64 go build -o aliyun .
scp aliyun root@106.14.204.11:~/aliyun_semi/