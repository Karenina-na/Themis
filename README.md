# Themis: 分布式记账调度中心


### [Author](https://www.wzxaugenstern.online/#/)
#### [Themis](https://www.wzxaugenstern.online/#/Article?ArticleId=818427233)

---

## API
### MessageAPI  通信api
* POST /api/message/register  服务注册
* PUT /api/message/beat  服务心跳
* PUT /api/message/election  新一轮选举发起，以Leader身份发起
* GET /api/message/getLeader  获取当前Leader，以非Leader身份发起
* POST /api/message/getServers  获取Leader领导的服务集群，以Leader身份发起

### OperatorAPI  操作api
* GET /api/operator/getInstances 获取全部服务实例
* DELETE /api/operator/deleteInstance 将指定服务拉入黑名单
* DELETE /api/operator/deleteColony 将指定集群内所有服务拉入黑名单
* GET /api/operator/getDeleteInstances 获取黑名单内的实例
* DELETE /api/operator/cancelDeleteInstance 将服务从黑名单删除  

### 具体API请运行项目并访问http://localhost:8080/swagger/index.html   查看swagger文档

## 项目结构
```
Themis:.
├─.idea
├─bin   //可执行文件
├─conf  //项目配置，需与bin中的exe文件放在同一个目录下
├─docs  //swagger文档
├─log   //日志
└─src   //源码
    ├─config    //配置加载器
    ├─controller    //前端控制器
    ├─entity    //数据模型
    │  └─util
    ├─router    //路由
    └─service   //业务逻辑层
        └─LeaderAlgorithm   //选举算法

```

### 版本
```
go version : 1.18.3 windows/amd64
gin : v1.8.1
viper : v1.12.0
swagger for go : v0.22.3
```