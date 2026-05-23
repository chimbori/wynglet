---
title: GitHub Repository Stats
description: Embed live GitHub repository statistics on your website
template: app-page
---

# GitHub Repository Stats

Embed live GitHub repository statistics (stars, forks, watchers, etc.) on your website without authentication.

## Overview

The GitHub Stats API allows you to fetch and display current information about any public GitHub repository.

Features:

- **No Authentication Required**: Works with public repositories
- **Cached Responses**: Results are cached for 1 hour to minimize API calls
- **CORS Enabled**: Can be called from any website
- **Multiple Fields**: Access stars, forks, watchers, open issues, and more
- **Lightweight**: Small, fast responses ideal for embedding

## Getting Started

### Fetch Repository Stats

Make a GET request to fetch repository statistics:

```text
GET https://wynglet.your-server.com/github/v1/{user}/{repo}/{field}
```

**URL Parameters:**

- `{user}` — GitHub username (e.g., `chimbori`)
- `{repo}` — Repository name (e.g., `wynglet`)
- `{field}` — What to fetch (see Available Fields below)

### Available Fields

- `name` — Repository name
- `description` — Repository description
- `stars` — Number of stars (stargazers_count)
- `forks` — Number of forks
- `issues` — Number of open issues
- `watchers` — Number of watchers (subscribers_count)

### Examples

**Get Star Count:**

```text
https://wynglet.your-server.com/github/v1/chimbori/wynglet/stars
```

Response:

```json
{
  "stars": 42
}
```

**Display in HTML:**

```html
<p>
  Wynglet has
  <span id="stars">--</span>
  stars on GitHub
</p>

<script>
fetch('https://wynglet.your-server.com/github/v1/chimbori/wynglet/stars')
  .then(r => r.json())
  .then(data => {
    document.getElementById('stars').textContent = data.stars.toLocaleString();
  });
</script>
```

## CORS

The GitHub Stats endpoint is fully CORS-enabled and can be called directly from any website.
No preflight requests are needed.

All requests receive permissive CORS headers:

```text
Access-Control-Allow-Origin: *
Access-Control-Allow-Headers: *
Access-Control-Allow-Methods: GET, OPTIONS
```

## Caching

Results are cached for 1 hour to minimize requests to GitHub’s API.
If you need fresh data before the cache expires, you can:

1. Wait for the cache to expire (1 hour)
2. Deploy a new instance with cleared caches
3. Query multiple times to get the most recent data

## Error Handling

If a request fails (invalid repository, network error, etc.), the API returns an error status with details.

**Example Error Response:**

```json
{
  "error": "Repository not found"
}
```

## Use Cases

- **Portfolio Sites**: Show your most popular projects
- **Documentation**: Display project health and community stats
- **Team Pages**: Showcase your open-source projects
- **Dashboards**: Monitor GitHub activity across repositories

## Limitations

- Only works with **public repositories** (GitHub API limitation)
- Results are **cached for 1 hour**, not real-time
- **Rate-limited** by GitHub’s public API limits (60 requests/hour per IP)

## Configuration

The GitHub stats endpoint requires no special configuration.
It uses a local disk cache to minimize GitHub API requests.

```yml
github:
  # Cache TTL is fixed at 1 hour for GitHub stats
  cache:
    ttl: 1h
    location: cache/github
```

## Best Practices

- **Cache on Your Side Too**: Don’t fetch on every page load; cache results in your site’s cache or in JavaScript sessionStorage
- **Handle Errors Gracefully**: Provide fallback values if GitHub stats are unavailable
- **Don’t Over-Request**: The endpoint is cached; requesting multiple times quickly won’t get you fresher data
- **Use Appropriate Fields**: Only request the fields you need

## Next Steps

- [Installation](/installation)
- [Home](/) — Start here
- [Security & Abuse Protection](/security)
- [Dashboard](/dashboard)
