# daemon-run-tabs-qa-1

evidence · `01M629Y6RBP93TAAR2QPBFX58S` · actor [claude-code-715b2246](../actor/claude-code-715b2246.md)

## result

FAIL

## statement

Built rhiza at f087c0e into $TMPDIR/rhiza-qa and drove it by hand in a scratch repo (RHIZA_* env unset, rex shimmed on PATH): the daemon's rex call was 'rex window new <cs[:9]> -c <repo> --json -- <rhiza> run start <cs>' with no --autonomous, as the acceptance says. With a rex shim that actually ran the tab's command, every started run exited 2 at once: 'step build/build is a leaves step whose builder runs with no permission mode: give --permission-mode or set permission_mode on the step' (also with --autonomous), writing no journal. A daemon --max-runs 1 left 80s then opened two tabs for ShipA, one per minute, both failing the same way. So the tab command matches the acceptance but no run the daemon starts on the built-in ship-change workflow can run.

## Links

### Out

- verifies daemon-run-tabs (retired)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
