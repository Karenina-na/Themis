# Themis: 分布式记账调度中心

<img src="icon/png/logo-white.png" alt="Themis分布式记账调度中心" width="500" />

### [Author](https://www.weizixiang.top)

---

## API
### MessageAPI  通信api
* POST /v1/message/leader/register  服务注册
* PUT /v1/message/leader/beat  服务心跳
* PUT /v1/message/leader/election  新一轮选举发起，以Leader身份发起
* 
* POST /v1/message/follow/getLeader  获取当前Leader，以非Leader身份发起
* POST /v1/message/follow/getServers  获取Leader领导的服务集群，以Leader身份发起
* POST /v1/message/follow/getServersNum	获取当前集群服务数量

### OperatorAPI  操作api
* GET /v1/operator/CURD/getNamespaces 获取全部命名空间名称
* GET /v1/operator/CURD/getColonies 获取指定命名空间全部集群名称
* GET /v1/operator/CURD/getColoniesInstances 获取指定命名空间全部集群名称和服务器名称
* GET /v1/operator/CURD/getInstances 获取全部服务实例
* POST /v1/operator/CURD/getInstance 获取指定服务实例
* POST /v1/operator/CURD/getInstancesByCondition 获取指定命名空间和区域的服务实例--返回list  模糊查询
* POST /v1/operator/CURD/getInstancesByConditionMap 获取指定命名空间和区域的服务实例--返回map  精确查询
*
* DELETE /v1/operator/CURD/blacklistInstance 将指定服务实例拉入黑名单
* DELETE /v1/operator/CURD/blacklistColony 将指定集群内所有服务实例拉入黑名单
* DELETE /v1/operator/CURD/blacklistNamespace 将指定命名空间内所有集群内所有服务实例拉入黑名单
*
* GET /v1/operator/CURD/getBlacklist 获取黑名单内的实例
* DELETE /v1/operator/CURD/deleteBlacklistInstance 将服务从黑名单删除

* GET /v1/operator/cluster/getStatus 获取调度中心服务状态
* GET /v1/operator/cluster/getClusterLeader 获取当前集群Leader
* GET /v1/operator/cluster/getClusterStatus 获取当前集群服务身份

* GET /v1/operator/manager/login 管理员登录

### 具体API请运行项目并访问http://localhost:8080/swagger/index.html   查看swagger文档
```
swag init -g swag.go -d .\src\swag\ #创建swag文档
```

## 项目结构
```     
Themis:.
├─bin   //可执行文件
├─conf  //项目配置，需与bin中的exe文件放在同一个目录下
├─db    //sqllit数据库文件
├─docs  //swagger文档
├─log   //日志
├─src   //源码
│  ├─config         //配置加载器
│  ├─controller     //前端控制器
│  │  ├─interception    //拦截器
│  │  └─util    //工具
│  ├─entity         //数据模型
│  ├─exception      //异常处理
│  ├─factory        //服务加载器
│  │  └─image       //图片生成器     
│  ├─mapper         //数据持久层
│  ├─pool           //线程池
│  ├─router         //后端路由
│  ├─service        //业务逻辑层层
│  │  ├─Bean        //逻辑状态Bean
│  │  └─LeaderAlgorithm  //选举算法
│  ├─swag           //swagger文档
│  │  ├─controller  //swagger控制器
│  │  ├─entity    //swagger数据模型
│  │  └─syncBean    //swagger数据模型
│  ├─sync           //集群同步包
│  │  ├─candidate   //candidate状态
│  │  ├─common      //公共代码
│  │  ├─follow      //follow状态
│  │  ├─leader      //leader状态  
│  │  └─syncBean    //同步数据模型
│  └─util           //工具类  
│     ├─token       //token工具  
│     └─encryption  //加密控制器
└─test              //测试
    ├─Base          //测试初始化加载器
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
mysql : v8.0.13
gowatch : v1.5.2
jwt-go : v3.2.0+incompatible
```