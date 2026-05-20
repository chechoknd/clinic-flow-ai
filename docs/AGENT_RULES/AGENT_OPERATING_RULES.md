# Agent Operating Rules

## Golden Rule

Before making any change, the agent must understand the current project context, read the relevant documentation, inspect the existing code, and produce a safe implementation plan.

No agent is allowed to make large, blind, destructive, or architectural changes without first validating the existing structure.

When in doubt, the agent must preserve the existing architecture and ask for clarification before changing project direction.

## Mandatory Behavior Before Changes

1. Read AGENTS.md, README.md, docs/DEVELOPMENT_STATUS.md.
2. Read task-specific documentation in docs/AGENT_RULES/.
3. Inspect the existing code and structure.
4. Understand project phases and scope.
5. Produce an implementation plan.

## How to Inspect the Repository

- `git status` — check dirty state.
- `git branch --show-current` — know current branch.
- `pwd` — confirm working directory.
- `ls -la` — list root contents.
- Read relevant files before editing.

## How to Avoid Destructive Changes

- Always read the file before editing.
- Never overwrite existing files without reporting it.
- Preserve existing architecture.
- Never rewrite the entire project without permission.
- Never delete documentation.

## How to Report Results

After finishing a task, respond with:

```md
## Summary
Short explanation.

## Files Changed
- `path`: explanation.

## Validation
Commands executed / not executed.

## Risks or Notes
Any concern.

## Next Suggested Step
One clear step.
```

## Forbidden Agent Behaviors

- Rewrite entire project without permission.
- Change the stack without permission.
- Add medical-record or diagnosis logic.
- Add payment systems in the MVP.
- Add WhatsApp Cloud API in the MVP.
- Add RAG as required MVP infrastructure.
- Commit secrets.
- Delete documentation.
- Ignore existing architecture.
- Push directly to main unless instructed.
- Mix frontend and backend responsibilities.
- Fake test results.
- Claim a feature is complete without validation.
- Create huge files with unrelated responsibilities.
- Add dependencies without justification.
- Use real-looking patient data examples.
- Store sensitive health information.
- Log private lead conversations unnecessarily.

## Preferred Agent Behaviors

- Work incrementally.
- Keep changes small.
- Explain trade-offs.
- Preserve project direction.
- Update docs continuously.
- Prefer simple implementation.
- Validate before refactoring.
- Create reusable but not over-engineered code.
- Keep UI clean.
- Optimize for receptionist workflows.
- Protect clinic data.
- Protect patient safety.
- Protect the repository from secrets.

## Definition of Done

A task is complete only when:

- Code compiles.
- Relevant tests pass or missing tests are clearly reported.
- No secrets are committed.
- Architecture rules are respected.
- Documentation is updated if needed.
- API contracts are updated if needed.
- Database migrations are included if needed.
- Frontend and backend remain aligned.
- Change does not introduce excluded MVP features.
- Agent reports exactly what changed.
