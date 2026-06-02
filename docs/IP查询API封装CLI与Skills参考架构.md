# IP查询API封装为 CLI + Skills 的项目参考分析

> 当前目录 `D:\my_work\ai_pros\aw_cli` 为空，本分析基于同级目录中最相关的项目：
>
> - `D:\my_work\ai_pros\baidu_ip_api\ip_query_plugin`
> - `D:\my_work\ai_pros\baidu_ip_api\demo2`
> - `D:\my_work\ai_pros\baidu_ip_api\埃文API.md`
> - `D:\my_work\ai_pros\fastapi_aiip_answer\table_know_job`

## 1. 现有项目概览

`baidu_ip_api/ip_query_plugin` 是一个把 IP 查询 API 封装成百度 Agents 插件服务的轻量项目。它不是 CLI 项目，也不是 Codex Skill 项目，但它已经具备后续封装 CLI + Skills 时最有价值的几个组成部分：

- API 调用封装：根据 IPv4/IPv6 选择不同上游接口。
- 输入校验：使用 Python `ipaddress` 校验 IP 地址合法性。
- 插件服务：通过 Flask 暴露 `/query_ip`。
- 插件注册文件：提供 `/.well-known/ai-plugin.json`、`/.well-known/openapi.yaml`、`/example.yaml`。
- 凭证管理雏形：`ip_query_plugin/server.py` 使用环境变量 `IPPLUS360_KEY`，比把 key 写进 OpenAPI 更安全。

`fastapi_aiip_answer/table_know_job` 是另一个可参考点。它和 IP 查询业务无关，但有 CLI 化、客户端类封装、日志、配置、上下文管理等模式，适合借鉴到 IP API CLI。

## 2. 技术栈分析

### 2.1 IP 查询插件服务

主要文件：

- `ip_query_plugin/server.py`
- `ip_query_plugin/.well-known/ai-plugin.json`
- `ip_query_plugin/.well-known/openapi.yaml`
- `ip_query_plugin/example.yaml`
- `ip_query_plugin/requirements.txt`

技术栈：

| 类别 | 使用技术 | 作用 |
|---|---|---|
| 语言 | Python | 服务端和 API 调用逻辑 |
| Web 框架 | Flask | 暴露插件 HTTP 服务 |
| 跨域 | flask_cors | 允许百度 Agents 和本地前端访问 |
| HTTP 客户端 | requests | 调用上游 IP 查询 API |
| IP 校验 | ipaddress | 区分 IPv4/IPv6 并校验格式 |
| 插件协议 | ai-plugin.json + OpenAPI YAML | 供 Agent 平台识别工具能力 |
| 配置 | 环境变量 | `IPPLUS360_KEY` 保存 API key |

依赖非常少：

```txt
flask
flask_cors
requests
```

### 2.2 CLI 可参考项目

`fastapi_aiip_answer/table_know_job/main.py` 使用：

- `argparse` 解析命令行参数。
- `logging` + `TimedRotatingFileHandler` 做日志。
- 独立 Client 类封装 HTTP 请求。
- 配置文件 `conf/config.py` 承载服务端地址、账号等配置。
- `with Client(...) as client` 管理 session 生命周期。

这些模式适合迁移到 IP 查询 CLI：

- CLI 层只负责参数解析和输出格式。
- Client 层负责认证、请求、异常、响应解析。
- Service/Command 层负责业务流程，例如单 IP 查询、批量查询、输出 CSV/JSON。

## 3. 现有系统架构

### 3.1 `ip_query_plugin` 调用链

当前架构可以抽象为：

```text
用户 / Agent
  |
  | POST /query_ip
  v
Flask server.py
  |
  | 1. 读取 JSON: ip, coordsys
  | 2. 校验 ip 是否为空
  | 3. 从环境变量读取 IPPLUS360_KEY
  | 4. 用 ipaddress 判断 IPv4 / IPv6
  | 5. 选择上游 URL
  v
requests.get(...)
  |
  v
ipplus360 API
  |
  v
原样返回上游 JSON
```

### 3.2 插件注册架构

项目同时暴露平台所需的元数据：

```text
/.well-known/ai-plugin.json  -> 插件名称、描述、OpenAPI 地址、logo、示例地址
/.well-known/openapi.yaml    -> 工具接口定义、operationId、请求响应 schema
/example.yaml                -> 正例和反例，帮助 Agent 判断何时调用工具
/logo.png                    -> 插件图标
```

这套结构适合参考，但如果目标是 Codex Skills，不需要完全照搬百度插件协议。Codex Skill 更关注：

- `SKILL.md` 中的触发描述和操作流程。
- `scripts/` 中可重复执行的确定性脚本。
- `references/` 中存放 API 文档、字段说明、错误码等。
- 可选 `agents/openai.yaml` 用于 UI 展示元信息。

## 4. 当前项目优点

### 4.1 值得复用的设计

1. IP 类型判断应该保留

`ipaddress.ip_address(ip)` 既能校验 IP，又能通过 `version` 判断 IPv4/IPv6。后续 CLI 和 Skill 都应该复用这个逻辑，避免让 Agent 或用户手动选择 API 类型。

2. API key 不写进插件描述

`ip_query_plugin` 通过 `IPPLUS360_KEY` 读取 key，这是正确方向。`demo2/openapi.yaml` 中出现过把 key 放进 schema 默认值的写法，这个不建议复用。

3. OpenAPI 的 `operationId` 很关键

`queryIpGeo` 这种清晰的 operationId 有助于 Agent 稳定触发工具。后续公司内部 API 如果有多个能力，建议命名为：

- `queryIpGeo`
- `queryIpBatchGeo`
- `queryIpRisk`
- `queryIpWhois`
- `queryIpOwner`
- `queryIpLocationHistory`

4. 示例文件包含反例

`example.yaml` 里不仅有“查询 8.8.8.8”的正例，也有“IP地址是什么”无需调用插件的反例。这个思路很适合迁移到 Skills：告诉 Agent 哪些问题需要调用 CLI，哪些只是常识回答。

5. 服务代理比直连上游更安全

`ip_query_plugin` 的模式是 Agent 调本地/服务端代理，代理再带 key 调上游。公司内部 API 如果包含鉴权、审计、额度、敏感字段，建议继续采用代理或 CLI 本地鉴权，不要让 Agent 直接持有密钥。

## 5. 当前项目不足

### 5.1 不足以直接作为 CLI

当前 `ip_query_plugin` 是 Flask 服务，没有：

- `argparse` / `typer` / `click` CLI 入口。
- 包结构和安装入口，例如 `console_scripts`。
- 输出格式选项，例如 `json`、`table`、`csv`。
- 批量查询能力。
- 本地配置初始化命令。
- 标准退出码。

### 5.2 不足以直接作为 Skill

当前项目没有 Skill 所需结构：

```text
skill-name/
  SKILL.md
  scripts/
  references/
  assets/
  agents/
```

也缺少面向 Agent 的操作流程说明，例如：

- 什么时候调用 IP 查询。
- 如何处理用户给多个 IP。
- 如何处理私网 IP、非法 IP、域名、CIDR。
- 查询失败时如何向用户解释。
- 哪些字段可以直接展示，哪些字段需要脱敏。

### 5.3 工程风险

需要注意的问题：

- `server.py` 默认只监听 `127.0.0.1:8081`，生产部署需要配置 host/port。
- 请求只有 timeout，没有 retry/backoff。
- 上游响应基本原样返回，缺少统一错误模型。
- 没有单元测试。
- 没有批量查询限流。
- 没有敏感信息过滤。
- `demo2` 中有 API key 硬编码痕迹，应避免复用，并检查是否需要废弃或轮换。
- 没有依赖锁定，生产建议增加 `pyproject.toml` 或锁文件。

## 6. 建议的新项目架构

如果要把公司 IP 查询 API 封装成 `CLI + Skills`，建议采用下面架构。

```text
aw_cli/
  pyproject.toml
  README.md
  src/
    aw_ip_cli/
      __init__.py
      cli.py
      client.py
      config.py
      models.py
      output.py
      errors.py
  skills/
    aw-cli-query/
      SKILL.md
      scripts/
        query_ip.py
        batch_query_ip.py
      references/
        api.md
        response-fields.md
        errors.md
      agents/
        openai.yaml
  tests/
    test_client.py
    test_cli.py
    test_ip_validation.py
```

### 6.1 Python 包分层

推荐职责划分：

| 模块 | 职责 |
|---|---|
| `cli.py` | 命令定义、参数解析、调用 service/client、输出结果 |
| `client.py` | HTTP 请求、认证、重试、响应解析 |
| `config.py` | 环境变量、配置文件、默认值 |
| `models.py` | 请求和响应数据结构 |
| `output.py` | JSON/table/CSV 输出 |
| `errors.py` | 统一异常和退出码 |

### 6.2 CLI 命令建议

```bash
aw-cli query 8.8.8.8
aw-cli query 2001:4860:4860::8888 --coordsys WGS84 --format json
aw-cli batch ips.txt --output result.csv
aw-cli validate 8.8.8.8
aw-cli config init
aw-cli config show
aw-cli serve-plugin --host 127.0.0.1 --port 8081
```

建议输出格式：

- `--format json`：给 Agent 和自动化流程使用。
- `--format table`：给人类终端使用。
- `--format csv`：给批量分析使用。

建议退出码：

| 退出码 | 含义 |
|---|---|
| 0 | 成功 |
| 1 | 普通运行错误 |
| 2 | 参数错误或 IP 非法 |
| 3 | 配置缺失 |
| 4 | 上游 API 错误 |
| 5 | 网络超时 |

## 7. Skill 设计建议

### 7.1 Skill 目录

建议 skill 名称使用短横线小写：

```text
aw-cli-query/
  SKILL.md
  scripts/
    query_ip.py
    batch_query_ip.py
  references/
    api.md
    response-fields.md
    errors.md
  agents/
    openai.yaml
```

### 7.2 `SKILL.md` 应包含的核心内容

`SKILL.md` 不应该写成长篇 API 文档。它应该只放 Agent 执行任务时必须知道的流程：

```yaml
---
name: aw-cli-query
description: Query company IP intelligence APIs through the aw-cli CLI for IPv4/IPv6 geolocation, ISP, owner, ASN, risk, and related network attributes. Use when users ask to look up one or more IP addresses, analyze IP ownership/location, export IP lookup results, or validate IP query data.
---
```

正文建议包含：

- 优先使用 `scripts/query_ip.py` 或已安装的 `aw-cli` CLI。
- 单 IP 查询用 `aw-cli query <ip> --format json`。
- 多 IP 查询用 `aw-cli batch <file> --format json/csv`。
- 对私网 IP、保留地址、非法 IP 的处理规则。
- 查询失败时返回错误原因，不编造地理位置。
- 需要字段解释时读取 `references/response-fields.md`。
- 需要 API 细节时读取 `references/api.md`。

### 7.3 Skill 脚本与 CLI 的关系

建议让 Skill 脚本调用同一个 Python 包，而不是复制一份业务逻辑。

```text
scripts/query_ip.py
  -> import aw_ip_cli.client
  -> 调用公司 API
  -> 输出 JSON
```

这样 CLI、Skill、插件服务可以共用核心逻辑：

```text
核心 client.py
  |
  +-- CLI: aw-cli query
  +-- Skill script: scripts/query_ip.py
  +-- Plugin server: /query_ip
```

### 7.4 Skill 参考资料

建议把详细资料放到 `references/`：

```text
references/api.md
  - API base URL
  - 认证方式
  - IPv4/IPv6 endpoint
  - 请求参数
  - 限流规则

references/response-fields.md
  - country/province/city
  - isp/owner/asnumber
  - accuracy/source/radius
  - lat/lng/coordsys

references/errors.md
  - 常见错误码
  - 上游失败解释
  - 是否需要重试
```

这样可以保持 `SKILL.md` 简洁，Agent 只有在需要解释字段或排错时才读取更详细文档。

## 8. 插件服务是否还需要

如果你们目标是 `CLI + Skills`，Flask 插件服务不是必须的。但建议保留一个可选 `serve-plugin` 命令，因为它有三个价值：

1. 兼容百度 Agents、OpenAI 插件风格或其他需要 OpenAPI 的平台。
2. 给无法直接执行本地 CLI 的 Agent 提供 HTTP 工具接口。
3. 统一鉴权和审计，把 API key 留在服务端。

推荐架构：

```text
aw-cli CLI
  |
  +-- query / batch / config
  |
  +-- serve-plugin
        |
        +-- /.well-known/ai-plugin.json
        +-- /.well-known/openapi.yaml
        +-- /query_ip
```

## 9. 配置与安全建议

公司内部 API 建议支持三层配置：

1. 环境变量

```powershell
$env:AW_IP_API_BASE_URL="https://api.example.com"
$env:AW_IP_API_KEY="your_key"
```

2. 本地配置文件

```text
~/.aw-cli/config.toml
```

3. 命令行参数临时覆盖

```bash
aw-cli query 8.8.8.8 --base-url https://api.example.com
```

安全原则：

- 不在 OpenAPI、Skill、示例 YAML、README 中写真实 key。
- CLI 输出默认不打印 key、token、完整请求头。
- 日志中脱敏认证信息。
- 批量查询增加最大数量限制和限速。
- 私网 IP、保留 IP、环回地址默认不调用付费 API，除非显式 `--include-private`。

## 10. 可直接参考的迁移路线

### 第一步：抽出核心 Client

从 `ip_query_plugin/server.py` 中抽出：

- IP 校验。
- IPv4/IPv6 endpoint 选择。
- requests 调用。
- 响应解析和错误处理。

形成：

```text
src/aw_ip_cli/client.py
```

### 第二步：实现 CLI

先实现最小命令：

```bash
aw-cli query <ip> --format json
```

再补：

- `batch`
- `config`
- `validate`
- `serve-plugin`

### 第三步：增加 Skill

创建：

```text
skills/aw-cli-query/SKILL.md
skills/aw-cli-query/scripts/query_ip.py
skills/aw-cli-query/references/api.md
```

Skill 脚本调用已实现的 CLI 或 Python client。

### 第四步：恢复插件能力

把 `ip_query_plugin` 里的：

- `ai-plugin.json`
- `openapi.yaml`
- `example.yaml`
- `logo.png`

迁移为可选服务资源，由 `aw-cli serve-plugin` 提供。

### 第五步：测试和验收

建议测试覆盖：

- IPv4 查询。
- IPv6 查询。
- 非法 IP。
- 私网 IP。
- API key 缺失。
- 上游超时。
- 上游返回非 JSON。
- 批量查询 CSV/JSON 输出。
- CLI 退出码。

## 11. 最小可行版本建议

如果先做一个 MVP，不建议一开始做太多能力。最小版本可以是：

```text
aw_cli/
  pyproject.toml
  src/aw_ip_cli/
    cli.py
    client.py
    config.py
  skills/aw-cli-query/
    SKILL.md
    scripts/query_ip.py
    references/api.md
```

MVP 支持：

```bash
aw-cli query 8.8.8.8 --format json
aw-cli batch ips.txt --format csv
```

Skill 支持：

- 用户问单个 IP，调用 query。
- 用户给多个 IP，写临时文件后调用 batch。
- 用户问字段含义，读取 `references/api.md` 或 `response-fields.md`。

## 12. 结论

这个项目最值得参考的是“API 代理 + OpenAPI 描述 + 示例触发”的 Agent 工具化思路，以及 `server.py` 里用 `ipaddress` 自动区分 IPv4/IPv6 的实现方式。

但如果要封装成公司可长期维护的 `CLI + Skills`，建议不要直接在 Flask 项目上叠功能，而是先抽出通用 Python client，再分别挂接 CLI、Skill 脚本和可选插件服务。这样核心逻辑只有一份，后续接入 Codex、百度 Agents、内部平台或批处理任务都会更稳。

