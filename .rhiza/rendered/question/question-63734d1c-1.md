# question-63734d1c-1

question · `01M6GCYES204GENK2QVG9ETB3D`

## session

63734d1c-2fd7-4c75-98ce-01a1a0cb9cca

## statement

Change set 01M6GCP (Fix: plan-run close-out fails on untracked .rhiza/log entries) serves no intent yet. Which landed intent should it serve? Recommendation: m24-on-the-vm, the intent the previous Fix (folder trust) served. This one is also a failure the VM's unattended daemon runs into. / How should close-out stop failing on .rhiza/log entries that sit untracked in the main checkout? Every change set opened in the main checkout (run start --new, change new) leaves its first entries there untracked, and the daemon reads drafts and accepted change sets from that checkout. Recommendation: close-out drops untracked log files that the incoming main has byte for byte, then pulls. This keeps the main checkout's view of in-flight change sets and fixes ship-change close-outs too. / `ship close-out --forget` deletes the ship state before it pulls main. If the pull fails, a retried close-out then fails with no ship state. How should a retry work? Recommendation: delete the ship state only after main is updated, so a retry finds it and finishes.

## Links

### In

- answers [close-out-untracked-log-and-retry](../decision/close-out-untracked-log-and-retry.md)

## History

- `01M6GCPSMFW8G6MXXNW6Y5J0Q6` Fix: Request plan runs fail at close-out: the daemon writes the new change set's first .rhiza/log entries untracked into the main checkout, so after the plan PR merges close-out's git pull --ff-only refuses to overwrite them (run 01M6FSQBBJYK). The plan run's log entries should be written only in its worktree, or close-out should drop untracked log files identical to the incoming ones before pulling; and a retried close-out must not fail because the first visit forgot the ship state
