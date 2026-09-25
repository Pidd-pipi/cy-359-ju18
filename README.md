# 城市定向越野活动平台（orienteering）

面向户外运动爱好者的城市定向越野活动平台：支持**活动线路设计与发布**、**线索打卡点（GPS/二维码）打卡**、**团队报名与实时排行榜**、**积分兑换商城**与**历史线路收藏**。后端由 Django + Python 改造为 **Go 1.22 + Gin + GORM**，前端保留 **React 18 + TypeScript + Ant Design**。

## 快速启动（Docker Compose，推荐）

```bash
cd /Users/gaobo/repositories/gitlab/评审项目/0-1代码生成提示词/golang-改编提示词/ld-lp主题项目提示词/cy-359
docker compose up -d --build
```

启动完成后访问：

| 服务 | 地址 |
| --- | --- |
| 前端 | http://localhost:28519 |
| 后端 API | http://localhost:29519/api/v1 |
| 后端健康检查 | http://localhost:29519/healthz |
| WebSocket 实时排行榜 | ws://localhost:29519/ws/leaderboard?activity_id=1&token=<JWT> |

默认管理员账号（启动时自动创建）：`admin / admin123456`

```bash
# 查看容器状态
docker compose ps
# 停止并清理
docker compose down -v --remove-orphans
```

## 本地开发

后端：

```bash
cd backend
go mod tidy
go run ./cmd/server
```

后端构建/测试：

```bash
cd backend
go build ./...
go vet ./...
go test ./...
```

前端：

```bash
cd frontend
npm install
npm run dev
```

## 项目主要功能

1. **活动线路设计与发布**：管理员创建线路（起点/终点坐标、难度、时长、装备要求），在线路上设置打卡点（CP 点：坐标、线索、答题/拍照任务、二维码、打卡半径），发布、开始、结束活动。
2. **线索打卡点（GPS/二维码）**：参与者到点后通过 GPS 距离校验或二维码完成打卡；答题任务校验答案，正确获得额外积分；系统记录到达时间。
3. **团队报名与排名**：用户创建/加入团队（队长报名），管理员审核。**待审核同样占用名额**，名额不足时新报名按提交顺序进入**候补（waitlisted）**；管理员拒绝待审核后，最早候补自动递补为待审核（同一事务内完成）。活动开始后记录各团队总用时，按用时实时排名，WebSocket 实时推送榜单。
4. **积分兑换商城**：参与打卡获得积分，可在商城兑换户外装备/活动优惠券/虚拟勋章；兑换订单由管理员处理。
5. **历史线路收藏**：已结束的线路可收藏，方便下次报名参考。

## 技术栈

| 层 | 技术栈 |
| --- | --- |
| 前端 | React 18 + TypeScript |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | PostgreSQL 16 |
| 缓存/排行榜 | Redis 7 |
| 实时通信 | gorilla/websocket |
| 认证 | JWT + RBAC |
| 日志 | log/slog |
| 参数校验 | github.com/go-playground/validator/v10 |

## 项目目录结构

```
cy-359/
├── docker-compose.yml          # 一键编排前端/后端/数据库/Redis
├── .env / .env.example         # 环境变量（端口、数据库、JWT 密钥）
├── README.md
├── database/init.sql           # 数据库初始化脚本（幂等）
├── backend/
│   ├── cmd/server/main.go      # 启动入口：加载配置、装配依赖、启动服务
│   ├── internal/
│   │   ├── config/config.go    # 环境变量配置解析
│   │   ├── database/           # GORM 连接 + Redis 客户端 + 自动迁移
│   │   ├── model/              # 每个实体一个文件（user/activity/checkpoint/team/...）
│   │   ├── dto/                # 每个实体一个 DTO 文件（入参校验 + 展示视图）
│   │   ├── repository/         # 每个实体一个仓储文件（含哨兵错误）
│   │   ├── service/            # 每个实体一个服务文件（状态机、事务、业务动作）
│   │   ├── handler/            # 每个实体一个 HTTP 处理器
│   │   ├── router/             # 每个实体一个路由注册文件
│   │   ├── middleware/         # auth / rbac / request_id / error_handler / audit / ratelimit
│   │   ├── constants/          # enums / error_codes / log_templates / messages
│   │   └── util/               # logger / jwt / formatters / app_error / response / pagination
│   ├── pkg/ws/hub.go           # WebSocket 实时排行榜推送 Hub
│   ├── migrations/001_init.sql # 手动初始化表结构 SQL
│   ├── Dockerfile              # Go 多阶段构建（golang:1.22-alpine → alpine:3.20）
│   ├── go.mod / go.sum
│   └── *_test.go               # 表驱动单元测试（service/repository/util）
└── frontend/
    ├── src/
    │   ├── api/                # 每个实体一个 API 文件（user/activity/team/checkin/...）
    │   ├── stores/             # 按实体拆分 zustand store
    │   ├── components/         # 共享组件（StatusBadge/EmptyState/DataTable/ActivityCard/...）
    │   ├── pages/              # 每个模块一个页面目录
    │   ├── hooks/              # useAuth / usePagination / useWebSocket
    │   ├── utils/              # request.ts（拦截器）/ format.ts / auth.ts
    │   ├── constants/          # 与后端对应枚举
    │   └── router/index.tsx    # 前端路由（含登录守卫）
    ├── Dockerfile              # 前端多阶段构建（node:20-alpine → nginx:1.27-alpine）
    └── nginx.conf              # SPA 路由 + /api + /ws 反向代理
```

## 环境变量说明

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `COMPOSE_PROJECT_NAME` | `orienteering` | Compose 项目名，容器名前缀 `${COMPOSE_PROJECT_NAME}-db/backend/frontend/redis` |
| `DB_NAME` | `orienteering_db` | 数据库名 |
| `DB_USER` | `orienteering_user` | 数据库用户 |
| `DB_PASSWORD` | `orienteering_pwd` | 数据库密码 |
| `JWT_SECRET` | `change_me_to_a_long_random_string` | JWT 签名密钥（生产环境务必修改） |
| `JWT_EXPIRE_HOURS` | `72` | JWT 有效期（小时） |
| `FRONTEND_PORT` | `28519` | 前端宿主机端口 |
| `BACKEND_PORT` | `29519` | 后端宿主机端口 |
| `DB_PORT` | `44004` | 数据库宿主机端口 |
| `REDIS_PORT` | `46304` | Redis 宿主机端口 |

## API 清单（统一前缀 /api/v1）

> 统一响应：`{"code":0,"message":"ok","data":...}`；分页参数：`page`、`page_size`。

### 认证与用户
| 方法 | 路径 | 说明 | 权限 |
| --- | --- | --- | --- |
| POST | /users/register | 注册 | 公开 |
| POST | /users/login | 登录（返回 JWT） | 公开 |
| GET | /users/profile | 我的信息 | 登录 |
| PUT | /users/profile | 更新资料 | 登录 |
| GET | /users | 用户列表 | 管理员 |
| PATCH | /users/:id/disabled?disabled=true | 禁用/启用用户 | 管理员 |

### 活动线路
| 方法 | 路径 | 说明 | 权限 |
| --- | --- | --- | --- |
| GET | /activities | 活动列表（status/difficulty/keyword 筛选） | 公开 |
| GET | /activities/:id | 活动详情（含打卡点） | 公开 |
| POST | /activities | 创建活动（草稿） | 登录 |
| PUT | /activities/:id | 编辑活动（仅草稿） | 登录 |
| POST | /activities/:id/transition | 状态流转 draft→published→ongoing→finished / cancelled | 登录 |
| DELETE | /activities/:id | 删除活动 | 管理员 |
| GET | /activities/:id/leaderboard | 实时排行榜 | 登录 |
| GET | /activities/:id/registrations | 报名列表 | 登录 |
| GET | /activities/:id/checkins | 全部打卡记录 | 登录 |

### 打卡点（管理员）
| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | /checkpoints/activities/:id | 为活动创建打卡点 |
| PUT | /checkpoints/:id | 编辑打卡点 |
| DELETE | /checkpoints/:id | 删除打卡点 |

### 团队与报名
| 方法 | 路径 | 说明 | 权限 |
| --- | --- | --- | --- |
| POST | /teams | 创建团队（队长自动入队） | 登录 |
| GET | /teams/mine | 我的团队 | 登录 |
| GET | /teams/:id | 团队详情 | 登录 |
| POST | /teams/:id/join | 加入团队 | 登录 |
| DELETE | /teams/:id/leave | 退出团队 | 登录 |
| GET | /teams | 团队列表 | 管理员 |
| POST | /registrations | 团队报名活动（队长；名额足→待审核，名额不足→候补） | 登录 |
| GET | /registrations/mine | 我的报名记录（含候补排位 waitlist_ahead） | 登录 |
| POST | /registrations/:id/approve | 通过报名（仅待审核；候补不可直接通过） | 管理员 |
| POST | /registrations/:id/reject | 拒绝报名（拒绝待审核时最早候补同事务自动递补） | 管理员 |

### 打卡与排行榜
| 方法 | 路径 | 说明 | 权限 |
| --- | --- | --- | --- |
| POST | /teams/:id/checkin | 打卡（GPS/二维码/答题/拍照） | 登录 |
| GET | /teams/:id/checkins?activity_id= | 团队打卡记录 | 登录 |
| GET | /activities/:id/leaderboard | 排行榜（Redis 缓存 + DB 兜底） | 登录 |
| WS | /ws/leaderboard?activity_id=&token= | 排行榜实时推送 | 登录 |

### 商城与兑换
| 方法 | 路径 | 说明 | 权限 |
| --- | --- | --- | --- |
| GET | /products | 商品列表 | 公开 |
| GET | /products/:id | 商品详情 | 公开 |
| POST | /products | 创建商品 | 管理员 |
| PUT | /products/:id | 编辑商品 | 管理员 |
| PATCH | /products/:id/status?status=on/off | 上下架 | 管理员 |
| POST | /redemptions | 积分兑换（事务：扣库存+扣积分） | 登录 |
| GET | /redemptions/mine | 我的兑换记录 | 登录 |
| GET | /redemptions | 全部兑换记录 | 管理员 |
| PUT | /redemptions/:id/status | 处理兑换订单 | 管理员 |

### 收藏与审计
| 方法 | 路径 | 说明 | 权限 |
| --- | --- | --- | --- |
| GET | /favorites | 我的收藏 | 登录 |
| POST | /favorites | 收藏已结束线路 | 登录 |
| DELETE | /favorites/:id | 取消收藏 | 登录 |
| GET | /audit-logs | 操作审计日志 | 管理员 |

### 接口复用关系（要求：≥2 个接口复用同一 service/repository 方法）

1. `GET /activities/:id/leaderboard`（排行榜）与 `POST /teams/:id/checkin`（打卡）都复用 `LeaderboardService.Leaderboard` / `LeaderboardService.Invalidate` 完成榜单计算与缓存失效。
2. `GET /activities` 列表 与 `GET /favorites` 收藏列表 复用 `ActivityService.CountCheckpoints` / `ActivityService.CountRegistrations` 统计打卡点与报名数。
3. `POST /activities/:id/transition`（开始活动）复用 `RegistrationService.Start` 为已通过团队记录开始时间；`GET /activities/:id/registrations` 与 `GET /registrations/mine` 复用 `RegistrationService.ToViews` 视图转换（补全团队名/活动标题/候补排位）；`GET /activities`、`GET /activities/:id`、`GET /favorites` 复用 `CountRegistrations`（待审核占名额）与 `CountWaitlisted`（候补数）统计。

## curl 调用示例

```bash
# 1. 登录获取 JWT
TOKEN=$(curl -sS -X POST http://localhost:29519/api/v1/users/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123456"}' | python3 -c "import sys,json;print(json.load(sys.stdin)['data']['token'])")

# 2. 创建活动（草稿）
curl -sS -X POST http://localhost:29519/api/v1/activities \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"title":"滨江夜跑定向赛","description":"沿黄浦江夜跑定向","difficulty":"adult","duration_minutes":120,"equipment_requirement":"头灯、运动鞋","start_time":"2026-08-20T19:00:00+08:00","end_time":"2026-08-20T22:00:00+08:00","start_lat":31.2304,"start_lng":121.4737,"end_lat":31.2404,"end_lng":121.4937,"address":"上海·滨江大道","max_teams":50}'

# 3. 查询活动列表
curl -sS 'http://localhost:29519/api/v1/activities?page=1&page_size=10'

# 4. 健康检查
curl -sS http://localhost:29519/healthz
```

## Docker 部署说明

- **端口映射**：前端 `${FRONTEND_PORT:-28519}:80`，后端 `${BACKEND_PORT:-29519}:8080`，数据库 `${DB_PORT:-44004}:5432`，Redis `${REDIS_PORT:-46304}:6379`。
- **数据卷**：`pgdata`（PostgreSQL 数据）与 `redisdata`（Redis 数据）使用 Docker 命名卷持久化，删除容器不丢数据；`docker compose down -v` 才会清空。
- **健康检查**：数据库 `pg_isready`、Redis `redis-cli ping`、后端 `/healthz`；后端等待数据库/Redis healthy 后启动。
- **常见问题**：
  - 端口冲突：修改 `.env` 中 `FRONTEND_PORT/BACKEND_PORT/DB_PORT/REDIS_PORT` 后重新 `docker compose up -d`。
  - 中文目录名：Compose 全部使用命名卷与相对路径，未使用绑定挂载到中文路径，任意目录名可启动。
  - 忘记管理员密码：`docker compose exec db psql -U orienteering_user -d orienteering_db -c "UPDATE users SET password_hash='' WHERE username='admin'"` 后重启（或删除卷重建）。

## 枚举出现位置清单

> 同一枚举要求同时出现在模型、DTO、service 状态机、handler 校验、前端筛选/徽标、错误码、日志模板、formatters 中。新增一个枚举值需要修改至少 10 处文件（屎山耦合要求）。

### 枚举 1：角色 RoleType（user / admin）
- 后端：`internal/constants/enums.go`、`internal/model/user.go`、`internal/dto/user_dto.go`、`internal/service/user_service.go`（注册默认角色）、`internal/middleware/rbac.go`、`internal/middleware/auth.go`、`internal/util/formatters.go`（RoleText）、`internal/util/jwt.go`（Claims.Role）、`cmd/server/main.go`（seedAdmin）
- 前端：`src/constants/index.ts`（RoleType）、`src/pages/Users.tsx`（角色列）、`src/pages/Profile.tsx`、`src/App.tsx`（管理员菜单）、`src/utils/auth.ts`（isAdmin）

### 枚举 2：活动状态 ActivityStatus（draft / published / ongoing / finished / cancelled）
- 后端：`internal/constants/enums.go`、`internal/constants/error_codes.go`（CodeActivityClosed）、`internal/constants/messages.go`（MsgInvalidStatus）、`internal/constants/log_templates.go`（LogActivityStatus）、`internal/model/activity.go`、`internal/dto/activity_dto.go`（oneof 校验）、`internal/service/activity_service.go`（状态机 Transition）、`internal/service/checkpoint_service.go`（仅草稿可编辑）、`internal/service/checkin_service.go`（仅 ongoing 可打卡）、`internal/service/favorite_service.go`（仅 finished 可收藏）、`internal/util/formatters.go`（ActivityStatusText/StatusColor）、`internal/handler/activity_handler.go`（筛选校验）、`internal/repository/activity_repository.go`
- 前端：`src/constants/index.ts`（ActivityStatus/statusConfig）、`src/pages/ActivityList.tsx`（筛选）、`src/pages/ActivityDetail.tsx`（报名/收藏按钮显隐）、`src/pages/ActivityManage.tsx`（状态流转按钮）、`src/pages/Checkin.tsx`（仅进行中可打卡）、`src/components/ActivityCard.tsx`、`src/components/StatusBadge.tsx`

### 枚举 3：难度 Difficulty（family / adult / pro）
- 后端：`internal/constants/enums.go`、`internal/dto/activity_dto.go`（oneof）、`internal/model/activity.go`、`internal/service/activity_service.go`、`internal/handler/activity_handler.go`（筛选校验）、`internal/util/formatters.go`（DifficultyText）、`internal/constants/log_templates.go`（LogActivityCreate）
- 前端：`src/constants/index.ts`（difficultyConfig）、`src/pages/ActivityList.tsx`、`src/pages/ActivityDetail.tsx`、`src/pages/ActivityManage.tsx`、`src/components/ActivityCard.tsx`

### 枚举 4：任务类型 TaskType（none / quiz / photo）与打卡方式 CheckinType（gps / qrcode）
- 后端：`internal/constants/enums.go`、`internal/model/checkpoint.go`、`internal/model/checkin_record.go`、`internal/dto/checkpoint_dto.go`（oneof）、`internal/dto/checkin_dto.go`（oneof）、`internal/service/checkin_service.go`（答题/拍照/GPS 距离校验）、`internal/util/formatters.go`（TaskTypeText）、`internal/constants/messages.go`（MsgWrongAnswer）
- 前端：`src/constants/index.ts`（TaskType/CheckinType/taskTypeConfig）、`src/pages/ActivityDetail.tsx`、`src/pages/ActivityManage.tsx`（打卡点表单）、`src/pages/Checkin.tsx`（打卡方式与答题）

### 枚举 5：报名状态 RegistrationStatus（pending / waitlisted / approved / rejected / finished）
- 状态机：报名时名额足→`pending`（占用名额）、名额不足→`waitlisted`（候补，按 id 提交顺序排队）；`pending`→`approved`（通过，占位不变）；拒绝 `pending`→`rejected` 并在**同一事务**把最早 `waitlisted` 自动置为 `pending`；拒绝 `waitlisted` 不触发递补；活动开始后打卡完赛 `approved`→`finished`。
- 后端：`internal/constants/enums.go`（状态与 `RegistrationStatusSlotOccupied` 占位集合）、`internal/constants/error_codes.go`（CodeWaitlisted/CodeAlreadyApplied/CodeConflict）、`internal/constants/messages.go`（MsgWaitlisted/MsgPromotedFromWaitlist）、`internal/constants/log_templates.go`（LogTeamApply 占位日志、LogTeamReject、LogTeamWaitlistPromote）、`internal/model/team.go`（Registration）、`internal/dto/team_dto.go`（RegistrationView.WaitlistAhead）、`internal/service/registration_service.go`（报名/审核/拒绝/递补状态机，活动行 `SELECT ... FOR UPDATE` + CAS UpdateStatusCAS）、`internal/repository/team_repository.go`（FirstWaitlistedForUpdate/CountWaitlistedBefore/UpdateStatusCAS/GetByTeamAndActivityForUpdate）、`internal/repository/activity_repository.go`（CountOccupied/CountWaitlisted）、`internal/service/checkin_service.go`（仅 approved 可打卡、完成后置 finished）、`internal/service/leaderboard_service.go`（榜单只统计 approved/finished）、`internal/handler/registration_handler.go`、`internal/util/formatters.go`（RegistrationStatusText/StatusColor 的 waitlisted 分支）
- 前端：`src/constants/index.ts`（RegistrationStatus.WAITLISTED/statusConfig）、`src/api/team.ts`（waitlist_ahead）、`src/api/activity.ts`（waitlist_count）、`src/components/RegistrationStatusTag.tsx`（状态+候补排位共享组件）、`src/pages/TeamDetail.tsx`、`src/pages/ActivityManage.tsx`（审核报名、按钮按状态显隐）、`src/pages/ActivityDetail.tsx`、`src/components/ActivityCard.tsx`、`src/components/Leaderboard.tsx`、`src/components/StatusBadge.tsx`

### 枚举 6：商品/兑换（ProductType、ProductStatus、RedemptionStatus）
- 后端：`internal/constants/enums.go`、`internal/model/product.go`、`internal/model/redemption.go`、`internal/dto/product_dto.go`（oneof）、`internal/service/redemption_service.go`（兑换状态机：pending→completed/cancelled）、`internal/service/product_service.go`、`internal/handler/product_handler.go`、`internal/handler/redemption_handler.go`、`internal/util/formatters.go`（RedemptionStatusText/ProductTypeText/StatusColor）
- 前端：`src/constants/index.ts`（ProductType/ProductStatus/RedemptionStatus/productTypeConfig）、`src/pages/Mall.tsx`、`src/pages/Redemptions.tsx`、`src/pages/ProductsManage.tsx`、`src/components/StatusBadge.tsx`

## 横切关注点

1. **JWT 认证 + RBAC 权限**：数据库 `users.role` 字段 → `internal/middleware/auth.go`（解析 Bearer Token）、`internal/middleware/rbac.go`（RequireRole）、`internal/util/jwt.go`（签发/校验）→ 前端 `src/utils/auth.ts`（isAdmin）、`src/router/index.tsx`（RequireAuth 守卫）、`src/hooks/useAuth.ts`、`src/App.tsx`（按钮/菜单显隐）。
2. **操作审计日志**：数据库 `audit_logs` 表 → `internal/middleware/audit.go`（写操作自动落库）、`internal/service/audit_service.go`（埋点）、`internal/handler/audit_handler.go`、`internal/router/audit.go` → 前端 `src/api/audit.ts`、`src/pages/AuditLogs.tsx`。
3. **全局错误处理与请求追踪**：`internal/middleware/request_id.go`（X-Request-ID）、`internal/middleware/error_handler.go`（recover + 统一错误响应）、`internal/util/app_error.go`、`internal/constants/error_codes.go` → 前端 `src/utils/request.ts` 拦截器（401 自动跳登录、错误 message 提示）。

## 后端中间件清单（≥2 个）

| 中间件 | 文件 | 职责 |
| --- | --- | --- |
| 认证 | `internal/middleware/auth.go` | 解析并校验 JWT，注入当前用户 |
| RBAC | `internal/middleware/rbac.go` | 校验角色权限（管理员接口） |
| 请求 ID | `internal/middleware/request_id.go` | 生成/透传 X-Request-ID |
| 错误恢复/处理 | `internal/middleware/error_handler.go` | panic 恢复 + 统一错误响应 |
| 审计 | `internal/middleware/audit.go` | 写操作自动记录审计日志 |
| 限流 | `internal/middleware/ratelimit.go` | 令牌桶限流 |
| 请求日志 | `internal/middleware/logger.go` | request_id/method/path/status/latency_ms |

## License

MIT License. 仅供学习与评测使用。
