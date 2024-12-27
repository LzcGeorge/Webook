# Webook

Webook小微书（仿小红书）

- DDD 框架：Domin-Drive Design

    ![image-20241226202631481](./assets/image-20241226202631481.png)

项目启动：
- 前端：在 webook-fe 目录下，执行 `npm run dev`
- 后端：在 webook 目录下，执行 `go run main.go`
- 数据库：在 webook 目录下，执行 `docker compose up`
  - 执行 `docker compose down` 会删除数据库，结束 `docker compose up` 进程不会

## 流程记录

### Week2

注册：

1. Bind 绑定请求参数，绑定到结构体 UserSignUpReq
2. 用正则表达式校验邮箱和密码格式
3. 确认密码和密码一致
4. 调用 service 层进行注册
5. 返回注册成功

跨域请求：

项目是前后端分离的，前端是 Axios，后端是Go，所以需要跨域请求。

- 跨域请求：协议、域名、端口有一个不同，就叫跨域
- Request Header 和 Response Header 中的字段要对应上

docker compose 安装数据库

- 静默启动；
   ```bash
    docker compose up -d
   ```

- `docker compose up` 初始化 docker compose 并启动
- `docker compose down` 删除 docker compose 里面创建的各种容器，数据库
- 只要不 down 数据库一直都在

DDD 框架：Domin-Drive Design

- Domain: 领域，存储对象
- Repository: 数据存储
- Service: 业务逻辑

