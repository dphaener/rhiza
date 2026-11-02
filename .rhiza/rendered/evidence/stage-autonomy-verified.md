# stage-autonomy-verified

evidence · `01M6CTTPKYWHZSQ99STGMYV2XZ` · actor [claude-code-d6305650](../actor/claude-code-d6305650.md)

## result

PASS

## statement

At fc0675ed I ran go test ./cmd/rhiza -run 'TestScripts/(stage-autonomy|describe-golden)$' -count=1 (ok) and go test ./... (all packages ok), and rhiza check passed on a clean clone of the commit. stage-autonomy.txtar uses a fake rex and covers each point. With rhiza.autonomy.ship=owner, daemon --once starts nothing, writes no rex.log, and daemon status --json shows reason ship-owner. With signoff=owner, the rex log has exactly one 'run start <cs> --permission-mode manual' line with no --autonomous, and with signoff=auto that line ends in --autonomous. With accept=owner, run gate on an autonomous plan run leaves the draft proposed, and inbox --json shows one accept item naming it ('waiting for acceptance'). With accept=auto, the same gate accepts it and the accept item is gone. A bad value is refused with an error naming the key. In the diff, autonomyStage re-reads git config on each use. The ship-owner branch skips every candidate. Signoff owner forces autonomous=false. The accept-owner gate commits without accepting. The inbox accept-item IDs are keyed per change set.

## Links

### Out

- verifies [stage-autonomy](../capability/rhiza/orchestration/accept-policy/stage-autonomy.md)

## History

- `01M6CRZQ0YKB2QV4V5TWY71G4V` The daemon reads its settings and each stage's autonomy from git config every poll, can be paused, and applies an acceptance policy
