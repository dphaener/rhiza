# scripts-free-disk-hermetic

capability · `01M6GGX8QP23PVT25E5XDZ7X49` · parent [self-check](../self-check.md)

## acceptance

The testscript Setup in cmd/rhiza/main_test.go sets RHIZA_TEST_FREE_BYTES=1099511627776 beside RHIZA_TEST_NOW; a script that sets no RHIZA_TEST_FREE_BYTES runs rhiza health --json and sees the 1 TiB figure for disk; no script clears the variable to empty: daemon-picks.txtar, which ran `env RHIZA_TEST_FREE_BYTES=` after its disk-guard block, resets it to 1099511627776 instead, so its later runs never read the host's disk; health, daemon-notify-disk-low and daemon-picks pass, and with that reset line pointed at a value below 10G daemon-picks fails, showing the rest of the script uses the variable rather than the host's disk

## statement

Every test script sees a fixed 1 TiB of free disk unless it sets RHIZA_TEST_FREE_BYTES itself, so the daemon's minFree guard never reads the host's disk and make check passes the same whatever free space TMPDIR has

## Links

### Out

- decided-by [test-setup-default-free-bytes](../../../../decision/test-setup-default-free-bytes.md)
- decided-by [decision-shell-portability](../../../../decision/decision-shell-portability.md)

### In

- verifies [scripts-free-disk-hermetic-verified](../../../../evidence/scripts-free-disk-hermetic-verified.md)
- verifies [scripts-free-disk-hermetic-qa-1](../../../../evidence/scripts-free-disk-hermetic-qa-1.md)
- verifies [pr-218-fix](../../../../evidence/pr-218-fix.md)

## History

- `01M6GG7VES74Y7T74DD8HKXFJ7` Fix: Test scripts depend on the host: daemon scripts start nothing when TMPDIR sits on a small tmpfs (the minFree guard reads the host's free disk), and on Linux sh is dash so background daemons need sh -c '... exec rhiza daemon ...'. Scripts should set RHIZA_TEST_FREE_BYTES (or minFree 0) and use portable shell so make check passes the same on macOS and Ubuntu
