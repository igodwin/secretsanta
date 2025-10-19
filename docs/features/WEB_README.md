# Secret Santa Web Application

A modern, user-friendly web interface for managing Secret Santa gift exchanges.

## Features

- ✅ **Create Participants** - Add participants one-by-one through an intuitive form
- ✅ **Upload Files** - Import participant lists from JSON files via drag-and-drop or file browser
- ✅ **Validation** - Real-time validation of participant configurations
- ✅ **Draw Execution** - Run the Secret Santa draw with a single click
- ✅ **Export Results** - Download draw results as JSON
- ✅ **Responsive Design** - Works on desktop and mobile devices

## Quick Start

### Run the Web Server

```bash
# Development mode
make run-web

# Or build and run
make build-web
./bin/secretsanta-web
```

The server will start on `http://localhost:8080` by default.

### Custom Port

```bash
./bin/secretsanta-web -addr :3000
```

## API Endpoints

The web application exposes the following REST API endpoints:

### `POST /api/validate`

Validate participant configuration without running a draw.

**Request Body:**
```json
[
  {
    "name": "Alice",
    "notification_type": "email",
    "contact_info": ["alice@example.com"],
    "exclusions": ["Bob"]
  }
]
```

**Response:**
```json
{
  "valid": true,
  "errors": [],
  "warnings": [],
  "min_compatibility": 2,
  "avg_compatibility": 3.5,
  "total_participants": 4
}
```

### `POST /api/draw`

Perform the Secret Santa draw.

**Request Body:** Same as validate endpoint

**Response:**
```json
{
  "success": true,
  "participants": [
    {
      "name": "Alice",
      "notification_type": "email",
      "contact_info": ["alice@example.com"],
      "exclusions": ["Bob"],
      "recipient": "Carol"
    }
  ]
}
```

### `POST /api/upload`

Upload a JSON file containing participant data.

**Form Data:**
- `file`: JSON file

**Response:**
```json
{
  "success": true,
  "participants": [...],
  "validation": {...}
}
```

### `POST /api/export`

Export draw results as JSON file.

## User Guide

### Creating Participants

1. Navigate to the **Create Participants** tab
2. Fill in the participant form:
   - **Name** (required): Participant's full name
   - **Notification Type**: email, slack, or stdout
   - **Contact Info** (required): Email addresses or usernames (comma-separated for multiple)
   - **Exclusions** (optional): Names of people this participant should NOT be assigned to
3. Click **Add Participant**
4. Repeat for all participants

### Uploading a File

1. Navigate to the **Upload File** tab
2. Either:
   - Drag and drop a file onto the upload area (supports JSON, YAML, TOML, CSV, TSV)
   - Click **Browse Files** to select a file
3. The file will be automatically validated
4. Valid participants will be loaded into the system
5. Switch between format tabs to see examples and download templates for each format

**Supported Formats:**

- **JSON** - Standard JSON array format
- **YAML** - Human-friendly YAML format
- **CSV** - Comma-separated values (great for Excel/Google Sheets)
- **TSV** - Tab-separated values (Excel-compatible)
- **TOML** - Configuration file format

**JSON Format Example:**
```json
[
  {
    "name": "Alice Johnson",
    "notification_type": "email",
    "contact_info": ["alice@example.com"],
    "exclusions": ["Bob Smith"]
  },
  {
    "name": "Bob Smith",
    "notification_type": "email",
    "contact_info": ["bob@example.com"],
    "exclusions": ["Alice Johnson"]
  }
]
```

**CSV Format Example:**
```csv
name,notification_type,contact_info,exclusions
Alice Johnson,email,alice@example.com,Bob Smith
Bob Smith,email,bob@example.com,Alice Johnson
Carol Davis,slack,@carol,
```

**Note:** For CSV/TSV, use semicolons to separate multiple values within a cell (e.g., `alice@work.com; alice@personal.com`)

### Validating Configuration

Before running a draw, you can validate your configuration:

1. Click **Validate Configuration** button
2. Review the validation results:
   - ✅ **Valid**: Configuration is ready for drawing
   - ❌ **Invalid**: Shows specific errors that need to be fixed
   - ⚠️ **Warnings**: Non-critical issues (can still proceed)

### Running the Draw

1. Navigate to the **Run Draw** tab
2. Review the participant count
3. Click **Run Draw**
4. View the results showing who draws whom
5. Optionally **Export Results** as JSON

## Validation Rules

The system validates:

### Errors (Must Fix)
- ❌ At least 2 participants required
- ❌ No duplicate names
- ❌ No participant can exclude everyone (must have at least one valid recipient)

### Warnings (Can Proceed)
- ⚠️ Missing contact information
- ⚠️ Excluding non-existent participants
- ⚠️ Very low compatibility (might be hard to find valid assignment)

## Common Use Cases

### Couples Exclusion

When organizing a Secret Santa with couples, exclude partners from each other:

```json
[
  {
    "name": "Alice",
    "exclusions": ["Bob"]  // Alice's partner
  },
  {
    "name": "Bob",
    "exclusions": ["Alice"]  // Bob's partner
  }
]
```

### Family Groups

Exclude family members from drawing each other:

```json
[
  {
    "name": "Parent1",
    "exclusions": ["Parent2", "Child1", "Child2"]
  },
  {
    "name": "Parent2",
    "exclusions": ["Parent1", "Child1", "Child2"]
  }
]
```

### Multiple Contact Methods

Support multiple notification methods:

```json
{
  "name": "Alice",
  "notification_type": "email",
  "contact_info": [
    "alice@work.com",
    "alice@personal.com"
  ]
}
```

## Notification Status

The web interface footer displays the currently available notification types based on your configuration:

- **Email** (📧) - Shown when SMTP is configured or external notifier has email enabled
- **Slack** (💬) - Shown when external notifier has Slack configured
- **Ntfy** (🔔) - Shown when external notifier has Ntfy configured
- **Stdout** (💻) - Always available (console output)

The status also indicates:
- **✓ notifier** - External notifier service is connected and healthy
- **✗ notifier** - External notifier service is configured but unavailable
- **✓ SMTP** - Using built-in SMTP configuration

This helps you know which notification types are available before creating participants.

## Development

### Project Structure

```
internal/
├── api/
│   └── handlers.go          # HTTP API handlers
└── web/
    └── static/
        ├── index.html       # Main HTML page
        ├── css/
        │   └── styles.css   # Application styles
        └── js/
            └── app.js       # Frontend JavaScript
```

### Building

```bash
# Build CLI
make build

# Build web server
make build-web

# Build both
make build-all
```

### Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage
```

## Integration with Notifier Service

The web application can integrate with the external notifier service for sending notifications:

1. Set the notifier service address in config:
   ```yaml
   notifier:
     service_addr: "localhost:50051"
     archive_email: "archive@example.com"
   ```

2. When drawing, notifications will be sent via the notifier service

### Multi-Account Support

The notifier service now supports multi-account configurations (added in latest protobuf update). This allows the notifier service to manage multiple email accounts or notification providers. The `account` field in the protobuf spec is optional - if not specified, the notifier service will use its default account configuration.

For secret santa, the account field is not required as we typically use a single notification configuration.

## Deployment

### Standalone

```bash
./bin/secretsanta-web -addr :8080
```

### Docker

```bash
docker build -t secretsanta-web -f docker/Dockerfile .
docker run -p 8080:8080 secretsanta-web
```

### Docker Compose

```bash
docker-compose up web
```

## Troubleshooting

### Port Already in Use

```bash
# Use a different port
./bin/secretsanta-web -addr :3000
```

### Static Files Not Loading

Make sure you're running the server from the project root directory where `internal/web/static` exists.

### CORS Issues

The server includes CORS middleware for development. For production, configure appropriate CORS settings.

## License

See LICENSE file in the project root.
