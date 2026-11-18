# conduct-no-wait-guide-verified

evidence · `01M7QESWZ0CZDZF9V005JDDMZV` · actor [claude-code-5ea4821b](../actor/claude-code-5ea4821b.md)

## result

PASS

## statement

Verified commit 7f2c598fe. go test -p 2 -trimpath ./internal/guide/... exit 0; TestConductNoWait run with HEAD~1's conduct.md and conduct-roles.md in a detached scratch worktree exit 1 (every no-wait/wake-up/end-turn phrase missing); go test -count=1 -p 2 -trimpath ./cmd/rhiza -run 'TestDocsPrimitives|TestScripts/guide-install' exit 0; rhiza ship generate exit 0 with 0 changes and rhiza render changed only this verifier's own actor record. Diff read: conduct-roles documents [--no-wait], the wake-up message text matching spawnWakeText in internal/api/run_lost.go, and keeps the waiting default for other callers; conduct Must do tells the conductor to spawn with --no-wait, start independent children together, end its turn with no sleep or poll and read rhiza run status --json on each wake-up; the shared 'Never sleep and never poll' rule remains and nothing says to poll run status; docs/primitives.md and docs/bring-your-own-agent.md describe --no-wait beside the blocking default. Only the leaf's listed files plus its .rhiza log/rendered records changed; no other role's guide touched, no TODOs.

## Links

### Out

- verifies [conduct-no-wait-guide](../capability/rhiza/workflows/workflow-runs/conduct-no-wait-guide.md)

## History

- `01M7M2FB5BNRSDSN7HFWQSC2WN` rhiza run spawn can return once the child starts, and the child's rhiza session done (or a blocked question) wakes the conductor with rhiza session message, so a conductor ends its turn instead of blocking on spawn; every child stays a Rhiza session on the run's ledger (idea from rex-teams, passed on by the orchestrator 2026-11-17)
