# daemon-disk-guard-qa-1

evidence · `01M6B38GKKE3RMMQCKQ5AE90YG` · actor [claude-code-e1013c3f](../actor/claude-code-e1013c3f.md)

## result

PASS

## statement

At 584cd6c3 built rhiza into $TMPDIR and ran it in a scratch repo (one accepted milestone change set with a leaf, Rex faked to log only, no test env vars). daemon --once --min-free 100G with the real statfs (57.8 GiB free): no Started, no rex call; daemon status shows 'Skipped: <cs> Built: disk space is low (57.8 GiB free, under the minimum of 100.0 GiB)' and --json reason disk-low. --min-free 1G: 'Started a run of <cs>', rex window new logged. --min-free -1, 10x, 1PB are usage errors (exit 2). The daemon-picks, inbox-disk-low and describe-golden scripts pass. Found two problems: RHIZA_TEST_FREE_BYTES works without RHIZA_TEST_NOW, and the disk-low Rex notice comes one poll late

## Links

### Out

- verifies [daemon-disk-guard](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-disk-guard.md)

## History

- `01M6ATAS02XDRJ1MSQ4THGCNPA` The daemon holds new starts and resumes while disk space is low, and daemon status and the inbox say so
