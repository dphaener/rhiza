# question-37d69292-1

question · `01M6H6X95DZDFPAC0FGD6Y4SV2`

## session

37d69292-5976-44cc-847c-01a1a26d0332

## statement

Daemon CPU fix (change set 01M6H6S): what should the daemon skip on a poll where nothing changed? Picks depend on the clock (start grace, 15-minute unjournaled hold, resume backoff) and on process liveness, which change without any stamp moving. / Daemon CPU fix (change set 01M6H6S): how is a poll's cost measured?

## Links

### In

- answers [daemon-poll-caches-fold-not-picks](../decision/daemon-poll-caches-fold-not-picks.md)

## History

- `01M6H6SGKR1N5GBG7XPECKN9Z3` Fix The daemon pins a full CPU core: on the 4-core VM rhiza daemon uses ~100% CPU continuously (120 CPU-minutes in 2 hours) because each poll (daemonPoll = 5s, internal/api/daemon.go) takes 10-15 s of work: it re-reads and re-folds the whole log and recomputes picks every time. A poll should be cheap when nothing changed: skip the fold/pick work when the log, git refs, run records and settings are unchanged since the last poll (cache keyed on their stamps), and measure a poll's cost in a test. Agents' rhiza apply calls on the same VM also take minutes under this load
