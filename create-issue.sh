#!/bin/bash

# Create GitHub issue from markdown file
# Usage: ./create-issue.sh <markdown-file> [labels] [assignee]
# Example: ./create-issue.sh ticket.md "enhancement,ui" "@me"

if [ $# -lt 1 ]; then
    echo "Usage: $0 <markdown-file> [labels] [assignee]"
    echo ""
    echo "Examples:"
    echo "  $0 ticket.md"
    echo "  $0 ticket.md \"enhancement,ui\""
    echo "  $0 ticket.md \"enhancement,ui\" \"@me\""
    exit 1
fi

MARKDOWN_FILE="$1"
LABELS="${2:-enhancement}"
ASSIGNEE="${3:-@me}"

# Check if file exists
if [ ! -f "$MARKDOWN_FILE" ]; then
    echo "Error: File '$MARKDOWN_FILE' not found"
    exit 1
fi

# Extract title from first H1 or use filename
TITLE=$(grep "^# " "$MARKDOWN_FILE" | head -1 | sed 's/^# //')
if [ -z "$TITLE" ]; then
    # Fallback: use filename without extension
    TITLE=$(basename "$MARKDOWN_FILE" .md)
fi

# Extract body (everything after first H1)
BODY=$(sed '1,/^# /d' "$MARKDOWN_FILE")
if [ -z "$BODY" ]; then
    # If no H1 found, use entire file
    BODY=$(cat "$MARKDOWN_FILE")
fi

echo "Creating issue..."
echo "Title: $TITLE"
echo "Labels: $LABELS"
echo "Assignee: $ASSIGNEE"
echo ""

# Create the issue
gh issue create \
  --title "$TITLE" \
  --body "$BODY" \
  --label "$LABELS" \
  --assignee "$ASSIGNEE"

if [ $? -eq 0 ]; then
    echo "✅ Issue created successfully!"
else
    echo "❌ Failed to create issue"
    exit 1
fi
