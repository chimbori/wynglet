---
title: Dashboard UI
description: Manage Wynglet from the web-based Dashboard
template: app-page
---

# Dashboard UI

Manage Wynglet from the Dashboard UI at `https://wynglet.your-server.com/dashboard`.

The Dashboard is available as an installable PWA (Progressive Web Application) that can be installed locally using any modern browser.

## Overview

The Dashboard provides centralized management for all Wynglet features, including:

- Viewing and managing generated content
- Authorizing domains
- Configuring per-domain settings
- Monitoring access and error logs
- Generating embed code for widgets

## Dashboard Sections

### Link Previews

Browse and manage all generated link preview images.

- **Browse images**: View all generated preview images.
- **Delete images**: Remove individual preview images to free up space.
- **View statistics**: See access statistics for each image.
- **Inspect user agents**: View a breakdown of which platforms.
  (Reddit, BlueSky, Facebook, Instagram, LinkedIn, Signal, Telegram, Slack, Discord, WhatsApp, iMessage, and many others)
  and bots have requested your images.

#### Sitemap Import

Bulk pre-generate link previews for an entire site by importing a sitemap URL.

- The import runs concurrently in the background
- Monitor progress in real-time
- Cancel imports at any time
- Configure concurrency and URL limits in your `wynglet.yml`

### QR Codes

Browse and manage all generated QR codes.

- **Browse codes**: View all generated QR codes
- **Delete codes**: Remove QR codes you no longer need
- **Generate embed code**: Get ready-to-use HTML for embedding QR codes on your pages

### Ratings

View and analyze all collected ratings.

- **View ratings by URL**: See all ratings collected for each page
- **Filter by timeline**: View ratings received in the last 1 day, or 7 days, or 28 days, etc.
- **Embed generator**: Generate ready-to-paste `<iframe>` code for rating widgets
- **Trend analysis**: Monitor rating trends over time

### Domains

Manage the Authorized Domains allowlist.

- **Add domains**: Authorize new domains to use Wynglet
- **Remove domains**: Revoke access for domains
- **Debug mode**: Toggle per-domain debug mode for detailed logging

### Logs

View recent error-level log entries.

- **Browse logs**: View all error logs in chronological order
- **Filter logs**: Search and filter logs by date, domain, or error type
- **Debug information**: See detailed error messages to troubleshoot issues

## PWA Installation

The Dashboard can be installed as a Progressive Web Application (PWA) on your device:

1. Open the Dashboard in your browser
2. Look for an “Install” option (usually in the address bar or menu)
3. Click “Install”
4. The Dashboard will be installed and appear in your applications

This allows you to:

- Access the Dashboard from your home screen
- Use it offline (cached data)
- Get a native app-like experience

## Dashboard Features

### Embed Code Generators

Several sections of the Dashboard include built-in embed code generators:

- **Rating Widget**: Generate `<iframe>` code for any authorized URL
- **QR Codes**: Generate `<img>` code for any authorized URL

Just select your URL and copy the generated code to your clipboard.

### Search and Filtering

Use the search boxes and filters in each section to quickly find what you’re looking for:

- Filter by URL
- Filter by date range
- Filter by domain
- Search by keyword

### Real-time Statistics

Monitor how your content is being used:

- View access statistics for link previews
- See which platforms are requesting your images
- Track rating submissions over time
- Monitor error rates and troubleshoot issues

## Security

The Dashboard is protected by HTTP Basic Authentication.
You configure your credentials in the `wynglet.yml` configuration file:

```yml
dashboard:
  username: admin
  password: "$2a$10$..."  # bcrypt-hashed password
```

Generate a new password hash using:

```shell
wynglet --bcrypt
```

## Next Steps

- [Home](/) — Start here
- [Installation](/installation)
- [OpenGraph Link Previews](/link-previews)
- [QR Codes](/qr-codes)
- [Rating Widget](/rating-widget)
