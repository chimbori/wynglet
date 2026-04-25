---
title: Butterfly Social
description: "The Dynamic Companion for your Static Sites: a self-hosted single-binary toolkit of embeddable services: OpenGraph preview images, QR codes, and Rating widgets"
template: app-page
nav:
- url: https://github.com/chimbori/butterfly
  title: Source
- url: https://github.com/chimbori/butterfly/issues
  title: Issues
- url: https://chimbori.com/feedback
  title: Feedback
---

# Butterfly Social

The dynamic companion for your static site: a self-hosted toolkit that brings dynamic capabilities to static websites without depending on paid third-party services.

Static site generators are great at producing fast, secure, cheap-to-host websites.
But they have one inherent limitation: they’re static.
Features like social preview images, QR codes, and rating widgets traditionally each require a separate paid SaaS subscription: each with its own account, pricing model, and data silo.

Butterfly bundles all of these into a single containerized open-source binary:

- **[OpenGraph Link Previews](#open-graph-link-previews)**:
  Auto-generate OpenGraph social preview images for every page, using your own HTML/CSS templates.

- **[QR Codes](#qr-codes)**:
  Generate and cache QR codes for any authorized URL.

- **[Rating Widget](#rating-widget)**:
  Embed a 👍/👎 or ⭐ rating widget on any page via a lightweight `<iframe>`.

**One instance serves many sites.** A single Butterfly deployment can power all your static sites simultaneously.
Every tool is gated by a domain allowlist: only explicitly authorized domains can embed or use your instance, preventing unauthorized use.

**Self-hosted.** Your data stays on your own infrastructure.
You can host Butterfly on a cheap VPS from any provider (here’s $200 in credits with a [DigitalOcean referral code](https://www.digitalocean.com/?refcode=e76ea0927117)).
No vendor lock-in, no per-request billing, no privacy trade-offs.

**Open source.** Apache 2.0 licensed, free forever.

## Open Graph Link Previews

Butterfly Social is a quick way to auto-generate OpenGraph link preview images in bulk for all your sites, without the use of a separate template editor or API integration at no cost.
The source of truth for the image data & design remains within your primary website, so you can use tools you are already familiar with & assets that

To prevent abuse and to conserve resources, Butterfly blocks all domains by default, until you explicitly authorize each domain you care about.

### Use the Default Template

Just one step: Paste the Butterfly `<meta>` tag into the original page, and you’re done!

```html
<meta property="og:image" content="https://butterfly.your-server.com/link-previews/v1?url=your-site.com/some/page">
```

### Use your Own Templates

1. Create a new hidden element inside your existing Web page, using whatever framework or template engine you use today.
   E.g.

    ```html
    <div id="link-preview" style="display: none; width: 1200px;">
      <h1>Butterfly Social</h1>
      <p>Self-hosted OpenGraph / social link preview image generation tool. Fully customizable, yet works out of the box with zero configuration.</p>
    </div>
    ```

2. Use Butterfly to craft a URL, and paste the `<meta>` tag into the original page.

    ```html
    <meta property="og:image" content="https://butterfly.your-server.com/link-previews/v1?url=your-site.com/some/page">
    ```

    If you can’t use the default selector (`#link-preview`) for any reason, you can provide an alternate one using the `&sel=` parameter.

3. There is no step 3.

### How it’s rendered

![Example](https://butterfly.chimbori.dev/example.png)

Test your Butterfly installation by posting your original page URL to any social platform.

### How it works

Butterfly fetches the URL you provide to it, using a Chrome Headless instance, runs JavaScript to un-hide the hidden element, takes a screenshot of it, and serves it (while also caching & compressing it).

Butterfly works well with static sites (using any static site generator) as well as dynamically-generated sites (using any CMS or platform).

### Can I use…

- Images? Yes.
- SVG backgrounds? Also, yes.
- Flexbox? Grid? Yes, of course.
- Custom fonts? Proprietary fonts? Absolutely.

Why limit yourself to the customization possible in a random WYSIWYG editor, when you have the entire Web platform available to you!

Anything you can design for the Web, you can use to create a link preview image.
The infinite is possible at Zombocom.
The unattainable is unknown at Zombocom.

## QR Codes

Butterfly generates QR Codes for your authorized URLs that you can embed on your site wherever appropriate.

Use this URL format:

```html
<img src="https://butterfly.your-server.com/qrcode/v1?url=your-site.com/some/page">
```

## Rating Widget

Butterfly includes an embeddable rating widget that your visitors can use to rate any page on your site.
Ratings are not public, but are visible on your Butterfly Dashboard.
It provides a lightweight mechanism to monitor the quality of content, e.g. for Help Articles.
It supports two UI styles:

- **Thumbs** (`ui=thumbs`): 👍 / 👎 buttons _(default)_
- **Stars** (`ui=stars`): ⭐ 1–5 star rating

Embed the widget on your page using an `<iframe>`:

```html
<iframe
  src="https://butterfly.your-server.com/rating/v1?ui=thumbs&url=https://your-site.com/some/page"
  style="border:none; width:200px; height:50px;">
</iframe>
```

Ratings are validated against the authorized domains allowlist.
To prevent abuse, a single IP is limited to submitting one rating for a URL per 24 hours, and 10 ratings across all URLs per hour).

Use the Dashboard’s **Ratings** section to view collected ratings and generate the correct `<iframe>` embed code for any authorized URL.

## CLI Utilities

Butterfly ships with a couple of built-in command-line helpers:

- **Generate a bcrypt password hash** (for use in `butterfly.yml`):

  Enter your desired password when prompted.
  The output hash can be pasted directly into the `dashboard.password` field.

  ```shell
  butterfly --bcrypt
  ```

- **Health check** (useful for container health checks):

  Pings the running service and exits `0` on success or `1` on failure.

  ```shell
  butterfly --healthcheck
  ```

## Install & Deploy

- We strongly recommend deploying using the official container image, which includes Chrome Headless for convenience.
  - Thanks to the [chromedp](https://github.com/chromedp/chromedp) project for making this possible!

- Butterfly is designed to be used behind a TLS reverse proxy for SSL termination (among other things).
  - We recommend [Caddy](https://caddyserver.com/); see sample Caddyfile below.
  - If you expect a lot of traffic, consider putting a CDN in front.

### Sample `compose.yml` (for Docker and Podman)

```yml
services:
  butterfly:
    container_name: butterfly
    image: ghcr.io/chimbori/butterfly:latest
    volumes:
      - $PWD/butterfly-data:/data
    restart: unless-stopped
    depends_on:
      - butterfly-db

  butterfly-db:
    container_name: butterfly-db
    image: postgres:18-alpine
    environment:
      POSTGRES_DB: butterfly
      POSTGRES_USER: chimbori
      POSTGRES_PASSWORD: chimbori
    volumes:
      - $PWD/butterfly-db-data:/var/lib/postgresql
    restart: unless-stopped

volumes:
  butterfly-data:
  butterfly-db-data:
```

#### If you prefer to install using the Go binary on raw metal, use `go install`

```shell
go install butterfly.chimbori.dev@latest
```

### Sample `butterfly.yml`

Butterfly requires basic configuration to be provided via a config file.

- **PostgreSQL Database URL** _(required)_

  ```yml
  database:
    url: postgresql://chimbori:chimbori@butterfly-db:5432/butterfly
  ```

- **Dashboard credentials** (encrypted via `bcrypt`) _(required)_

  ```yml
  dashboard:
    username: admin
    password: "$2a$10$a8LnUkK1UiB.9yQrUp3wyuGsH1AAHhlHVy1cjIaaIUVAwCtGvaX7q" # "test"
  ```

- Web config _(optional)_

  Assuming there’s a reverse proxy in front of Butterfly Social, there should be no need to change the port here; just configure the reverse proxy to forward requests to port 9999.

  ```yml
  web:
    port: 9999
  ```

- Link Previews config _(optional)_

  Configure Cache TTL if needed for testing (e.g. set to a low value like `1m`).

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

- QR Codes config _(optional)_

  Configure Cache TTL if needed for testing (e.g. set to a low value like `1m`).

  ```yml
  qr-codes:
    cache:
      ttl: 720h0m0s
      max_size_bytes: 1073741824
  ```

- Ratings config _(optional)_

  ```yml
  ratings:
    retention: 8760h  # How long to keep rating data (default: 365 days)
  ```

- Logs config _(optional)_

  ```yml
  logs:
    retention: 720h  # How long to keep error logs (default: 30 days)
    pagination:
      limit: 10
  ```

- Dashboard pagination _(optional)_

  ```yml
  dashboard:
    pagination:
      limit: 10
  ```

- Debug Mode _(optional)_

  Turn on additional logging in Debug Mode.
  Debug mode can also be toggled per-domain from within the Dashboard UI.

  ```yml
  debug: true
  ```

### Sample `Caddyfile`

```caddy
butterfly.your-server.com {
  reverse_proxy butterfly:9999
  encode zstd gzip
}
```

# Dashboard UI

You can manage Butterfly from the Dashboard UI at `https://butterfly.your-server.com/dashboard`.
The Dashboard is available as an installable PWA (Progressive Web Application) that can be "installed" locally using any modern browser.

The Dashboard includes the following sections:

- **Link Previews**:
  Browse all generated link preview images, delete individual ones, view access statistics, and inspect a breakdown of user agents (social platforms, bots, etc.)
  that have requested them.

  - **Sitemap Import**:
    Bulk pre-generate link previews for an entire site by importing a sitemap URL.
    The import runs concurrently in the background; you can monitor progress or cancel it at any time.

- **QR Codes**:
  Browse and delete generated QR codes.

- **Ratings**:
  View collected ratings for any authorized URL.
  Use the built-in embed builder to generate ready-to-paste `<iframe>` code for your pages.

- **Domains**:
  Manage the authorized domains allowlist.
  Add or remove domains and toggle per-domain debug mode.

- **Logs**:
  View recent error-level log entries.

<img src="screenshot-pwa.webp">

# License

This project is licensed under the Apache License 2.0.

# Comparison with Alternatives

There are a lot of paid SaaS tools in this space.

- BannerBear
- RenderForm
- Templated.io
- Imejis.io
- Pablle
- Orshot
- Abyssale

They all work roughly the same way: you design a template using their custom tools, then provide them your data (title, description, etc.), and pay them per-request (or per-render) to create & serve those images for you.

This model works great if you do not have access to the source of the page, or have no influence over the developers who build your website.

But now,

- You’ve got to learn a whole new tool.
- That tool exposes a certain amount of design expressiveness, but nowhere near what the Web platform offers natively.
- Anytime you need to change the preview image, you have to visit a completely separate website.
- Anytime your own webpage changes, you have to remember to update the templates to match the theme.
- There’s no way to share themes between your website & these third-party tools: colors, gradients, logos must be copy/pasted manually.
- You have to rely on these companies being around long enough, and not disappearing completely after running out of money or being bought over by a VC.
- And you have to pay, based on volume.

Butterfly is none of those things. All you need is the ability to write some HTML/CSS (no JavaScript necessary!) to design your preview image. And it’s free in perpetuity.
