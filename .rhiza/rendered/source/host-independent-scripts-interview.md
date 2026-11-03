# host-independent-scripts-interview

source · `01M6GGTYW2E9JCR33JAZVB7YNF`

## quotes

```json
[]
```

## raw_sha256

687a2cd2d5db13169295c13f9f00f0c9b5612e58c6c0d13130c2b3126c26e1f0

## session

9e9e2812-e28f-4033-a5ee-01a1a10429d2

## summary

Autonomous plan run for the Fix change set on test scripts that depend on the host. Code reading found daemonDiskLow reads the repo filesystem's real free space unless RHIZA_TEST_FREE_BYTES is set alongside RHIZA_TEST_NOW, and only a handful of daemon scripts set it, so on a host whose TMPDIR has under 10G free the daemon scripts start nothing. On the Ubuntu VM (sh is dash) every daemon script already passes, and the fake rex scripts that background a daemon already use exec. Two questions went to the arbiter. It decided that the testscript Setup sets RHIZA_TEST_FREE_BYTES to 1 TiB for every script, with low-disk scripts overriding it, and that this change set makes no shell-portability rewrite: the dash half of the intent is treated as already satisfied, and a script's shell changes only if it is shown to fail under dash. Nothing is left open.

## title

Host-independent test scripts, autonomous plan interview

## Links

### In

- derived-from [test-setup-default-free-bytes](../decision/test-setup-default-free-bytes.md)
- derived-from [decision-shell-portability](../decision/decision-shell-portability.md)

## History

- `01M6GG7VES74Y7T74DD8HKXFJ7` Fix: Test scripts depend on the host: daemon scripts start nothing when TMPDIR sits on a small tmpfs (the minFree guard reads the host's free disk), and on Linux sh is dash so background daemons need sh -c '... exec rhiza daemon ...'. Scripts should set RHIZA_TEST_FREE_BYTES (or minFree 0) and use portable shell so make check passes the same on macOS and Ubuntu
