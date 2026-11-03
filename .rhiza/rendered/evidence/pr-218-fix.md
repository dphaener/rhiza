# pr-218-fix

evidence · `01M6HGTMJTQTNVYDFQNN8QKBH5` · actor [dphaener](../actor/dphaener.md)

## head

b417e1775cfff71995a230c757f8c03327c03ba2

## pr

https://github.com/dphaener/rhiza/pull/218

## result

1 check passed: check SUCCESS

## statement

Pull request #218 (https://github.com/dphaener/rhiza/pull/218) passed CI at b417e17

## status

verified

## Links

### Out

- verifies [scripts-free-disk-hermetic](../capability/rhiza/integrity/self-check/scripts-free-disk-hermetic.md)

## History

- `01M6GG7VES74Y7T74DD8HKXFJ7` Fix: Test scripts depend on the host: daemon scripts start nothing when TMPDIR sits on a small tmpfs (the minFree guard reads the host's free disk), and on Linux sh is dash so background daemons need sh -c '... exec rhiza daemon ...'. Scripts should set RHIZA_TEST_FREE_BYTES (or minFree 0) and use portable shell so make check passes the same on macOS and Ubuntu
