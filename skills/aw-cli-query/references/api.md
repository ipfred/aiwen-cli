# API Reference

## Base URL

`https://api.ipplus360.com`

## Authentication

All requests require the `key` query parameter with your API key.

## Common Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `key` | string | Yes | API access key |
| `channel` | string | Yes | Channel identifier (default: `aw_cli`) |
| `ip` | string | Yes* | IP address to query (*not needed for `current`) |
| `coordsys` | string | No | Coordinate system: WGS84, GCJ02, BD09 |
| `lang` | string | No | Response language: cn, en |
| `accuracy` | string | No | Location accuracy: city, district, street |

## Endpoints

### IP Geolocation (`loc`)

| IP Type | Accuracy | Endpoint |
|---|---|---|
| IPv4 | city | `/ip/geo/v1/city/` |
| IPv4 | district | `/ip/geo/v1/district/` |
| IPv4 | street | `/ip/geo/v1/street/psi/` |
| IPv6 | city | `/ip/geo/v1/ipv6/` |
| IPv6 | district | `/ip/geo/v1/ipv6/district/` |
| IPv6 | street | `/ip/geo/v1/ipv6/street/biz/` |

### Current IP (`current`)

First calls `GET https://www.ipuu.net/ipuu/user/getIP` to detect public IP, then queries `loc`.

### Usage Scene (`scene`)

| IP Type | Endpoint |
|---|---|
| IPv4 | `/ip/info/v1/scene/` |
| IPv6 | `/ip/info/v1/ipv6Scene/` |

### IPv4-Only Endpoints

| Action | Endpoint |
|---|---|
| whois | `/ip/info/v1/ipWhois` |
| asn | `/as/info/v1/asWhois` |
| host | `/ip/geo/v1/host/` |
| risk | `/ip/info/v3/portrait/` |
| identity | `/ip/info/v1/person/` |
| industry | `/ip/info/v1/industry/` |

## Rate Limits

Contact AIWEN/IPPlus360 for rate limit details for your API plan.
