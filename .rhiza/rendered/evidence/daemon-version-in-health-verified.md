# daemon-version-in-health-verified

evidence · `01M6FKWGQCZPE76N0PD7Y0QSXW` · actor [claude-code-2403cd29](../actor/claude-code-2403cd29.md)

## result

PASS

## statement

At 933178f7: go test ./cmd/rhiza -run 'TestScripts/(daemon-version-in-health|describe-golden)$' -count=1 passes; make fmt-check vet test exits 0. The scenario starts a daemon with RHIZA_TEST_COMMIT set to the built commit; daemon status prints 'Commit: <40-hex>' and --json carries "commit". With origin/main at that commit, health's daemon part is ok; after an extra commit moved origin/main ahead, it is warn: '...; built from <built7>, origin/main is at <head7>'. I re-ran a temporary copy of the scenario that printed health --json (deleted afterwards): the daemon part was state "ok" when the commits matched, and state "warn", summary 'built from a10ad21, origin/main is at 4919c15', next "make update" when origin/main was ahead. Diff read: Commit comes from the build info's vcs.revision (RHIZA_TEST_COMMIT only under RHIZA_TEST_NOW), is stored in daemon.State and DaemonStatusResponse, and healthOriginAhead warns only when origin/main is a strict descendant (merge-base --is-ancestor), with no fetch. make self-check's only failures were rendered actor files for actors registered after this commit (including this verifier), so not caused by the leaf.

## Links

### Out

- verifies [daemon-version-in-health](../capability/rhiza/orchestration/daemon/daemon-version-in-health.md)

## History

- `01M6FEN1BRQRHP9RTJKN2YPTH6` The daemon pulls main and runs make update when main moves, and health shows the version it runs
