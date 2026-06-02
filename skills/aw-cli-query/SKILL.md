---
name: aw-cli-query
version: 1.0.0
description: "AIWEN/IPPlus360 IP intelligence query skill. Query IPv4/IPv6 geolocation, current network IP, usage scene, WHOIS, AS mapping, host info, risk portrait, identity check, and industry classification via the aw-cli CLI. Use when users ask to look up IP location, ISP, ownership, risk scoring, bot/human detection, or batch IP intelligence."
metadata:
  requires:
    bins: ["aw-cli"]
  cliHelp: "aw-cli --help"
---

# aw-cli-query Skill

Use the `aw-cli` CLI to query IP intelligence. Never guess or fabricate IP query results.

## Quick Reference

| User Intent | Command |
|---|---|
| IP geolocation (city, district, street) | `aw-cli loc <ip>` |
| Current machine's public IP location | `aw-cli current` |
| IP usage scene (residential, datacenter, CDN) | `aw-cli scene <ip>` |
| IP WHOIS registration info | `aw-cli whois <ip>` |
| AS number mapping | `aw-cli asn <ip>` |
| AS name, ISP, organization | `aw-cli host <ip>` |
| VPN, proxy, Tor, risk score | `aw-cli risk <ip>` |
| Real human vs bot traffic | `aw-cli identity <ip>` |
| Industry classification | `aw-cli industry <ip>` |
| Multiple IPs from file | `aw-cli batch <file> --action <action>` |

## Required Setup

Set your API key before first use:

```bash
aw-cli config set api_key YOUR_KEY
# or
export AIWEN_API_KEY=YOUR_KEY
```

## Common Flags

| Flag | Purpose |
|---|---|
| `--format json/table/csv/ndjson` | Output format (default: json) |
| `--dry-run` | Preview request without calling upstream |
| `--jq .data.country` | Filter JSON output |

## IPv4-Only Actions

`whois`, `asn`, `host`, `risk`, `identity`, and `industry` only support IPv4. Passing an IPv6 address returns a validation error.

## Batch Query

```bash
aw-cli batch ips.txt --action loc --format ndjson
aw-cli batch ips.csv --ip-column ip --action risk --output results.csv --format csv
```

## Error Handling

| Exit Code | Meaning |
|---|---|
| 0 | Success |
| 1 | Internal error |
| 2 | Validation error (invalid IP, unsupported action) |
| 3 | Config error (missing API key) |
| 4 | API error (upstream returned error) |
| 5 | Network error (timeout, connection failure) |

## Important Rules

- Always call the CLI; never guess IP query results.
- For multiple IPs, use `batch` instead of sequential queries.
- IPv4-only actions will fail on IPv6 inputs — explain this to the user.
- For field definitions, read `references/response-fields.md`.
- For API details, read `references/api.md`.
- For error details, read `references/errors.md`.
