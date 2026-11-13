# land-conflict-to-fixer-qa-1

evidence · `01M793HGF9TTJNR9W5D6K0FJ7P` · actor [claude-code-6b7d70a1](../actor/claude-code-6b7d70a1.md)

## result

PASS

## statement

QA claude-code-6b7d70a1 at ccf01f58d, hands-on with a built rhiza binary on PATH (fake gh/rex/claude from run-ship-land-conflict.txtar, bare origin, hermetic HOME). HEAD's binary cannot get a ship-change run past build/check (actor reuse, see the blocker on actor-record-no-duplicate), so these runs used HEAD built with only that reuse loop removed; the code under test here is HEAD's unchanged. Case 1: head build refuses the first land with exit 3 'cannot land change set X: node N changed' after main moved feature.txt; rhiza run status shows land fail -> land-triage -> land-conflict-triage -> land-main-conflict-triage pass -> conflict-fix (1 fixer session) -> land visit 2 pass, never conflict-fix-landed; two head-build land calls, gh pr merge ran, change set landed. Case 2 (exit 3, other message): run fails at land-main-conflict-triage, no fixer session.

## Links

### Out

- verifies [land-conflict-to-fixer](../capability/rhiza/orchestration/shipping/land-conflict-to-fixer.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).
