# daemon-settings-qa-3

evidence · `01M6CZZWKV1TPPXZCN3H19QEF7` · actor [claude-code-91d47930](../actor/claude-code-91d47930.md)

## result

PASS

## statement

At 3c66a166 I built go build -o $TMPDIR/rhiza-qa3 ./cmd/rhiza and used it in a scratch repo with a logging fake rex and every RHIZA_* var unset. daemon settings on a fresh repo listed all 8 keys, including rhiza.autonomy.accept = auto (default), and --json gave value auto, source default. daemon set rhiza.autonomy.accept owner then settings showed owner (git config). daemon set rhiza.autonomy.accept maybe exited 2 and left the value at owner. The cycle 1 finding is fixed: with accept unset, an autonomous plan's clean gate accepted its change set, which matches the listed default.

## Links

### Out

- verifies [daemon-settings](../capability/rhiza/orchestration/daemon/daemon-settings.md)

## History

- `01M6CRZQ0YKB2QV4V5TWY71G4V` The daemon reads its settings and each stage's autonomy from git config every poll, can be paused, and applies an acceptance policy
