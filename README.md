# Themis: 分布式记账调度中心


### [Author](https://www.wzxaugenstern.online/#/)
#### [Themis](https://www.wzxaugenstern.online/#/Article?ArticleId=818427233)

---

## API
### MessageAPI  通信api
* POST /api/v1/message/register  服务注册
* PUT /api/v1/message/beat  服务心跳
* PUT /api/v1/message/election  新一轮选举发起，以Leader身份发起
* POST /api/v1/message/getLeader  获取当前Leader，以非Leader身份发起
* POST /api/v1/message/getServers  获取Leader领导的服务集群，以Leader身份发起
* POST /api/v1/message/getServersNum	获取当前集群服务数量

### OperatorAPI  操作api
* GET /api/v1/operator/getInstances 获取全部服务实例
* POST /api/v1/operator/getInstances 获取指定命名空间和区域的服务实例
* DELETE /api/v1/operator/deleteInstance 将指定服务拉入黑名单
* DELETE /api/v1/operator/deleteColony 将指定集群内所有服务拉入黑名单
* GET /api/v1/operator/getDeleteInstances 获取黑名单内的实例
* DELETE /api/v1/operator/cancelDeleteInstance 将服务从黑名单删除  
* GET /api/v1/operator/getStatus 获取调度中心服务状态

### 具体API请运行项目并访问http://localhost:8080/swagger/index.html   查看swagger文档

## 项目结构
```     
Themis:.
├─bin   //可执行文件
├─conf  //项目配置，需与bin中的exe文件放在同一个目录下
├─db    //数据库文件
├─docs  //swagger文档
├─log   //日志
├─src   //源码
│  ├─config         //配置加载器
│  ├─controller     //前端控制器
│  ├─entity         //数据模型
│  ├─exception      //异常处理模块
│  ├─Init           //初始化操作
│  ├─mapper         //数据持久层
│  ├─router         //路由
│  ├─service        //业务逻辑层
│  │  └─LeaderAlgorithm  //选举算法
│  └─util           //工具类
└─test              //测试
    ├─Base          //初始化测试
    └─ServerTest    //服务测试用例
```

### 运行
#### 启动项目
```
cd Themis

go build -o Themis.exe main.go      //windows
./Themis.exe                       //windows

go build -o Themis main.go          //linux
./Themis                            //linux

debug模式:    (windows为例)
./Themis.exe -mode=debug

release模式:
./Themis.exe -mode=release

test模式:
./Themis.exe -mode=test
```

### 版本
```
go version : 1.18.3 windows/amd64
gin : v1.8.1
swagger for go : v0.22.3
gorm : v1.3.6
sqllit : v1.3.6
gowatch : v1.5.2
go-callvis :  v0.6.1
```