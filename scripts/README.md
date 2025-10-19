# Secret Santa Scripts

This directory contains helper scripts for automating Secret Santa draws via the API.

## secretsanta-draw.sh

A simple example script that demonstrates how to use the Secret Santa API for automated draws.

### Usage

```bash
./secretsanta-draw.sh <participants.json> [server-url]
```

**Arguments:**
- `participants.json` - Path to JSON file with participant data (required)
- `server-url` - Optional: API server URL (default: http://localhost:8080)

### Examples

```bash
# Draw with local server
./secretsanta-draw.sh participants.json

# Draw with remote server
./secretsanta-draw.sh participants.json http://secretsanta.example.com

# Draw with custom port
./secretsanta-draw.sh participants.json http://localhost:3000
```

### Features

- ✅ Validates participants before drawing
- ✅ Shows validation errors and warnings
- ✅ Displays compatibility statistics
- ✅ Runs the draw and shows results
- ✅ Color-coded output for better readability
- ✅ JSON formatting with jq (if available)

### Requirements

- **curl** - Required for API calls
- **jq** - Optional but recommended for formatted output
  - macOS: `brew install jq`
  - Ubuntu/Debian: `apt-get install jq`
  - Arch: `pacman -S jq`

### What it Does

1. **Validation**: Calls `/api/validate` to check participant configuration
   - Verifies all participants have valid recipient options
   - Checks for errors (missing data, impossible configurations)
   - Shows warnings (non-critical issues)
   - Displays compatibility metrics

2. **Drawing**: If validation passes, calls `/api/draw` to run the draw
   - Performs the Secret Santa assignment
   - Sends notifications (if configured)
   - Returns the complete results

3. **Output**: Displays results in a friendly format
   - Shows who draws whom
   - Indicates notification status
   - Reports any errors

### Example Output

```
==================================================
Secret Santa Draw Tool
==================================================

ℹ Participants file: participants.json
ℹ API server: http://localhost:8080

ℹ Step 1: Validating participants...
✓ Validation passed!
  Participants: 6
  Min compatibility: 4
  Avg compatibility: 4.5

ℹ Step 2: Running Secret Santa draw...
✓ Draw completed successfully!

ℹ Results:

  Alice → Charlie
  Bob → Diana
  Charlie → Eve
  Diana → Frank
  Eve → Alice
  Frank → Bob

==================================================
✓ Secret Santa draw complete!
==================================================
```

## Creating Custom Scripts

This script serves as an example. You can create your own automation by calling the Secret Santa API directly:

### API Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/validate` | POST | Validate participants without drawing |
| `/api/draw` | POST | Run the draw and send notifications |
| `/api/status` | GET | Check notification service status |
| `/api/export` | POST | Export results as JSON |
| `/api/template` | GET | Download participant file templates |

### Example with curl

```bash
# Validate participants
curl -X POST http://localhost:8080/api/validate \
  -H "Content-Type: application/json" \
  -d @participants.json

# Run the draw
curl -X POST http://localhost:8080/api/draw \
  -H "Content-Type: application/json" \
  -d @participants.json

# Check notification status
curl http://localhost:8080/api/status
```

### Example with Python

```python
import requests
import json

# Load participants
with open('participants.json') as f:
    participants = json.load(f)

# Validate
response = requests.post(
    'http://localhost:8080/api/validate',
    json=participants
)
validation = response.json()

if validation['valid']:
    # Run draw
    response = requests.post(
        'http://localhost:8080/api/draw',
        json=participants
    )
    results = response.json()

    for p in results['participants']:
        print(f"{p['name']} → {p['recipient']}")
```

## Integration Examples

### Cron Job

Run Secret Santa draw automatically:

```bash
# Add to crontab
0 9 1 12 * /path/to/secretsanta-draw.sh /path/to/participants.json
```

### CI/CD Pipeline

```yaml
# GitHub Actions example
- name: Run Secret Santa Draw
  run: |
    ./scripts/secretsanta-draw.sh participants.json http://secretsanta-api:8080
```

### Docker

```bash
# Run script in container
docker run --rm -v $(pwd)/participants.json:/data/participants.json \
  alpine sh -c "apk add --no-cache curl && \
  wget -O /tmp/draw.sh https://raw.githubusercontent.com/igodwin/secretsanta/main/scripts/secretsanta-draw.sh && \
  chmod +x /tmp/draw.sh && \
  /tmp/draw.sh /data/participants.json http://secretsanta-api:8080"
```

## Notes

- The web UI provides a more user-friendly experience for most users
- These scripts are useful for automation, CI/CD, and testing
- Always validate before drawing to catch configuration errors early
- Notification configuration is handled server-side (not in the script)
