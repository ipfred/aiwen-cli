---
name: aw-geoip
version: 1.1.0
description: "AIWEN/IPPlus360 IP geolocation query skill. Query IP location, ISP, risk, scene, WHOIS, ASN, identity, or industry via aw-cli. Only query the data the user actually asked for — each call is billable."
metadata:
  requires:
    bins: ["aw-cli"]
  cliHelp: "aw-cli --help"
---

# aw-geoip Skill

Use the `aw-cli` CLI to query IP intelligence. **Each CLI call is billable — only query the specific data the user asked for. Never guess or fabricate results.**

## Cost Warning

- Every `aw-cli` invocation triggers a billable API request (except `--dry-run`).
- **Only query what the user asked for.** If they ask "where is this IP", just call `loc` — don't also fetch `scene`, `risk`, and `whois`.
- **It's fine to call multiple commands** when the user genuinely asks for multiple data points (e.g., "check the risk and scene for 1.2.3.4"). Just don't blindly fetch everything.
- When unsure which data point the user needs, ask them rather than guessing and calling multiple commands.
- Use `--dry-run` to preview requests without consuming API quota.

## Intent-to-Command Mapping

Match the user's question to the right command(s):

| If the user asks about... | Run this |
|---|---|
| Where is this IP? Country, city, province, coordinates | `aw-cli loc <ip>` |
| My own public IP and its location | `aw-cli current` |
| Is it residential, datacenter, CDN, or Anycast? | `aw-cli scene <ip>` |
| Who owns this IP? Registration, net range, contacts | `aw-cli whois <ip>` |
| What AS does this IP belong to? AS number mapping | `aw-cli asn <ip>` |
| What ISP / organization / AS name is behind this IP? | `aw-cli host <ip>` |
| Is this IP risky? VPN, proxy, Tor, risk score | `aw-cli risk <ip>` |
| Is this real human traffic or a bot? | `aw-cli identity <ip>` |
| What industry does this IP belong to? | `aw-cli industry <ip>` |
| Query many IPs from a file with one action | `aw-cli batch <file> --action <action>` |

**Ambiguous queries**: If a user says "look up 1.2.3.4" or "check this IP", default to `aw-cli loc <ip>` (geolocation is the most commonly expected answer). If they ask about "security" or "reputation", use `aw-cli risk <ip>`.

## Command Reference

### Geolocation

```bash
aw-cli loc <ip>                        # city-level (default)
aw-cli loc <ip> --accuracy district    # district-level
aw-cli loc <ip> --accuracy street      # street-level (IPv4 only)
aw-cli loc <ip> --lang en              # English response
```

### Current IP

```bash
aw-cli current                         # detect public IP, then loc
aw-cli current --accuracy district
```

### Usage Scene

```bash
aw-cli scene <ip>
```

### IPv4-Only Commands

These commands support **IPv4 only**. IPv6 inputs return exit code 2.

```bash
aw-cli whois <ip>      # WHOIS registration info
aw-cli asn <ip>        # AS number mapping
aw-cli host <ip>       # ISP, organization, AS name
aw-cli risk <ip>       # VPN/proxy/Tor detection, risk score
aw-cli identity <ip>   # human vs bot classification
aw-cli industry <ip>   # industry classification
```

### Batch

```bash
aw-cli batch ips.txt --action loc --format ndjson
aw-cli batch ips.csv --ip-column ip --action risk --output results.csv --format csv
```

## Common Flags

| Flag | Purpose |
|---|---|
| `--format json/table/csv/ndjson` | Output format (default: json) |
| `--dry-run` | Preview request without consuming API quota |
| `--jq .data.country` | Filter JSON output with jq expression |
| `--timeout 30s` | Request timeout (default: 10s) |
| `--lang en/cn` | Response language |

## Setup

```bash
aw-cli config set api_key YOUR_KEY
# or
export AIWEN_API_KEY=YOUR_KEY
```

## Error Handling

| Exit Code | Meaning |
|---|---|
| 0 | Success |
| 1 | Internal error |
| 2 | Validation error (invalid IP, IPv6 on IPv4-only action) |
| 3 | Config error (missing API key) |
| 4 | API error (upstream returned error) |
| 5 | Network error (timeout, connection failure) |

## Rules Summary

1. **Match the ask, don't over-fetch.** Each CLI call is billable. Only query the data points the user actually asked for — don't blindly run all available commands against an IP.
2. **Never guess results.** Always call the CLI.
3. **Multiple commands are fine** when the user genuinely needs multiple data points (e.g., "tell me the risk and WHOIS for 1.2.3.4").
4. **IPv4-only actions** (`whois`, `asn`, `host`, `risk`, `identity`, `industry`) fail on IPv6 — warn the user before calling.
5. **Batch for multiple IPs.** Use `batch` instead of looping individual queries.
6. **Ask, don't assume.** If intent is unclear, ask the user which data they need rather than calling everything.
7. For field definitions, see `references/response-fields.md`.
8. For troubleshooting, see `references/errors.md`.
