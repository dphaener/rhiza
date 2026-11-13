# ship-gate-reuse-code-tree-qa-1

evidence · `01M792XMPDTVJTZY0M68NY1907` · actor [claude-code-7ac56481](../actor/claude-code-7ac56481.md)

## result

PASS

## statement

QA with the head build (go build ./cmd/rhiza at 4610130) in a scratch git repo with origin/main, a change set, a CI stability record with TestPond flaking on 2 heads, and rhiza.ship.localGate set to a command that appends to a log: rhiza ship gate ran the gate once and wrote gated, gated_tree and gated_command to the ship state; the same head passed at once ('already passed on 4c0f0cb with the same code outside .rhiza/, so the gate was not run again'); a commit adding a .rhiza/ file and a .rhiza/log entry passed at once with the new head and that reason; a commit changing pkg/a.go ran it again; a dirty worktree ran it and recorded nothing (the next clean gate reused the previous record); a changed gate command ran it again on the same tree; a ship state with gated only (gated_tree and gated_command removed) ran it once, then reused.

## Links

### Out

- verifies [ship-gate-reuse-code-tree](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-gate-reuse-code-tree.md)

## History

- `01M78T0G869AE4NQ9YGNB2DDGE` A passing local gate is reused while the code it checked hasn't changed, and the gate checks only what the change touches. Owner, 2026-11-12, agreed to file after the performance report. Today run 01M787SV ran the full make check 4 times (local-gate, gate-rebase-merge, after conflict-fix, after a rendered-only merge), ~20 min each serialized on the test lock; most reruns followed changes that touched no Go code (rendered model files, .rhiza/log entries). Wanted: (1) record the tree hash (or the hash of the non-.rhiza/ tree) a gate passed on and skip the gate when the next push has the same code tree; (2) the gate runs the change's scoped checks (go test for the packages the branch changes vs origin/main, their reverse dependencies, the affected testscripts, make self-check) instead of the whole suite, leaving the full suite to CI; (3) the project setting stays language-neutral (the gate command and its scoping come from project/workflow config, per the cleanup milestone). Estimated saving ~20 min per run on days the gate is on. Evidence: ~/orchestrator/perf-2026-11-12/report-prospective.md.
