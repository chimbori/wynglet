---
title: Wynglet
description: "The Dynamic Companion for your Static Sites: a self-hosted single-binary toolkit of embeddable services: OpenGraph preview images, form submissions, QR codes, Rating widgets, and GitHub stats"
template: app-page
nav:
- url: https://github.com/chimbori/wynglet
  title: Source
- url: https://github.com/chimbori/wynglet/issues
  title: Issues
- url: https://chimbori.com/feedback
  title: Feedback
---

The dynamic companion for your static site:
a self-hosted toolkit that brings dynamic capabilities to static websites without depending on paid third-party services.

**One instance serves many sites.** A single Wynglet deployment powers all your static sites simultaneously,
with a domain allow-list ensuring only authorized sites can use it.

**Self-hosted.** Your data stays on your infrastructure.
No vendor lock-in, no per-request billing, no privacy compromises.

**Open source.** Apache 2.0 licensed, free forever.

# Wynglet Features

## Auto-generated OpenGraph Link Preview Images

[Learn more →](/link-previews)

- Auto-generate beautiful social media preview images using your own HTML/CSS templates.
- Design once in your website’s theme, then Wynglet generates previews automatically for every page.
- No need to learn proprietary tools or maintain separate design workflows.

## Form Submissions

[Learn more →](/forms)

- Collect form data from your site visitors with zero configuration required.
- Built-in security: CSRF token protection, automatic honeypot spam detection, and per-IP rate limiting.
- All submissions are validated, logged, and accessible through your dashboard.

## QR Codes

[Learn more →](/qr-codes)

- Generate scannable QR codes for any URL.
- Perfect for linking from print materials, posters, business cards, and physical signage to your online content.
- Generated codes are automatically cached to minimize server load.

## Rating Widget

[Learn more →](/rating-widget)

- Embed a lightweight rating interface (👍/👎 or ⭐ 1–5) to collect visitor feedback directly on your site.
- Privacy-focused, rate-limited to prevent abuse, and fully logged in your dashboard for analysis.

## GitHub Repository Stats

[Learn more →](/github)

- Display live GitHub repository statistics (stars, forks, issues, watchers) on your site without external JavaScript.
- Handles CORS appropriately, so you can use it from any of your authorized domains.
- Perfect for showcasing your open-source projects or highlighting popular repositories you maintain.

# Getting Started

Ready to get started? Follow these steps:

1. **[Install Wynglet](/installation)** on your server
2. **Configure** your `wynglet.yml` file
3. **Authorize domains** in the Dashboard
4. **Choose a feature** to start with:
   - [OpenGraph Link Previews](/link-previews)
   - [Form Submissions](/forms)
   - [QR Codes](/qr-codes)
   - [Rating Widget](/rating-widget)
   - [GitHub Stats](/github)
5. **Manage everything** from the [Dashboard UI](/dashboard)

## Quick Example

The simplest way to get started is with the default link preview template. Just paste this meta tag into your page:

```html
<meta property="og:image" content="https://wynglet.your-server.com/link-previews/v1?url=your-site.com/some/page">
```

Then test it by sharing your page on social media. For more details, see [OpenGraph Link Previews](/link-previews).

## Setup & Administration

- [Installation & Deployment](/installation) — Deploy Wynglet on your infrastructure
- [Dashboard UI](/dashboard) — Manage all features from the web interface
- [Security & Abuse Protection](/security) — Domain authorization, rate limiting, CSRF tokens

## Domain-Based Authorization

Wynglet uses a **deny-by-default** security model to prevent unauthorized use:

1. **All domains are blocked by default**: You control which sites can use your instance
2. **Allow-list-based access**: Add authorized domains in the Dashboard
3. **Every request is validated**: Each API call checks the request origin

This means:

- You can run one instance for multiple sites
- Malicious sites cannot use your instance
- You maintain complete control over who has access

Learn more: [Security & Abuse Protection](/security)

## Support

- **[GitHub Issues](https://github.com/chimbori/wynglet/issues)**: Report bugs and request features
- **[Feedback](https://chimbori.com/feedback)**: Send feedback & suggestions
- **[Source Code](https://github.com/chimbori/wynglet)**: View and contribute

## License

Wynglet is licensed under the [Apache License, Version 2.0](https://github.com/chimbori/wynglet/blob/main/LICENSE.md).
Copyright 2025 onwards, Chimbori.
