# daemon-fold-from-projection-verified

evidence · `01M7H7SCD6MTTGHZVZXYK9AQZZ` · actor [claude-code-82f9bd45](../actor/claude-code-82f9bd45.md)

## result

PASS

## statement

Verified commit d802ae4a6 (only internal/api/daemon_fold.go and its test changed outside .rhiza/). Re-ran under the capped scope: go test -trimpath -count=1 -p 2 ./internal/api -run 'Daemon|TestPropertyProjectionParity' ok; -run 'Daemon' ok (TestDaemonFoldIndexedMatchesLog, TestDaemonIdlePollSkipsFold, TestDaemonPollCacheInvalidates PASS); go test -trimpath -count=1 -p 2 ./cmd/rhiza -run TestScripts/daemon ok (34 daemon*.txtar, RHIZA_/SHIP_/REX_ unset); go vet ./internal/api exit 0; rhiza check exit 0 (only pre-existing unclaimed warnings). Diff: with ctx.IndexPath set daemonFold takes change sets from p.changeSets(ws.pending) and models from p.model; with it empty it still calls changeset.FoldLog(ws.log()) once and newModelBase(list).model (newModelBase kept). The derivation is moved unchanged into daemonDerive; schedule, journal and --json types untouched. Cache tests: daemonFoldCache.folds increments exactly on a cache miss in daemonFoldOf, so idle polls asserting unchanged folds plus zero log reads and zero planRead still prove reuse, and each invalidation (log entry, accept, setting via pickFolds before the pick, run start, ref move, resume) asserts a new fold; no coverage lost and the zero-log-read assertion is added. A projection failure is returned after the derive and before caching, the same path as a planRead/workspace error today: daemonPoll1 returns the error, notify/demand/health degrade as on any fold error, nothing is cached so the next poll retries. TestDaemonFoldIndexedMatchesLog covers proposed with/without leaves, an accepted one and two in-flight change sets setting done-0 (overlap), zero log reads, candidates/proposed/planSkips/overlaps/moves/states equal to daemonDerive over FoldLog, and the vanished propose entry case (stamp change, rebuild) still matching and dropping it from proposed.

## Links

### Out

- verifies [daemon-fold-from-projection](../capability/rhiza/views/fast-reads/daemon-fold-from-projection.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
