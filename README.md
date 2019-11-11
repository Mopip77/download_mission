### 一站式下载部署
本程序提供aria2(下载HTTP，FTP，磁力等)，和you-get的视频解析下载服务的结合。  
并且视频解析下载采用的是任务卡片式，可以中断或重启任务、查看输出。  
![屏幕快照 2019-11-11 下午7.23.14.png](https://i.loli.net/2019/11/11/YJwLcXRNlmIu9O7.png)

程序使用go的后台 + redis + Vue的前端，并且这三者分别使用三个不同的docker image

一站式部署脚本为`allrun.emample.sh`(只在Ubuntu18下使用，其他的Linux系统可能不支持)，需要自己填写的就是onedrive的token值  
可以在本机上执行`rclone authorize "onedrive"`来获取

由于aira2 rpc端口为6800，所以vps防火墙除了需要开启80，也需要开启6800

#### 可配置的环境变量
分为api后台的配置和aria，onedrive应用的配置，分别可以在`.env.example`和`Dockerfile`中找到  
都可以在开启docker时传入用于替代，重点的几个配置项是：  
(由于使用场景的单一性，所以这里只使用了`BasicAuth`，所以整个应用只有一个用户名密码)
- `APP_ID` 应用用户名 默认guest
- `APP_PWD` 应用密码 默认guest
- `ARIA_RPC_PWD` Aria Rpc的密码 默认123