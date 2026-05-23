---
title: QR Codes
description: Generate and cache QR codes for any URL
template: app-page
---

# QR Codes

Generate and cache QR codes for any authorized URL.
Ideal for print media, business cards, product packaging, and anywhere you want to direct users to your website.

## Getting Started

Wynglet generates QR codes for your authorized URLs that you can embed on your site wherever appropriate.

Use this URL format:

```html
<img src="https://wynglet.your-server.com/qrcode/v1?url=your-site.com/some/page">
```

That’s it! The QR code will be generated, cached, and served to your users.

## How It Works

1. You request a QR code for a URL
2. Wynglet checks if the domain is in your authorized allowlist
3. Generates the QR code (or serves a cached version if available)
4. The image is cached for fast subsequent requests
5. Users can scan the QR code to visit the URL

## Security

To prevent abuse and to conserve resources, Wynglet blocks all domains by default, until you explicitly authorize each domain you care about.

[Manage your authorized domains in the Dashboard](/dashboard).

## Use Cases

- **Print Media**: Add QR codes to brochures, flyers, and ads
- **Business Cards**: Link to your contact information or website
- **Product Packaging**: Direct customers to product pages or manuals
- **Email Campaigns**: Add scannable links to email newsletters
- **Events**: Use QR codes on event materials for registration or information

## Dashboard Integration

Use the Dashboard’s **QR Codes** section to:

- Browse all generated QR codes
- Delete QR codes you no longer need
- Generate ready-to-use embed code for any authorized URL

## Caching

QR codes are cached for 30 days by default. You can adjust this in the configuration.

Since QR codes for the same URL are always identical, caching provides significant performance benefits with minimal storage overhead.

## Configuration

You can configure QR code caching in your `wynglet.yml`:

```yml
qr-codes:
  cache:
    ttl: 720h0m0s  # How long to cache generated QR codes
    max_size_bytes: 1073741824  # Maximum cache size (1GB)
```

## Next Steps

- [Installation](/installation)
- [Home](/) — Start here
- [View Dashboard](/dashboard)
