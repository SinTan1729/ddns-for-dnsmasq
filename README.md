# DDNS for Dnsmasq
This is a tool to automatically update the hosts defined in an additional hosts file for `Dnsmasq`.

# Installation
Use the provided `compose.yaml` file to install the Docker image. The options are described inside it.

The hostfile denoted by `HOSTFILE_PATH` should be a valid hostfile readable by `Dnsmasq`, recommended
location is inside a shared directory like `/srv/hostfiles`. Add `hostsdir=/srv/hostfiles` to
`/etc/dnsmasq.conf` so that `Dnsmasq` automatically reads the file and updated its state.

# Config
Take a look at `example_config.yaml`. The entries `server1`, `server2` etc. are the hostnames used in the
hostfile. The `api-key` entries are `Argon2id` hashed of the intended API keys.  
Recommended way of generating API keys:
```bash
tr -dc A-Za-z0-9 </dev/urandom | head -c 64 
```
Make sure to put them in quotes, or validation would not work properly.

If using any reverse-proxy, please set `ip-header` accordingly. If you use a custom port, update your
docker compose file accordingly.

# Usage
## `/update`
The main endpoint is the `/update` path. Use it like this.
```bash
curl -X PUT -H 'X-API-Key: <api-key>' -d '{"host":"<hostname>"}' <your-server>/update
```
You can optionally add `"ip":"<some-ip>"` to the payload to ask the server to use a specific IP address.
By default, it tries to detect your IP and use that.

The reply will be in `json` of either of two types.
```json
{
  "hsot": "<hostname>",
  "ip": "<used-ip>"
}
```
or
```json
{
  "error": true,
  "reason": "<some-reason>"
```

## `/whoami`
It will simply reply with the detected IP address. Use it like this.
```bash
curl <your-server>/whoami
```

## `getinfo`
It will reply with the info stored about a specific host. Use it like this.
```bash
curl -X POST -H 'X-API-Key: <api-key>' -d '{"host":"<hostname>"}' <your-server>/getinfo
```
The reply will look exactly like that of `/update`.

## `version`
It will simply reply with the server version. Use it like this.
```bash
curl <your-server>/version
```
