# daemon-settings-verified

evidence · `01M6CWWZPKX64NP0A5AP34HYG2` · actor [claude-code-f549fe0b](../actor/claude-code-f549fe0b.md)

## result

PASS

## statement

At 965d6d13 (attempt 2): go test ./cmd/rhiza -run 'TestScripts/(daemon-settings|describe-golden|daemon-pause)$' -count=1 passed. daemon-settings.txtar checks every key with its default and source default; daemon set rhiza.daemon.maxRuns 1 then settings shows 1 (git config); bad key/values fail and write nothing; with maxRuns 1 a daemon --once poll of two ready change sets starts exactly one run; --max-runs 5 wins and daemon status shows source flag, and the next poll without the flag reads config again. describe-golden passes against the regenerated golden. By hand on a scratch repo, daemon set with an unknown key, maxRuns zero, ship sometimes, minFree lots and paused maybe each exited 2 and left no rhiza.daemon/rhiza.autonomy keys in git config. Attempt 1's finding (daemon-lock flaky) is fixed: the fix installs the interrupt handler before daemon.TryLock; TestScripts/daemon-lock passed 20 of 20 runs with -count=1. make check exited 0.

## Links

### Out

- verifies [daemon-settings](../capability/rhiza/orchestration/daemon/daemon-settings.md)

## History

- `01M6CRZQ0YKB2QV4V5TWY71G4V` The daemon reads its settings and each stage's autonomy from git config every poll, can be paused, and applies an acceptance policy
