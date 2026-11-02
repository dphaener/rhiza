# stage-autonomy-qa-3

evidence · `01M6D0024DM260H3V40JBETQ6X` · actor [claude-code-91d47930](../actor/claude-code-91d47930.md)

## result

PASS

## statement

At 3c66a166, with the built binary and a logging fake rex, after the autonomyStage refactor: with rhiza.autonomy.ship owner, daemon --once started nothing and never called rex, and daemon status said 'Built: ship is set to owner: only the owner starts its run'. With ship auto and signoff owner, daemon --autonomous --once started one run, and rex was called with run start <cs> --permission-mode manual, without --autonomous. accept owner/auto behaved as in accept-policy-qa-3.

## Links

### Out

- verifies [stage-autonomy](../capability/rhiza/orchestration/accept-policy/stage-autonomy.md)

## History

- `01M6CRZQ0YKB2QV4V5TWY71G4V` The daemon reads its settings and each stage's autonomy from git config every poll, can be paused, and applies an acceptance policy
