# question-9e9e2812-1

question · `01M6GGG628MZGQBJ0FREHQQCVK`

## session

9e9e2812-e28f-4033-a5ee-01a1a10429d2

## statement

Daemon test scripts read the host's real free disk unless they set RHIZA_TEST_FREE_BYTES (only 5 of ~17 daemon scripts do), so on a host whose TMPDIR has <10G free the daemon starts nothing. How should scripts stop depending on the host's disk?

## Links

### In

- answers [test-setup-default-free-bytes](../decision/test-setup-default-free-bytes.md)

## History

- `01M6GG7VES74Y7T74DD8HKXFJ7` Fix: Test scripts depend on the host: daemon scripts start nothing when TMPDIR sits on a small tmpfs (the minFree guard reads the host's free disk), and on Linux sh is dash so background daemons need sh -c '... exec rhiza daemon ...'. Scripts should set RHIZA_TEST_FREE_BYTES (or minFree 0) and use portable shell so make check passes the same on macOS and Ubuntu
