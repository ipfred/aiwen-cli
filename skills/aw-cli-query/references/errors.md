# Error Reference

## Exit Codes

| Code | Type | Description |
|---|---|---|
| 0 | Success | Command completed successfully |
| 1 | `internal` | Unexpected internal error |
| 2 | `validation` | Invalid IP, unsupported action for IP version, bad flag |
| 3 | `config` | Missing API key, malformed config, invalid timeout |
| 4 | `api_error` / `parse_error` | Upstream returned non-200 or non-JSON |
| 5 | `network` | Connection failure, timeout, DNS error |

## Error Response Format

```json
{
  "ok": false,
  "error": {
    "type": "validation",
    "message": "action risk only supports IPv4"
  }
}
```

## Common Errors

### Validation (exit 2)

- `invalid IP address: <ip>` — The provided string is not a valid IPv4 or IPv6 address.
- `action <name> only supports IPv4` — An IPv6 address was passed to an IPv4-only action.
- `action <name> does not support IPv4` — An IPv4 address was passed to an IPv6-only action.
- `invalid accuracy <level>; valid options are city, district, street` — Invalid accuracy level for `loc`.
- `unsupported format: <format>` — Invalid output format.
- `unsupported action: <action>` — Unknown action name.

### Config (exit 3)

- `AIWEN_API_KEY is required` — No API key found in flags, env, or config file.
- `failed to parse config file` — Malformed JSON in config.
- `invalid IPV4_ACCURACY` / `invalid IPV6_ACCURACY` — Accuracy value not in city/district/street.
- `invalid timeout duration` — Timeout string is not a valid Go duration.

### API Error (exit 4)

- `upstream API returned an error` — HTTP 4xx from upstream.
- `upstream returned non-JSON response` — Response body could not be parsed as JSON.

### Network (exit 5)

- `upstream server error` — HTTP 5xx from upstream.
- Network timeout, DNS failure, connection refused.

## Troubleshooting

1. **Missing API key**: Set it via `aw-cli config set api_key YOUR_KEY`, environment variable `AIWEN_API_KEY`, or `--api-key` flag.
2. **IPv6 on IPv4-only action**: Use IPv4 addresses with `whois`, `asn`, `host`, `risk`, `identity`, and `industry` commands.
3. **Network timeout**: Increase with `--timeout 30s` or `aw-cli config set timeout 30s`.
4. **Config file location**: Default is `~/.aw-cli/config.json`. Override with `--config /path/to/config.json`.
