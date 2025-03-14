# Webook

Webook小微书（仿小红书）

- DDD 框架：Domin-Drive Design

    ![image-20241226202631481](./img/image-20241226202631481.png)

项目启动：
- 前端：在 webook-fe 目录下，执行 `npm run dev`
- 后端：在 webook 目录下，执行 `go run . --config config/dev.yaml`
  - 配置文件：config/dev.yaml
  - `go run .` ：在当前目录下运行，包含 wire 生成的代码
- 第三方依赖：在 webook 目录下，执行 `docker compose up`
  - 执行 `docker compose down` 会删除数据库，结束 `docker compose up` 进程不会
  - 包含 mysql ，redis，viper，etcd
- 依赖注入：在 Webook 目录下，执行 `wire`
- mock: 在 Webook 目录下，执行 `make mock`

功能详细介绍见 **[博客: webook](https://docs.selfknow.cn/projects/golang/webook/)**