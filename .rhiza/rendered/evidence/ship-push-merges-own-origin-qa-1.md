# ship-push-merges-own-origin-qa-1

evidence · `01M793HGJXFFA1DYPJ4EFRFFC5` · actor [claude-code-6b7d70a1](../actor/claude-code-6b7d70a1.md)

## result

PASS

## statement

QA claude-code-6b7d70a1 at ccf01f58d, hands-on with a built rhiza binary on PATH (fake gh/rex/claude from run-ship-land-conflict.txtar, bare origin, hermetic HOME). HEAD's binary cannot get a ship-change run past build/check (actor reuse, see the blocker on actor-record-no-duplicate), so these runs used HEAD built with only that reuse loop removed; the code under test here is HEAD's unchanged. In the run worktree, pushed an origin-only commit to origin's cs/ship-it and made a local-only commit; RHIZA_RUN=... rhiza ship open-pr merged it ('Merge origin/cs/ship-it into cs/ship-it') and pushed plainly: origin's graph keeps e8ffb16 Origin-only commit. Then a conflicting origin commit on clash.txt: open-pr refused 'origin's cs/ship-it has commits that conflict with ...: b9e3c70 Origin clash', no MERGE_HEAD, local HEAD unchanged, origin untouched.

## Links

### Out

- verifies [ship-push-merges-own-origin](../capability/rhiza/orchestration/shipping/ship-push-merges-own-origin.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).
