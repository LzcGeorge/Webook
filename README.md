# Webook

## Webook小微书

- DDD 框架：Domin-Drive Design

    ![image-20241226202631481](./img/image-20241226202631481.png)

项目架构

```bash
.
├── config
│   ├── dev.go
│   ├── dev.yaml
│   ├── k8s.go
│   ├── redisExpire.go
│   ├── test.yaml
│   └── types.go
├── internal
│   ├── domain
│   ├── integration
│   ├── job
│   ├── repository
│   ├── service
│   └── web
├── ioc
│   ├── db.go
│   ├── job.go
│   ├── log.go
│   ├── redis.go
│   ├── sms.go
│   ├── web.go
│   └── wechat.go
├── main.go
├── pkg
│   ├── ginx
│   ├── limiter
│   └── logger
├── script
│   └── mysql
├── webook
├── Dockerfile
├── docker-compose.yaml
├── app.go
├── Makefile
├── wire.go
└── wire_gen.go
```

项目启动：
- 前端：在 webook-fe 目录下，执行 `npm run dev`
- 后端：在 webook 目录下，执行 `go run . --config config/dev.yaml`
  - 配置文件：config/dev.yaml
  - `go run .` ：在当前目录下运行，包含 wire 生成的代码
- 第三方依赖：在 webook 目录下，执行 `docker compose up`
  - 执行 `docker compose down` 会删除数据库，结束 `docker compose up` 进程不会
  - 包含 mysql ，redis，viper，etcd
- 依赖注入：在 webook 目录下，执行 `wire`
- mock: 在 Webook 目录下，执行 `make mock`



## 页面预览

### 主页

![image-20250315190925018](./img/image-20250315190925018.png)

![image-20250315191253565](./img/image-20250315191253565.png)

### 登录/注册

![image-20250315191117324](./img/image-20250315191117324.png)

![image-20250315191025972](./img/image-20250315191025972.png)

![image-20250315191055835](./img/image-20250315191055835.png)

### 文章广场

![image-20250315191441243](./img/image-20250315191441243.png)

### 热榜

![image-20250315191523379](./img/image-20250315191523379.png)

### 我的文章

![image-20250315191926233](./img/image-20250315191926233.png)

### 编辑文章

![image-20250315191702846](./img/image-20250315191702846.png)

### 个人资料

![image-20250315192020511](./img/image-20250315192020511.png)

## 功能介绍

功能详细介绍见 **[博客: webook](https://docs.selfknow.cn/projects/golang/webook/)**