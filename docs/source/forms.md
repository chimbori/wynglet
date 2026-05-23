---
title: Form Submissions
description: Collect form submissions from your static site with CSRF protection and spam detection
template: app-page
---

# Form Submissions

Accept form submissions from your static site with built-in CSRF token protection, rate limiting, and spam detection.

## Overview

Wynglet provides a secure form submission API that allows your static site visitors to submit data without requiring a backend.

Features:

- **CSRF Protection**: Token-based anti-CSRF mechanism to prevent cross-site request forgery
- **Rate Limiting**: Prevent spam with per-IP rate limits
- **Honeypot Fields**: Detect automated form submissions
- **CORS Validation**: Only accept requests from authorized domains
- **Redirect on Success**: Customize where users are sent after submission

## Getting Started

### Step 1: Generate a Token

Before users can submit a form, you need to generate a CSRF token.
This is typically done when the page loads.

```javascript
fetch('https://wynglet.your-server.com/forms/v1/token', {
  method: 'GET'
})
.then(response => response.json())
.then(data => {
  document.getElementById('token').value = data.token;
})
```

### Step 2: Create Your Form

Create a form on your page that includes:

- The CSRF token from Step 1
- A `_form_id` to identify which form this is
- Your custom form fields
- Optional: `_subject` for email subject line
- Optional: `_redirect` for redirect URL after submission
- Optional: `_honeypot` for spam detection

```html
<form method="POST" action="https://wynglet.your-server.com/forms/v1/submit">
  <!-- Required: CSRF token -->
  <input type="hidden" name="_token" id="token" value="">

  <!-- Required: Form ID (identifies this form) -->
  <input type="hidden" name="_form_id" value="contact-form">

  <!-- Optional: Email subject line -->
  <input type="hidden" name="_subject" value="New Contact Form Submission">

  <!-- Optional: Redirect after successful submission -->
  <input type="hidden" name="_redirect" value="https://your-site.com/thank-you">

  <!-- Optional: Honeypot field (hidden from users) -->
  <input type="hidden" name="_honeypot" value="">

  <!-- Your form fields -->
  <input type="text" name="name" placeholder="Your name" required>
  <input type="email" name="email" placeholder="Your email" required>
  <textarea name="message" placeholder="Your message" required></textarea>

  <button type="submit">Send</button>
</form>
```

### Step 3: Handle Form Data

Form submissions are validated and stored in the Wynglet database.
You can view them in the Dashboard.

## Security

To prevent abuse and to conserve resources, Wynglet blocks all domains by default,
until you explicitly authorize each domain you care about.

[Manage your authorized domains in the Dashboard](/dashboard).

## Rate Limiting

To prevent spam:

- **20 submissions per IP per form per hour**: Prevents rapid-fire submissions from the same IP on the same form

These limits help ensure quality submissions while preventing spam.

## Spam Detection

### Honeypot Fields

The `_honeypot` field is a hidden field that should always be empty.
If it contains data, the submission is likely automated spam and is rejected.

```html
<!-- Hidden from real users, visible to bots -->
<input type="hidden" name="_honeypot" value="">
```

### CORS Validation

Requests must come from authorized domains. Requests from unrecognized origins are rejected.

### Form ID Validation

Each form submission must include a valid `_form_id` that has been previously registered.

## Configuration

Configure form handling in your `wynglet.yml`:

```yml
forms:
  rate-limit:
    per-ip-hour: 20
```

**Configuration options:**

- `per-ip-hour` — Maximum form submissions per IP address per hour (default: 20)

## Best Practices

- **Always use CSRF tokens**: Never accept form submissions without token validation
- **Use honeypot fields**: Add hidden fields that bots are likely to fill in
- **Validate on both sides**: Validate form data in Wynglet and on the client side
- **Use appropriate redirects**: Redirect users to a thank-you page after successful submission
- **Monitor submissions**: Regularly check the Dashboard for spam patterns

## Troubleshooting

- **Invalid token**: Make sure you generate a fresh token for each page load
- **Domain not authorized**: Ensure your domain is authorized in the Dashboard
- **Rate limit exceeded**: Wait before resubmitting; rate limits reset after the time window

## Next Steps

- [Installation](/installation)
- [Home](/) — Start here
- [Security & Abuse Protection](/security)
- [Dashboard](/dashboard)
