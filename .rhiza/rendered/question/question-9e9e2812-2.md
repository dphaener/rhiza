# question-9e9e2812-2

question · `01M6GGP08H2Z6RCM3PJAXCNG7Y`

## session

9e9e2812-e28f-4033-a5ee-01a1a10429d2

## statement

The intent's second half says background daemons on Linux need sh -c '... exec rhiza daemon ...' because sh is dash. On this Ubuntu VM (sh -> dash) every daemon script already passes, and the fake rex scripts that background a daemon already use exec "$@". What should this change set do about shell portability?

## Links

### In

- answers [decision-shell-portability](../decision/decision-shell-portability.md)

## History

- `01M6GG7VES74Y7T74DD8HKXFJ7` Fix: Test scripts depend on the host: daemon scripts start nothing when TMPDIR sits on a small tmpfs (the minFree guard reads the host's free disk), and on Linux sh is dash so background daemons need sh -c '... exec rhiza daemon ...'. Scripts should set RHIZA_TEST_FREE_BYTES (or minFree 0) and use portable shell so make check passes the same on macOS and Ubuntu
