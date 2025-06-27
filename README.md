# Webook

## 整体架构

项目采用典型的领域驱动设计(DDD)分层架构:

1. **表示层(Web层)** - `internal/web`
   - 负责HTTP请求处理和响应
   - 使用Gin框架作为Web框架
   - 包含各种Handler，如UserHandler、ArticleHandler等

2. **业务逻辑层(Service层)** - `internal/service`
   - 实现核心业务逻辑
   - 定义接口和实现类，如UserService、CodeService等
   - 不直接与数据库交互，通过Repository层

3. **数据访问层(Repository层)** - `internal/repository`
   - 定义数据访问接口
   - 包含缓存逻辑实现(cache目录)
   - 通过DAO层与数据库交互

4. **数据访问对象层(DAO层)** - `internal/repository/dao`
   - 直接操作数据库(如MongoDB、Redis等)

5. **领域模型层(Domain层)** - `internal/domain`
   - 定义核心业务实体，如User、Article等
   - 不包含业务逻辑，只包含数据结构

6. **定时任务层(Job层)** - `internal/job`
   - 处理定时和周期性任务
   - 使用cron库实现

### 技术栈

- **Web框架**: Gin
- **数据库**: MySQL 
- **缓存**: Redis
- **依赖注入**: Wire
- **配置管理**: Viper
- **日志**: Zap
- **定时任务**: cron
- **认证**: JWT, Session
- **容器化**: Docker

### 业务功能

1. **用户管理**:
   - 注册/登录(邮箱密码、手机验证码)
   - 个人资料管理
   - OAuth2集成(微信登录)

2. **文章系统**:
   - 文章创建、编辑、查询
   - 作者和读者视图

3. **互动功能**:
   - 点赞、收藏、评论等

4. **排行榜功能**:
   - 文章排名
   - 使用Redis缓存实现

### 架构特点

1. **依赖注入**: 使用Wire框架管理依赖，简化组件初始化
2. **缓存设计**: 多级缓存(本地缓存+Redis)
3. **接口分离**: 每层都定义接口，遵循依赖倒置原则
4. **领域驱动**: 基于DDD设计思想
5. **微服务准备**: 目录结构适合未来拆分为微服务
6. **测试**: 各层都有单元测试和mock实现

## 项目部署

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