# Improvement Plan

## Baseline
- Test commands: to be filled after first run (e.g., `go test ./... -count=1`)
- Current status: pass/fail summary and durations (to be recorded)

## Milestone 1 (short-term)
- Stabilize tests; fix failing or flaky cases uncovered by baseline.
- Enforce formatting (go fmt) and lint rules if available (golangci-lint or `go vet`).
- Add missing tests for critical paths and error cases.

## Milestone 2 (medium-term)
- Refactor high-complexity areas; improve developer workflow (Makefile targets for test/lint/format).

## Milestone 3 (long-term)
- Performance/resiliency improvements and documentation hardening.

## Priorities
Security > Correctness > Reliability > Developer Experience > Performance > Aesthetics

## Constraints
- Small, incremental changes with tests and docs updates.
