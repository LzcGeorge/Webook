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

### 注册功能

1. Bind 绑定请求参数，绑定到结构体 UserSignUpReq
2. 用正则表达式校验邮箱和密码格式
3. 确认密码和密码一致
4. 调用 service 层进行注册
5. 返回注册成功

> 跨域请求：
>
> 项目是前后端分离的，前端是 Axios，后端是Go，所以需要跨域请求。
>
> - 跨域请求：协议、域名、端口有一个不同，就叫跨域
> - Request Header 和 Response Header 中的字段要对应上
> - 采用 middleware 中间件进行跨域请求
>
> docker compose 安装数据库
>
> - 静默启动；
>
>     ```bash
>      docker compose up -d
>     ```
>
> - `docker compose up` 初始化 docker compose 并启动
>
> - `docker compose down` 删除 docker compose 里面创建的各种容器，数据库
>
> - 只要不 down 数据库一直都在
>
> DDD 框架：Domin-Drive Design
>
> - Domain: 领域，存储对象
> - Repository: 数据存储
> - Service: 业务逻辑

### 登录功能

登录功能分为两件事：
- 实现登录功能
- 登录状态的校验

登录功能：

1. 绑定请求参数，绑定到结构体 UserLoginReq
2. 在 service 层中，根据邮箱查询用户是否存在，密码是否正确
3. 返回登录结果

登录状态的校验：
- 利用 Gin 的 session 插件，从 cookie 中获取 sessionID，校验登录状态
- 采用 Cookie 和 Session 进行登录状态的保持

> Cookie：
> - Domain：Cookie 可以在什么域名下使用
> - Path：Cookie 可以在什么路径下使用
> - Expires/Max-Age：Cookie 的过期时间
> - HttpOnly：Cookie 是否可以通过 JS 访问
> - Secure：Cookie 是否只能通过 HTTPS 访问
> - SameSite：Cookie 是否只能在同一个站点下使用

> Session：
> - 存储在服务器端
> - 通过 SessionID 来识别用户
> - 一般通过 Cookie 来传递 SessionID
>
> Redis：
> - 用户数据存储在 Redis 中
>
> LoginMiddlewareBuilder：
> - 登录中间件，用于校验登录状态
> - 通过 IgnorePaths 方法，设置不校验登录状态的路径
> - 通过 Build 方法，构建中间件: 链式调用
>
> Debug 定位问题：
> 倒排确定：http 发送请求，中间件，业务逻辑，数据库
> F12 查看错误信息
> 后端看日志
>
> Session 的过期时间：
> - 通过中间件 LoginMiddlewareBuilder 设置，当访问不在 IgnorePaths 的路径时，会更新 Session 的 update_time 字段
> - 同时更新 Session 的过期时间 MaxAge
> - 但每次访问都要从 Redis 中获取 Session，性能较差（所以后面引入 JWT）
> 
