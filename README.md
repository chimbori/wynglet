# Wynglet

The dynamic companion for your static site: a self-hosted toolkit that brings dynamic capabilities to static websites without depending on paid third-party services.

Static site generators are great at producing fast, secure, cheap-to-host websites.
But they have one inherent limitation: they’re static.
Features like auto-generated OpenGraph social preview images, forms, QR codes, etc. each require a back-end or a separate paid SaaS subscription: each with its own account, pricing model, and data silo.

Wynglet bundles all of these into a single containerized open-source binary:

## Features

### OpenGraph Link Previews

Auto-generate social media preview images for every page. No template editor needed—just write HTML/CSS that you control.
[Learn more →](https://wynglet.chimbori.dev/link-previews)

### QR Codes

Generate and cache QR codes for any authorized URL. Ideal for print media, business cards, and product packaging.
[Learn more →](https://wynglet.chimbori.dev/qr-codes)

### Rating Widget

Embed a simple rating widget (👍/👎 or ⭐) on your pages to collect visitor feedback. Ratings are private but visible in the Dashboard.
[Learn more →](https://wynglet.chimbori.dev/rating-widget)

### Dashboard

Manage all features from the web-based Dashboard at `https://wynglet.your-server.com/dashboard`. Available as a Progressive Web Application (PWA).
[Learn more →](https://wynglet.chimbori.dev/dashboard)

## Quick Start

For detailed installation and deployment instructions, see
**[Installation & Deployment Guide](https://wynglet.chimbori.dev/installation)**.

Supported deployment methods:

- **Docker Compose** (recommended) — includes Chrome Headless
- Binary: Go binary installation

**One instance serves many sites.** A single Wynglet deployment can power all your static sites simultaneously.
Every feature is gated by a domain allowlist: only explicitly authorized domains can embed or use your instance, preventing unauthorized use.

**Self-hosted.** Your data stays on your own infrastructure.
You can host Wynglet on a cheap VPS from any provider (here’s $200 in credits with a [DigitalOcean referral code](https://www.digitalocean.com/?refcode=e76ea0927117)).
No vendor lock-in, no per-request billing, no privacy trade-offs.

**Open source.** Apache 2.0 licensed, free forever.

## Getting Help

- **[Documentation](https://wynglet.chimbori.dev/)**: Full documentation and guides
- **[GitHub Issues](https://github.com/chimbori/wynglet/issues)**: Report bugs and request features
- **[Feedback](https://chimbori.com/feedback)**: Send general feedback
- **[Source Code](https://github.com/chimbori/wynglet)**: View and contribute to the source

# License

Copyright 2025, Chimbori

Licensed under the [Apache License, Version 2.0](LICENSE.md) (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
