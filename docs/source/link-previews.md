---
title: OpenGraph Link Preview Images
description: Auto-generate OpenGraph social preview images for every page
template: app-page
---

# OpenGraph Link Preview Images

Generate OpenGraph link preview images in bulk for all your sites, without the use of a separate template editor or API integration.

## Overview

When you share a link on social media and messaging platforms—like Reddit, BlueSky, Facebook, Instagram, LinkedIn, Signal, Telegram, Slack, Discord, WhatsApp, iMessage, and many others—a preview image is displayed. This preview is generated from the OpenGraph meta tags on your page.

Wynglet automatically generates these preview images by:

1. Fetching your page
2. Rendering a designated HTML element using Chrome Headless
3. Taking a screenshot of it
4. Serving and caching the image

The source of truth for the image data & design remains within your primary website, so you can use tools you are already familiar with & assets that you already have.

## Getting Started

### Option 1: Use the Default Template

The easiest way to get started is with Wynglet’s default template. Just paste this meta tag into the `<head>` of your page:

```html
<meta property="og:image" content="https://wynglet.your-server.com/link-previews/v1?url=your-site.com/some/page">
```

Then test it by posting your page URL to any social platform.

### Option 2: Use Your Own Custom Template

1. Create a hidden element in your page with the design you want for the preview image:

    ```html
    <div id="link-preview" style="display: none; width: 1200px;">
      <h1>Wynglet</h1>
      <p>Self-hosted OpenGraph / social link preview image generation tool.</p>
    </div>
    ```

2. Paste the meta tag into your page (using the default `#link-preview` selector):

    ```html
    <meta property="og:image" content="https://wynglet.your-server.com/link-previews/v1?url=your-site.com/some/page">
    ```

3. If you can’t use the default selector for any reason, specify an alternate one:

    ```html
    <meta property="og:image" content="https://wynglet.your-server.com/link-previews/v1?url=your-site.com/some/page&sel=.custom-selector">
    ```

## Design Capabilities

Wynglet uses the full power of the Web platform for rendering preview images. You can use:

- Images and backgrounds
- SVG backgrounds
- Flexbox and CSS Grid
- Custom fonts
- Proprietary fonts
- Any CSS you can write

Anything you can design for the Web, you can use to create a link preview image. You’re not limited to a WYSIWYG editor or a predefined set of templates.

## Security

To prevent abuse and to conserve resources, Wynglet blocks all domains by default, until you explicitly authorize each domain you care about.

[Manage your authorized domains in the Dashboard](/dashboard).

## How It Works

Here’s what happens when someone requests a link preview:

1. Wynglet fetches the URL you provide
2. Uses Chrome Headless to render the page and execute JavaScript
3. Un-hides your designated preview element
4. Takes a screenshot of that element
5. Compresses and caches the image
6. Serves it to the requesting platform

## Dashboard Integration

Use the Dashboard’s **Link Previews** section to:

- Browse all generated link preview images
- Delete individual preview images
- View access statistics
- Inspect a breakdown of user agents (social platforms, bots, etc.) that have requested them

### Bulk Import with Sitemap

Pre-generate link previews for an entire site by importing a sitemap URL. The import runs concurrently in the background; you can monitor progress or cancel it at any time.

## Example

Here’s an example of a preview image design:

![Example](https://wynglet.chimbori.dev/example.png)

Test your Wynglet installation by posting your original page URL to any social platform and seeing the preview appear.

## Tips & Best Practices

- **Cache TTL**: Link preview images are cached for 720 hours (30 days) by default.
  You can adjust this in the configuration.
- **Performance**: Consider using Wynglet behind a CDN for better performance at scale.
- **Dynamic Content**: Wynglet fully supports dynamically-generated content as well as static sites.
- **Testing**: Use your platform’s built-in link preview tools to debug your previews.
  Most platforms have preview debuggers to test how your links will appear when shared.

## Comparison with Alternatives

There are many paid SaaS tools for generating preview images:
BannerBear, RenderForm, Templated.io, Imejis.io, Pablle, Orshot, Abyssale, etc.

They all work roughly the same way: you design a template using their custom tools, then provide your data (title, description, etc.), and **pay them per request** to generate and serve images.

This model comes with several drawbacks:

- **Have to learn new tools**: You have to learn to use a proprietary template editor instead of using HTML/CSS you already know.
- **Limited design expressiveness**: These proprietary tools expose only a fraction of what the Web platform offers.
- **Context switching**: Anytime you need to change a layout, you visit a separate website.
- **Theme maintenance**: Your preview templates can’t easily reference your website’s theme or keep up with changes
  colors, gradients, and logos must be manually copied.
- **Vendor lock-in**: You depend on these companies staying in business; mergers and acquisitions have historically shut down such services.
  Pricing can vary arbitrarily after you’ve signed up.
- **Per-request billing**: Costs scale with volume.
  If you have a large site with new, fresh content all the time, you’ll spend a fortune on these tools.

Wynglet is different.

- **Use HTML/CSS**: You can author your templates in languages you already know. Use all the Web design tools you use today.
- **Full design control**: Access the entire Web platform’s capabilities (flexbox, grid, custom fonts, SVG, etc.)
- **Co-located templates**: Update preview templates alongside your website where they live.
- **Shared themes**: Reference the same colors, gradients, and logos as your website.
  When they change site-wide, your template automatically receives the same changes
- **Open source**: Wynglet won’t go away randomly one day.
  It’s yours to run on your own infrastructure.
- **Free forever**: One-time deployment cost, no ongoing per-request charges.

## Next Steps

- [Installation](/installation)
- [Home](/) — Start here
- [View Dashboard](/dashboard)
