# ship-gate-reuse-code-tree-verified

evidence · `01M7926XSX87DFAHF6ZTVV0GCC` · actor [claude-code-859d14f4](../actor/claude-code-859d14f4.md)

## result

PASS

## statement

Verifier at 669309152: go test -trimpath -count=1 ./internal/api -run TestShipLocalGate passes (TestShipLocalGate and TestShipLocalGateOncePerCodeTree), and -run TestShip passes. TestShipLocalGateOncePerCodeTree, with CI unstable and a gate appending to a log to count runs, shows: first head runs the gate once; a new head whose git ls-tree differs only in the .rhiza entry passes with no run and a reason naming beef000 and 'not run again'; a changed non-.rhiza blob runs it again; a changed rhiza.ship.localGate on the same tree runs it again; a dirty worktree runs it and records nothing (the next clean gate reuses the earlier record); a state of only {"gated":...} reads and reruns once, then reuses. Diff read: shipState gains gated_tree and gated_command; shipCodeTree sha256-hashes git ls-tree HEAD lines minus the top-level .rhiza entry; reuse requires a clean worktree, same tree and same trimmed gate command. Checked by hand that real git ls-tree HEAD lists .rhiza as a tab-separated top-level entry. go vet and gofmt clean on internal/api.

## Links

### Out

- verifies [ship-gate-reuse-code-tree](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-gate-reuse-code-tree.md)

## History

- `01M78T0G869AE4NQ9YGNB2DDGE` A passing local gate is reused while the code it checked hasn't changed, and the gate checks only what the change touches. Owner, 2026-11-12, agreed to file after the performance report. Today run 01M787SV ran the full make check 4 times (local-gate, gate-rebase-merge, after conflict-fix, after a rendered-only merge), ~20 min each serialized on the test lock; most reruns followed changes that touched no Go code (rendered model files, .rhiza/log entries). Wanted: (1) record the tree hash (or the hash of the non-.rhiza/ tree) a gate passed on and skip the gate when the next push has the same code tree; (2) the gate runs the change's scoped checks (go test for the packages the branch changes vs origin/main, their reverse dependencies, the affected testscripts, make self-check) instead of the whole suite, leaving the full suite to CI; (3) the project setting stays language-neutral (the gate command and its scoping come from project/workflow config, per the cleanup milestone). Estimated saving ~20 min per run on days the gate is on. Evidence: ~/orchestrator/perf-2026-11-12/report-prospective.md.
