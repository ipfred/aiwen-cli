# aw-cli

AIWEN/IPPlus360 IP 情报查询 CLI，支持 IPv4/IPv6 地理定位、当前网络 IP、应用场景、WHOIS、AS 映射、宿主信息、风险画像、真人/机器流量判断和行业分类。

## 安装

### 从源码构建

```bash
git clone https://github.com/aiwen/aw-cli.git
cd aw_cli
go build -o aw-cli .
```

### 依赖

- Go 1.23+

## 快速开始

### 1. 配置 API Key

```bash
# 方式一：配置文件设置（推荐）
aw-cli config set api_key YOUR_API_KEY

# 方式二：环境变量
export AIWEN_API_KEY=YOUR_API_KEY

# 方式三：命令行临时覆盖
aw-cli loc 8.8.8.8 --api-key YOUR_API_KEY
```

### 2. 初始化配置

```bash
aw-cli config init
```

将在 `~/.aw-cli/config.json` 生成默认配置文件。

### 3. 查看当前配置

```bash
aw-cli config show
```

输出（API Key 自动脱敏）：

```json
{
  "base_url": "https://api.ipplus360.com",
  "api_key": "",
  "channel": "aw_cli",
  "timeout": "10s",
  "ipv4_accuracy": "city",
  "ipv6_accuracy": "city"
}
```

## 命令一览

| 命令 | 说明 | IP 支持 |
|---|---|---|
| `aw-cli loc <ip>` | IP 地理定位（城市/区县/街道） | IPv4 / IPv6 |
| `aw-cli current` | 当前网络出口 IP 定位 | IPv4 / IPv6 |
| `aw-cli scene <ip>` | IP 应用场景（住宅/数据中心/CDN 等） | IPv4 / IPv6 |
| `aw-cli whois <ip>` | IP WHOIS 注册信息 | 仅 IPv4 |
| `aw-cli asn <ip>` | AS 号 / AS WHOIS 映射 | 仅 IPv4 |
| `aw-cli host <ip>` | IP 宿主归属信息 | 仅 IPv4 |
| `aw-cli risk <ip>` | IP 风险画像（VPN/代理/Tor 等） | 仅 IPv4 |
| `aw-cli identity <ip>` | 真人/机器流量判断 | 仅 IPv4 |
| `aw-cli industry <ip>` | IP 行业分类 | 仅 IPv4 |
| `aw-cli batch <file>` | 批量 IP 查询 | — |
| `aw-cli config init` | 初始化配置文件 | — |
| `aw-cli config show` | 显示当前配置 | — |
| `aw-cli config set <key> <value>` | 设置配置值 | — |
| `aw-cli completion <shell>` | 生成 Shell 补全脚本 | — |

## 使用示例

### IP 地理定位

```bash
# 默认城市级定位
aw-cli loc 8.8.8.8

# 区县级定位
aw-cli loc 8.8.8.8 --accuracy district

# 街道级定位
aw-cli loc 8.8.8.8 --accuracy street

# IPv6 定位
aw-cli loc 2001:4860:4860::8888

# 指定坐标系
aw-cli loc 8.8.8.8 --coordsys GCJ02

# 仅查看请求（不调用上游，密钥脱敏）
aw-cli loc 8.8.8.8 --dry-run
```

输出格式：

```json
{
  "ok": true,
  "action": "loc",
  "ip": "8.8.8.8",
  "data": {
    "country": "US",
    "province": "California",
    "city": "Mountain View",
    ...
  }
}
```

### 当前网络 IP 定位

```bash
# 查询当前出口 IP 的地理位置
aw-cli current

# 指定精度
aw-cli current --accuracy district
```

### IP 应用场景

```bash
# 查询 IP 使用场景
aw-cli scene 8.8.8.8

# 英文返回
aw-cli scene 8.8.8.8 --lang en
```

### IPv4 专用查询

```bash
# WHOIS 信息
aw-cli whois 1.1.1.1

# AS 映射
aw-cli asn 1.1.1.1

# 宿主信息
aw-cli host 1.1.1.1

# 风险画像
aw-cli risk 1.1.1.1

# 真人/机器判断
aw-cli identity 1.1.1.1

# 行业分类
aw-cli industry 1.1.1.1
```

> 注意：`whois`、`asn`、`host`、`risk`、`identity`、`industry` 仅支持 IPv4。传入 IPv6 地址会返回验证错误。

### 批量查询

```bash
# 从文本文件查询（每行一个 IP）
aw-cli batch ips.txt --action loc

# 从 CSV 文件查询
aw-cli batch ips.csv --ip-column ip --action risk

# 指定输出文件和并发数
aw-cli batch ips.txt --action loc --output result.ndjson --format ndjson --concurrency 10

# 查询所有支持的 action
aw-cli batch ips.txt --action all --format csv --output result.csv

# 包含私网/保留地址（默认跳过）
aw-cli batch ips.txt --action loc --include-private

# txt 文件格式（# 开头为注释）
# 8.8.8.8
# 1.1.1.1
# 2001:4860:4860::8888

# csv 文件格式
# ip,name
# 8.8.8.8,Google DNS
# 1.1.1.1,Cloudflare
```

### 输出格式

```bash
# JSON（默认）
aw-cli loc 8.8.8.8 --format json

# NDJSON（批量查询适合流式处理）
aw-cli batch ips.txt --action loc --format ndjson

# 表格格式
aw-cli loc 8.8.8.8 --format table

# CSV 格式
aw-cli loc 8.8.8.8 --format csv

# 使用 jq 过滤表达式提取字段
aw-cli loc 8.8.8.8 --jq .data.country
```

### 配置管理

```bash
# 初始化配置文件
aw-cli config init

# 查看当前配置（密钥自动脱敏）
aw-cli config show

# 设置配置值
aw-cli config set api_key YOUR_KEY
aw-cli config set base_url https://api.ipplus360.com
aw-cli config set timeout 30s
aw-cli config set channel my_app
aw-cli config set ipv4_accuracy district
aw-cli config set ipv6_accuracy city
```

### Shell 补全

```bash
# Bash
aw-cli completion bash > /etc/bash_completion.d/aw-cli

# Zsh
aw-cli completion zsh > "${fpath[1]}/_aw-cli"

# Fish
aw-cli completion fish > ~/.config/fish/completions/aw-cli.fish

# PowerShell
aw-cli completion powershell >> $PROFILE
```

## 全局标志

| 标志 | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `--config` | string | `~/.aw-cli/config.json` | 配置文件路径 |
| `--base-url` | string | `https://api.ipplus360.com` | API 基础地址 |
| `--api-key` | string | — | 临时覆盖 API Key |
| `--timeout` | duration | `10s` | HTTP 超时时间 |
| `--format` | string | `json` | 输出格式：`json` / `ndjson` / `table` / `csv` |
| `--jq` / `-q` | string | — | JSON 过滤表达式（如 `.data.country`） |
| `--dry-run` | bool | `false` | 仅打印请求，不调用上游 |
| `--verbose` | bool | `false` | 调试输出，密钥自动脱敏 |

## 命令专属标志

### `loc`

| 标志 | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `--accuracy` | string | `city` | 定位精度：`city` / `district` / `street` |
| `--coordsys` | string | `WGS84` | 坐标系：`WGS84` / `GCJ02` / `BD09` |

### `current`

| 标志 | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `--accuracy` | string | `city` | 定位精度 |
| `--coordsys` | string | `WGS84` | 坐标系 |

### `scene`

| 标志 | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `--lang` | string | `cn` | 返回语言 |

### `batch`

| 标志 | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `--action` | string | `loc` | 查询类型：`loc` / `scene` / `whois` / `asn` / `host` / `risk` / `identity` / `industry` / `all` |
| `--ip-column` | string | `ip` | CSV 中 IP 列名 |
| `-o, --output` | string | stdout | 输出文件路径 |
| `--concurrency` | int | `5` | 并发请求数 |
| `--retries` | int | `2` | 网络错误重试次数 |
| `--include-private` | bool | `false` | 包含私网/保留地址 |

## 配置优先级

命令行 flag > 环境变量 > 配置文件 > 默认值

### 环境变量

| 变量 | 说明 |
|---|---|
| `AIWEN_API_KEY` | API 密钥 |
| `AIWEN_API_BASE_URL` | API 基础地址 |
| `AIWEN_CHANNEL` | 渠道标识 |
| `AIWEN_TIMEOUT` | HTTP 超时 |
| `IPV4_ACCURACY` | IPv4 默认定位精度 |
| `IPV6_ACCURACY` | IPv6 默认定位精度 |

## 退出码

| 码 | 类型 | 说明 |
|---|---|---|
| 0 | 成功 | 命令执行成功 |
| 1 | `internal` | 程序内部错误 |
| 2 | `validation` | 参数错误、IP 非法、IPv6 传入了 IPv4-only 接口 |
| 3 | `config` | 缺少 API Key、配置文件格式错误 |
| 4 | `api_error` / `parse_error` | 上游返回错误或非 JSON |
| 5 | `network` | 网络错误、超时 |

错误输出示例：

```json
{
  "ok": false,
  "error": {
    "type": "validation",
    "message": "action risk only supports IPv4"
  }
}
```

## 项目结构

```
aw_cli/
├── main.go                  # 入口
├── cmd/                     # Cobra 命令定义
│   ├── root.go              # 根命令 & 全局标志 & 错误处理
│   ├── loc.go               # loc 子命令
│   ├── current.go           # current 子命令
│   ├── scene.go             # scene 子命令
│   ├── ip_commands.go        # whois/asn/host/risk/identity/industry
│   ├── batch.go              # batch 子命令
│   ├── config.go             # config 子命令组
│   ├── completion.go         # completion 子命令
│   └── ip/
│       └── ip.go             # IP 查询共享逻辑
├── errs/
│   └── types.go              # 类型化错误 & 退出码
├── internal/
│   ├── batch/
│   │   └── batch.go          # 批量查询 & 并发 worker
│   ├── build/
│   │   └── build.go          # 版本信息
│   ├── client/
│   │   └── aiwen.go         # HTTP 请求构造 & 响应处理
│   ├── cmdutil/
│   │   ├── factory.go        # Factory 依赖注入
│   │   └── iostreams.go      # IO 流抽象
│   ├── core/
│   │   ├── config.go         # 配置加载/解析/写入
│   │   └── secret.go         # 密钥脱敏
│   ├── endpoint/
│   │   └── endpoint.go       # Action 元数据 & Endpoint 路径映射
│   ├── iputil/
│   │   └── ip.go             # IP 校验 & 版本判断
│   └── output/
│       └── format.go         # JSON/NDJSON/Table/CSV 格式化输出
├── skills/
│   └── aw-cli-query/
│       ├── SKILL.md           # Agent Skill 描述
│       └── references/
│           ├── api.md          # API 参考文档
│           ├── response-fields.md  # 响应字段说明
│           └── errors.md       # 错误码参考
├── go.mod
└── Makefile
```

## 开发

```bash
# 运行测试
make test

# 格式化
make fmt

# 构建
make build

# 直接运行
go run . --help
go run . loc 8.8.8.8 --dry-run
```

## License

Private — AIWEN/IPPlus360


# TODO 
1. 通过npx @larksuite/cli@latest install的方式安装 skills 和cli 