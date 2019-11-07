#!/bin/bash

# 该文件为自动部署脚本，但是测试使用的vps为ubuntu 18，其他种类或版本的linux系统不一定支持

# env
# 其中，除了ONEDRIVE_TOKEN必须设置，其他都有默认值，详情参见README.md

# 应用ID
APP_ID=
APP_PWD=
# 应用密码
# aria_prc的密码
ARIA_RPC_PWD=
# onedrive云上的同步根文件夹
ONEDRIVE_BASE_PATH=
# onedrive的to
# ken，在本地使用rclone authorize "onedrive" 获取
ONEDRIVE_TOKEN='必须配置'

# version 默认不指定版本号，缺省即使用latest
#dl_api_version=":v0.0.1"
#dl_vue_version=":v0.0.1"

# 由于程序自带了环境变量的默认值，所以如果上方的配置项留空，那么就不传入该配置项
# 该方法接收配置项项名，如果该配置项不为空(以APP_ID=me为例，传入APP_ID)，那么就返回 "-e APP_ID=me", 如果(APP_ID=)，那么就不返回，docker也就不会传入该配置项
fmt_docker_env() {
  env_name=$1
  eval env_var=`echo '$'$env_name`
  if [ "$env_var" != "" ]
  then
    echo "-e ${env_name}=${env_var}"
  fi
}

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
docker pull redis
docker run -d \
        --network=dl \
        --hostname=redis \
        --name=redis redis
# 开启 api docker
docker pull mopip77/dl_api${dl_api_version}
docker run -d \
        -e REDIS_ADDR=redis:6379 \
        -e ONEDRIVE_TOKEN=${ONEDRIVE_TOKEN} \
        `fmt_docker_env ARIA_RPC_PWD` \
        `fmt_docker_env APP_ID` \
        `fmt_docker_env APP_PWD` \
        `fmt_docker_env ONEDRIVE_BASE_PATH` \
        --network=dl \
        --hostname=dl_api \
        -p 3000:3000 \
        -p 6800:6800 \
        --name=dl_api mopip77/dl_api${dl_api_version}
#开启 vue docker
docker pull mopip77/dl_vue${dl_vue_version}
docker run -d --network=dl \
        --hostname=vue \
        --name=vue \
        -p 3001:80 mopip77/dl_vue${dl_vue_version}

# 由于api容器安装较慢，所以最后输出一下其日志，用于判断部署完成
docker logs -f dl_api