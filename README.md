# proxz
Socket proxy with zstd compression

Proxz behaves roughly like socat with built-in zstd compression on either the left or the right side of the connection. TCP and Unix domain sockets are supported. Proxz favors speed over compression ratio by using zstd's fastest compression level.

Proxz uses [DataDog's Zstd wrapper for Go](https://github.com/DataDog/zstd) which wraps the original zstd C library. Building requires `gcc` and the compiled binary requires `glibc`.

## Usage

```
proxz [-c|--compress] <(tcp|unix):address> <(tcp|unix):address>
```

## Examples

```bash
# Tunnel SSH through proxz
# Server
proxz tcp::22022 tcp:127.0.0.1:22

# Client
proxz -c tcp:127.0.0.1:22022 tcp:server:22022 &
ssh -p 22022 user@127.0.0.1

# Tunnel MySQL through proxz
# Server
proxz tcp::33066 unix:/var/lib/mysql/mysql.sock

# Client
proxz -c unix:/opt/mysql.sock tcp:server:33066 &
mysql -S /opt/mysql.sock
```
