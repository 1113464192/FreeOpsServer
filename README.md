# 欢迎来到FreeOps!
> 这里是后端代码

[前端代码入口](https://github.com/1113464192/FreeOpsClient)

FreeOps是一个功能齐全的运维自动化平台，只需要接入运维入口脚本即可使用，涵盖有基本的RBAC、每个用户的操作记录、项目管理、服务器管理、云平台同步操作、工单审批、游戏服管理、操作模板管理、模板参数管理、任务管理、任务日志等
基本上一个公司所需要的运维自动化业务，这里都能涵盖~
>精力有限，如果后续有人关注到，再补CI等功能

## 功能细节补充
> 普通用户(只有get权限且无法查看操作记录)
> 
> 账密
> 
> normal1
> 
> normal1_123456

[项目展示](http://106.52.66.254:81/)

[后端接口展示](http://106.52.66.254:9080/swagger/index.html#/)

> RBAC等常见功能就忽略了，只介绍运维功能
### 用户操作记录
包含所有请求与返回，防止单表过大做了水平分表处理。
### 项目管理
有字段限制每个服务器限制部署的 服数/Mem，接入运维脚本后，可以在该页面直接 创建云项目 以及 同步云平台 的项目信息
单项目单平台也可以多项目单平台
> 如果单项目单平台，可以将对应服务器的运维管理机设置为127.0.0.1
### 服务器管理
接入运维脚本后，可以在该页面选择服务器数量进行购买
### 模板参数管理
根据运营从PHP处生成的执行内容，取出关键字对应的内容，映射到关联模板命令的变量，避免了人工修改变量出错风险
### 任务管理
1. 可以关联模板顺序，选择主机(可选择内网执行or外网IP执行)，可选择按顺序执行or多模板同时执行
2. 在执行前还可以选择增删用审批人与命令，可以实时查看增删而致新生成的命令，如果不设置审批人与预设时间，则直接执行。
3. 在执行前可以预设执行时间，以凌晨12.举例：只要审批人(无审批人则到点自动执行)在12.前审批，则会等到12.自动执行。超过12.审批，则不予执行。
4. 有任务执行时可以点开modal查看任务实时进度。
### 审批任务
包含提交人、运营从PHP处生成的执行内容、检查脚本返回的检查内容、待执行的命令、执行的管理机IP、预设执行时间(如有设置) 等关键信息
### 任务日志
任务日志除了基本的每个命令的返回内容与SSH返回码，还包含提交人、运营从PHP处生成的执行内容、运维检查脚本返回的检查内容、审批人等关键信息
### WebSSH
> 这里key与passphrase示范的是通过配置文件读取使用，为了安全起见demo代码注释了写的功能

可以在网页实现与服务器的终端连接交互，这里的按钮权限理论上只给运维(便于运维有快捷操作的时候迅速登录服务器进行操作，免得进入xshell等应用)

这里的key与passphrase通过读配置文件(有jumpserver的情况下一般这样子做)，也可以直接每个用户录入密钥与passphrase字段(后端有AESCBC加密与解密的函数，即拿即用，加个字段和接口调用就好了)
## 部署方式
### 环境依赖
#### 前端

 - node-v18.20.3

#### 后端

 - mysql8.0+/mariadb:11.2
 - go-1.21.7

> P.S：go用1.19试过也可以，没有特殊库绝大多数版本应该通用，Mysql也没有使用特殊字段类型


> 部署上述所需环境，然后假设git clone到/data目录下
### 后端
#### 数据库
> 此处以mariadb举例，使用docker便于示例
> 
> 可以将mariadb换成mysql
##### vim /data/mariadb-docker.yaml 

    version: '3'
    services:
      db:
        image: "mariadb:11.2"
        ports:
          - "3306:3306"
        volumes:
          - /data/mariadb-11.2/data:/var/lib/mysql
        environment:
          TIME_ZONE: Asia/Shanghai
          MYSQL_ROOT_PASSWORD: "yourPassword"
##### 启动数据库
    docker-compose -f mariadb-docker.yaml up -d
##### 健康检查
    docker ps
##### 创建数据库
    docker exec -it yourDockerContainerID mariadb -uyourUser -p'yourPassword' -e "CREATE DATABASE yourDatabaseName CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
#### 启动后端服务
> 配置好configs/config.yaml，尤其是mysql/mariadb路径

    cd /data/FreeOpsServer/ && go run main.go
    // 如果能启动，则Crtrl+C关闭，然后执行编译
    go build -o FreeOpsServer main.go
    后台启动(二选一)：nohup ./FreeOpsServer > /dev/null 2>&1 &
    前台启动(二选一)：./FreeOpsServer
    docker cp /data/FreeOpsServer/init.sql yourDockerContainerID:/tmp/init.sql
    docker exec -it yourDockerContainerID bash
    mariadb -uroot -p'yourDBPassword' yourDatabaseName < /tmp/init.sql

### 前端
> 配置好.env*的对应后端url

    cd /data/FreeOpsClient && npm install -g pnpm
    pnpm i
    pnpm run build
### Nginx
#### 部署
    apt install nginx
#### vim /etc/nginx/conf.d/freeOps.conf 

    server { 
            listen 80;
            server_name yourAdmin;
            proxy_buffer_size 64k;
            proxy_buffers   32 32k;
            proxy_busy_buffers_size 128k;
            access_log /var/log/nginx/freeops.log;
            error_log /var/log/nginx/freeops_error.log;            
            location / { 
                root /data/FreeOpsClient/dist;
                index  index.html index.htm;
                try_files $uri $uri/ /index.html;
            } 
    	location /api/ {
                proxy_pass http://127.0.0.1:9080;
                proxy_http_version 1.1;
                proxy_set_header Upgrade $http_upgrade;    # WebSocket 协议升级
                proxy_set_header Connection "Upgrade";
                proxy_set_header Host $host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            
                # CORS 支持
                add_header 'Access-Control-Allow-Origin' $http_origin;
                add_header 'Access-Control-Allow-Credentials' 'true';
                add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS';
                add_header 'Access-Control-Allow-Headers' 'DNT,web-token,app-token,Authorization,Accept,Origin,Keep-Alive,User-Agent,X-Mx-ReqToken,X-Data-Type,X-Auth-Token,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,x-request-id';
                add_header 'Access-Control-Expose-Headers' 'Content-Length,Content-Range';
            }
    }
#### 刷新
    nginx -t && nginx -s reload

## P.S
1. 可以定时任务检查后台输出是否有日志写入失败,定时任务检查日志十分钟内是否有ERROR
2. 数据库如果量级不大且对性能要求不高可以改为使用主外键关联，可以省去做约束的时间
3. 如果新增/删除Model，记得从consts和model的init中更改常量/变量

    
    



 

## 联系我
 如果在二开/学习过程中有疑问，可以联系邮箱：fqh1113464192@gmail.com
 因为这也是我初次从0手搓运维平台，如果有意见或者批评也可以联系我，看到就会回。
 辛苦看到这里！！再次感谢！
 

