# Score Server

## Running in production

```bash
docker run --rm -d -v ./certs:/certs -p 443:443 -p 80:80 "docker image location with hash"
```