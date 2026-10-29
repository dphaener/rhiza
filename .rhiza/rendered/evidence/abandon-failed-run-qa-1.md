# abandon-failed-run-qa-1

evidence · `01M634BMB9T4MQF5V5CDQJHBW9` · actor [claude-code-5b91cdbb](../actor/claude-code-5b91cdbb.md)

## result

PASS

## statement

At fd3f4b3 I built rhiza to a scratch path and, in a scratch repo, ran a one-step workflow that fails: run start exited 1 and rhiza inbox listed the failed run with next 'rhiza run resume --retry'. rhiza run abandon <run> exited 0; run status showed 'Run ... failed' with 'Detail: abandoned by darin-haener', and --json gave state failed, detail 'abandoned by darin-haener'. rhiza inbox then said 'Nothing waits on you' and inbox --json gave items []. run resume and run resume --retry each exited 2 with 'run <id> was abandoned'; a second abandon exited 2 with 'is finished: it failed'. After making the step pass, a succeeded run's abandon exited 2 with 'is finished: it succeeded'.

## Links

### Out

- verifies [abandon-failed-run](../capability/rhiza/workflows/workflow-runs/abandon-failed-run.md)

## History

- `01M62MKZFRPZ2C7QTYAF815XAG` The inbox leaves out failed runs whose change set has since landed, been withdrawn or rejected, and rhiza run abandon also takes a failed run, so every failed item can be cleared; no item's next command is one that can only fail again.
