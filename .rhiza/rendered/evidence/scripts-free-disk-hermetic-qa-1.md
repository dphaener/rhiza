# scripts-free-disk-hermetic-qa-1

evidence · `01M6HG54WRM33GW9KRZKV0EEJK` · actor [claude-code-ec2a1e14](../actor/claude-code-ec2a1e14.md)

## result

PASS

## statement

QA cycle 1 on the Ubuntu VM. Built the binary (go build -o $TMPDIR/rhiza-x ./cmd/rhiza); in a scratch repo, rhiza health shows the host's 63.2 GiB disk, 1.0 TiB with RHIZA_TEST_NOW+RHIZA_TEST_FREE_BYTES=1099511627776, bad at 1.0 GiB, and the host figure when RHIZA_TEST_NOW is unset. Ran go test -count=1 ./cmd/rhiza with TMPDIR on a 3.7G tmpfs (/dev/shm): at the merge-base 15 scripts fail (daemon-picks, daemon-status, daemon-pause, run-start-queue, stage-autonomy and others); at HEAD the whole package passes (ok, 474s). Pointed daemon-picks' reset line at 1073741824: it fails at line 78 (no 'Started a run of'); reverted, daemon-picks, health, daemon-notify-disk-low and free-bytes-default pass. No script sets RHIZA_TEST_FREE_BYTES to empty.

## Links

### Out

- verifies [scripts-free-disk-hermetic](../capability/rhiza/integrity/self-check/scripts-free-disk-hermetic.md)

## History

- `01M6GG7VES74Y7T74DD8HKXFJ7` Fix: Test scripts depend on the host: daemon scripts start nothing when TMPDIR sits on a small tmpfs (the minFree guard reads the host's free disk), and on Linux sh is dash so background daemons need sh -c '... exec rhiza daemon ...'. Scripts should set RHIZA_TEST_FREE_BYTES (or minFree 0) and use portable shell so make check passes the same on macOS and Ubuntu
