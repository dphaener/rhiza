# inbox-disk-low-qa-2

evidence · `01M6B52Z2EZBSHWQ8FBWWHJADK` · actor [claude-code-2be1dffb](../actor/claude-code-2be1dffb.md)

## result

PASS

## statement

At 5218260b, same scratch repo and built binary, real statfs with --min-free 100G then 1G, each daemon --once: poll 1 (low) emits the OSC 7501 state=blocked notice 'disk space is low: the daemon holds new starts and resumes' on the very first --once poll, and rhiza inbox --json lists one item of kind disk-low; poll 2 (still low) sends no notice; poll 3 (enough space) starts the run, sends no disk-low notice and the inbox has no disk-low item; poll 4 (low again) sends a fresh notice. The cycle-1 one-poll-late notice is fixed.

## Links

### Out

- verifies [inbox-disk-low](../capability/rhiza/orchestration/inbox/inbox-list/inbox-disk-low.md)

## History

- `01M6ATAS02XDRJ1MSQ4THGCNPA` The daemon holds new starts and resumes while disk space is low, and daemon status and the inbox say so
