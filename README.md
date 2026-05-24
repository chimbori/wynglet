# Wynglet

The dynamic companion for your static site:
a self-hosted toolkit that brings dynamic capabilities to static websites without depending on paid third-party services.

Static site generators are great at producing fast, secure, cheap-to-host websites.
But they have one inherent limitation: they’re static.
Features like auto-generated OpenGraph social preview images, forms, QR codes, etc. each require a back-end or a separate paid SaaS subscription: each with its own account, pricing model, and data silo.

Wynglet bundles all of these features into a single containerized open-source binary:

## Auto-generated OpenGraph Link Preview Images

[Learn more →](https://wynglet.chimbori.dev/link-previews)

- Auto-generate beautiful social media preview images using your own HTML/CSS templates.
- Design once in your website’s theme, then Wynglet generates previews automatically for every page.
- No need to learn proprietary tools or maintain separate design workflows.

## Form Submissions

[Learn more →](https://wynglet.chimbori.dev/forms)

- Collect form data from your site visitors with zero configuration required.
- Built-in security: CSRF token protection, automatic honeypot spam detection, and per-IP rate limiting.
- All submissions are validated, logged, and accessible through your dashboard.

## QR Codes

[Learn more →](https://wynglet.chimbori.dev/qr-codes)

- Generate scannable QR codes for any URL.
- Perfect for linking from print materials, posters, business cards, and physical signage to your online content.
- Generated codes are automatically cached to minimize server load.

## Rating Widget

[Learn more →](https://wynglet.chimbori.dev/rating-widget)

- Embed a lightweight rating interface (👍/👎 or ⭐ 1–5) to collect visitor feedback directly on your site.
- Privacy-focused, rate-limited to prevent abuse, and fully logged in your dashboard for analysis.

## GitHub Repository Stats

[Learn more →](https://wynglet.chimbori.dev/github)

- Display live GitHub repository statistics (stars, forks, issues, watchers) on your site without external JavaScript.
- Handles CORS appropriately, so you can use it from any of your authorized domains.
- Perfect for showcasing your open-source projects or highlighting popular repositories you maintain.

# Getting Started

For detailed installation and deployment instructions, see
**[Installation & Deployment Guide](https://wynglet.chimbori.dev/installation)**.

Supported deployment methods:

- **Docker Compose** (recommended) — includes Chrome Headless
- Binary: Go binary installation

**One instance serves multiple sites.** A single Wynglet deployment can power all your static sites simultaneously.
Every feature is gated by a domain allowlist: only explicitly authorized domains can embed or use your instance, preventing unauthorized use.

**Self-hosted.** Your data stays on your own infrastructure.
You can host Wynglet on a cheap VPS from any provider (here’s $200 in credits with a [DigitalOcean referral code](https://www.digitalocean.com/?refcode=e76ea0927117)).
No vendor lock-in, no per-request billing, no privacy trade-offs.

**Open source.** Apache 2.0 licensed, free forever.

## Quick Example

The simplest way to get started is with the default link preview template. Just paste this meta tag into your page:

```html
<meta property="og:image" content="https://wynglet.your-server.com/link-previews/v1?url=your-site.com/some/page">
```

Then test it by sharing your page on social media. For more details, see [OpenGraph Link Previews](https://wynglet.chimbori.dev/link-previews).

## Getting Help

- **[Documentation](https://wynglet.chimbori.dev/)**: Full documentation and guides
- **[GitHub Issues](https://github.com/chimbori/wynglet/issues)**: Report bugs and request features
- **[Feedback](https://chimbori.com/feedback)**: Send general feedback
- **[Source Code](https://github.com/chimbori/wynglet)**: View and contribute to the source

# License

Copyright 2025 onwards, Chimbori

Licensed under the [Apache License, Version 2.0](https://github.com/chimbori/wynglet/blob/main/LICENSE.md) (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
