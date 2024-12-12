# 欢迎来到FreeOps!
> 这里是后端代码

[前端代码](https://github.com/1113464192/FreeOpsClient)

FreeOps是一个功能齐全的运维自动化平台，只需要接入运维入口脚本即可使用，涵盖有基本的RBAC、每个用户的操作记录、项目管理、服务器管理、云平台同步操作、工单审批、游戏服管理、操作模板管理、模板参数管理、任务管理、任务日志等
基本上一个中小型公司所需要的运维自动化业务，这里都能涵盖~
>精力有限，如果后续有人关注到，再补webssh、CI等功能(后端已实现)

## 功能细节补充

> RBAC等常见功能就忽略了，具体可以到下面查看
> 
[项目展示](http://106.52.66.254:81/)
[后端接口展示](http://106.52.66.254:9080/swagger/index.html#/)
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
可以关联模板顺序，选择主机(可选择内网执行or外网IP执行)，可选择按顺序执行or多模板同时执行，在执行前还可以选择增删用审批人与命令，如果不设置审批人，则直接执行。
有任务执行时可以点开modal查看任务实时进度。
### 审批任务
包含提交人、运营从PHP处生成的执行内容、检查脚本返回的检查内容、待执行的命令、执行的管理机IP 等关键信息
### 任务日志
任务日志除了基本的每个命令的返回内容与SSH返回码，还包含提交人、运营从PHP处生成的执行内容、运维检查脚本返回的检查内容、审批人等关键信息

## 联系我
 如果在二开/学习过程中有疑问，可以联系邮箱：fqh1113464192@gmail.com
 因为这也是我初次从0手搓运维平台，如果有意见或者批评也可以联系我，看到就会回。
 辛苦大家看到这里了！！再次感谢！
 

