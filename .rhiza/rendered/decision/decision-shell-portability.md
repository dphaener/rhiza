# decision-shell-portability

decision · `01M6GGRTK77A4WVXXV5BBZVYR0`

## reasoning

The question reports that on the Ubuntu VM (sh is dash) every daemon script already passes. The scripts confirm why: daemon-lock.txtar backgrounds the daemon with testscript's own `exec rhiza daemon ... &name&`, which runs no shell at all; the sh -c lines there are foreground and use only POSIX constructs; the fake rex scripts that background a daemon already use exec "$@". So the intent's premise that background daemons need sh -c '... exec rhiza daemon ...' does not hold for the current code, and rewriting passing scripts would be churn with no failing case to check against. The intent's goal is that make check passes the same on macOS and Ubuntu; the free-disk guard is the real host dependency, so the change set's work belongs there.

## statement

This change set makes no shell-portability rewrite. It fixes only the host free-disk dependency, through the testscript Setup default decided in test-setup-default-free-bytes. The dash half of the intent is treated as already satisfied: change a script's shell only if one is shown to fail under dash, and then only that script.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-9e9e2812-2](../question/question-9e9e2812-2.md)
- derived-from [host-independent-scripts-interview](../source/host-independent-scripts-interview.md)

### In

- decided-by [scripts-free-disk-hermetic](../capability/rhiza/integrity/self-check/scripts-free-disk-hermetic.md)

## History

- `01M6GG7VES74Y7T74DD8HKXFJ7` Fix: Test scripts depend on the host: daemon scripts start nothing when TMPDIR sits on a small tmpfs (the minFree guard reads the host's free disk), and on Linux sh is dash so background daemons need sh -c '... exec rhiza daemon ...'. Scripts should set RHIZA_TEST_FREE_BYTES (or minFree 0) and use portable shell so make check passes the same on macOS and Ubuntu
