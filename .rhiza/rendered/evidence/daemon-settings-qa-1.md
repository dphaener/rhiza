# daemon-settings-qa-1

evidence · `01M6CX4VKARNPFT2JHQHWP3MKA` · actor [claude-code-63750355](../actor/claude-code-63750355.md)

## result

PASS

## statement

At 1a7a3961 I built go build -o $TMPDIR/rhiza-x ./cmd/rhiza and used it in a scratch repo with a logging fake rex on PATH and every RHIZA_* var unset. daemon settings listed all 8 keys as default. daemon set with an unknown key, maxRuns 0, ship sometimes, minFree lots, paused maybe or permissionMode 'a b' each exited 2 and wrote no rhiza.* git config key; --dry-run wrote nothing. After daemon set maxRuns 1 and minFree 1G, settings showed both as (git config). With maxRuns 1 in config, daemon --once on two ready change sets started one and skipped the other as waiting for a slot. --max-runs 3 started 3 and status showed maxRuns = 3 (flag). A long-running daemon re-read config every 5s poll: raising maxRuns from 1 to 2 in config started the second run on the next poll. A bad hand-edited value (paused maybe) was logged as 'rhiza: daemon: git config rhiza.daemon.paused is "maybe"' and the daemon kept running. daemon --once with a bad value exited 2.

## Links

### Out

- verifies [daemon-settings](../capability/rhiza/orchestration/daemon/daemon-settings.md)

## History

- `01M6CRZQ0YKB2QV4V5TWY71G4V` The daemon reads its settings and each stage's autonomy from git config every poll, can be paused, and applies an acceptance policy
