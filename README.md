# Themis: 分布式记账调度中心


### [Author](https://www.wzxaugenstern.online/#/)

---

## API
### MessageAPI  通信api
* POST /api/message/register  服务注册  
` {
  "Port":"服务端口",
  "IP": "服务IP",
  "Name":"服务名",
  "Time":"服务启动时间",
  "Colony":"服务集群",
  "Namespace":"命名空间"
  }`


* PUT /api/message/beat  服务心跳  
  ` {
  "Port":"服务端口",
  "IP": "服务IP",
  "Name":"服务名",
  "Time":"服务启动时间",
  "Colony":"服务集群",
  "Namespace":"命名空间"
  }`


* PUT /api/message/election  新一轮选举发起，以Leader身份发起    
  ` {
  "Port":"Leader服务端口",
  "IP": "Leader服务IP",
  "Name":"Leader服务IP服务名",
  "Time":"Leader服务IP服务启动时间",
  "Colony":"Leader服务IP服务集群",
  "Namespace":"Leader服务IP命名空间"
  }`


* GET /api/message/getLeader  获取当前Leader，以非Leader身份发起  


* POST /api/message/getServers  获取Leader领导的服务集群，以Leader身份发起  
  ` {
  "Port":"Leader服务端口",
  "IP": "Leader服务IP",
  "Name":"Leader服务IP服务名",
  "Time":"Leader服务IP服务启动时间",
  "Colony":"Leader服务IP服务集群",
  "Namespace":"Leader服务IP命名空间"
  }`

### OperatorAPI  操作api


* GET /api/operator/getInstances 获取全部服务实例  


* DELETE /api/operator/deleteInstance 将指定服务拉入黑名单  
  ` {
  "Port":"服务端口",
  "IP": "服务IP",
  "Name":"服务名",
  "Time":"服务启动时间",
  "Colony":"服务集群",
  "Namespace":"命名空间"
  }`


* DELETE /api/operator/deleteColony 将指定集群内所有服务拉入黑名单  
  ` {
  "Port":"服务端口",
  "IP": "服务IP",
  "Name":"服务名",
  "Time":"服务启动时间",
  "Colony":"服务集群",
  "Namespace":"命名空间"
  }`


* GET /api/operator/getDeleteInstances 获取黑名单内的实例  


* DELETE /api/operator/cancelDeleteInstance 将服务从黑名单删除  
  ` {
  "Port":"服务端口",
  "IP": "服务IP",
  "Name":"服务名",
  "Time":"服务启动时间",
  "Colony":"服务集群",
  "Namespace":"命名空间"
  }`


## 项目结构
```
Themis:.
│  .gitignore
│  go.mod
│  go.sum
│  README.md
│
└─src                               源码
    │  main.go
    │
    ├─config                        配置文件
    │      Config.go
    │
    ├─controller                    前端控制器
    │      Interception.go
    │      ManagerController.go
    │      ServerController.go
    │
    ├─entity                        数据模型
    │  │  ResultModel.go
    │  │  ServerModel.go
    │  │
    │  └─util
    │          GoroutinePool.go
    │          LinkList.go
    │          Logger.go
    │          Strval.go
    │
    ├─router                        后端路由
    │      MessageAPI.go
    │      OperatorAPI.go
    │
    └─service                       业务层
        │  init.go
        │  MessageService.go
        │  OperatorService.go
        │
        └─LeaderAlgorithm           选举算法
                init.go
                RandomAlgorithmCreateLeader.go
```
