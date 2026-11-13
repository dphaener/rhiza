# ship-gate-reuse-code-tree

capability · `01M78TAJ9VCKZHPQZN5WS0A6DF` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/api -run 'TestShipLocalGate' shows, with CI unstable and a gate that counts its runs: the gate runs once on a head; a new commit changing only .rhiza/ files (a rendered file and a log entry) passes at once with a reason naming the earlier passing head and the gate not run again; a commit changing a file outside .rhiza/ runs it again; a changed rhiza.ship.localGate runs it again on the same tree; a dirty worktree runs it and records nothing; a ship state written before this change (Gated head only, no tree) still reads and reruns the gate once

## statement

A passing local gate records, in the ship state beside the head it passed on, the gate command and the hash of that head's tree with the top-level .rhiza/ entry left out, and rhiza ship gate passes at once without running the gate when the worktree is clean and its HEAD has the same command and the same non-.rhiza/ tree, so render, log-entry and model-only merge commits never rerun it while any change outside .rhiza/ or a changed gate command does

## Links

### Out

- depends-on [ship-local-gate](ship-local-gate.md)
- decided-by [local-gate-reuse-and-scope](../../../../../decision/local-gate-reuse-and-scope.md)

### In

- verifies [ship-gate-reuse-code-tree-verified](../../../../../evidence/ship-gate-reuse-code-tree-verified.md)
- verifies [ship-gate-reuse-code-tree-qa-1](../../../../../evidence/ship-gate-reuse-code-tree-qa-1.md)
- verifies [pr-513-a-passing-local-gate-is-reused](../../../../../evidence/pr-513-a-passing-local-gate-is-reused.md)

## History

- `01M78T0G869AE4NQ9YGNB2DDGE` A passing local gate is reused while the code it checked hasn't changed, and the gate checks only what the change touches. Owner, 2026-11-12, agreed to file after the performance report. Today run 01M787SV ran the full make check 4 times (local-gate, gate-rebase-merge, after conflict-fix, after a rendered-only merge), ~20 min each serialized on the test lock; most reruns followed changes that touched no Go code (rendered model files, .rhiza/log entries). Wanted: (1) record the tree hash (or the hash of the non-.rhiza/ tree) a gate passed on and skip the gate when the next push has the same code tree; (2) the gate runs the change's scoped checks (go test for the packages the branch changes vs origin/main, their reverse dependencies, the affected testscripts, make self-check) instead of the whole suite, leaving the full suite to CI; (3) the project setting stays language-neutral (the gate command and its scoping come from project/workflow config, per the cleanup milestone). Estimated saving ~20 min per run on days the gate is on. Evidence: ~/orchestrator/perf-2026-11-12/report-prospective.md.
