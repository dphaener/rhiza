# daemon-idle-poll-no-fold-verified

evidence · `01M6HCBKKXCN9HGNQ7B292QYQA` · actor [claude-code-e6ffa3cf](../actor/claude-code-e6ffa3cf.md)

## result

PASS

## statement

At bb59f3b9 I read the diff and ran the test and benchmark. The diff adds internal/api/daemon_fold.go, a fold cache keyed by the stat fingerprint of .rhiza/log (index.LogStamp, which hashes folder names, mtimes and inodes, and is distrusted when racy), git for-each-ref, and the current change set. daemonPoll1, daemonNotify and daemonEvents now read the plan result, candidates, change-set lookups, overlaps and lifecycle moves from the cache. Run journals, liveness, the disk check, autonomy settings, holds, backoff and picks still run on every poll, and the daemon's own starts invalidate the cache. The planReadCount and planFoldCount counters are in planRead and planGather (around FoldCanonical). `go test ./internal/api -run TestDaemonIdlePollSkipsFold -count=1 -v` PASS. The test builds about 105 change sets (6 in flight) with IndexPath set and an event command, and runs the full poll closure. The second poll makes 0 readLog calls, 0 planRead calls and 0 in-flight folds. With the clock moved past daemonUnjournaledHold and no stamp change, the held change set is picked. After a new change set is created and proposed, the next poll reads the log, runs planRead again and writes one more event line, `proposed`. `go test -bench BenchmarkDaemonIdlePoll -benchtime=200x` reported 7935126 ns/op, and no test asserts a wall-clock threshold. go vet and gofmt are clean, and make self-check is ok. In make check, only the wall-clock budget tests failed, with the VM at load average ~9 on 4 cores: TestReadBudget (what-if medians over 50ms) and modelapp TestGeneratedModel. TestGeneratedModel passed when run alone. TestReadBudget fails the same way at the parent commit ebee728a under the same load, so neither failure comes from this leaf.

## Links

### Out

- verifies [daemon-idle-poll-no-fold](../capability/rhiza/orchestration/daemon/daemon-idle-poll-no-fold.md)

## History

- `01M6H6SGKR1N5GBG7XPECKN9Z3` Fix The daemon pins a full CPU core: on the 4-core VM rhiza daemon uses ~100% CPU continuously (120 CPU-minutes in 2 hours) because each poll (daemonPoll = 5s, internal/api/daemon.go) takes 10-15 s of work: it re-reads and re-folds the whole log and recomputes picks every time. A poll should be cheap when nothing changed: skip the fold/pick work when the log, git refs, run records and settings are unchanged since the last poll (cache keyed on their stamps), and measure a poll's cost in a test. Agents' rhiza apply calls on the same VM also take minutes under this load
