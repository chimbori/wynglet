---
title: Rating Widget
description: Embed a rating widget on any page to collect user feedback
template: app-page
---

# Rating Widget

Embed an interactive rating widget on any page to collect user feedback.
The widget supports both thumbs (👍/👎) and star rating (⭐) interfaces.

## Overview

The Rating Widget is a lightweight mechanism to monitor the quality of content and collect feedback from your visitors.

Features:

- **Private ratings**: Ratings are not public, but are visible on your Wynglet Dashboard.
- **Two UI styles**: Thumbs (👍/👎) or Stars (⭐ 1–5)
- **Lightweight**: Embedded via a simple `<iframe>`, no heavy JavaScript required.
- **Rate-limited**: Prevents abuse by limiting submissions per user.

## Getting Started

Embed the widget on your page using an `<iframe>`:

### Thumbs UI (Default)

```html
<iframe
  src="https://wynglet.your-server.com/rating/v1?ui=thumbs&url=https://your-site.com/some/page"
  style="border:none; width:200px; height:50px;">
</iframe>
```

### Stars UI

```html
<iframe
  src="https://wynglet.your-server.com/rating/v1?ui=stars&url=https://your-site.com/some/page"
  style="border:none; width:200px; height:50px;">
</iframe>
```

## UI Styles

### Thumbs (Default)

Displays a simple 👍 (like) and 👎 (dislike) button pair.
Best for binary feedback.

```html
<iframe
  src="https://wynglet.your-server.com/rating/v1?ui=thumbs&url=https://your-site.com/some/page"
  style="border:none; width:200px; height:50px;">
</iframe>
```

### Stars

Displays a 1–5 star rating interface.
Best for detailed feedback.

```html
<iframe
  src="https://wynglet.your-server.com/rating/v1?ui=stars&url=https://your-site.com/some/page"
  style="border:none; width:200px; height:50px;">
</iframe>
```

## Rate Limiting

To prevent abuse, a single IP is limited to:

- **1 rating per URL per 24 hours**: Users can’t spam the same page repeatedly.
- **10 ratings across all URLs per hour**: Abusive users can’t overload the system.

These limits help ensure quality feedback while preventing abuse.

## Security

Ratings are validated against the authorized domains allowlist.
Only ratings from authorized domains are accepted.

[Manage your authorized domains in the Dashboard](/dashboard).

## Dashboard Integration

Use the Dashboard’s **Ratings** section to:

- View collected ratings for any authorized URL
- Monitor feedback trends
- Generate ready-to-paste `<iframe>` embed code for any authorized URL

The embed builder in the Dashboard makes it easy to copy the correct code for your pages without having to remember the URL format.

## Use Cases

- **Help Articles**: Monitor the usefulness of documentation
- **Blog Posts**: Collect reader feedback on content quality
- **Product Pages**: Gauge customer satisfaction with product information
- **FAQs**: Identify which answers are most helpful
- **Tutorials**: Get feedback on tutorial clarity and usefulness

## Configuration

You can configure rating data retention in your `wynglet.yml`:

```yml
ratings:
  retention: 8760h  # How long to keep rating data (default: 365 days)
```

## Best Practices

- **Place strategically**: Put the widget at the end of your content where readers have finished consuming it
- **Collect feedback early**: Use thumbs for quick "Was this helpful?" feedback
- **Follow up**: For detailed feedback (5-star ratings), consider adding a comment form or link to your feedback system
- **Act on feedback**: Regularly review ratings and update your content based on feedback

## Next Steps

- [Installation](/installation)
- [Home](/) — Start here
- [View Dashboard](/dashboard)
