# stage-autonomy-qa-1

evidence · `01M6CX4ZCVR9PKDM65P9HH1C9C` · actor [claude-code-63750355](../actor/claude-code-63750355.md)

## result

PASS

## statement

At 1a7a3961 I used the built binary in a scratch repo with a fake rex. With ship owner, daemon --once started nothing, made no rex call, and status said 'ship is set to owner: only the owner starts its run'. With signoff owner, daemon --autonomous --once sent rex 'run start <cs> --permission-mode manual' without --autonomous; with signoff auto the line ended in --autonomous. On a real autonomous plan run (hold workflow) with clean findings, accept owner made run gate leave the drafted change set proposed with next 'rhiza accept --rule owner'. inbox showed one accept item, 'waiting for acceptance'. Setting accept auto and gating again accepted it, and the accept item was gone.

## Links

### Out

- verifies [stage-autonomy](../capability/rhiza/orchestration/accept-policy/stage-autonomy.md)

## History

- `01M6CRZQ0YKB2QV4V5TWY71G4V` The daemon reads its settings and each stage's autonomy from git config every poll, can be paused, and applies an acceptance policy
