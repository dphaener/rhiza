# daemon-pause-qa-1

evidence · `01M6CX4XGZNYRT2ZVD9E7RH3Z8` · actor [claude-code-63750355](../actor/claude-code-63750355.md)

## result

PASS

## statement

At 1a7a3961 I used the built binary in a scratch repo with a fake rex. With rhiza.daemon.paused true, daemon --once on two accepted ready change sets started nothing and never called rex. daemon status printed 'Paused: yes' and skipped both as paused. After daemon set paused false, the next --once started a run. Live: a looping daemon started paused polled with both candidates skipped as paused and no rex call. After daemon set paused false it started a run within one poll, and status showed 'Paused: no'. SIGINT stopped it cleanly, exit 0.

## Links

### Out

- verifies [daemon-pause](../capability/rhiza/orchestration/daemon/daemon-pause.md)

## History

- `01M6CRZQ0YKB2QV4V5TWY71G4V` The daemon reads its settings and each stage's autonomy from git config every poll, can be paused, and applies an acceptance policy
