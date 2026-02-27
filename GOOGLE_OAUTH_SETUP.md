# Google OAuth Setup Guide

## Prerequisites

- Google Account
- Backend server running on `http://localhost:8080`

---

## Step-by-Step Setup

### 1. Create Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Click **Select a project** → **New Project**
3. Enter project name: `Finance Hub`
4. Click **Create**
5. Wait for project creation to complete

### 2. Enable Google+ API

1. In the left sidebar, click **APIs & Services** → **Library**
2. Search for `Google+ API`
3. Click on **Google+ API**
4. Click **Enable**

### 3. Configure OAuth Consent Screen

1. Go to **APIs & Services** → **OAuth consent screen**
2. Select **External** (unless you have a Google Workspace)
3. Click **Create**

**Fill in the form:**

**App information:**

- App name: `Finance Hub`
- User support email: `your-email@gmail.com`
- App logo: (Optional)

**App domain:**

- Application home page: `http://localhost:5173`
- Application privacy policy link: `http://localhost:5173/privacy` (optional for testing)
- Application terms of service link: `http://localhost:5173/terms` (optional for testing)

**Developer contact information:**

- Email addresses: `your-email@gmail.com`

4. Click **Save and Continue**

**Scopes:**

- Click **Add or Remove Scopes**
- Select these scopes:
    - `openid`
    - `profile`
    - `email`
- Click **Update**
- Click **Save and Continue**

**Test users (for development):**

- Click **Add Users**
- Add your test email addresses
- Click **Add**
- Click **Save and Continue**

5. Click **Back to Dashboard**

### 4. Create OAuth 2.0 Credentials

1. Go to **APIs & Services** → **Credentials**
2. Click **Create Credentials** → **OAuth client ID**

**Configure OAuth client:**

- Application type: **Web application**
- Name: `Finance Hub - Development`

**Authorized JavaScript origins:**

```
http://localhost:5173
http://localhost:8080
```

**Authorized redirect URIs:**

```
http://localhost:8080/api/v1/auth/google/callback
http://localhost:5173/auth/callback
```

3. Click **Create**

### 5. Copy Credentials

A dialog will appear with your credentials:

```
Client ID: 123456789-xxxxxxxxx.apps.googleusercontent.com
Client Secret: GOCSPX-xxxxxxxxxxxxxxxxxxxxxx
```

**⚠️ IMPORTANT:** Copy these credentials immediately!

---

## Configure Backend

### Update `.env` file

```bash
cd server
nano .env  # or open with your editor
```

Add/update these lines:

```env
# Google OAuth Configuration
GOOGLE_CLIENT_ID=123456789-xxxxxxxxx.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-xxxxxxxxxxxxxxxxxxxxxx
GOOGLE_REDIRECT_URI=http://localhost:8080/api/v1/auth/google/callback

# Frontend URL
FRONTEND_URL=http://localhost:5173
```

Save and exit.

### Restart Server

```bash
# Stop the running server (Ctrl+C)
# Start again
go run cmd/api/main.go
```

---

## Test OAuth Flow

### Method 1: Browser

1. Open browser
2. Visit: `http://localhost:8080/api/v1/auth/google?redirect_uri=http://localhost:5173`
3. You should be redirected to Google login
4. Select your Google account
5. Click **Continue** on consent screen
6. You should be redirected back with tokens in URL

### Method 2: Using Google Sign-In Button (Frontend)

#### HTML

```html
<!DOCTYPE html>
<html>
    <head>
        <script
            src="https://accounts.google.com/gsi/client"
            async
            defer
        ></script>
    </head>
    <body>
        <div
            id="g_id_onload"
            data-client_id="YOUR_CLIENT_ID"
            data-callback="handleCredentialResponse"
        ></div>
        <div class="g_id_signin" data-type="standard"></div>

        <script>
            function handleCredentialResponse(response) {
                console.log("ID Token: " + response.credential);

                // Send to your backend
                fetch("http://localhost:8080/api/v1/auth/google/token", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ id_token: response.credential }),
                })
                    .then((res) => res.json())
                    .then((data) => {
                        console.log("Auth response:", data);
                        localStorage.setItem(
                            "access_token",
                            data.data.access_token,
                        );
                        // Redirect to dashboard
                    });
            }
        </script>
    </body>
</html>
```

---

## Production Setup

### Update Google Cloud Console

1. Go to **Credentials** → Edit your OAuth client
2. Add production URIs:

**Authorized JavaScript origins:**

```
https://your-domain.com
https://api.your-domain.com
```

**Authorized redirect URIs:**

```
https://api.your-domain.com/api/v1/auth/google/callback
https://your-domain.com/auth/callback
```

3. Click **Save**

### Update Production `.env`

```env
GOOGLE_CLIENT_ID=123456789-xxxxxxxxx.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-xxxxxxxxxxxxxxxxxxxxxx
GOOGLE_REDIRECT_URI=https://api.your-domain.com/api/v1/auth/google/callback
FRONTEND_URL=https://your-domain.com
```

---

## Troubleshooting

### Error: "redirect_uri_mismatch"

**Cause:** The redirect URI in your request doesn't match the authorized URIs in Google Console.

**Solution:**

1. Go to Google Cloud Console → Credentials
2. Edit your OAuth client
3. Make sure exact URI is listed in "Authorized redirect URIs":
    - `http://localhost:8080/api/v1/auth/google/callback`
4. Wait 5 minutes for changes to propagate

### Error: "invalid_client"

**Cause:** Client ID or Secret is incorrect.

**Solution:**

1. Double-check `.env` file
2. Make sure no extra spaces or quotes
3. Restart server after updating `.env`

### Error: "access_denied"

**Cause:** User denied consent or not added as test user.

**Solution:**

1. Go to OAuth consent screen
2. Add user email to "Test users"
3. Try again with that Google account

### OAuth Screen Shows "This app isn't verified"

**Normal for development!**

**To continue:**

1. Click "Advanced"
2. Click "Go to Finance Hub (unsafe)"

**For production:**

1. Complete OAuth consent screen verification
2. Submit app for Google verification
3. This takes several days

---

## Security Best Practices

✅ **Never commit `.env` to Git**

```bash
# Make sure .gitignore has:
.env
.env.local
```

✅ **Use HTTPS in production**

- OAuth requires HTTPS for production
- Use Let's Encrypt for free SSL

✅ **Rotate secrets regularly**

- Generate new Client Secret every 90 days
- Update in production immediately

✅ **Limit redirect URIs**

- Only add URIs you control
- Remove development URIs in production

✅ **Verify tokens on backend**

- Never trust tokens from client
- Always verify signature with Google

---

## Testing Checklist

- [ ] Server starts without errors
- [ ] Can visit `/api/v1/auth/google` in browser
- [ ] Redirects to Google login page
- [ ] Can select Google account
- [ ] Consent screen appears
- [ ] Redirects back to frontend with tokens
- [ ] Tokens are valid (can access protected routes)
- [ ] User info is correct in database
- [ ] Can login again with same Google account
- [ ] Account linking works (if email exists)

---

## Getting Help

**Google OAuth Documentation:**

- [OAuth 2.0 Overview](https://developers.google.com/identity/protocols/oauth2)
- [OpenID Connect](https://developers.google.com/identity/protocols/oauth2/openid-connect)

**Common Issues:**

- Check redirect URI matches exactly (including http/https)
- Make sure API is enabled in Google Cloud Console
- Wait 5 minutes after making changes
- Clear browser cookies if issues persist

---

## Quick Reference

**OAuth URLs:**

- Authorization: `https://accounts.google.com/o/oauth2/v2/auth`
- Token Exchange: `https://oauth2.googleapis.com/token`
- User Info: `https://www.googleapis.com/oauth2/v3/userinfo`
- Public Keys: `https://www.googleapis.com/oauth2/v3/certs`

**Scopes Used:**

- `openid` - OpenID Connect
- `profile` - Basic profile info (name, picture)
- `email` - Email address

**Token Validation:**

- Audience (`aud`): Must match Client ID
- Issuer (`iss`): Must be `accounts.google.com` or `https://accounts.google.com`
- Expiry (`exp`): Must be in future
- Subject (`sub`): Google User ID

---

**Setup Complete!** 🎉

Your Google OAuth is now configured and ready to use.
