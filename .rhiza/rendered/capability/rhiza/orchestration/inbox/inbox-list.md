# inbox-list

capability · `01M61R7Q0JX5BXJP6M3P1SJXTG` · parent [inbox](../inbox.md)

## acceptance

A testscript with fixture journals of a run blocked on a question, a run waiting at sign-off, a failed run, a run interrupted again after a daemon resume and a done run lists exactly the first four in order in text and --json; once the question is answered, the sign-off approved, the given-up run abandoned and a new fake run started and finished on the failed run's change set, rhiza inbox says nothing waits and lists no failed run for the abandoned one.

## statement

rhiza inbox lists, oldest first across every run of the repo, each item waiting on the owner with a short item ID: a question (its run, change set, session tab and the question), a run waiting at sign-off, a failed run (its failed step and detail, until it is retried with rhiza run resume --retry or a newer run starts on its change set), and a run the daemon gave up on (interrupted again after the daemon's one resume, its step, until it is resumed or abandoned), in text and --json, and says so when nothing waits. A run failed by rhiza run abandon is the owner's own act and is never listed as a failed run.

## Links

### Out

- depends-on [daemon-resume-once](../daemon/daemon-resume-once.md)
- decided-by [answer-8cd68bc7-4](../../../../decision/answer-8cd68bc7-4.md)
- decided-by [answer-8cd68bc7-16](../../../../decision/answer-8cd68bc7-16.md)
- decided-by [answer-8cd68bc7-15](../../../../decision/answer-8cd68bc7-15.md)

### In

- depends-on [inbox-open](inbox-open.md)
- depends-on [inbox-signoff](inbox-signoff.md)
- depends-on [inbox-advisor](inbox-advisor.md)
- depends-on [inbox-notify](inbox-notify.md)
- verifies [inbox-list-verified](../../../../evidence/inbox-list-verified.md)
- verifies [inbox-list-qa-1](../../../../evidence/inbox-list-qa-1.md)
- verifies [pr-53-the-inbox](../../../../evidence/pr-53-the-inbox.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
