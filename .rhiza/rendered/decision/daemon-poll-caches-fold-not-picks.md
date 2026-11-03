# daemon-poll-caches-fold-not-picks

decision · `01M6H6ZC969219DGD2H3FHVX3V`

## reasoning

Read daemonPoll1 in internal/api/daemon.go: the 10-15 s cost is planRead and openCtxWorkspace, which re-read and fold the whole log; everything that depends on the clock or process liveness (Resumable with rs.Live and ctx.Now, daemonHeld with grace and the 15-minute hold, Pick/Skips over run journals) is cheap computation over run journals plus the folded workspace. Caching only the fold-derived inputs keyed on the stamps the change set's intent names keeps every time- and liveness-driven pick correct without inventing time-bucketed invalidation, while removing the expensive part. Run journals are small and are where liveness and backoff live, so they stay read every poll rather than becoming a cache key. For measurement, wall-clock thresholds are flaky on the shared VM and CI runners (test timings swing with load), so the test asserts the work skipped (fold count) rather than elapsed time, which is deterministic and directly proves the intent's 'skip the fold/pick work when unchanged'.

## statement

On a poll where the log, git refs and settings stamps are unchanged, the daemon skips only the fold-derived work: planRead (open reader, planGather, planCore), openCtxWorkspace (in daemonPoll1 and in the events pass daemonEvents), daemonCandidates, the change-set lookups (accepted state, ship Since, titles, has-leaves) and the lifecycle moves daemonEvents collects are reused from the last poll. The clock- and liveness-dependent steps run every poll against that cached result: rs.List of run journals, daemon.Resumable, daemonHeld (start grace, 15-minute unjournaled hold), resume backoff, disk check, autonomy stages, Pick/Skips, queued continues/starts, plan requests, the notify pass, and the event command for moves state does not list as handled. Any write the daemon itself makes (start, resume journal) invalidates the cache. A poll's cost is measured by a deterministic Go test, not wall-clock: build a repo with a large log and an event command, run the whole poll closure (daemonPoll1, daemonNotify, daemonEvents) twice with nothing changed, and assert through a fold counter hook that the second poll performs zero folds/plan reads, and that changing a stamp (new log entry, ref move, settings edit) makes the next poll fold again; an optional go benchmark reports ns/op for the unchanged poll but no test asserts a timing threshold.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-37d69292-1](../question/question-37d69292-1.md)
- derived-from [daemon-cpu-plan-autonomous](../source/daemon-cpu-plan-autonomous.md)

### In

- decided-by [daemon-idle-poll-no-fold](../capability/rhiza/orchestration/daemon/daemon-idle-poll-no-fold.md)
- decided-by [daemon-poll-cache-invalidates](../capability/rhiza/orchestration/daemon/daemon-poll-cache-invalidates.md)

## History

- `01M6H6SGKR1N5GBG7XPECKN9Z3` Fix The daemon pins a full CPU core: on the 4-core VM rhiza daemon uses ~100% CPU continuously (120 CPU-minutes in 2 hours) because each poll (daemonPoll = 5s, internal/api/daemon.go) takes 10-15 s of work: it re-reads and re-folds the whole log and recomputes picks every time. A poll should be cheap when nothing changed: skip the fold/pick work when the log, git refs, run records and settings are unchanged since the last poll (cache keyed on their stamps), and measure a poll's cost in a test. Agents' rhiza apply calls on the same VM also take minutes under this load
