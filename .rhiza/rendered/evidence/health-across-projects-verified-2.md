# health-across-projects-verified-2

evidence · `01M6SYFBTQCAMMZPXVTNMSP6W2` · actor [claude-code-ea498c12](../actor/claude-code-ea498c12.md)

## result

PASS

## statement

Verifier claude-code-ea498c12 at 52fcbbae: go test -trimpath ./cmd/rhiza -run 'TestScripts/(health|daemon-status)' passes. health-projects.txtar sets up a root with two projects, alpha and beta. At the root, rhiza health --json gives daemon, disk, usage and machine once each, inbox and check twice, and 5 project-tagged items per project. A junk log in beta makes rhiza health exit 1 with 'bad check beta:' and no bad alpha. --project alpha and running inside beta each give one project's parts, untagged. --project nope reports not registered. The attempt-1 finding is fixed: with the root daemon locked, home beta, and runs in both projects' state, health reports '3 live runs' and beta's usage window. With the commit's test applied to the parent 1cf77bb7, the script fails with '1 live runs', so the test catches the bug. Read the diff of health.go, daemon.go and main_test.go. internal/api: in this checkout under ~/code, 3 usage tests (TestUsageWindowsHitKinds, TestLimitHitsBecomeLearnedLimits, TestSessionRunRecordsLimitHit) fail. They pass at this same commit and at its parent in scratch worktrees outside the root, so the failure comes from the environment, not this change. go vet and gofmt are clean.

## Links

### Out

- verifies [health-across-projects](../capability/rhiza/orchestration/health-across-projects.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
