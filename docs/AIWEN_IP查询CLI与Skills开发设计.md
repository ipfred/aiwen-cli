# AIWEN IP 查询 CLI + Skills 开发设计（Go / Cobra 架构版）

本文基于当前目录的 `aiwen_loc.py` 和参考项目 `E:\my_work\github_pro\cli` 的技术栈与架构进行设计。目标是后续可以直接按本文落地开发一个 Go 语言 CLI，并配套 Codex Skills。

## 1. 技术栈选型

参考项目 `E:\my_work\github_pro\cli` 的核心技术栈：

| 类别 | 参考项目技术 | 本项目采用 |
|---|---|---|
| 语言 | Go 1.23 | Go 1.23+ |
| CLI 框架 | `github.com/spf13/cobra` | 同样使用 Cobra |
| Flag 框架 | `github.com/spf13/pflag` | Cobra 内置 pflag |
| 输出格式 | `internal/output`，支持 json / ndjson / table / csv | 同样抽象 `internal/output` |
| HTTP 客户端 | `net/http` + client 封装 | `net/http` + `internal/client` |
| 配置 | `internal/core` + config 文件 + env | `internal/core` + config 文件 + env |
| 错误处理 | `errs/` + root 统一处理退出码 | 简化版 typed error + root 统一处理 |
| 依赖注入 | `internal/cmdutil.Factory` | 同样使用 Factory |
| Skills | 根目录 `skills/<skill-name>/SKILL.md` | 同样放 `skills/aw-cli-query` |
| 测试 | 单测 + `tests/cli_e2e` | 单测 + CLI e2e |

不建议继续使用 Python 作为主实现。原 `aiwen_loc.py` 只作为接口和能力来源，后续 MCP 服务也应复用 Go CLI 的核心 client，而不是维护两套请求逻辑。

## 2. 现有 MCP 能力清单

`aiwen_loc.py` 当前 MCP 服务：

- MCP 服务名：`aiwen_ip_geo`
- 上游 Host：`https://api.ipplus360.com`
- 认证环境变量：`AIWEN_API_KEY`
- 公共参数：
  - `key`: API key
  - `channel`: 当前为 `py_mcp`，Go CLI 建议默认 `aw_cli`
- 定位精度环境变量：
  - `IPV4_ACCURACY`: 默认 `city`
  - `IPV6_ACCURACY`: 默认 `city`
  - 可选值：`city`、`district`、`street`

当前 MCP 暴露了 9 类能力：

| Action | MCP 工具名 | 功能 | IP 支持 |
|---|---|---|---|
| `loc` | `aiwen_ip_location` | IP 地理定位 | IPv4 / IPv6 |
| `current` | `user_network_ip` | 当前网络出口 IP 定位 | IPv4 / IPv6 |
| `scene` | `ip_usage_scene` | IP 应用场景 | IPv4 / IPv6 |
| `whois` | `ip_whois_info` | IP WHOIS 注册信息 | IPv4 |
| `asn` | `ip_as_mapping` | AS WHOIS / IP 到 AS 映射 | IPv4 |
| `host` | `ip_host_info` | IP 宿主归属信息 | IPv4 |
| `risk` | `ip_risk_portrait` | IP 风险画像 | IPv4 |
| `identity` | `ip_identity_check` | 真人 / 机器流量判断 | IPv4 |
| `industry` | `ip_industry_classify` | IP 行业分类 | IPv4 |

## 3. API 接口与参数

### 3.1 公共参数

所有 `api.ipplus360.com` 请求自动追加：

| 参数 | 类型 | 必填 | 来源 | 说明 |
|---|---|---:|---|---|
| `key` | string | 是 | `AIWEN_API_KEY` 或配置文件 | API 访问密钥 |
| `channel` | string | 是 | 默认 `aw_cli` | 调用渠道标识 |

### 3.2 IP 定位：`loc`

功能：查询大洲、国家、省份、城市、区县、街道、经纬度、机构、运营商、精度等信息。

参数：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---:|---|---|
| `ip` | string | 是 | 无 | IPv4 或 IPv6 |
| `coordsys` | string | 否 | `WGS84` | 坐标系 |
| `accuracy` | enum | 否 | `city` | `city` / `district` / `street` |

Endpoint 映射：

| IP 类型 | 精度 | Endpoint |
|---|---|---|
| IPv4 | `city` | `/ip/geo/v1/city/` |
| IPv4 | `district` | `/ip/geo/v1/district/` |
| IPv4 | `street` | `/ip/geo/v1/street/psi/` |
| IPv6 | `city` | `/ip/geo/v1/ipv6/` |
| IPv6 | `district` | `/ip/geo/v1/ipv6/district/` |
| IPv6 | `street` | `/ip/geo/v1/ipv6/street/biz/` |

### 3.3 当前网络 IP：`current`

功能：先获取当前公网 IP，再调用 `loc`。

调用链：

```text
GET https://www.ipuu.net/ipuu/user/getIP
  -> response.data
  -> loc(ip)
```

参数：无。

### 3.4 IP 应用场景：`scene`

功能：查询保留 IP、未分配 IP、组织机构、移动网络、家庭宽带、数据中心、企业专线、CDN、卫星通信、交换中心、Anycast 等场景。

参数：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---:|---|---|
| `ip` | string | 是 | 无 | IPv4 或 IPv6 |
| `lang` | string | 否 | `cn` | 返回语言 |

Endpoint：

| IP 类型 | Endpoint |
|---|---|
| IPv4 | `/ip/info/v1/scene/` |
| IPv6 | `/ip/info/v1/ipv6Scene/` |

### 3.5 IPv4-only 接口

| Action | 功能 | Endpoint | 参数 |
|---|---|---|---|
| `whois` | 查询 IP 注册网段、机构、技术联系人、管理员 | `/ip/info/v1/ipWhois` | `ip: string` |
| `asn` | 查询 AS 号 / AS Whois | `/as/info/v1/asWhois` | `ip: string` |
| `host` | 查询 AS Number、AS 名称、运营商、所属机构 | `/ip/geo/v1/host/` | `ip: string` |
| `risk` | 查询 VPN、代理、Tor、数据中心、风险评分 | `/ip/info/v3/portrait/` | `ip: string` |
| `identity` | 判断真人概率、秒拨概率、机器流量 | `/ip/info/v1/person/` | `ip: string` |
| `industry` | 查询行业分类 | `/ip/info/v1/industry/` | `ip: string` |

Go CLI 必须在调用前校验 IPv4-only 能力。如果用户传 IPv6，返回 validation error，不调用上游。

## 4. 目标项目结构

按参考项目风格设计：

```text
aw_cli/
  go.mod
  go.sum
  main.go
  Makefile
  README.md
  cmd/
    root.go
    root_test.go
    global_flags.go
    ip/
      ip.go
      loc.go
      current.go
      scene.go
      whois.go
      asn.go
      host.go
      risk.go
      identity.go
      industry.go
      batch.go
      commands_test.go
    config/
      config.go
      init.go
      show.go
      set.go
    completion/
      completion.go
  errs/
    types.go
    category.go
    problem.go
  internal/
    build/
      build.go
    client/
      aiwen.go
      request.go
      response.go
      errors.go
      aiwen_test.go
    cmdutil/
      factory.go
      iostreams.go
      json.go
      dryrun.go
      completion.go
    core/
      config.go
      paths.go
      secret.go
      errors.go
    endpoint/
      endpoint.go
      registry.go
    iputil/
      ip.go
      ip_test.go
    output/
      format.go
      json.go
      ndjson.go
      table.go
      csv.go
      errors.go
    batch/
      batch.go
      input.go
      worker.go
      batch_test.go
    validate/
      input.go
      path.go
  skills/
    aw-cli-query/
      SKILL.md
      references/
        api.md
        response-fields.md
        errors.md
  tests/
    cli_e2e/
      loc_workflow_test.go
      batch_workflow_test.go
      errors_test.go
```

与参考项目保持一致的关键点：

- `main.go` 只负责 `os.Exit(cmd.Execute())`。
- `cmd/root.go` 负责构建根命令、全局 flag、统一错误处理。
- 每个命令使用 `Options` 结构体承载参数。
- 每个命令提供 `NewCmdXxx(f *cmdutil.Factory, runF func(*Options) error)`，方便测试注入。
- `internal/cmdutil.Factory` 统一提供配置、HTTP client、IOStreams。
- `internal/client` 只做 API 请求，不处理终端展示。
- `internal/output` 统一格式化 `json / ndjson / table / csv`。
- `skills/` 与 CLI 在同一仓库维护。

## 5. Go 模块与依赖

建议 `go.mod`：

```go
module github.com/your-org/aw-cli

go 1.23.0

require (
    github.com/spf13/cobra v1.10.2
    github.com/spf13/pflag v1.0.9
    github.com/stretchr/testify v1.11.1
    github.com/tidwall/gjson v1.18.0
    github.com/itchyny/gojq v0.12.17
)
```

可选依赖：

- `github.com/charmbracelet/lipgloss`：如果要做更好的 table / 颜色输出。
- `gopkg.in/yaml.v3`：如果后续要生成 OpenAPI / MCP 配置。

MVP 阶段不引入过多 UI 库，优先保证 CLI 可被 Agent 稳定调用。

## 6. 命令设计

命令名建议：`aw-cli`

根命令：

```bash
aw-cli <command> [flags]
```

### 6.1 单 IP 查询

```bash
aw-cli loc 8.8.8.8
aw-cli loc 8.8.8.8 --accuracy district --coordsys WGS84 --format json
aw-cli loc 2001:4860:4860::8888 --accuracy city

aw-cli current
aw-cli scene 8.8.8.8
aw-cli whois 1.1.1.1
aw-cli asn 1.1.1.1
aw-cli host 1.1.1.1
aw-cli risk 1.1.1.1
aw-cli identity 1.1.1.1
aw-cli industry 1.1.1.1
```

全局 flags：

| Flag | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `--config` | string | 空 | 指定配置文件 |
| `--base-url` | string | `https://api.ipplus360.com` | 覆盖 API host |
| `--api-key` | string | 空 | 临时覆盖 key；不建议在日常命令中使用 |
| `--timeout` | duration | `10s` | HTTP 超时 |
| `--format` | enum | `json` | `json` / `ndjson` / `table` / `csv` |
| `--jq`, `-q` | string | 空 | JSON 查询过滤 |
| `--dry-run` | bool | false | 只打印请求，不发起调用 |
| `--verbose` | bool | false | 调试输出，密钥必须脱敏 |

`loc` flags：

| Flag | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `--accuracy` | enum | config 或 `city` | `city` / `district` / `street` |
| `--coordsys` | string | `WGS84` | 坐标系 |

`scene` flags：

| Flag | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `--lang` | string | `cn` | 语言 |

### 6.2 批量查询

```bash
aw-cli batch ips.txt --action loc --accuracy city --format ndjson
aw-cli batch ips.csv --ip-column ip --action risk --output risk.csv --format csv
aw-cli batch ips.txt --action all --output result.ndjson
```

批量 flags：

| Flag | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `--action` | enum | `loc` | `loc` / `scene` / `whois` / `asn` / `host` / `risk` / `identity` / `industry` / `all` |
| `--ip-column` | string | `ip` | CSV 输入 IP 列 |
| `--output`, `-o` | string | stdout | 输出路径 |
| `--concurrency` | int | `5` | 并发数 |
| `--retries` | int | `2` | 网络错误重试次数 |
| `--include-private` | bool | false | 是否允许私网 / 保留地址 |

批量输出每条记录建议：

```json
{
  "ip": "8.8.8.8",
  "action": "loc",
  "ok": true,
  "status_code": 200,
  "data": {},
  "error": null
}
```

## 7. 核心包设计

### 7.1 `cmd/root.go`

职责：

- 创建 root Cobra command。
- 注册全局 flags。
- 创建 `cmdutil.Factory`。
- 注册 `loc/current/scene/.../batch/config/completion` 子命令。
- 统一捕获 error 并转成退出码。

参考结构：

```go
func Execute() int {
    f, rootCmd := buildInternal(context.Background(), os.Args[1:])
    if err := rootCmd.Execute(); err != nil {
        return handleRootError(f, err)
    }
    return 0
}
```

### 7.2 `cmd/ip/*.go`

每个命令保持参考项目模式：

```go
type LocOptions struct {
    Factory  *cmdutil.Factory
    Ctx      context.Context
    IP       string
    Accuracy string
    CoordSys string
    Format   string
    JqExpr   string
    DryRun   bool
}

func NewCmdLoc(f *cmdutil.Factory, runF func(*LocOptions) error) *cobra.Command
```

命令执行流程：

```text
解析 args/flags
  -> 读取 config
  -> 校验 IP 和 action
  -> dry-run 分支
  -> client.Query(...)
  -> output.FormatValue(...)
```

### 7.3 `internal/client`

核心类型：

```go
type AiwenClient struct {
    BaseURL    string
    APIKey     string
    Channel    string
    HTTP       *http.Client
    ErrOut     io.Writer
}

type QueryRequest struct {
    Action   string
    IP       string
    Accuracy string
    CoordSys string
    Lang     string
}

type QueryResult struct {
    StatusCode int
    Raw        []byte
    JSON       any
}
```

核心方法：

```go
func (c *AiwenClient) Query(ctx context.Context, req QueryRequest) (*QueryResult, error)
func (c *AiwenClient) Current(ctx context.Context, req QueryRequest) (*QueryResult, error)
func (c *AiwenClient) BuildRequest(ctx context.Context, req QueryRequest) (*http.Request, error)
func (c *AiwenClient) CheckResponse(result any) error
```

要求：

- 使用 `net/http`。
- 对 API key 日志脱敏。
- JSON 解析失败时保留 raw body 前 4KB 作为错误 detail。
- 业务错误码不在 client 内直接打印，由 error 返回给 root 统一处理。

### 7.4 `internal/endpoint`

用结构化元数据替代散落的 if/else：

```go
type ActionSpec struct {
    Name         string
    SupportsIPv4 bool
    SupportsIPv6 bool
    IPv4Only     bool
    Params       []string
    Paths        map[string]any
}
```

需要覆盖：

- `loc`
- `scene`
- `whois`
- `asn`
- `host`
- `risk`
- `identity`
- `industry`

### 7.5 `internal/iputil`

职责：

- 使用 `net/netip` 或 `net.ParseIP` 校验 IP。
- 判断 IPv4 / IPv6。
- 判断 private / loopback / multicast / unspecified。
- 校验 action 是否支持当前 IP 版本。

建议优先用 Go 标准库 `net/netip`：

```go
addr, err := netip.ParseAddr(ip)
addr.Is4()
addr.Is6()
addr.IsPrivate()
addr.IsLoopback()
```

### 7.6 `internal/core`

配置文件建议：

```text
%USERPROFILE%\.aw-cli\config.json
~/.aw-cli/config.json
```

配置结构：

```go
type CliConfig struct {
    BaseURL      string `json:"base_url"`
    APIKey       string `json:"api_key"`
    Channel      string `json:"channel"`
    Timeout      string `json:"timeout"`
    IPv4Accuracy string `json:"ipv4_accuracy"`
    IPv6Accuracy string `json:"ipv6_accuracy"`
}
```

配置优先级：

```text
命令行 flag > 环境变量 > config.json > 默认值
```

环境变量：

| 环境变量 | 说明 |
|---|---|
| `AIWEN_API_KEY` | API key |
| `AIWEN_API_BASE_URL` | API host |
| `AIWEN_CHANNEL` | channel |
| `IPV4_ACCURACY` | IPv4 默认定位精度 |
| `IPV6_ACCURACY` | IPv6 默认定位精度 |

### 7.7 `internal/output`

参考项目风格支持：

- `json`
- `ndjson`
- `table`
- `csv`

建议输出默认包一层稳定 envelope，方便 Agent 和自动化解析：

```json
{
  "ok": true,
  "action": "loc",
  "ip": "8.8.8.8",
  "data": {}
}
```

错误输出写 stderr：

```json
{
  "ok": false,
  "error": {
    "type": "validation",
    "message": "action risk only supports IPv4",
    "ip": "2001:4860:4860::8888",
    "action": "risk"
  }
}
```

## 8. 错误与退出码

沿用参考项目“root 统一处理错误”的思想，但本项目可以先做简化版：

| 错误类型 | 场景 | 退出码 |
|---|---|---:|
| `validation` | 参数错误、IP 非法、IPv6 调 IPv4-only 接口 | 2 |
| `config` | 缺少 API key、配置非法 | 3 |
| `api_error` | 上游业务错误码非成功 | 4 |
| `network` | 网络错误、超时、HTTP 5xx | 5 |
| `parse_error` | 上游返回非 JSON 或结构异常 | 4 |
| `internal` | 程序内部错误 | 1 |

`errs/` 包建议包含：

```go
type Problem struct {
    Type    string         `json:"type"`
    Message string         `json:"message"`
    Hint    string         `json:"hint,omitempty"`
    Detail  map[string]any `json:"detail,omitempty"`
}

type ExitError struct {
    Code    int
    Problem Problem
}
```

## 9. Skill 设计

Skill 路径：

```text
skills/aw-cli-query/
  SKILL.md
  references/
    api.md
    response-fields.md
    errors.md
```

参考项目的 Skill frontmatter 风格：

```yaml
---
name: aw-cli-query
version: 1.0.0
description: "AIWEN/IPPlus360 IP 情报查询 skill。用于通过 aw-cli CLI 查询 IPv4/IPv6 地理定位、当前网络 IP、应用场景、WHOIS、AS 映射、宿主信息、风险画像、真人概率和行业分类。当用户要求查询 IP 位置、运营商、归属机构、风险、VPN/代理/Tor、数据中心、真人/机器流量、行业分类或批量导出 IP 情报时使用。"
metadata:
  requires:
    bins: ["aw-cli"]
  cliHelp: "aw-cli --help"
---
```

`SKILL.md` 正文只写高频决策：

| 用户意图 | CLI 命令 |
|---|---|
| 查 IP 位置、省市、经纬度、运营商 | `aw-cli loc <ip>` |
| 查当前机器出口 IP | `aw-cli current` |
| 查住宅宽带 / 数据中心 / CDN / Anycast | `aw-cli scene <ip>` |
| 查 WHOIS 注册信息 | `aw-cli whois <ip>` |
| 查 AS 号 | `aw-cli asn <ip>` |
| 查宿主、AS 名称、运营商、机构 | `aw-cli host <ip>` |
| 查 VPN、代理、Tor、风险分 | `aw-cli risk <ip>` |
| 查真人概率、机器流量、秒拨 | `aw-cli identity <ip>` |
| 查行业分类 | `aw-cli industry <ip>` |
| 多 IP 查询 | `aw-cli batch <file> --action <action>` |

Skill 注意事项：

- 不要猜测 IP 查询结果，必须调用 CLI。
- 多个 IP 使用 `batch`。
- IPv4-only 能力遇到 IPv6 时，说明该能力不支持 IPv6。
- 字段解释读取 `references/response-fields.md`。
- API 细节读取 `references/api.md`。
- 错误处理读取 `references/errors.md`。

不需要在 Skill 内放 Python 脚本。参考项目的 Skills 主要指导 Agent 使用已安装 CLI，因此本项目也让 Skill 依赖 `aw-cli` 二进制。

## 10. MCP 服务后续处理

当前 `aiwen_loc.py` 可以保留为历史参考。若必须继续提供 MCP，建议后续单独做 Go MCP 或 HTTP wrapper，但核心逻辑仍复用 Go client。

可选方案：

```text
aw-cli mcp serve
```

或者：

```text
cmd/mcp/
internal/mcp/
```

MCP 工具名保持兼容：

| 旧 MCP 工具名 | 内部 action |
|---|---|
| `aiwen_ip_location` | `loc` |
| `user_network_ip` | `current` |
| `ip_usage_scene` | `scene` |
| `ip_whois_info` | `whois` |
| `ip_as_mapping` | `asn` |
| `ip_host_info` | `host` |
| `ip_risk_portrait` | `risk` |
| `ip_identity_check` | `identity` |
| `ip_industry_classify` | `industry` |

MVP 阶段优先开发 CLI + Skill，MCP 服务放到后续阶段。

## 11. 测试设计

### 11.1 单元测试

| 路径 | 测试内容 |
|---|---|
| `internal/iputil` | IPv4/IPv6/非法 IP/private/loopback |
| `internal/endpoint` | action 到 endpoint 映射 |
| `internal/client` | 请求参数、key/channel 注入、错误解析 |
| `internal/output` | json/ndjson/table/csv 输出 |
| `cmd/ip` | flags、参数校验、dry-run |
| `internal/batch` | 输入解析、并发、失败记录 |

### 11.2 E2E 测试

放在：

```text
tests/cli_e2e/
```

用例：

```bash
aw-cli loc 8.8.8.8 --format json
aw-cli loc 2001:4860:4860::8888 --format json
aw-cli scene 8.8.8.8 --format json
aw-cli risk 1.1.1.1 --format json
aw-cli current --format json
aw-cli batch ips.txt --action loc --format ndjson
```

异常用例：

```bash
aw-cli loc not-an-ip
aw-cli risk 2001:4860:4860::8888
aw-cli loc 8.8.8.8 --accuracy invalid
aw-cli loc 8.8.8.8
```

最后一个用例在未设置 `AIWEN_API_KEY` 时应返回 config error。

## 12. 开发阶段规划

### Phase 1：Go CLI 骨架

- `go.mod`
- `main.go`
- `cmd/root.go`
- `internal/cmdutil.Factory`
- `internal/core` 配置读取
- `internal/output` JSON 输出
- `internal/iputil` IP 校验

### Phase 2：核心 API Client

- `internal/endpoint`
- `internal/client.AiwenClient`
- `aw-cli loc`
- `aw-cli scene`
- `aw-cli risk`
- `--dry-run`
- 基础单测

### Phase 3：补齐所有命令

- `current`
- `whois`
- `asn`
- `host`
- `identity`
- `industry`
- IPv4-only 校验
- `table/csv/ndjson`

### Phase 4：批量查询

- `aw-cli batch`
- txt/csv/jsonl 输入
- 并发 worker
- retry
- JSONL/CSV 输出

### Phase 5：Skills

- `skills/aw-cli-query/SKILL.md`
- `references/api.md`
- `references/response-fields.md`
- `references/errors.md`
- 验证 Agent 能根据用户意图选择正确命令

### Phase 6：可选 MCP

- `aw-cli mcp serve`
- 保留旧 MCP 工具名
- 去掉 Python MCP 中的 key 泄露风险

## 13. 现有 `aiwen_loc.py` 迁移注意事项

必须修正：

- 不再用 `":" in ip` 判断 IPv6，改用 Go `net/netip`。
- 不再在 import 阶段读取 key，改成运行时从 config/env 解析。
- 不打印请求参数，尤其不能打印 key。
- 不保留 `return response.text` 后的不可达代码。
- IPv4-only action 明确校验。
- API key 缺失返回 config error，退出码 3。
- 上游非 JSON 返回 parse error。
- 上游业务错误码返回 api_error。

## 14. 结论

后续开发应以 Go CLI 为核心，而不是扩展 Python MCP。整体架构对齐 `E:\my_work\github_pro\cli`：

```text
main.go -> cmd.Execute() -> Cobra commands -> cmdutil.Factory -> internal/client -> internal/output
```

Skills 也按参考项目方式维护在仓库 `skills/` 下，让 Agent 学会稳定调用 `aw-cli` 二进制。这样 CLI、人类终端、Agent Skills、未来 MCP wrapper 都共享同一套 Go client 和 endpoint 元数据，维护成本最低。

