# Response Fields

## Geolocation (`loc`)

| Field | Description |
|---|---|
| `continent` | Continent name |
| `country` | Country name |
| `province` | Province/state |
| `city` | City name |
| `district` | District/county |
| `street` | Street-level detail (street accuracy) |
| `lng` | Longitude |
| `lat` | Latitude |
| `radius` | Accuracy radius in meters |
| `owner` | Network owner/organization |
| `isp` | Internet Service Provider |
| `asnumber` | AS number |
| `accuracy` | Accuracy level queried |

## Usage Scene (`scene`)

| Field | Description |
|---|---|
| `scene` | Usage scene: residential, datacenter, CDN, etc. |
| `is_residential` | Whether IP is residential |
| `is_datacenter` | Whether IP belongs to a data center |
| `is_cdn` | Whether IP is a CDN node |
| `is_anycast` | Whether IP uses Anycast |

## IPv4-Only Actions

### WHOIS (`whois`)

| Field | Description |
|---|---|
| `asnumber` | AS number |
| `org` | Organization |
| `net` | Network range |
| `country` | Country |
| `tech_contact` | Technical contact |
| `admin_contact` | Administrative contact |

### AS Mapping (`asn`)

| Field | Description |
|---|---|
| `asnumber` | AS number |
| `asname` | AS name |
| `org` | Organization |

### Host Info (`host`)

| Field | Description |
|---|---|
| `asnumber` | AS number |
| `asname` | AS name |
| `isp` | ISP name |
| `org` | Owning organization |

### Risk Portrait (`risk`)

| Field | Description |
|---|---|
| `is_vpn` | VPN detected |
| `is_proxy` | Proxy detected |
| `is_tor` | Tor exit node |
| `is_datacenter` | Data center IP |
| `risk_score` | Risk score (0-100) |

### Identity (`identity`)

| Field | Description |
|---|---|
| `is_human` | Probability of real human |
| `is_bot` | Probability of bot traffic |
| `second_use` | Second-use probability |

### Industry (`industry`)

| Field | Description |
|---|---|
| `industry` | Industry classification |
| `sub_industry` | Sub-industry category |
