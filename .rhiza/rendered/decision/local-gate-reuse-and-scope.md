# local-gate-reuse-and-scope

decision · `01M78T8Q83CG5M2526J2YGVHN7`

## reasoning

Milestone 35's statement names a 40-minute ship-change run, and m35-ship-40m is its measurable done point; the perf report (~/orchestrator/perf-2026-11-12/report-prospective.md, item 1) puts gate and gate-fix at ~25m per run and says no m35 change set covers it, so the change set serves that done point rather than m4-runs-hold-up (about hand repair, which the sibling's trigger change fits). internal/api/ship_gate.go already skips the gate when HEAD equals the recorded Gated head on a clean tree; the intent's 4 reruns in run 01M787SV followed rendered-model and .rhiza/log changes, which change HEAD but not code, so keying the reuse on the non-.rhiza/ tree hash (option 1 of the intent) removes them while any code change still reruns it. Leaving .rhiza/ out is safe because the pre-push check runs rhiza check and make self-check before every push regardless. The intent's point (3) and m32-go-in-own-workflows require Go behaviour to live in Rhiza's own config, so the product only hands the command its base and the per-language scoping lives in this repo's gate command, with a full-suite fallback so a scoping gap never lets an untested change through.

## statement

Change set 01M78T0 serves m35-ship-40m (milestone 35: a ship-change run takes 40 minutes or less). Reuse: ship gate records, beside the head it passed on (ship state Gated), the hash of that head's tree with the top-level .rhiza/ entry left out, and passes at once without running the gate when the branch's current HEAD, with a clean worktree, has the same non-.rhiza/ tree hash; any change outside .rhiza/ (a conflict fix, a merge of origin/main that brings code) reruns it. .rhiza/ changes are already checked on every push by the pre-push check (rhiza check, then rhiza.ship.prePushCheck). Scope: Rhiza stays language-neutral. The gate command remains the project's setting rhiza.ship.localGate; ship gate runs it with the base it compares against in the environment (RHIZA_GATE_BASE, the merge-base with origin/main, and RHIZA_GATE_HEAD), and the project's command does the scoping. This repo sets its gate to a script/make target of its own that tests the Go packages the branch changes since that base, their reverse dependencies, the affected testscripts and make self-check, falling back to the full suite when it cannot work out the scope (or the base is missing); the full suite stays with CI. No Go-specific scoping goes into the shipped product.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-aea6d8fc-1](../question/question-aea6d8fc-1.md)
- derived-from [local-gate-reuse-interview](../source/local-gate-reuse-interview.md)

### In

- decided-by [ship-gate-reuse-code-tree](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-gate-reuse-code-tree.md)
- decided-by [ship-gate-base-env](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-gate-base-env.md)
- decided-by [rhiza-scoped-gate](../capability/rhiza/integrity/self-check/rhiza-scoped-gate.md)

## History

- `01M78T0G869AE4NQ9YGNB2DDGE` A passing local gate is reused while the code it checked hasn't changed, and the gate checks only what the change touches. Owner, 2026-11-12, agreed to file after the performance report. Today run 01M787SV ran the full make check 4 times (local-gate, gate-rebase-merge, after conflict-fix, after a rendered-only merge), ~20 min each serialized on the test lock; most reruns followed changes that touched no Go code (rendered model files, .rhiza/log entries). Wanted: (1) record the tree hash (or the hash of the non-.rhiza/ tree) a gate passed on and skip the gate when the next push has the same code tree; (2) the gate runs the change's scoped checks (go test for the packages the branch changes vs origin/main, their reverse dependencies, the affected testscripts, make self-check) instead of the whole suite, leaving the full suite to CI; (3) the project setting stays language-neutral (the gate command and its scoping come from project/workflow config, per the cleanup milestone). Estimated saving ~20 min per run on days the gate is on. Evidence: ~/orchestrator/perf-2026-11-12/report-prospective.md.
