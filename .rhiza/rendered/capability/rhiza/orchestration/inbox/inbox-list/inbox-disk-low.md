# inbox-disk-low

capability · `01M6ATBT8PTDQ8JQ49CV33ZG28` · parent [inbox-list](../inbox-list.md)

## acceptance

An inbox testscript with a daemon state recording a disk-low poll shows one item of kind disk-low in rhiza inbox --json, and none once a later poll records enough space

## statement

rhiza inbox lists one disk-low item while the daemon's last poll found free space under its minimum, naming the free space, and the inbox app shows it

## Links

### Out

- decided-by [disk-guard-10-gib](../../../../../decision/disk-guard-10-gib.md)

### In

- verifies [inbox-disk-low-verified](../../../../../evidence/inbox-disk-low-verified.md)
- verifies [inbox-disk-low-qa-1](../../../../../evidence/inbox-disk-low-qa-1.md)
- verifies [inbox-disk-low-verified-2](../../../../../evidence/inbox-disk-low-verified-2.md)
- verifies [inbox-disk-low-qa-2](../../../../../evidence/inbox-disk-low-qa-2.md)
- verifies [pr-162-disk-guard](../../../../../evidence/pr-162-disk-guard.md)

## History

- `01M6ATAS02XDRJ1MSQ4THGCNPA` The daemon holds new starts and resumes while disk space is low, and daemon status and the inbox say so
