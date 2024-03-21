# Score Server

## Running in production

```bash
VERSION=0.0.6
docker pull ghcr.io/soockee/terminal-games-scoreserver:${VERSION}-scoreserver
docker run --rm -d -v ./certs:/certs -p 13337:443 -p 13338:80 ghcr.io/soockee/terminal-games-scoreserver:${VERSION}-scoreserver
```