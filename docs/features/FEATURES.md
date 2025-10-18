# Secret Santa - Feature Overview

## Project Components

### 1. CLI Application (`cmd/cli`)
Command-line interface for running Secret Santa draws from JSON files.

**Usage:**
```bash
./bin/secretsanta participants.json
```

### 2. Web Application (`cmd/web`)
Modern web interface with REST API for managing Secret Santa events.

**Usage:**
```bash
make run-web
# Visit http://localhost:8080
```

## Core Features

### Participant Management
- ✅ Add participants individually via web UI
- ✅ Upload participant lists from JSON files
- ✅ Support for multiple contact methods per participant
- ✅ Exclusion rules (e.g., couples, family members)
- ✅ Validation with detailed error reporting

### Drawing Algorithm
- ✅ **Fast random retry algorithm** - 3-26× faster than backtracking
- ✅ Handles 500+ participants efficiently
- ✅ Respects all exclusion constraints
- ✅ Guaranteed fairness (everyone gives and receives)

### Validation System
- ✅ **Separate validation** - Check configurations before drawing
- ✅ Detailed error messages and warnings
- ✅ Compatibility analysis
- ✅ Quick validation mode for boolean checks

### Notification Integration
- ✅ **Multiple notification types**: email, Slack, stdout
- ✅ **External notifier service** integration via gRPC
- ✅ **Multi-account support** via notifier service (optional account field in protobuf)
- ✅ **Archive BCC** support for record-keeping
- ✅ Fallback to built-in SMTP
- ✅ Support for multiple recipients per participant

### Web Interface
- ✅ Modern, responsive UI
- ✅ Three tabs: Create, Upload, Draw
- ✅ Drag-and-drop file upload
- ✅ **Multiple file format support**: JSON, YAML, TOML, CSV, TSV
- ✅ **Format-specific templates** with download functionality
- ✅ Real-time validation feedback
- ✅ Export results as JSON
- ✅ Toast notifications
- ✅ Validation modal with detailed results

## Technical Highlights

### Performance
- **Draw Algorithm**: 18-197 µs for 10-500 participants
- **Validation**: 2.6-3,547 µs for 10-500 participants
- **Total (validation + draw)**: < 4 ms even for 500 participants

### Architecture
- **Multi-tier**: CLI, Web UI, REST API
- **gRPC integration** with external notifier service
- **Docker support** with multi-architecture builds
- **Comprehensive testing** with benchmarks

### Best Practices
- ✅ Clean architecture with `internal/` packages
- ✅ Separation of concerns (api, draw, notification, validation)
- ✅ Comprehensive error handling
- ✅ Input validation and sanitization
- ✅ CORS support for API
- ✅ Security (non-root Docker user, input escaping)

## Configuration

### Application Config (`secretsanta.config.toml`)
```toml
[smtp]
host = "smtp.example.com"
port = "587"
username = "your-email@example.com"
password = "YOUR_PASSWORD"
from_address = "your-email@example.com"
from_name = "Secret Santa"

[notifier]
service_addr = "localhost:50051"  # Optional: external notifier service
archive_email = "archive@example.com"  # Optional: BCC for all notifications
```

### Participant Data Format
```json
[
  {
    "name": "Alice Johnson",
    "notification_type": "email",
    "contact_info": ["alice@example.com"],
    "exclusions": ["Bob Smith"]
  }
]
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/validate` | POST | Validate participant configuration |
| `/api/draw` | POST | Run Secret Santa draw |
| `/api/upload` | POST | Upload participant file (JSON, YAML, TOML, CSV, TSV) |
| `/api/export` | POST | Export results as JSON |
| `/api/template` | GET | Download template file (query param: `format=json\|yaml\|toml\|csv\|tsv`) |
| `/api/status` | GET | Get notification configuration status (available types, notifier health) |
| `/` | GET | Serve web interface |
| `/static/*` | GET | Static assets (CSS, JS) |

## Docker & Deployment

### Multi-Architecture Support
- ✅ linux/amd64 (x86_64)
- ✅ linux/arm64 (ARM 64-bit)
- ✅ linux/arm/v7 (ARM 32-bit)

### Deployment Options
1. **Standalone binary**
2. **Docker container**
3. **Docker Compose** (with notifier service)
4. **Kubernetes** (via Helm charts - roadmap)

## Development

### Building
```bash
make build          # Build CLI
make build-web      # Build web server
make build-all      # Build both
make docker-buildx-build  # Multi-arch Docker images
```

### Testing
```bash
make test           # Run all tests
make test-coverage  # Coverage report
make lint          # Run linter
```

### Running
```bash
# CLI
./bin/secretsanta participants.json

# Web server
make run-web
# or
./bin/secretsanta-web -addr :8080
```

## Future Roadmap

### Planned Features
- [ ] Kubernetes Operator with CRDs
- [ ] gRPC API alongside REST
- [ ] User authentication & multi-tenancy
- [ ] Event history & audit log
- [ ] Email template customization
- [ ] Scheduling for future draws
- [ ] Reminder notifications
- [ ] Gift suggestions integration

### Improvements
- [ ] WebSocket for real-time updates
- [ ] Progressive Web App (PWA)
- [ ] Database persistence
- [ ] Admin dashboard
- [ ] Analytics & reporting

## Contributing

### Code Organization
```
secretsanta/
├── cmd/
│   ├── cli/          # CLI application
│   └── web/          # Web server
├── internal/
│   ├── api/          # HTTP handlers
│   ├── draw/         # Draw algorithms
│   ├── notification/ # Notification logic
│   └── web/          # Web assets
├── pkg/
│   ├── config/       # Configuration
│   └── participant/  # Participant types
├── api/grpc/         # gRPC proto files
├── configs/          # Config templates
└── docker/           # Docker files
```

## Documentation

- **README.md** - Project overview and setup
- **docs/features/WEB_README.md** - Web application guide
- **docs/algorithms/VALIDATION.md** - Validation system documentation
- **docs/algorithms/ALGORITHM_ANALYSIS.md** - Performance analysis and complexity
- **docs/performance/BENCHMARK_COMPARISON.md** - Benchmark comparisons
- **docs/performance/VALIDATION_PERFORMANCE.md** - Validation performance metrics
- **docs/algorithms/BUGFIX_HALLS_THEOREM.md** - Hall's theorem bugfix documentation
- **docs/features/FEATURES.md** - This file

## License

See LICENSE file for details.
