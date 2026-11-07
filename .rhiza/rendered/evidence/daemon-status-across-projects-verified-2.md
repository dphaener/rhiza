# daemon-status-across-projects-verified-2

evidence · `01M6SXPTEFYA3WANGN831XZEKK` · actor [claude-code-4e4a3ac7](../actor/claude-code-4e4a3ac7.md)

## result

PASS

## statement

Verifier at 7f571ac0 (attempt 2, fixing review finding f1). Read the diff: daemon-status-projects.txtar now gives repo and second each an accepted no-milestone change set (Orphan) and a fixture live-run journal (dead pid, change set the daemon does not resume). Ran go test -trimpath ./cmd/rhiza -run TestScripts/daemon-status -count=1: ok. The script asserts at the root, in --json, each project's part with its runs (LIVEREPO001 / LIVESECOND1 interrupted), skipped (Orphan no-milestone) and last_started (count=2), the bare project's 'not initialised' once, and the same per project in text. Re-ran the built binary by hand in the kept $WORK: root --json shows the daemon fields once, projects[] with repo and second each carrying runs, skipped, last_started, and skipped_projects bare: not initialised; --project repo --json and from inside second give projects: [] with only that project's runs/skipped/last_started at top level (LIVEREPO001 only, LIVESECOND1 only respectively); --project nope fails 'not registered'.

## Links

### Out

- verifies [daemon-status-across-projects](../capability/rhiza/orchestration/daemon-status-across-projects.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
