# Testing Rules

Add or update tests when practical. Never fake test results.

## Backend Tests

Include tests for:
- Auth service.
- Password hashing.
- JWT generation/validation.
- Clinic service logic.
- Lead service logic.
- AI prompt builder.
- AI safety validator.
- Repository behavior where practical.

Run:
```bash
go test ./...
go vet ./...
go fmt ./...
```

## Frontend Tests

Include tests for:
- Auth service.
- Guards.
- Form validation.
- Lead status transitions.
- Copy-to-clipboard behavior where practical.
- API service methods.

Run:
```bash
npm install
npm run build
npm test
npm run lint
```

## Manual Validation

After implementation, run relevant commands. If a command does not exist, report it clearly. Do not pretend it ran.
