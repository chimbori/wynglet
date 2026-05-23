---
title: Security & Abuse Protection
description: How Wynglet protects against abuse and unauthorized use
template: app-page
---

# Security & Abuse Protection

Wynglet is designed with security and abuse prevention as core features.
Every API and service is protected by multiple layers of validation and rate limiting.

## Domain-Based Authorization

### How It Works

Wynglet operates on a **deny-by-default** security model:

- **All domains are blocked by default**: No one can use your Wynglet instance until you explicitly authorize them.
- **Allow-list-based access**: You maintain an allow-list of authorized domains in the Dashboard.
- **Domain validation on every request**: Each request checks that its origin is on the allow-list.

### Setting Up Domain Authorization

1. Go to the Dashboard at `https://wynglet.your-server.com/dashboard`
2. Navigate to the **Domains** section
3. Click **Add Domain**
4. Enter your domain (e.g., `example.com` or `subdomain.example.com`)
5. Save

If your instance has recently received requests from particular domains, you will see these listed there (with the denied-by-default status).
This allows you to quickly allow-list domains while bringing up new sites.
If you do not explicitly authorize a domain, it will stop showing up in the list after a few days.

### Domain Matching Rules

- By default, only the **exact domain** is allow-listed: `example.com` and `www.example.com` are treated as different domains,
  and authorizing one will not authorize the other one.
- **Sub-domains can be optionally allow-listed**:
  You can authorize all subdomains (without having to list each one individually) if you prefer.
  E.g. allow-listing all subdomains for `example.com` will match `www.example.com`, `api.example.com`, `blog.example.com`, etc.
- **Case-insensitive** (like all domains in general): `Example.com` matches `example.com`.
- **Protocols don’t matter**: both `http://example.com` and `https://example.com` are authorized if `example.com` is on the allow-list.
- **Ports are ignored**: `example.com:8080` is treated as `example.com`.

### Why This Matters

Without domain authorization, anyone could use your Wynglet instance:

- Use your QR code generator to link to their site.
- Embed your rating widget and collect fake ratings.
- Generate preview images at your expense.
- Overwhelm your instance with requests.

## Rate Limiting

In addition to Domain Authorization, Wynglet uses IP-based rate limiting to prevent abuse and protect server resources.

### Global Rate Limits

Each IP address is limited across all features:

- **Form Submissions**:
  - 1 submission per form per 60 seconds
  - 10 total submissions per hour

- **Ratings**:
  - 1 rating per URL per 24 hours
  - 10 ratings across all URLs per hour

### What Happens When Rate Limits Are Exceeded

When a rate limit is exceeded, the request is rejected with a `429 Too Many Requests` status code.

### Bypassing Rate Limits

Rate limits apply per-IP. Multiple clients from different IP addresses can submit simultaneously. However:

- **Shared Networks**: All users on the same corporate network, school, or ISP may share a rate limit bucket
- **Residential ISPs**: Most residential ISPs assign static IPs per household
- **Mobile Networks**: Mobile ISPs often assign IPs per cell tower

## CSRF Protection

### What Is CSRF?

Cross-Site Request Forgery (CSRF) is an attack where a malicious site tricks your browser into making unwanted requests to another site on your behalf.

### How Wynglet Protects Against It

Wynglet uses **CSRF tokens** for form submissions:

1. Before rendering a form, you request a token from Wynglet
2. You include this token in any form submission
3. Wynglet validates that the token is present and valid
4. Malicious sites don’t have access to your valid tokens

### Token Validation

- Tokens are **unique per form**: each new form render gets a new token
- Tokens **expire**: old tokens become invalid
- Tokens are **single-use**: after submission, they cannot be reused

## Honeypot Fields

### How It Works

A honeypot field is a hidden form field that humans won’t see, but automated bots will fill in.

Example:

```html
<!-- Hidden from real users, visible to bots -->
<input type="hidden" name="_honeypot" value="">
```

When a form is submitted:

- **Humans** leave the field empty (they can’t see it)
- **Bots** fill it in with something (they fill every field)

Wynglet automatically rejects submissions where honeypot fields are filled in.

## CORS Validation

### What Is CORS?

Cross-Origin Resource Sharing (CORS) is a browser security feature that controls which sites can make requests to an API.

### How Wynglet Uses CORS

- **Strict CORS headers**: Wynglet sets CORS headers that only allow requests from authorized domains.
  This prevents unauthorized sites from making requests to your Wynglet instance.
- **Reject file:// URLs**: Requests from `file://` are always rejected
- **GitHub stats endpoint**: Allows any origin (intentionally permissive for public stats)

In practice, this means:

- Only domains you’ve explicitly authorized in the Dashboard can make cross-origin requests to Wynglet
- Unauthorized domains will be blocked by the browser’s CORS policy
- This provides an additional layer of security beyond domain allow-listing

## Debug Mode

### Toggling Debug Mode

You can enable per-domain debug mode in the Dashboard:

1. Go to **Domains** section
2. Find your domain
3. Toggle **Debug Mode**

When enabled, Wynglet logs detailed information about requests from that domain, helping you troubleshoot issues.

## Best Practices

- **Authorize only your domains**: Add all your sites’ domains, but no others
- **Monitor the Logs**: Check the Dashboard Logs section for suspicious activity
- **Review Ratings & Submissions**: Look for spam patterns in your collected data
- **Use Honeypots**: Always include honeypot fields in your contact forms
- **Validate Client-Side**: Validate form data on your website too

## Monitoring Abuse

### Signs of Abuse

- Unusually high number of ratings from a single IP
- Many form submissions from a single IP
- Requests from unexpected domains
- Requests from suspicious origins

### What To Do

1. Check the **Logs** section in the Dashboard for errors
2. Review **Ratings** and **Submissions** for spam
3. Check if the attacking IP is on your domain allow-list
4. Consider removing offending domains if they’re being abused
5. Enable **Debug Mode** for detailed information

## Reporting Security Issues

If you discover a security vulnerability in Wynglet:

1. **Do not** post it publicly
2. **Do not** include it in a GitHub issue
3. Visit the [Security Policy](https://github.com/chimbori/wynglet/security/policy) for reporting procedures

## Next Steps

- [Installation](/installation)
- [Form Submissions](/forms)
- [Dashboard](/dashboard)
- [Installation](/installation)
