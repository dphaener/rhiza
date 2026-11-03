# test-setup-default-free-bytes

decision · `01M6GGJXKR6VPWTH7BVJN87SRP`

## reasoning

repoFreeBytes already honours RHIZA_TEST_FREE_BYTES only alongside RHIZA_TEST_NOW, which Setup always sets, so a Setup default makes every daemon script hermetic in one place instead of patching ~12 scripts and relying on future scripts to remember. This matches how Setup already makes HOME, git config, time and IDs hermetic (CLAUDE.md: deterministic scripts). Setting minFree 0 per script was rejected: it is again per-script, and it would skip the guard's code path rather than feed it a known figure. The scripts that test low disk already set the variable explicitly, so they keep working unchanged.

## statement

The testscript Setup in cmd/rhiza/main_test.go sets RHIZA_TEST_FREE_BYTES to a large default (1 TiB) for every script, next to RHIZA_TEST_NOW and RHIZA_TEST_SEED; scripts that exercise the low-disk guard (health, daemon-notify-disk-low) override it with env, and no script reads the host's real free disk

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-9e9e2812-1](../question/question-9e9e2812-1.md)
- derived-from [host-independent-scripts-interview](../source/host-independent-scripts-interview.md)

### In

- decided-by [scripts-free-disk-hermetic](../capability/rhiza/integrity/self-check/scripts-free-disk-hermetic.md)

## History

- `01M6GG7VES74Y7T74DD8HKXFJ7` Fix: Test scripts depend on the host: daemon scripts start nothing when TMPDIR sits on a small tmpfs (the minFree guard reads the host's free disk), and on Linux sh is dash so background daemons need sh -c '... exec rhiza daemon ...'. Scripts should set RHIZA_TEST_FREE_BYTES (or minFree 0) and use portable shell so make check passes the same on macOS and Ubuntu
