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
- **[Feedback](https://wynglet.chimbori.dev/feedback)**: Send Feedback
- **[Source Code](https://github.com/chimbori/wynglet)**: View and contribute to the source

# Tech Stack

Wynglet is primarily Golang + PostgreSQL + HTMX + Templ (+ some custom TypeScript compiled to JavaScript).

## Why Golang + PostgreSQL

Go’s standard library offers excellent APIs for Web services, and is extremely lean on resources at runtime.

- Memory-efficient: Go’s compiled binary and PostgreSQL’s optimized database engine require a fraction of the memory footprint of interpreted/JVM languages.
- Easy to deploy: Wynglet requires Chrome as a dependency for OpenGraph rendering, but besides that has no other runtime dependencies, thanks to Go compiling it to a single static binary. No `CGO`.
- Battle-tested: PostgreSQL has decades of production stability and reliability. I’ve tried Sqlite in other projects, and the consistently high throughput offered by PostgreSQL, even on cheap hardware, made it a better choice.
- Cheap: A single Wynglet instance can run reliably on minimal hardware (even a $5/month VPS), making it ideal for self-hosted deployments.

## Templ, Tailwind, & HTMX

- Instead of Go stdlib’s `html/template`, I prefer [Templ](https://templ.guide/) for its type-safety, build-time validation, & runtime speed. Its language server integration offers a pretty awesome developer experience.
- I used to be a [Tailwind](https://tailwindcss.com/) skeptic and had a stash of copy/paste-able SASS/CSS lying around, but after using Tailwind in several side projects, I’ve come to appreciate the developer experience. I still define some reusable classes (using `@apply`) but being able to apply consistent grid-based styles quickly has been a game-changer.
- Although I have significant experience authoring production-grade JavaScript at scale (e.g. for Google.com and Gmail Web), the current Web ecosystem is needlessly complex. I wanted something simple for this project, and [HTMX](https://htmx.org) fit the bill perfectly.
- I’m not dogmatic about it, though; there are valid use cases for custom JavaScript. This project uses [`esbuild`](https://esbuild.github.io/) to build TypeScript to JavaScript (`esbuild` itself is built in Go!), which makes it possible to avoid `npm` and `node_modules`.

## sqlc and Goose

- Wynglet uses [`sqlc`](https://sqlc.dev/) to use SQL as the canonical schema & query language. Instead of ORMs, `sqlc` generates Go code from SQL schemas & queries, and the rest of the project references the generated Go code.
- Migrations are handled using [Goose](https://pressly.github.io/goose/). When updating a Wynglet installation, any required migrations are applied automatically at startup. No user intervention or separate database management is needed.
- Only upgrades are supported; the database cannot be migrated backwards for restores. Take regular backups before upgrading Wynglet.

If this interests you, I’d love for you to contribute to Wynglet!

## AI Usage

As a professional software engineer for the last 25+ years (and 17+ years at Google),
I value immensely the craft & joy of software engineering.
I sweat the big parts & the small parts—the high-level architecture as well as each and every line of code—to make sure it’s clean, bug-free, and maintainable.

That said, as all effective engineers should, I use the best tools available for the job.
In the 1990s, that was Turbo C++; in the 2000s, Visual Studio; in 2010s, JetBrains; in 2025, that’s LLMs.
I use AI unapologetically to ship features faster, locate tricky bugs, and ensure a high standard of security across my projects.

I do not let LLMs commit any code on my behalf; I review every line of code before it is merged.
If a function was LLM-generated, I typically fix up things by hand before committing anything.

And irrespective of how that code got authored—either in a text editor, IDE, via LLMs, or by an external human contributor—you should feel free to hold me personally responsible for everything in this project.

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
