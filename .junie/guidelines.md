# Guidelines

## Code style
- Use `go fmt` for formatting; adhere to idiomatic Go.
- Keep functions small and cohesive; use early returns; choose explicit names; write clear error messages with context.
- If `golangci-lint` is configured, ensure a clean run; otherwise, at minimum, run `go vet ./...`.

## Testing
- Add or extend unit tests for each change, including edge and error paths.
- Keep tests deterministic and fast; avoid external dependencies by using mocks/fakes where applicable.
- Maintain or improve coverage where coverage tooling exists.

## Documentation
- Update README/CHANGELOG when behavior changes or when commands are added.
- After completing a task from docs/tasks.md, check it off by changing `[ ]` to `[x]` and include brief notes if needed.

## Commits/PRs
- One logical change per commit; reference the corresponding task and plan milestone.
- Include rationale and verification steps in commit messages or PR descriptions.

## CI
- Ensure CI runs formatting, vet/lint, and tests consistently with local commands.
