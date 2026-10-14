
---
name: verify-change
description: Verify completed Autofort code changes, including diffs, tests, regressions and integration risks. Use after implementation or before accepting a patch.
---

# Verify Change

## Procedure

1. Inspect Git status and identify changed files.
2. Distinguish pre-existing changes from task changes.
3. Read the relevant diff.
4. Identify affected modules and API contracts.
5. Run applicable verification.
6. Review correctness and potential regressions.
7. Report findings with file references.

## Backend verification

Working directory: backend/go-app

Run:
- go build ./...
- go test ./...
- go vet ./...

Check:
- context propagation
- error handling
- API compatibility
- SQL and transaction behavior
- migration compatibility, when relevant

## Frontend verification

Working directory: frontend

Inspect package.json before running scripts.

Check:
- Vue/Nuxt component contracts
- async state handling
- API compatibility
- form behavior
- type checks and tests, when available

Do not treat static generation as proof of
correct browser behavior.

## Integration

For cross-service changes:
- Check Docker configuration.
- Verify affected service contracts.
- Distinguish compilation from runtime verification.

Never run migrations against production.
Use disposable test databases.

## Report

Provide:
- Findings ordered by severity
- Checks executed and outcomes
- Checks not executed
- Remaining risks
- Final recommendation: PASS / PASS WITH RISKS / FAIL

Do not modify the reviewed code unless explicitly asked.
