# Documentation Rules

Documentation is part of the product. Keep it updated when making changes.

## Required Docs

```txt
docs/PROJECT_PLAN.md
docs/ARCHITECTURE.md
docs/DEVELOPMENT_STATUS.md
docs/API_CONTRACTS.md
docs/DECISIONS_LOG.md
```

## DEVELOPMENT_STATUS.md

Must track:
- Current phase.
- Completed tasks.
- In-progress tasks.
- Pending tasks.
- Known issues.
- Technical debt.
- Next recommended step.

## DECISIONS_LOG.md

Format for each entry:

```md
## YYYY-MM-DD — Decision title

### Context
Why the decision was needed.

### Decision
What was decided.

### Consequences
Trade-offs and implications.
```

## API_CONTRACTS.md

Update when: adding an endpoint, changing payload/response shape, changing auth requirements, changing error response formats.

## Docs Must Not Lie

If something is not implemented, documentation must not claim it is implemented.

Use labels: `Planned`, `In Progress`, `Implemented`, `Deprecated`, `Blocked`.
