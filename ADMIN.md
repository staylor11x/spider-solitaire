# Admin Documentation

Maintainer and administrator guide for Spider Solitaire project management tasks.

## Table of Contents

- [Issue & Ticket Management](#issue--ticket-management)
- [Release Management](#release-management)
- [Maintenance](#maintenance)

---

## Issue & Ticket Management

### Creating GitHub Issues from Tickets

This project uses a streamlined workflow for creating GitHub issues from markdown ticket files.

#### Setup

1. **Ensure GitHub CLI is installed and authenticated:**
   ```bash
   gh auth login
   ```
   
   When prompted, ensure you grant the `repo` scope for full repository access.

2. **Verify authentication:**
   ```bash
   gh auth status
   ```

#### Ticket Structure

Store all ticket markdown files in the `tickets/` directory. Each ticket should follow this structure:

```markdown
# Issue Title

## Description
Brief overview of the feature/bug/task.

## Technical Details
Implementation approach, affected files, architecture notes.

## Acceptance Criteria
- [ ] Specific, testable requirement 1
- [ ] Specific, testable requirement 2
- [ ] Update documentation

## Testing Checklist
- [ ] Manual test case 1
- [ ] Manual test case 2
```

#### Creating an Issue

**Basic usage:**
```bash
./create-issue.sh tickets/ticket-name.md
```

**With custom labels:**
```bash
./create-issue.sh tickets/ticket-name.md "enhancement,ui"
```

**With custom labels and assignee:**
```bash
./create-issue.sh tickets/ticket-name.md "enhancement,ui" "@username"
```

#### Available Labels

- `bug` - Bug fixes
- `enhancement` - New features or improvements
- `documentation` - Documentation updates

#### Script Details

The `create-issue.sh` script:
- Reads markdown file from `tickets/` directory
- Extracts title from first `# Header` (or uses filename)
- Uses everything after the header as issue body
- Creates the issue with specified labels and assignee
- Provides success/failure feedback

#### Example Workflow

1. **Create ticket markdown:**
   ```bash
   # Create tickets/my-feature.md
   ```

2. **Write ticket content** with description, technical details, acceptance criteria

3. **Create issue on GitHub:**
   ```bash
   ./create-issue.sh tickets/my-feature.md "enhancement,ui" "@me"
   ```

4. **Work on the issue** - create branch, implement, test

5. **Update CHANGELOG.md** when merging PR

## Notes

- Always use the ticket system for tracking work
- Keep issue descriptions detailed and clear
- Include acceptance criteria in all tickets
- Test thoroughly before marking issues as complete
