# daemon-poll-cache-invalidates-qa-1

evidence · `01M6HF8ZEQRQ23MNFQMDAMG6PZ` · actor [claude-code-797a9572](../actor/claude-code-797a9572.md)

## result

PASS

## statement

At dd9192ca, with the branch binary's real `rhiza daemon` running paused on a scratch clone of this repo, I sampled its CPU every 5 s while changing things from another process. `git branch qa-scratch-ref` (a ref move) cost 361 ticks in the next window, which is one re-fold, and then it went back to idle. `rhiza daemon set rhiza.daemon.maxRuns 3` (a setting) cost 346 ticks, then idle. `rhiza change new` (a log entry) caused re-folds over the next polls, then idle. `rhiza propose` caused re-folds, and the event command ran exactly once, writing one line {"event":"proposed","change":"01M6HF4SM0QW..."}. No moves already in the log were replayed. Before each change the daemon sat at idle CPU, so the cached fold was in use. I did not exercise daemon start/resume invalidation hands-on, because that would launch real agents. The daemon logs had no errors.

## Links

### Out

- verifies [daemon-poll-cache-invalidates](../capability/rhiza/orchestration/daemon/daemon-poll-cache-invalidates.md)

## History

- `01M6H6SGKR1N5GBG7XPECKN9Z3` Fix The daemon pins a full CPU core: on the 4-core VM rhiza daemon uses ~100% CPU continuously (120 CPU-minutes in 2 hours) because each poll (daemonPoll = 5s, internal/api/daemon.go) takes 10-15 s of work: it re-reads and re-folds the whole log and recomputes picks every time. A poll should be cheap when nothing changed: skip the fold/pick work when the log, git refs, run records and settings are unchanged since the last poll (cache keyed on their stamps), and measure a poll's cost in a test. Agents' rhiza apply calls on the same VM also take minutes under this load
