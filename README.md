# ![Logo](docs/logos/readme.png)

[![License][1]][2] [![Release][3]][4] [![Telegram][5]][6]

[1]: https://img.shields.io/github/license/tobyxdd/hysteria?style=flat-square

[2]: LICENSE.md

[3]: https://img.shields.io/github/v/release/tobyxdd/hysteria?style=flat-square

[4]: https://github.com/tobyxdd/hysteria/releases

[5]: https://img.shields.io/badge/chat-Telegram-blue?style=flat-square

[6]: https://t.me/hysteria_github

[中文](README.zh.md)

Hysteria is a feature-packed network utility optimized for networks of poor quality (e.g. satellite connections,
congested public Wi-Fi, connecting from China to servers abroad) powered by a custom version of QUIC protocol. It
currently has the following features: (still growing!)

- SOCKS5 proxy (TCP & UDP)
- HTTP/HTTPS proxy
- TCP/UDP relay
- TCP/UDP TPROXY (Linux only)
- TUN (TAP on Windows)

## Installation

### Windows, Linux, macOS CLI

- Download pre-built binaries from https://github.com/tobyxdd/hysteria/releases
  - Linux builds are available as `hysteria` (with tun support) and `hysteria-notun` (without tun support). Builds 
    without tun support are statically linked and do not depend on glibc. **If you use a non-standard distribution that 
    can't run `hysteria` properly, try `hysteria-notun` instead.**
- Use Docker or Docker Compose: https://github.com/HyNetwork/hysteria/blob/master/Docker.md
- Use our Arch Linux AUR: https://aur.archlinux.org/packages/hysteria/
- Build from source with `go build ./cmd`

### OpenWrt LuCI app

- [openwrt-passwall](https://github.com/xiaorouji/openwrt-passwall)

### Android

- [SagerNet](https://github.com/SagerNet/SagerNet) with [hysteria-plugin](https://github.com/SagerNet/SagerNet/releases/tag/hysteria-plugin-0.9.1)

### iOS

- [Shadowrocket](https://apps.apple.com/us/app/shadowrocket/id932747118)

## Quick Start

Note: This is only a bare-bones example to get the server and client running. Go to [Advanced usage](#advanced-usage)
for all the available options.

### Server

Create a `config.json` under the root directory of the program:

```json
{
  "listen": ":36712",
  "acme": {
    "domains": [
      "your.domain.com"
    ],
    "email": "hacker@gmail.com"
  },
  "obfs": "fuck me till the daylight",
  "up_mbps": 100,
  "down_mbps": 100
}
```

Hysteria requires a TLS certificate. You can either get a trusted TLS certificate from Let's Encrypt automatically using
the built-in ACME integration, or provide it yourself. It does not have to be valid and trusted, but in that case the
clients need additional configuration. To use your own existing TLS certificate, refer to this config:

```json
{
  "listen": ":36712",
  "cert": "/home/ubuntu/my.crt",
  "key": "/home/ubuntu/my.key",
  "obfs": "fuck me till the daylight",
  "up_mbps": 100,
  "down_mbps": 100
}
```

The (optional) `obfs` option obfuscates the protocol using the provided password, so that it is not apparent that this
is Hysteria/QUIC, which could be useful for bypassing DPI blocking or QoS. If the passwords of the server and client do
not match, no connection can be established. Therefore, this can also serve as a simple password authentication. For
more advanced authentication schemes, see `auth` below.

`up_mbps` and `down_mbps` limit the maximum upload and download speed of the server for each client. Feel free to remove
them if you don't need.

To launch the server, simply run

```
./hysteria-linux-amd64 server
```

If your config file is not named `config.json` or is in a different path, specify it with `-config`:

```
./hysteria-linux-amd64 -config blah.json server
```

### Client

Same as the server side, create a `config.json` under the root directory of the program:

```json
{
  "server": "example.com:36712",
  "obfs": "fuck me till the daylight",
  "up_mbps": 10,
  "down_mbps": 50,
  "socks5": {
    "listen": "127.0.0.1:1080"
  },
  "http": {
    "listen": "127.0.0.1:8080"
  }
}
```

This config enables a SOCKS5 proxy (with both TCP & UDP support), and an HTTP proxy at the same time. There are many
other modes in Hysteria, be sure to check them out in [Advanced usage](#advanced-usage)! To enable or disable a mode,
simply add or remove its entry in the config file.

If your server certificate is not issued by a trusted CA, you need to specify the CA used
with `"ca": "/path/to/file.ca"` on the client or use `"insecure": true` to ignore all certificate errors (not
recommended).

`up_mbps` and `down_mbps` are mandatory on the client side. Try to fill in these values as accurately as possible
according to your network conditions, as they are crucial for Hysteria to work optimally.

Some users may attempt to forward other encrypted proxy protocols such as Shadowsocks with relay. While this technically
works, it's not optimal from a performance standpoint - Hysteria itself uses TLS, considering that the proxy protocol
being forwarded is also encrypted, and the fact that almost all sites are now using HTTPS, it essentially becomes triple
encryption. If you need a proxy, just use our proxy modes.

## Comparison

![Bench](docs/bench/bench.png)

## Advanced usage

### Server

```json5
{
  "listen": ":36712", // Listen address
  "protocol": "faketcp", // Blank or "udp" for UDP mode, "faketcp" for TCP "masquerade", see below for details
  "acme": {
    "domains": [
      "your.domain.com",
      "another.domain.net"
    ], // Domains for the ACME cert
    "email": "hacker@gmail.com", // Registration email, optional but recommended
    "disable_http": false, // Disable HTTP challenges
    "disable_tlsalpn": false, // Disable TLS-ALPN challenges
    "alt_http_port": 8080, // Alternate port for HTTP challenges
    "alt_tlsalpn_port": 4433 // Alternate port for TLS-ALPN challenges
  },
  "cert": "/home/ubuntu/my_cert.crt", // Cert file, mutually exclusive with the ACME options above
  "key": "/home/ubuntu/my_key.crt", // Key file, mutually exclusive with the ACME options above
  "up_mbps": 100, // Max upload Mbps per client
  "down_mbps": 100, // Max download Mbps per client
  "disable_udp": false, // Disable UDP support
  "acl": "my_list.acl", // See ACL below
  "obfs": "AMOGUS", // Obfuscation password
  "auth": { // Authentication
    "mode": "password", // Mode, supports "password" "none" and "external" for now
    "config": {
      "password": "yubiyubi"
    }
  },
  "alpn": "ayaya", // QUIC TLS ALPN
  "prometheus_listen": ":8080", // Prometheus HTTP metrics server listen address (at /metrics)
  "recv_window_conn": 15728640, // QUIC stream receive window
  "recv_window_client": 67108864, // QUIC connection receive window
  "max_conn_client": 4096, // Max concurrent connections per client
  "disable_mtu_discovery": false, // Disable Path MTU Discovery (RFC 8899)
  "ipv6_only": false, // Only resolve domains to IPv6 address
  "resolver": "1.1.1.1:53" // DNS resolver address
}
```

#### ACME

Only HTTP and TLS-ALPN challenges are currently supported (no DNS challenges). Make sure your TCP ports 80/443 are
accessible respectively.

#### External authentication integration

If you are a commercial proxy provider, you may want to connect Hysteria to your own authentication backend.

```json5
{
  // ...
  "auth": {
    "mode": "external",
    "config": {
      "http": "https://api.example.com/auth" // Both HTTP and HTTPS are supported
    }
  }
}
```

For the above config, Hysteria sends a POST request to `https://api.example.com/auth` upon each client's connection:

```json5
{
  "addr": "111.222.111.222:52731",
  "payload": "[BASE64]", // auth or auth_str of the client
  "send": 12500000, // Negotiated server send speed for this client (Bps)
  "recv": 12500000 // Negotiated server recv speed for this client (Bps)
}
```

The endpoint must return results with HTTP status code 200 (even if the authentication failed):

```json5
{
  "ok": false,
  "msg": "No idea who you are"
}
```

`ok` indicates whether the authentication passed. `msg` is a success/failure message.

#### Prometheus Metrics

You can make Hysteria expose a Prometheus HTTP client endpoint for monitoring traffic usage with `prometheus_listen`.
If configured on port 8080, the endpoint would be at `http://example.com:8080/metrics`.

```text
hysteria_active_conn{auth="55m95auW5oCq"} 32
hysteria_active_conn{auth="aGFja2VyISE="} 7

hysteria_traffic_downlink_bytes_total{auth="55m95auW5oCq"} 122639
hysteria_traffic_downlink_bytes_total{auth="aGFja2VyISE="} 3.225058e+06

hysteria_traffic_uplink_bytes_total{auth="55m95auW5oCq"} 40710
hysteria_traffic_uplink_bytes_total{auth="aGFja2VyISE="} 37452
```

`auth` is the auth payload sent by the clients, encoded in Base64.

### Client

```json5
{
  "server": "example.com:36712", // Server address
  "protocol": "faketcp", // Blank or "udp" for UDP mode, "faketcp" for TCP "masquerade", see below for details
  "up_mbps": 10, // Max upload Mbps
  "down_mbps": 50, // Max download Mbps
  "socks5": {
    "listen": "127.0.0.1:1080", // SOCKS5 listen address
    "timeout": 300, // TCP timeout in seconds
    "disable_udp": false, // Disable UDP support
    "user": "me", // SOCKS5 authentication username
    "password": "lmaolmao" // SOCKS5 authentication password
  },
  "http": {
    "listen": "127.0.0.1:8080", // HTTP listen address
    "timeout": 300, // TCP timeout in seconds
    "user": "me", // HTTP authentication username
    "password": "lmaolmao", // HTTP authentication password
    "cert": "/home/ubuntu/my_cert.crt", // Cert file (HTTPS proxy)
    "key": "/home/ubuntu/my_key.crt" // Key file (HTTPS proxy)
  },
  "tun": {
    "name": "tun-hy", // TUN interface name
    "timeout": 300, // Timeout in seconds
    "address": "192.0.2.2", // TUN interface address, not applicable for Linux
    "gateway": "192.0.2.1", // TUN interface gateway, not applicable for Linux
    "mask": "255.255.255.252", // TUN interface mask, not applicable for Linux
    "dns": [ "8.8.8.8", "8.8.4.4" ], // TUN interface DNS, only applicable for Windows
    "persist": false // Persist TUN interface after exit, only applicable for Linux
  },
  "relay_tcps": [
    {
      "listen": "127.0.0.1:2222", // TCP relay listen address
      "remote": "123.123.123.123:22", // TCP relay remote address
      "timeout": 300 // TCP timeout in seconds
    },
    {
      "listen": "127.0.0.1:13389", // TCP relay listen address
      "remote": "124.124.124.124:3389", // TCP relay remote address
      "timeout": 300 // TCP timeout in seconds
    }
  ],
  "relay_udps": [
    {
      "listen": "127.0.0.1:5333", // UDP relay listen address
      "remote": "8.8.8.8:53", // UDP relay remote address
      "timeout": 60 // UDP session timeout in seconds
    },
    {
      "listen": "127.0.0.1:11080", // UDP relay listen address
      "remote": "9.9.9.9.9:1080", // UDP relay remote address
      "timeout": 60 // UDP session timeout in seconds
    }
  ],
  "tproxy_tcp": {
    "listen": "127.0.0.1:9000", // TCP TProxy listen address
    "timeout": 300 // TCP timeout in seconds
  },
  "tproxy_udp": {
    "listen": "127.0.0.1:9000", // UDP TProxy listen address
    "timeout": 60 // UDP session timeout in seconds
  },
  "acl": "my_list.acl", // See ACL below
  "obfs": "AMOGUS", // Obfuscation password
  "auth": "[BASE64]", // Authentication payload in Base64
  "auth_str": "yubiyubi", // Authentication payload in string, mutually exclusive with the option above
  "alpn": "ayaya", // QUIC TLS ALPN
  "server_name": "real.name.com", // TLS hostname used to verify the server certificate
  "insecure": false, // Ignore all certificate errors 
  "ca": "my.ca", // Custom CA file
  "recv_window_conn": 15728640, // QUIC stream receive window
  "recv_window": 67108864, // QUIC connection receive window
  "disable_mtu_discovery": false, // Disable Path MTU Discovery (RFC 8899)
  "resolver": "1.1.1.1:53" // DNS resolver address
}
```

#### Fake TCP / TCP masquerade

Certain networks may impose various restrictions on UDP traffic or block it altogether. Hysteria offers a "faketcp" mode
that allows servers and clients to communicate using a protocol that looks like TCP but does not actually go through the
system TCP stack. This tricks whatever middleboxes into thinking it's actually TCP traffic, rendering UDP-specific
restrictions useless.

This mode is currently only supported on Linux (both client and server) and requires root privileges.

If your server is behind a firewall, open the corresponding TCP port instead of UDP.

#### Transparent proxy

TPROXY modes (`tproxy_tcp` & `tproxy_udp`) are only available on Linux.

References:
- https://www.kernel.org/doc/Documentation/networking/tproxy.txt
- https://powerdns.org/tproxydoc/tproxy.md.html

## Optimization tips

### Optimizing for extreme transfer speeds

If you want to use Hysteria for very high speed transfers (e.g. 10GE, 1G+ over inter-country long fat pipes), consider
increasing your system's UDP receive buffer size.

```shell
sysctl -w net.core.rmem_max=4000000
```

This would increase the buffer size to roughly 4 MB on Linux.

You may also need to increase `recv_window_conn` and `recv_window` (`recv_window_client` on server side) to make sure
they are at least no less than the bandwidth-delay product. For example, if you want to achieve a transfer speed of 500
MB/s on a line with an RTT of 200 ms, you need a minimum receive window size of 100 MB (500*0.2).

### Routers and other embedded devices

For devices with very limited computing power and RAM, turning off obfuscation can bring a slight performance boost.

The default receive window size for both Hysteria server and client is 64 MB. Consider lowering them if it's too large
for your device. Keeping a ratio of one to four between stream receive window and connection receive window is
recommended.

## ACL

[ACL File Format](ACL.md)

ACL is available on both client & server. On the server side it can be used to restrict what the clients can access, and
is valid for any mode on the client side. On the client side, it's only supported in SOCKS5 & HTTP proxy modes, and has
no effect in other modes (all traffic will go through the proxy)

## URI Scheme

Third party clients looking to implement a "share by link" feature are advised to follow the following URI scheme 
(initially introduced by Shadowrocket):

    hysteria://host:port?protocol=udp&auth=123456&peer=sni.domain&insecure=1&upmbps=100&downmbps=100&alpn=hysteria&obfs=xplus&obfsParam=123456#remarks

    - host: hostname or IP address of the server to connect to (required)
    - port: port of the server to connect to (required)
    - protocol: protocol to use ("udp" or "faketcp") (optional, default: "udp")
    - auth: authentication payload (string) (optional)
    - peer: SNI for TLS (optional)
    - insecure: ignore certificate errors (optional)
    - upmbps: upstream bandwidth in Mbps (required)
    - downmbps: downstream bandwidth in Mbps (required)
    - alpn: QUIC ALPN (optional)
    - obfs: Obfuscation mode (optional, empty or "xplus")
    - obfsParam: Obfuscation password (optional)
    - remarks: remarks (optional)

## Logging

The program outputs `DEBUG` level, text format logs via stdout by default.

To change the logging level, use `LOGGING_LEVEL` environment variable. The available levels are `panic`, `fatal`
, `error`, `warn`, `info`, ` debug`, `trace`

To print JSON instead, set `LOGGING_FORMATTER` to `json`

To change the logging timestamp format, set `LOGGING_TIMESTAMP_FORMAT`

## Custom CA

  1. Suppose the server address is `123.123.123.123`, UDP port `5678` is not blocked by firewall
  2. openssl is already installed
  3. hysteria is already installed in `/root/hysteria/` directory
<details>
  <summary>4. Generate custom CA certificate</summary>

- Run below shell in `/root/hysteria/` folder

``` shell
#!/usr/bin/env bash

domain=$(openssl rand -hex 8)
password=$(openssl rand -hex 16)
obfs=$(openssl rand -hex 6)
path="/root/hysteria"

openssl genrsa -out hysteria.ca.key 2048

openssl req -new -x509 -days 3650 -key hysteria.ca.key -subj "/C=CN/ST=GD/L=SZ/O=Hysteria, Inc./CN=Hysteria Root CA" -out hysteria.ca.crt

openssl req -newkey rsa:2048 -nodes -keyout hysteria.server.key -subj "/C=CN/ST=GD/L=SZ/O=Hysteria, Inc./CN=*.${domain}.com" -out hysteria.server.csr

openssl x509 -req -extfile <(printf "subjectAltName=DNS:${domain}.com,DNS:www.${domain}.com") -days 3650 -in hysteria.server.csr -CA hysteria.ca.crt -CAkey hysteria.ca.key -CAcreateserial -out hysteria.server.crt

cat > ./client.json <<EOF
{
    "server": "123.123.123.123:5678",
    "alpn": "h3",
    "obfs": "${obfs}",
    "auth_str": "${password}",
    "up_mbps": 30,
    "down_mbps": 30,
    "socks5": {
        "listen": "0.0.0.0:1080"
    },
    "http": {
        "listen": "0.0.0.0:8080"
    },
    "server_name": "www.${domain}.com",
    "ca": "${path}/hysteria.ca.crt"
}
EOF


cat > ./server.json <<EOF
{
    "listen": ":5678",
    "alpn": "h3",
    "obfs": "${obfs}",
    "cert": "${path}/hysteria.server.crt",
    "key": "${path}/hysteria.server.key" ,
    "auth": {
        "mode": "password",
        "config": {
            "password": "${password}"
        }
    }
}
EOF
```
</details>

5. Server side: copy `server.json`、 `hysteria.server.crt`、 `hysteria.server.key` to `/root/hysteria/` directory, run `/root/hysteria/hysteria -c /root/hysteria/server.json server` command

6. Client side: Assuming that the client directory is also`/root/hysteria`, copy `client.json`、`hysteria.ca.crt` to `/root/hysteria/` directory, run `/root/hysteria/hysteria -c /root/hysteria/client.json` cmmand

7. After generating CA certificate, modify the server address, port and certificate file path according to your own situation, add obfs and alpn to prevent the first time to be walled in some environment, after the first test passed in full parameters, you can remove the unnecessary parameters such as obfs and alpn in your own network environment.

8. If you are using shadowrocket on IOS, you can airdrop the file `hysteria.ca.crt` to your iPhone and install it, then you can use custom CA certificate.
