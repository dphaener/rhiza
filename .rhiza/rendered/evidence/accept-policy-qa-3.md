# accept-policy-qa-3

evidence · `01M6CZZZDG4S8VGMYJWCKX6BG7` · actor [claude-code-91d47930](../actor/claude-code-91d47930.md)

## result

PASS

## statement

At 3c66a166, by hand with the built binary on an autonomous plan run (hold workflow, --no-wait) whose review round was clean: with rhiza.autonomy.accept unset, run gate accepted Serve B one (by the gate: proposed, accepted). After daemon set rhiza.autonomy.accept owner, the same gate left Serve B two proposed, and inbox listed it as waiting for acceptance with next: rhiza accept --rule owner.

## Links

### Out

- verifies [accept-policy](../capability/rhiza/orchestration/accept-policy.md)

## History

- `01M6CRZQ0YKB2QV4V5TWY71G4V` The daemon reads its settings and each stage's autonomy from git config every poll, can be paused, and applies an acceptance policy
