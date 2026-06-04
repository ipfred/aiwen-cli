# aw-cli

[埃文科技](https://www.ipplus360.com/) 官方的IP地址信息查询 CLI，支持 IPv4/IPv6定位、当前网络IP位置、IP场景、WHOIS、AS域名、宿主信息、IP风险画像、IP真假人和IPv4行业分类等信息查询。


## 安装

### 方式一：npx 一键安装+配置+Skills（推荐）

适合普通用户首次安装。该命令会启动交互式安装向导，并完成三件事：

```bash
npx aiwen-geoip-cli@latest install
```

- 安装或升级全局 CLI：`aiwen-geoip-cli`
- 配置 API Key
- 从独立的 [aiwen-skills](https://github.com/ipfred/aiwen-skills) 仓库拉取并安装 AI Skills

安装完成后验证：

```bash
aw-cli --version
aw-cli loc 8.8.8.8 --dry-run
```

如果安装向导中跳过了 API Key，可以稍后手动配置：

```bash
aw-cli config set api_key YOUR_API_KEY
```

### 方式二：只安装 CLI 工具

适合只需要使用 `aw-cli` 命令，或希望自己管理 API Key 和 Skills 的场景。

```shell
npm install -g aiwen-geoip-cli
aw-cli --version
#  配置key
aw-cli config set api_key YOUR_API_KEY
```

### 方式三：从 GitHub Release 手动安装 CLI

适合不能直接通过 npm 下载二进制包，或需要手动分发到内网机器的场景。

1. 打开 Release 页面并下载对应系统和架构的压缩包：

   <https://github.com/ipfred/aiwen-cli/releases/latest>

   | 系统 | 架构 | 文件名 |
   |---|---|---|
   | Linux | amd64 | `aw-cli-<version>-linux-amd64.tar.gz` |
   | Linux | arm64 | `aw-cli-<version>-linux-arm64.tar.gz` |
   | macOS | amd64 | `aw-cli-<version>-darwin-amd64.tar.gz` |
   | macOS | arm64 | `aw-cli-<version>-darwin-arm64.tar.gz` |
   | Windows | amd64 | `aw-cli-<version>-windows-amd64.zip` |
   | Windows | arm64 | `aw-cli-<version>-windows-arm64.zip` |

2. 解压并把可执行文件所在目录加入 `PATH`。

   Linux / macOS 示例：

   ```bash
   mkdir -p ~/.local/aw-cli
   tar -xzf aw-cli-<version>-linux-amd64.tar.gz -C ~/.local/aw-cli
   export PATH="$HOME/.local/aw-cli:$PATH"
   ```

   Windows PowerShell 示例：

   ```powershell
   New-Item -ItemType Directory -Force C:\Tools\aw-cli
   Expand-Archive .\aw-cli-<version>-windows-amd64.zip -DestinationPath C:\Tools\aw-cli -Force
   ```

   Windows 需要把 `C:\Tools\aw-cli` 加入系统或用户 `PATH`。

3. 配置 API Key 并验证：

   ```bash
   aw-cli config set api_key YOUR_API_KEY
   aw-cli --version
   ```


### AI Skills安装

本仓库主要维护 CLI 的安装、配置和功能说明。AI Skills 独立维护在 [ipfred/aiwen-skills](https://github.com/ipfred/aiwen-skills)，安装方式、更新方式、作用域说明和 Skill 使用说明请以该仓库为准。

### 环境依赖

- Node.js 16+：用于 `npm` / `npx` 安装。
- Go 1.23+：仅源码构建需要。

## 卸载

```shell
# 卸载cli工具
npm uninstall -g aiwen-geoip-cli
```

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
| `aw-cli scene <ip>` | IP 应用场景（家庭宽带/数据中心/企业专线 等） | IPv4 / IPv6 |
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

This project is licensed under the MIT License. When running, it calls AIWEN/IPPlus360 APIs. To use these APIs, you must comply with the following agreements and privacy policies:

- [埃文科技隐私政策 (AIWEN Privacy Policy)](https://aidata.ipplus360.com/privacy-policy/)
