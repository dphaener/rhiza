# daemon-disk-guard-qa-2

evidence · `01M6B52XAW90B2HK3EZEVKGNY8` · actor [claude-code-2be1dffb](../actor/claude-code-2be1dffb.md)

## result

PASS

## statement

At 5218260b built rhiza into $TMPDIR and drove it in a scratch repo (one accepted milestone change set with a leaf, Rex faked to log only, RHIZA_* run vars unset). RHIZA_TEST_FREE_BYTES=1 GiB without RHIZA_TEST_NOW is now ignored: with the real statfs (55 GiB free, default 10 GiB minimum) daemon --once started the run and logged rex window new. With RHIZA_TEST_NOW also set, the same 1 GiB is honoured: nothing started, daemon status shows 'Skipped: <cs> Built: disk space is low (1.0 GiB free, under the minimum of 10.0 GiB)'. --min-free 100G against the real statfs skips with '54.3 GiB free, under the minimum of 100.0 GiB'; --min-free 1G starts the run. go test of daemon-picks, daemon-notify-disk-low, inbox-disk-low and describe-golden: ok.

## Links

### Out

- verifies [daemon-disk-guard](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-disk-guard.md)

## History

- `01M6ATAS02XDRJ1MSQ4THGCNPA` The daemon holds new starts and resumes while disk space is low, and daemon status and the inbox say so
