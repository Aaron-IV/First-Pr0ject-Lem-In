#!/bin/bash

cd cmd
go build -o ../lemin .
cd ..
./lemin graph01.txt 
# start.sh - Script to start the lem-in project
# Author: rsatimov
# Date: $(date +%Y-%m-%d)
rm -rf lemin

set -euo pipefail
IFS=$'\n\t'