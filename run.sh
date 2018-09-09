#!/usr/bin/env bash

cd cmd
go build
mv cmd ../test
cd ../
./test