# code-check-tolerates-branch-newer-format-verified

evidence · `01M73501M9FQWQPH7WYXE2RR85` · actor [claude-code-196a2d45](../actor/claude-code-196a2d45.md)

## result

PASS

## statement

Verifier claude-code-196a2d45 at 876ef85a, without editing code. Scoped tests: go test -trimpath ./internal/check/... ok; ./internal/api/... ok (544s, run from a detached worktree outside ~/code: the first run inside ~/code had no test failures but tripped the guard on writes to the enclosing root's .rhiza-root/usage.live.json, a known leak this change does not touch); TestCheckBranchNewerFormat PASS: a branch-added newer-format entry is a warning, one at the base fails, and with the option off it fails; cmd/rhiza TestScripts review-code-check, describe-golden and check* ok (describe-golden has branch_newer_format). By hand, with a scratch build in a hermetic HOME, a real run worktree and v=5 (FormatVersion+1): (a) entry committed on the run branch: rhiza check exits 1 with 'fail log ...: newer format: field v: 5, this Rhiza reads 4'; rhiza run code-check exits 0 'code-check passed' with no cause naming it, and --json shows it at severity warn; (b) the same entry committed before the merge-base: rhiza run code-check exits 1 with the cause 'rhiza check: fail log ... newer format', and check --branch-newer-format still fails it. Diff: the option only downgrades errors.Is(ErrNewerFormat) problems whose path is missing from the base listing, so other bad entries still fail; only run code-check sets it, not plain check or self-check. Gap: the testscript writes v:99 untracked rather than committing v=5, and it leaves the merge-base case to the unit test; the manual runs above cover both cases end to end.

## Links

### Out

- verifies [code-check-tolerates-branch-newer-format](../capability/code-check-tolerates-branch-newer-format.md)

## History

- `01M732D2EZSS20MWDHZGG66BJ5` Fix A run whose change bumps the log entry format can never pass review/code. Run 01M72YSG33J9JB66SYG9AD0Z0F (Docs sign-off workflow, change set 01M71G46Q5378CYYVEHFF2EMP3) failed at review/code with three causes like: 'rhiza check: fail log .rhiza/log/01M71G4BKFHP1E172RR6H82WFT/01M731FHDC43E8XR3XAPNAFE1Y.json: entry: newer format: field "v": 5, this Rhiza reads 4'. The builder's own binary wrote v5 entries (pinning a change set's workflow), but run code-check runs rhiza check in-process with the installed daemon binary (internal/api/run_codecheck.go runCheck; entry.ErrNewerFormat at internal/entry/entry.go:43), which reads only v4. #417 (01M6ZT7P) skipped the reference rules for the same root cause but not this one. Expected: code-check checks a run's worktree with that worktree's own build (or otherwise accepts the format its change introduces), so format bumps can ship through the daemon. Unstick: once fixed and merged (daemon self-updates), rhiza run resume --retry 01M72YSG33J9JB66SYG9AD0Z0F.
