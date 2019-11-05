FROM python:3.5.8-slim-stretch

COPY dl_api .env  /root/api/
COPY ./script /root/api/script

ENV ARIA_RPC_PWD 123
ENV ONEDRIVE_BASE_PATH /share

RUN apt update -y && \
    apt install -y wget curl unzip expect procps bc && \
    pip3 install you-get

WORKDIR /root/api/

ENTRYPOINT "./script/deploy/run.sh"