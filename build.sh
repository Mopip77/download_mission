#!/bin/bash

version="v0.1.0"

GOOS=linux GOARCH=amd64 go build -o dl_api
docker build -t mopip77/dl_api:${version} ./
docker push mopip77/dl_api:${version}
docker tag mopip77/dl_api:${version} mopip77/dl_api:latest
docker push mopip77/dl_api:latest
docker rmi mopip77/dl_api:latest
rm dl_api