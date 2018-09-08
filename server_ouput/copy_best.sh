#!/usr/bin/env bash

SERVER_IP=${SERVER_IP}

scp root@${SERVER_IP}:~/_output/c/best.csv ./_output/c/
scp root@${SERVER_IP}:~/_output/c/best_summary.csv ./_output/c/

scp root@${SERVER_IP}:~/_output/d/best.csv ./_output/d/
scp root@${SERVER_IP}:~/_output/d/best_summary.csv ./_output/d/