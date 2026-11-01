# daemon-disk-guard-verified-2

evidence · `01M6B4XQ0N72JK8TZ4NDP4N70B` · actor [claude-code-6b128d4d](../actor/claude-code-6b128d4d.md)

## result

PASS

## statement

Verifier, attempt 2, at 50de518c (fix for review finding f1 on top of bed13ffa). go test ./cmd/rhiza -run 'TestScripts/(daemon-picks|describe-golden)$' -count=1 -v: both PASS. daemon-picks covers RHIZA_TEST_FREE_BYTES=1 GiB -> daemon --once starts nothing, opens no rex tab, daemon status --json shows reason disk-low with '1.0 GiB free, under the minimum of 10.0 GiB'; 20 GiB -> one run started; --min-free 0 at 1 GiB -> run started; --min-free lots rejected. Read the code: daemonPoll1 skips daemon.Resumable when low, so resumes are held too. describe-golden includes --min-free. Fix diff: daemonDiskLow now reads RHIZA_TEST_FREE_BYTES through ctxGetenv and only when RHIZA_TEST_NOW is set, matching testRenderStale and Decision disk-guard-10-gib. Built the binary and ran daemon --once in a scratch repo with RHIZA_TEST_FREE_BYTES=5: without RHIZA_TEST_NOW, no disk_low in daemon.state.json and no blocked status. With RHIZA_TEST_NOW set, it emits the blocked 'disk space is low' status. make fmt-check, vet and go test ./... pass; one make test run tripped the read-timing benchmark under load and passed on three reruns. make self-check fails only on stale .rhiza/rendered actor files from the review run's chore commits, which this leaf's commits do not touch.

## Links

### Out

- verifies [daemon-disk-guard](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-disk-guard.md)

## History

- `01M6ATAS02XDRJ1MSQ4THGCNPA` The daemon holds new starts and resumes while disk space is low, and daemon status and the inbox say so
