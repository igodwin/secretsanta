#!/bin/bash
#
# Secret Santa CLI - Example script for running draws via API
#
# Usage:
#   ./secretsanta-draw.sh <participants.json> [server-url]
#
# Example:
#   ./secretsanta-draw.sh participants.json
#   ./secretsanta-draw.sh participants.json http://localhost:8080
#

set -e

# Configuration
PARTICIPANTS_FILE="${1:-}"
SERVER_URL="${2:-http://localhost:8080}"
API_VALIDATE="${SERVER_URL}/api/validate"
API_DRAW="${SERVER_URL}/api/draw"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_error() {
    echo -e "${RED}ERROR: $1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Check arguments
if [ -z "$PARTICIPANTS_FILE" ]; then
    echo "Usage: $0 <participants.json> [server-url]"
    echo ""
    echo "Arguments:"
    echo "  participants.json  Path to JSON file with participant data"
    echo "  server-url         Optional: API server URL (default: http://localhost:8080)"
    echo ""
    echo "Example:"
    echo "  $0 participants.json"
    echo "  $0 participants.json http://secretsanta.example.com"
    exit 1
fi

# Check if file exists
if [ ! -f "$PARTICIPANTS_FILE" ]; then
    print_error "File not found: $PARTICIPANTS_FILE"
    exit 1
fi

# Check if jq is available (optional but recommended)
if ! command -v jq &> /dev/null; then
    print_warning "jq is not installed. Output will not be formatted."
    print_info "Install jq for better output: brew install jq (macOS) or apt-get install jq (Linux)"
    JQ_AVAILABLE=false
else
    JQ_AVAILABLE=true
fi

# Check if curl is available
if ! command -v curl &> /dev/null; then
    print_error "curl is required but not installed"
    exit 1
fi

echo "=================================================="
echo "Secret Santa Draw Tool"
echo "=================================================="
echo ""
print_info "Participants file: $PARTICIPANTS_FILE"
print_info "API server: $SERVER_URL"
echo ""

# Step 1: Validate participants
print_info "Step 1: Validating participants..."
VALIDATE_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_VALIDATE" \
    -H "Content-Type: application/json" \
    -d @"$PARTICIPANTS_FILE")

HTTP_CODE=$(echo "$VALIDATE_RESPONSE" | tail -n 1)
RESPONSE_BODY=$(echo "$VALIDATE_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" != "200" ]; then
    print_error "Validation request failed with HTTP $HTTP_CODE"
    if [ "$JQ_AVAILABLE" = true ]; then
        echo "$RESPONSE_BODY" | jq '.'
    else
        echo "$RESPONSE_BODY"
    fi
    exit 1
fi

# Parse validation result
if [ "$JQ_AVAILABLE" = true ]; then
    VALID=$(echo "$RESPONSE_BODY" | jq -r '.valid')
    ERRORS=$(echo "$RESPONSE_BODY" | jq -r '.errors // [] | length')
    WARNINGS=$(echo "$RESPONSE_BODY" | jq -r '.warnings // [] | length')
    TOTAL=$(echo "$RESPONSE_BODY" | jq -r '.total_participants // 0')
else
    VALID=$(echo "$RESPONSE_BODY" | grep -o '"valid":[^,}]*' | cut -d':' -f2 | tr -d ' ')
    ERRORS=0
    WARNINGS=0
    TOTAL=0
fi

if [ "$VALID" != "true" ]; then
    print_error "Validation failed!"
    echo ""
    if [ "$JQ_AVAILABLE" = true ]; then
        echo "Errors:"
        echo "$RESPONSE_BODY" | jq -r '.errors[]' | sed 's/^/  - /'
        echo ""
        if [ "$WARNINGS" -gt 0 ]; then
            echo "Warnings:"
            echo "$RESPONSE_BODY" | jq -r '.warnings[]' | sed 's/^/  - /'
        fi
    else
        echo "$RESPONSE_BODY"
    fi
    exit 1
fi

print_success "Validation passed!"
if [ "$JQ_AVAILABLE" = true ]; then
    echo "  Participants: $TOTAL"
    echo "  Min compatibility: $(echo "$RESPONSE_BODY" | jq -r '.min_compatibility // "N/A"')"
    echo "  Avg compatibility: $(echo "$RESPONSE_BODY" | jq -r '.avg_compatibility // "N/A"')"

    if [ "$WARNINGS" -gt 0 ]; then
        echo ""
        print_warning "Warnings found (draw can still proceed):"
        echo "$RESPONSE_BODY" | jq -r '.warnings[]' | sed 's/^/  - /'
    fi
fi
echo ""

# Step 2: Run the draw
print_info "Step 2: Running Secret Santa draw..."
DRAW_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_DRAW" \
    -H "Content-Type: application/json" \
    -d @"$PARTICIPANTS_FILE")

HTTP_CODE=$(echo "$DRAW_RESPONSE" | tail -n 1)
RESPONSE_BODY=$(echo "$DRAW_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" != "200" ]; then
    print_error "Draw request failed with HTTP $HTTP_CODE"
    if [ "$JQ_AVAILABLE" = true ]; then
        echo "$RESPONSE_BODY" | jq '.'
    else
        echo "$RESPONSE_BODY"
    fi
    exit 1
fi

# Parse draw result
if [ "$JQ_AVAILABLE" = true ]; then
    SUCCESS=$(echo "$RESPONSE_BODY" | jq -r '.success')
else
    SUCCESS=$(echo "$RESPONSE_BODY" | grep -o '"success":[^,}]*' | cut -d':' -f2 | tr -d ' ')
fi

if [ "$SUCCESS" != "true" ]; then
    print_error "Draw failed!"
    echo "$RESPONSE_BODY"
    exit 1
fi

print_success "Draw completed successfully!"
echo ""

# Display results
print_info "Results:"
echo ""
if [ "$JQ_AVAILABLE" = true ]; then
    echo "$RESPONSE_BODY" | jq -r '.participants[] | "  \(.name) → \(.recipient)"'

    # Check notification status
    NOTIFICATION_STATUS=$(echo "$RESPONSE_BODY" | jq -r '.notification_status // "unknown"')
    if [ "$NOTIFICATION_STATUS" != "null" ] && [ "$NOTIFICATION_STATUS" != "unknown" ]; then
        echo ""
        print_info "Notification status: $NOTIFICATION_STATUS"
    fi
else
    echo "$RESPONSE_BODY"
fi

echo ""
echo "=================================================="
print_success "Secret Santa draw complete!"
echo "=================================================="
