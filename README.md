### 一站式下载部署
本程序提供aria2(下载http，磁力)，和you-get的视频解析下载服务的结合。

程序使用go的后台 + redis + Vue的前端

一站式部署脚本为`allrun.sh`，需要传入`rclone authorize "onedrive"`的token值  
除此之外还可以选择性设置aira2 rpc的密码，只需要设置环境变量`ARIA_PWD`即可

由于aira2 rpc端口为6800，所以vps防火墙除了需要开启80，也需要开启6800

#### 待添加
- 显示vps剩余容量