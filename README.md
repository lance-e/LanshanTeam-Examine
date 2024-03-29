# LanshanTeam winter examine

##### 基础功能：

* [X]  用户注册
* [X]  用户登录
* [X]  用户个人主页
* [X]  添加游戏好友
* [X]  大厅创建房间
* [X]  大厅加入房间，多位玩家准备后开始游戏
* [X]  基本游戏逻辑（任意选择一款游戏进行实现）：五子棋、象棋、坦克大战、贪吃蛇等等
* [X]  积分排行榜/段位

##### 加分项：

* [X]  用户密码加盐加密
* [X]  用户登录有短信登录、邮箱登录、第三方登录多种形式
* [X]  验证码（登录，注册，修改密码）
* [X]  用户状态保存使用 JWT 或 Session
* [X]  实现对局的断线重连
* [X]  实现对局内聊天
* [ ]  实现好友邀请进行对局
* [ ]  实现观战功能
* [X]  用户对局记录，要求能够复现每一步下法
* [ ]  将项目部署上线（包括前端和后端的项目，也就是登录你的网站能够像正常的网站一样访问）
* [ ]  使用 https 加密
* [X]  合理的缓存策略（提高访问速度）
* [X]  考虑服务端安全性（xxs，sql注入，cors，csrf 等）
* [ ]  实现游戏性能的优化：容纳更多人同时在线、降低对局延迟等等，如果可以做测试进行优化前后对比就更好了
* [ ]  其他任何你想加的让我们耳目一新的功能...

### tech stack

web framework:gin

rpc framework:grpc

relational database:mysql

none-relational database:redis

service register and discovery:etcd

deployment:docker

configuration:viper

logger:zap

### structure

![structures.png](assets/structures.png)

~~~bash
LanshanTeam-Examine
├── LICENSE
├── README.md
├── assets
│   └── structures.png
├── caller
│   ├── api
│   │   ├── middleware
│   │   │   ├── cors.go
│   │   │   └── jwt.go
│   │   └── router
│   │       └── route.go
│   ├── config
│   │   └── config.yaml
│   ├── handle
│   │   ├── addFriend.go
│   │   ├── game.go
│   │   ├── history.go
│   │   ├── homepage.go
│   │   ├── login.go
│   │   ├── rank.go
│   │   ├── register.go
│   │   ├── thirdPart.go
│   │   └── verify.go
│   ├── main.go
│   ├── model
│   │   └── model.go
│   ├── pkg
│   │   ├── consts
│   │   │   └── code.go
│   │   └── utils
│   │       └── logger.go
│   ├── rpc
│   │   ├── discovery
│   │   │   └── discovery.go
│   │   ├── gameModule
│   │   │   ├── game.go
│   │   │   └── pb
│   │   │       ├── game.pb.go
│   │   │       ├── game.proto
│   │   │       └── game_grpc.pb.go
│   │   └── userModule
│   │       ├── pb
│   │       │   ├── user.pb.go
│   │       │   ├── user.proto
│   │       │   └── user_grpc.pb.go
│   │       └── user.go
│   └── ws
│       └── serve.go
├── docker-compose.yaml
├── go.mod
├── go.sum
├── gobang-client
│   └── main.go
└── server
    ├── game
    │   ├── config
    │   │   └── gameConfig.yaml
    │   ├── dao
    │   │   ├── Init
    │   │   │   └── init.go
    │   │   ├── cathe
    │   │   │   └── model.go
    │   │   └── db
    │   │       ├── migrate.go
    │   │       └── model.go
    │   ├── handle
    │   │   └── handle.go
    │   ├── main.go
    │   ├── pb
    │   │   ├── game.pb.go
    │   │   ├── game.proto
    │   │   └── game_grpc.pb.go
    │   ├── serveRegister
    │   │   └── register.go
    │   └── utils
    │       └── logger.go
    └── user
        ├── config
        │   └── userConfig.yaml
        ├── dao
        │   ├── Init
        │   │   └── init.go
        │   ├── cathe
        │   │   └── models.go
        │   └── db
        │       ├── migrate.go
        │       └── models.go
        ├── handle
        │   ├── friendShip.go
        │   ├── homepage.go
        │   ├── registerAndLogin.go
        │   └── scoreAndRank.go
        ├── main.go
        ├── pb
        │   ├── user.pb.go
        │   ├── user.proto
        │   └── user_grpc.pb.go
        ├── pkg
        │   └── utils
        │       ├── bcrypt.go
        │       └── logger.go
        └── serveRegister
            └── register.go

~~~
