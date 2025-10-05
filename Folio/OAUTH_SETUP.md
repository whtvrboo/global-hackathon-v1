# Google OAuth Setup Guide

This guide will walk you through setting up Google OAuth for Folio.

## Step 1: Create a Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Click **Select a project** → **New Project**
3. Enter project name: "Folio" (or your preferred name)
4. Click **Create**

## Step 2: Enable Google+ API

1. In the Google Cloud Console, navigate to **APIs & Services** → **Library**
2. Search for "Google+ API" or "People API"
3. Click **Enable**

## Step 3: Configure OAuth Consent Screen

1. Navigate to **APIs & Services** → **OAuth consent screen**
2. Select **External** (unless you have a Google Workspace account)
3. Click **Create**

### Fill in the required information:

**App information:**

- App name: `Folio`
- User support email: Your email
- App logo: (Optional) Upload a logo

**App domain:**

- Application home page: `http://localhost` (for development)
- Application privacy policy link: (Optional)
- Application terms of service link: (Optional)

**Authorized domains:**

- Add: `localhost` (for development)

**Developer contact information:**

- Email addresses: Your email

4. Click **Save and Continue**
5. On the **Scopes** page, click **Add or Remove Scopes**
6. Select these scopes:
   - `.../auth/userinfo.email`
   - `.../auth/userinfo.profile`
7. Click **Update** → **Save and Continue**
8. On **Test users**, add your email address for testing
9. Click **Save and Continue**

## Step 4: Create OAuth Client ID

1. Navigate to **APIs & Services** → **Credentials**
2. Click **+ Create Credentials** → **OAuth client ID**
3. Application type: **Web application**
4. Name: `Folio Web Client`

### Authorized JavaScript origins:

```
http://localhost
http://localhost:8080
```

### Authorized redirect URIs:

```
http://localhost:8080/api/auth/google/callback
```

5. Click **Create**
6. **IMPORTANT**: Copy your **Client ID** and **Client Secret** - you'll need these!

## Step 5: Configure Folio

### Option 1: Using Environment Variables (Recommended for Production)

Create a `.env` file in the Folio directory:

```bash
# Copy the example file
cp env.example .env

# Edit with your values
nano .env  # or use your preferred editor
```

Update with your credentials:

```env
GOOGLE_CLIENT_ID=your-actual-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-actual-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
JWT_SECRET=your-random-secret-key-here
```

### Option 2: Update docker-compose.yml Directly (Development Only)

Edit `docker-compose.yml` and update the `api` service environment variables:

```yaml
environment:
  GOOGLE_CLIENT_ID: your-actual-client-id.apps.googleusercontent.com
  GOOGLE_CLIENT_SECRET: your-actual-client-secret
  JWT_SECRET: your-random-secret-key
```

⚠️ **Never commit real credentials to version control!**

## Step 6: Generate JWT Secret

Generate a secure random string for JWT_SECRET:

### Windows (PowerShell):

```powershell
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Minimum 0 -Maximum 256 }))
```

### Linux/Mac:

```bash
openssl rand -base64 32
```

## Step 7: Start Folio

```bash
# Start all services
docker-compose down
docker-compose up --build -d

# Check logs
docker-compose logs -f api
```

## Step 8: Test OAuth Flow

1. Visit: http://localhost
2. You should be redirected to `/login`
3. Click **Sign in with Google**
4. You'll be redirected to Google's consent screen
5. Grant permissions
6. You'll be redirected back to Folio
7. You should now be logged in!

## Production Deployment

When deploying to production:

1. **Update OAuth Consent Screen:**

   - Change from "Testing" to "In Production"
   - Add your production domain to authorized domains

2. **Update OAuth Client:**

   - Add production URLs to Authorized JavaScript origins:
     ```
     https://yourdomain.com
     ```
   - Add production callback to Authorized redirect URIs:
     ```
     https://yourdomain.com/api/auth/google/callback
     ```

3. **Update Environment Variables:**

   ```env
   FRONTEND_URL=https://yourdomain.com
   GOOGLE_REDIRECT_URL=https://yourdomain.com/api/auth/google/callback
   JWT_SECRET=use-a-strong-random-secret-in-production
   ```

4. **Security Best Practices:**
   - Use HTTPS in production
   - Store secrets in your platform's secret management (Railway Secrets, Fly.io Secrets, etc.)
   - Rotate JWT_SECRET periodically
   - Enable rate limiting

## Troubleshooting

### "redirect_uri_mismatch" Error

This means the redirect URI doesn't match what's configured in Google Cloud Console.

**Solution:**

1. Check the exact URL in the error message
2. Add that exact URL to Authorized redirect URIs in Google Cloud Console
3. Wait a few minutes for changes to propagate

### "invalid_client" Error

This means the Client ID or Client Secret is incorrect.

**Solution:**

1. Double-check your credentials in Google Cloud Console
2. Verify no extra spaces or characters when copying
3. Ensure you're using the Web Application credentials (not Android/iOS)

### "access_denied" Error

The user denied permission or your app is not approved.

**Solution:**

1. Make sure you added your email as a test user
2. Try again and grant all requested permissions

### OAuth Works But User Not Saved

Check backend logs:

```bash
docker-compose logs api | grep -i "error"
```

Common issues:

- Database connection failed
- Missing `username` field (check migrations)
- Unique constraint violation

### Token Expired Immediately

Check JWT_SECRET is consistent:

```bash
docker-compose exec api env | grep JWT_SECRET
```

### CORS Errors

Ensure CORS is enabled in the Echo middleware (already configured in main.go):

```go
e.Use(middleware.CORS())
```

## Testing Without Real Google OAuth

For quick testing without setting up OAuth:

1. Use the development JWT secret
2. Create a user manually in the database
3. Generate a JWT token:
   ```bash
   # Use an online JWT generator with HS256
   # Payload: {"user_id": "some-uuid", "email": "test@example.com"}
   # Secret: dev-secret-change-in-production
   ```
4. Use the token in API requests:
   ```bash
   curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8080/api/me
   ```

## Additional Resources

- [Google OAuth 2.0 Documentation](https://developers.google.com/identity/protocols/oauth2)
- [Google Cloud Console](https://console.cloud.google.com/)
- [OAuth 2.0 Playground](https://developers.google.com/oauthplayground/)

---

**Need Help?** Check the backend logs: `docker-compose logs -f api`
