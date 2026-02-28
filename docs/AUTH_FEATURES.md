# Authentication Features - Finance Hub API

## Overview

This document describes the complete authentication system implemented in Finance Hub API, including email verification, password reset, and security features.

## Features Implemented

### 1. Email-Based Authentication

#### Registration & Login

- User registration with email and password
- Password strength validation (min 8 chars, uppercase, lowercase, numbers)
- Secure password hashing using bcrypt
- JWT-based authentication (access + refresh tokens)
- Auto email verification link sent upon registration

#### Email Verification

- Automatic verification email sent after registration
- Secure token generation using crypto/rand (32-byte hex)
- Token stored in MongoDB with TTL (24 hours)
- One-time use tokens (marked as used after verification)
- Resend verification email functionality
- Beautiful HTML email templates

**Endpoints:**

- `POST /api/v1/auth/send-verification-email` - Send verification email
- `POST /api/v1/auth/verify-email` - Verify email with token
- `POST /api/v1/auth/resend-verification-email` - Resend verification email

### 2. Password Reset

Complete password reset flow with secure token system:

1. User requests password reset with email
2. System generates secure reset token
3. Email sent with reset link
4. User clicks link and enters new password
5. Password updated, token invalidated

**Endpoints:**

- `POST /api/v1/auth/forgot-password` - Request password reset
- `POST /api/v1/auth/reset-password` - Reset password with token

**Security Features:**

- Reset tokens expire after 24 hours
- One-time use tokens
- Old tokens invalidated when new one created
- Secure token generation (crypto/rand)

### 3. Google OAuth Integration

Dual OAuth flow support:

- Redirect-based OAuth flow (traditional)
- Direct ID token verification (for SPA)

**Endpoints:**

- `GET /api/v1/auth/google` - Initiate OAuth flow
- `GET /api/v1/auth/google/callback` - OAuth callback
- `POST /api/v1/auth/google/token` - Verify Google ID token

### 4. Rate Limiting

Two-tier rate limiting system:

#### Strict Rate Limiting (5 requests/minute)

Applied to sensitive authentication endpoints:

- `POST /auth/login`
- `POST /auth/register`
- `POST /auth/forgot-password`
- `POST /auth/reset-password`
- `POST /auth/send-verification-email`
- `POST /auth/resend-verification-email`

#### Moderate Rate Limiting (60 requests/minute)

Applied to all other API endpoints.

**Implementation:**

- In-memory rate limiter
- IP-based tracking
- Automatic cleanup of old visitor records
- Returns 429 Too Many Requests when limit exceeded

### 5. Email Service

Full-featured email service with:

- SMTP configuration support
- HTML email templates
- Development mode (logs instead of sending)
- Templates for:
    - Email verification
    - Password reset

**Configuration:**

```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=noreply@financehub.com
SMTP_FROM_NAME=Finance Hub
FRONTEND_URL=http://localhost:5173
```

**Development Mode:**
When SMTP credentials are not configured, emails are logged to console instead of being sent. Perfect for development!

## Database Schema

### verification_tokens Collection

```json
{
    "_id": "ObjectId",
    "user_id": "string",
    "token": "string (hex, 32 bytes)",
    "type": "email_verification | password_reset",
    "expires_at": "timestamp",
    "used": "boolean",
    "created_at": "timestamp"
}
```

**Indexes:**

- `token` - unique, for fast lookup
- `user_id` + `type` - for finding user's tokens
- `expires_at` - TTL index for auto-cleanup

## Security Best Practices

✅ **Implemented:**

- Password hashing with bcrypt
- Secure token generation (crypto/rand)
- JWT access & refresh tokens
- Rate limiting on sensitive endpoints
- Token expiration (24 hours)
- One-time use tokens
- Email verification before full access
- IP-based rate limiting

✅ **Token Security:**

- 32-byte cryptographically secure random tokens
- Tokens stored hashed in database
- Automatic expiration and cleanup
- Tokens invalidated after use
- Old tokens removed when new one created

✅ **Rate Limiting:**

- Prevents brute force attacks
- Prevents email spam
- IP-based tracking
- Memory-efficient implementation

## Frontend Integration

### Email Verification Flow

```typescript
// After registration, user receives email with link:
// http://localhost:5173/auth/verify-email?token=abc123...

// Frontend extracts token and calls:
await fetch("/api/v1/auth/verify-email", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ token: "abc123..." }),
});
```

### Password Reset Flow

```typescript
// 1. Request reset
await fetch("/api/v1/auth/forgot-password", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email: "user@example.com" }),
});

// 2. User receives email with link:
// http://localhost:5173/auth/reset-password?token=xyz789...

// 3. Frontend submits new password
await fetch("/api/v1/auth/reset-password", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
        token: "xyz789...",
        new_password: "NewSecurePass123!",
    }),
});
```

### Resend Verification Email

```typescript
await fetch("/api/v1/auth/resend-verification-email", {
    method: "POST",
    headers: {
        Authorization: `Bearer ${accessToken}`,
    },
});
```

## Email Templates

### Verification Email

Vietnamese language template with:

- Clean, professional design
- Prominent verification button
- Security notice
- Token expiration information
- Support links

### Password Reset Email

Vietnamese language template with:

- Clear reset instructions
- Security warnings
- Prominent reset button
- Token expiration notice
- "Didn't request this?" notice

## Configuration

### Required Environment Variables

```env
# Server
PORT=8080
FRONTEND_URL=http://localhost:5173

# Database
MONGODB_URI=mongodb+srv://...
MONGODB_DATABASE=finance_hub

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h

# Email (Optional - logs if not configured)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=noreply@financehub.com
SMTP_FROM_NAME=Finance Hub

# Google OAuth
GOOGLE_CLIENT_ID=...
GOOGLE_CLIENT_SECRET=...
GOOGLE_REDIRECT_URI=http://localhost:8080/api/v1/auth/google/callback
```

## Testing

### Manual Testing

1. **Registration:**

    ```bash
    curl -X POST http://localhost:8080/api/v1/auth/register \
      -H "Content-Type: application/json" \
      -d '{
        "email": "test@example.com",
        "password": "Test1234!",
        "full_name": "Test User"
      }'
    ```

2. **Verify Email:**

    ```bash
    curl -X POST http://localhost:8080/api/v1/auth/verify-email \
      -H "Content-Type: application/json" \
      -d '{"token": "your-token-here"}'
    ```

3. **Request Password Reset:**

    ```bash
    curl -X POST http://localhost:8080/api/v1/auth/forgot-password \
      -H "Content-Type: application/json" \
      -d '{"email": "test@example.com"}'
    ```

4. **Reset Password:**

    ```bash
    curl -X POST http://localhost:8080/api/v1/auth/reset-password \
      -H "Content-Type: application/json" \
      -d '{
        "token": "your-reset-token",
        "new_password": "NewPass1234!"
      }'
    ```

5. **Test Rate Limiting:**
    ```bash
    # Send 6 requests quickly - 6th should return 429
    for i in {1..6}; do
      curl -X POST http://localhost:8080/api/v1/auth/login \
        -H "Content-Type: application/json" \
        -d '{"email":"test@test.com","password":"test"}';
    done
    ```

## Architecture

### Service Layer

- **AuthService**: Business logic for authentication
- **EmailService**: Email sending with SMTP
- **TokenGenerator**: Secure token generation

### Repository Layer

- **UserRepository**: User CRUD operations
- **VerificationTokenRepository**: Token management

### Handler Layer

- **AuthHandler**: HTTP request handlers

### Middleware Layer

- **RateLimiter**: Request rate limiting
- **Auth**: JWT validation

## Future Enhancements

Potential improvements:

- [ ] Two-factor authentication (2FA)
- [ ] Email change with verification
- [ ] Session management & device tracking
- [ ] OAuth with more providers (Facebook, GitHub)
- [ ] Redis-based rate limiting (distributed)
- [ ] Email queue with background workers
- [ ] Account lockout after failed attempts
- [ ] IP whitelist/blacklist
- [ ] Security audit logs

## Troubleshooting

### Emails Not Sending

**Development Mode:**

- Check console logs - emails are logged if SMTP not configured
- Look for "📧 [DEV MODE]" in logs

**Production:**

- Verify SMTP credentials in .env
- Check SMTP_HOST and SMTP_PORT
- For Gmail: Use App Password, not regular password
- Check firewall/network allows SMTP port (587)

### Rate Limiting Issues

- Clear by waiting 1 minute
- Restart server to clear in-memory limiter
- Check IP address being tracked correctly

### Token Errors

- Tokens expire after 24 hours - request new one
- Tokens are one-time use - can't reuse
- Check MongoDB for token record and expiry

## Production Considerations

Before deploying to production:

1. **Environment Variables:**
    - Use strong JWT_SECRET (32+ random characters)
    - Configure real SMTP credentials
    - Set FRONTEND_URL to production domain
    - Use HTTPS URLs

2. **Rate Limiting:**
    - Consider Redis instead of in-memory
    - Adjust limits based on actual usage
    - Monitor for abuse

3. **Email:**
    - Use professional email service (SendGrid, AWS SES)
    - Configure SPF, DKIM, DMARC records
    - Monitor delivery rates

4. **Security:**
    - Enable CORS only for your domain
    - Use HTTPS everywhere
    - Implement CSP headers
    - Regular security audits
    - Monitor for suspicious activity

5. **Database:**
    - Ensure TTL indexes are working
    - Monitor token collection size
    - Regular backups

## Summary

Complete, production-ready authentication system with:

- ✅ Email/password authentication
- ✅ Email verification
- ✅ Password reset
- ✅ Google OAuth
- ✅ Rate limiting
- ✅ Secure token management
- ✅ Beautiful email templates
- ✅ Development-friendly (no SMTP needed)

The system is secure, scalable, and ready for production use.
