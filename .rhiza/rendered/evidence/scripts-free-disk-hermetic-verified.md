# scripts-free-disk-hermetic-verified

evidence · `01M6HF3WF90EMZJJ08GHQ41XG4` · actor [claude-code-07f8e3e5](../actor/claude-code-07f8e3e5.md)

## result

PASS

## statement

At 84cb674d: cmd/rhiza/main_test.go setup sets RHIZA_TEST_FREE_BYTES=1099511627776 right after RHIZA_TEST_NOW/SEED. grep finds no script that sets RHIZA_TEST_FREE_BYTES to empty; daemon-picks.txtar:74 now resets it to 1099511627776. go test ./cmd/rhiza -run 'TestScripts/(health|daemon-notify-disk-low|daemon-picks|free-bytes-default)$' -count=1 passes. A temporary copy of free-bytes-default that sets nothing and runs `rhiza health --json` showed part disk with summary "1.0 TiB free, minimum 10.0 GiB". With daemon-picks line 74 changed to 1073741824 (under 10G), daemon-picks failed at line 78 (no 'Started a run of'), so the rest of the script reads the variable and not the host disk. Both temporary edits were reverted. Full cmd/rhiza package passes. Two things outside this leaf: make check self-check reports stale .rhiza/rendered for the builder's actor record, and internal/api TestReadBudget (a timing test this diff does not touch) failed on a heavily loaded machine.

## Links

### Out

- verifies [scripts-free-disk-hermetic](../capability/rhiza/integrity/self-check/scripts-free-disk-hermetic.md)

## History

- `01M6GG7VES74Y7T74DD8HKXFJ7` Fix: Test scripts depend on the host: daemon scripts start nothing when TMPDIR sits on a small tmpfs (the minFree guard reads the host's free disk), and on Linux sh is dash so background daemons need sh -c '... exec rhiza daemon ...'. Scripts should set RHIZA_TEST_FREE_BYTES (or minFree 0) and use portable shell so make check passes the same on macOS and Ubuntu
