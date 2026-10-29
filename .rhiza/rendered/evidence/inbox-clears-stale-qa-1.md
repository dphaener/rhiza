# inbox-clears-stale-qa-1

evidence · `01M634BMVYC4X3DFQN7W4XF5AZ` · actor [claude-code-5b91cdbb](../actor/claude-code-5b91cdbb.md)

## result

PASS

## statement

At fd3f4b3, with the built binary in a scratch repo, I failed a real run on each of four change sets, so rhiza inbox listed four failed runs. I then landed one change set, withdrew one and rejected one, and left the fourth as a draft. rhiza inbox and inbox --json then listed only the failed run on the draft change set. rhiza inbox open on the withdrawn change set's item exited 2 with 'no inbox item'. I did not reproduce the gave-up and superseded cases by hand. They and the question/sign-off cases pass in inbox-stale.txtar, which I ran with run-abandon-failed and the other inbox and run-abandon scripts (go test -count=1, ok).

## Links

### Out

- verifies [inbox-clears-stale](../capability/rhiza/orchestration/inbox/inbox-clears-stale.md)

## History

- `01M62MKZFRPZ2C7QTYAF815XAG` The inbox leaves out failed runs whose change set has since landed, been withdrawn or rejected, and rhiza run abandon also takes a failed run, so every failed item can be cleared; no item's next command is one that can only fail again.
