#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o dl_api
docker build -t mopip77/dl_api:v0.0.1 ./
docker push mopip77/dl_api:v0.0.1
rm dl_api