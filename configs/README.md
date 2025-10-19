# Configuration

## Quick Start

The Secret Santa application can run **without any configuration** - it will default to stdout (console) notifications.

To enable email or other notification types, create a config file.

## Config File Locations

The application looks for `secretsanta.config` in these locations (in order):

1. `/etc/secretsanta/secretsanta.config` (system-wide)
2. `~/.secretsanta/secretsanta.config` (user-specific)
3. `./secretsanta.config` (current working directory)
4. `<binary-directory>/secretsanta.config` (same directory as the executable)

**Note:** The file has no extension but uses YAML format internally.

## Example Configs

### Minimal (No Config Needed)

Just run the application - notifications will be printed to console (stdout).

```bash
./bin/secretsanta-web
# Available notifications: stdout only
```

### Example 1: Gmail SMTP

Create `secretsanta.config`:

```yaml
smtp:
  host: "smtp.gmail.com"
  port: "587"
  username: "your-email@gmail.com"
  password: "your-app-password"  # Generate at https://myaccount.google.com/apppasswords
  from_address: "your-email@gmail.com"
  from_name: "Secret Santa"
```

**Important for Gmail users:**
- You must use an [App Password](https://myaccount.google.com/apppasswords), not your regular password
- Enable 2-Factor Authentication first
- Create an App Password specifically for this application

### Example 2: Generic SMTP Server

```yaml
smtp:
  host: "smtp.example.com"
  port: "587"
  username: "notifications@example.com"
  password: "your-smtp-password"
  from_address: "notifications@example.com"
  from_name: "Holiday Gift Exchange"
```

### Example 3: External Notifier Service

For advanced setups with multiple notification types (Slack, Ntfy, etc.):

```yaml
notifier:
  service_addr: "localhost:50051"
  archive_email: "archive@example.com"  # Optional: BCC all notifications
```

See the [notifier service documentation](https://github.com/igodwin/notifier) for setup.

### Example 4: SMTP with Archiving

```yaml
smtp:
  host: "smtp.gmail.com"
  port: "587"
  username: "santa@example.com"
  password: "your-app-password"
  from_address: "santa@example.com"
  from_name: "Secret Santa"

notifier:
  archive_email: "archive@example.com"  # BCC all emails to this address
```

## Configuration Reference

### SMTP Section

| Field | Required | Description | Example |
|-------|----------|-------------|---------|
| `host` | Yes* | SMTP server hostname | `smtp.gmail.com` |
| `port` | Yes* | SMTP server port | `587` (TLS) or `465` (SSL) |
| `username` | Yes* | SMTP username | `your-email@gmail.com` |
| `password` | Yes* | SMTP password | `your-app-password` |
| `from_address` | Yes* | From email address | `santa@example.com` |
| `from_name` | No | From display name | `Secret Santa` (default) |
| `identity` | No | SMTP identity (rarely needed) | Usually empty |

*Required only if you want to use email notifications

### Notifier Section

| Field | Required | Description | Example |
|-------|----------|-------------|---------|
| `service_addr` | No | External notifier gRPC address | `localhost:50051` |
| `archive_email` | No | BCC address for all notifications | `archive@example.com` |

## Testing Your Configuration

After creating your config file:

1. Start the web server:
   ```bash
   ./bin/secretsanta-web
   ```

2. Check the logs for:
   ```
   Loaded config from: /path/to/secretsanta.config
   ```

3. Open http://localhost:8080 in your browser

4. Check the footer - it will show available notification types:
   - 📧 Email (if SMTP configured)
   - 💬 Slack (if notifier with Slack configured)
   - 💻 Stdout (always available)

## Common Issues

### Gmail: "Username and Password not accepted"

**Solution:** You need an App Password, not your regular Gmail password.
1. Enable 2FA on your Google account
2. Go to https://myaccount.google.com/apppasswords
3. Generate a new App Password for "Mail"
4. Use that 16-character password in your config

### "Config file not found"

**This is fine!** The application will use defaults (stdout notifications only).

To add email support, create `secretsanta.config` in one of these locations:
- Current directory: `./secretsanta.config`
- Binary directory: `bin/secretsanta.config` (next to the executable)
- Your home: `~/.secretsanta/secretsanta.config`
- System-wide: `/etc/secretsanta/secretsanta.config`

### External notifier service unreachable

Check the footer - it will show "✗ notifier" if the service is configured but unreachable.

1. Make sure the notifier service is running
2. Verify the `service_addr` in your config
3. Test connection: `curl http://localhost:50051` (should connect even if response is garbled)

## Environment Variables

You can also use environment variables (they override config file values):

```bash
export SMTP_HOST=smtp.gmail.com
export SMTP_PORT=587
export SMTP_USERNAME=your-email@gmail.com
export SMTP_PASSWORD=your-app-password
export SMTP_FROM_ADDRESS=your-email@gmail.com
export NOTIFIER_SERVICE_ADDR=localhost:50051

./bin/secretsanta-web
```

## Security Best Practices

1. **Never commit config files with real credentials to git**
2. Use App Passwords instead of real passwords when possible
3. Restrict file permissions:
   ```bash
   chmod 600 secretsanta.config
   ```
4. For production, use the external notifier service with proper secret management
5. The config file in the binary directory is convenient for deployment
