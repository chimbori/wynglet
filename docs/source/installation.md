---
title: Installation & Deployment
description: How to install and deploy Wynglet
template: app-page
---

# Installation & Deployment

Wynglet is designed to be easy to deploy and can run on any server or container platform.

## Container Deployment (Recommended)

We strongly recommend deploying using the official container image, which includes Chrome Headless for convenience.

### Docker Compose Example

Create a `compose.yml` file with the following content:

```yml
services:
  wynglet:
    container_name: wynglet
    image: ghcr.io/chimbori/wynglet:latest
    volumes:
      - $PWD/wynglet-data:/data
    restart: unless-stopped
    depends_on:
      - wynglet-db

  wynglet-db:
    container_name: wynglet-db
    image: postgres:18-alpine
    environment:
      POSTGRES_DB: wynglet
      POSTGRES_USER: chimbori
      POSTGRES_PASSWORD: chimbori
    volumes:
      - $PWD/wynglet-db-data:/var/lib/postgresql
    restart: unless-stopped

volumes:
  wynglet-data:
  wynglet-db-data:
```

Then deploy with:

```shell
docker-compose up -d
```

## Binary Installation

If you prefer to install using the Go binary on raw metal, use `go install`:

```shell
go install wynglet.chimbori.dev@latest
```

## Reverse Proxy Setup

Wynglet is designed to be used behind a TLS reverse proxy for SSL termination (among other things).
We recommend [Caddy](https://caddyserver.com/).

### Sample Caddyfile

```caddy
wynglet.your-server.com {
  reverse_proxy wynglet:9999
  encode zstd gzip
}
```

If you expect a lot of traffic, consider putting a CDN in front.

## Configuration

Wynglet requires basic configuration to be provided via a `wynglet.yml` config file.

### Minimal Configuration

```yml
# PostgreSQL Database URL (required)
database:
  url: postgresql://chimbori:chimbori@wynglet-db:5432/wynglet

# Dashboard credentials (encrypted via bcrypt) (required)
dashboard:
  username: admin
  password: "$2a$10$a8LnUkK1UiB.9yQrUp3wyuGsH1AAHhlHVy1cjIaaIUVAwCtGvaX7q" # "test"
```

### Generate a Dashboard Password

Use the bcrypt utility to generate a password hash:

```shell
wynglet --bcrypt
```

Enter your desired password when prompted.
The output hash can be pasted directly into the `dashboard.password` field.

### Optional Configuration

**Web config:**

Assuming there’s a reverse proxy in front of Wynglet, there should be no need to change the port here.
Just configure the reverse proxy to forward requests to port 9999.

```yml
web:
  port: 9999
```

**Link Previews config:**

```yml
link-previews:
  screenshot:
    timeout: 20s
  sitemap:
    concurrent_urls: 4  # Number of URLs processed in parallel during sitemap import
    max_urls: 1000      # Maximum URLs fetched from a sitemap
  cache:
    ttl: 720h0m0s
    max_size_bytes: 1073741824
```

**QR Codes config:**

```yml
qr-codes:
  cache:
    ttl: 720h0m0s
    max_size_bytes: 1073741824
```

**Ratings config:**

```yaml
ratings:
  retention: 8760h  # How long to keep rating data (default: 365 days)
```

**Logs config:**

```yml
logs:
  retention: 720h  # How long to keep error logs (default: 30 days)
  pagination:
    limit: 10
```

**Dashboard pagination:**

```yml
dashboard:
  pagination:
    limit: 10
```

**Debug Mode:**

Turn on additional logging in Debug Mode.
Debug mode can also be toggled per-domain from within the Dashboard UI.

```yml
debug: true
```

## Health Check

Use the health check utility for container health checks:

```shell
wynglet --healthcheck
```

This pings the running service and exits `0` on success or `1` on failure.

## Next Steps

Once Wynglet is installed and running:

1. Access the [Dashboard](/dashboard) at `https://wynglet.your-server.com/dashboard`
2. [Authorize your domains](/dashboard)
3. [Start here](/) with one of the features
