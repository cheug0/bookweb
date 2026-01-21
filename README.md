# BookWeb

BookWeb 是一个使用 Go 语言开发的现代化小说阅读 Web 程序，旨在替代杰奇 1.7 系统，提供更高的性能、更好的可维护性和更丰富的功能。

## ✨ 特性

- **高性能**：基于 Go 语言和 httprouter，提供出色的并发处理能力
- **Redis 缓存**：可选的 Redis 缓存支持，大幅提升页面响应速度
- **热重载**：支持路由配置热重载，无需重启即可生效
- **插件系统**：灵活的插件架构，支持广告管理、数据库优化、长尾词采集等
- **多模板支持**：支持多套前端模板，轻松切换站点风格
- **完整后台**：功能齐全的管理后台，包含文章、用户、配置等管理
- **用户系统**：支持用户注册、登录、书架、书签等功能
- **SEO 友好**：可配置的 URL 路由和 SEO 规则
- **GZIP 压缩**：可选的 GZIP 压缩，减少传输带宽
- **ID 转换**：支持 ID 算术转换，便于多站点共享数据库

## 📁 项目结构

```
bookweb/
├── admin/              # 后台管理模块
│   ├── admin.go        # 后台控制器
│   ├── auth.go         # 后台认证
│   └── template/       # 后台模板
├── config/             # 配置文件目录
│   ├── config.conf     # 主配置文件
│   ├── router.conf     # 路由配置
│   ├── seo.conf        # SEO 规则配置
│   ├── link.conf       # 友情链接配置
│   └── plugins.conf    # 插件配置
├── controller/         # 前台控制器
├── dao/                # 数据访问层
├── model/              # 数据模型
├── plugin/             # 插件目录
│   ├── ads/            # 广告管理插件
│   ├── db_optimizer/   # 数据库优化插件
│   └── langtail/       # 长尾词采集插件
├── router/             # 路由管理
├── service/            # 服务层
├── sql/                # SQL 脚本
├── static/             # 静态资源
├── template/           # 前台模板
├── utils/              # 工具函数
├── files/              # 小说文件存储
├── main.go             # 程序入口
└── go.mod              # Go 模块定义
```

## 🚀 快速开始

### 环境要求

- Go 1.24+
- MySQL 5.7+
- Redis（可选）

### 安装步骤

1. **克隆项目**
   ```bash
   git clone https://github.com/cheug0/bookweb.git
   cd bookweb
   ```

2. **初始化数据库**
   ```bash
   # 导入数据库结构和初始数据
   mysql -u root -p your_database < sql/data.sql
   mysql -u root -p your_database < sql/admin.sql
   ```

3. **修改配置文件**
   
   编辑 `config/config.conf`：
   ```json
   {
     "db": {
       "driver": "mysql",
       "host": "localhost",
       "port": 3306,
       "user": "your_user",
       "password": "your_password",
       "dbname": "your_database"
     },
     "server": {
       "host": "localhost",
       "port": 8080
     },
     "site": {
       "sitename": "我的小说站",
       "domain": "localhost:8080",
       "template": "html"
     }
   }
   ```

4. **编译运行**
   ```bash
   go build -o bookweb
   ./bookweb
   ```

5. **访问站点**
   - 前台：`http://localhost:8080`
   - 后台：`http://localhost:8080/admin`（默认账号：admin / admin123）

## ⚙️ 配置说明

### 主配置 (config.conf)

| 配置项 | 说明 |
|--------|------|
| `db` | 数据库连接配置 |
| `server` | HTTP 服务器配置 |
| `site.sitename` | 站点名称 |
| `site.domain` | 站点域名 |
| `site.template` | 使用的模板目录 |
| `site.admin_path` | 后台管理路径 |
| `site.force_domain` | 强制域名跳转 |
| `site.gzip_enabled` | 启用 GZIP 压缩 |
| `site.id_trans_rule` | ID 转换规则（如 `+1000`） |
| `redis` | Redis 缓存配置 |
| `storage` | 存储配置（local/oss） |

### 路由配置 (router.conf)

支持自定义 URL 模式：
```json
{
  "routes": {
    "book": "/book_:aid.html",
    "read": "/book/:aid/:cid.html",
    "sort": "/sort/:sid/:page/"
  }
}
```

### SEO 配置 (seo.conf)

可配置各页面的 Title、Keywords、Description 模板。

## 🔌 插件系统

### 内置插件

| 插件 | 说明 |
|------|------|
| `ads` | 广告位管理，支持自定义广告槽位 |
| `db_optimizer` | 数据库连接优化 |
| `langtail` | 长尾词采集支持 |

### 插件配置 (plugins.conf)

```json
{
  "ads": {
    "enabled": true,
    "slots": {
      "header": {
        "name": "顶部广告",
        "content": "<script>...</script>",
        "enabled": true
      }
    }
  }
}
```

## 🎨 模板开发

模板文件位于 `template/` 目录，支持多套模板切换。

### 模板变量

模板中可用的常用变量和函数请参考 `utils/template_funcs.go`。

### 切换模板

修改 `config.conf` 中的 `site.template` 为对应模板目录名。

## 📝 后台功能

- **仪表板**：站点统计概览
- **小说管理**：小说增删改查
- **用户管理**：用户列表、编辑、书架书签管理
- **友情链接**：链接管理
- **模块设置**：路由、分类、SEO 配置
- **插件管理**：插件启用/配置
- **系统设置**：站点、数据库、Redis、存储配置
- **安全设置**：修改密码、后台入口

## 🔧 开发说明

### 技术栈

- **Web 框架**：httprouter
- **数据库**：MySQL + go-sql-driver
- **缓存**：Redis (go-redis/v9)
- **模板**：Go html/template
- **其他**：UUID、GZIP 中间件

### 代码规范

- 使用 MVC 分层架构
- DAO 层使用预编译 SQL 语句
- 支持缓存自动失效

## 📄 License

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！
