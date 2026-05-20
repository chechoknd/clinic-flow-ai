# GitFlow Rules

## Branch Naming

Never work directly on `main` unless explicitly instructed.

Preferred branches:

```
feature/<short-description>
fix/<short-description>
chore/<short-description>
docs/<short-description>
refactor/<short-description>
```

Examples: `feature/backend-auth-jwt`, `fix/frontend-login`, `chore/docker-compose-postgres`.

## Before Starting Work

```bash
git status
git branch --show-current
git pull
```

If there are uncommitted changes not created by the agent, do not overwrite them.

## Commit Message Format

```
type(scope): short description
```

Allowed types: `feat`, `fix`, `docs`, `chore`, `refactor`, `test`, `style`, `build`, `ci`.

Examples:
- `feat(auth): add JWT login endpoint`
- `fix(frontend): protect private routes with auth guard`
- `docs(project): add MVP roadmap`
- `chore(docker): add postgres service`

## Pull Request Rules

If the environment supports PRs, create a PR instead of pushing directly to protected branches.

PR description must include: summary, files changed, how it was tested, risks, follow-up tasks.

## Push Rules

After committing, push:

```bash
git push -u origin HEAD
```

## Dirty Working Tree

If the working tree has uncommitted changes that are not yours, report them. Do not overwrite.
