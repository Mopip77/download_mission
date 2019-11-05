#!/bin/bash

# 该文件为自动部署脚本，但是测试使用的vps为ubuntu 18，其他种类或版本的linux系统不一定支持

# env
# aria_prc的密码
ARIA_RPC_PWD='123'
# onedrive云上的同步根文件夹
ONEDRIVE_BASE_PATH='/share'
# onedrive的to
# ken，在本地使用rclone authorize "onedrive" 获取
ONEDRIVE_TOKEN=''

# version 默认不指定版本号，缺省即使用latest
#dl_api_version=":v0.0.1"
#dl_vue_version=":v0.0.1"

apt update

# nginx 部分
apt install -y nginx
systemctl start nginx
echo -e '
server {
    listen 80;

location / {
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_pass http://127.0.0.1:3001;
    }

    location /api {
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_pass http://127.0.0.1:3000;
    }
}
' > /etc/nginx/sites-enabled/default
nginx -s reload

# docker 部分
apt install -y docker.io

# 创建一个私有的局域网用于通信
docker network create dl
# 开启 redis docker
docker run -d \
        --network=dl \
        --hostname=redis \
        --name=redis redis
# 开启 api docker
docker run -d \
        -e REDIS_ADDR=redis:6379 \
        -e ARIA_RPC_PWD=${ARIA_RPC_PWD} \
        -e ONEDRIVE_BASE_PATH=${ONEDRIVE_BASE_PATH} \
        -e ONEDRIVE_TOKEN=${ONEDRIVE_TOKEN} \
        --network=dl \
        --hostname=dl_api \
        -p 3000:3000 \
        -p 6800:6800 \
        --name=dl_api mopip77/dl_api${dl_api_version}
#开启 vue docker
docker run -d --network=dl \
        --hostname=vue \
        --name=vue \
        -p 3001:80 mopip77/dl_vue${dl_vue_version}

# 由于api容器安装较慢，所以最后输出一下其日志，用于判断部署完成
docker logs -f dl_api