# daemon-disk-guard-verified

evidence · `01M6B1MCY9EHQ4T75JFQXG1XSX` · actor [claude-code-ef797ba0](../actor/claude-code-ef797ba0.md)

## result

PASS

## statement

At bed13ffa ran go test ./cmd/rhiza -run 'TestScripts/(daemon-picks|describe-golden)$' -count=1: ok. daemon-picks.txtar with RHIZA_TEST_FREE_BYTES=1 GiB runs daemon --once: no 'Started', no rex.log, daemon status --json shows reason disk-low with '1.0 GiB free, under the minimum of 10.0 GiB'; at 20 GiB the same fixture starts one run; at 1 GiB with --min-free 0 it starts one run; --min-free lots is a usage error. Regenerating describe-golden with -update left the tree clean (min_free and disk-low present). Read the diff: daemonPoll1 skips both Resumable and Pick when low, marks every candidate disk-low with the free space; default 10 GiB; FreeBytes via statfs Bavail*Bsize. make check exit 0.

## Links

### Out

- verifies [daemon-disk-guard](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-disk-guard.md)

## History

- `01M6ATAS02XDRJ1MSQ4THGCNPA` The daemon holds new starts and resumes while disk space is low, and daemon status and the inbox say so
