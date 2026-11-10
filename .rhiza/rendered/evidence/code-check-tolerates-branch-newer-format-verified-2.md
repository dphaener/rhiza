# code-check-tolerates-branch-newer-format-verified-2

evidence · `01M7368R7784QKSTW6TKT2VYS8` · actor [claude-code-7280e7a9](../actor/claude-code-7280e7a9.md)

## result

PASS

## statement

Verifier at 53b5c153: read the branch diff (CheckRequest.BranchNewerFormat flag; runCheck marks log problems not in the merge-base ls-tree as NewerOK; check.log downgrades only store.ErrNewerFormat on NewerOK paths to warn; run code-check sets it, plain check and self-check do not; describe-golden and docs/reference regenerated). review-code-check.txtar now commits an entry with v 5 (entry.FormatVersion is 4) on the run branch: plain rhiza check fails naming it, rhiza run code-check passes with exit 0 and no cause naming it, --json shows severity warn; then the same entry committed on main and merged into the branch (so it sits at the merge-base) makes run code-check fail with 1 cause naming the file. TestCheckBranchNewerFormat shows option off still fails the branch-added entry. Ran: go test -trimpath ./internal/api -run 'TestCheck|TestRunCodeCheck' -v (all PASS incl. TestCheckBranchNewerFormat), go test -trimpath ./internal/check/... (ok), go test -trimpath ./cmd/rhiza -run 'TestScripts/(review-code-check|describe-golden)$' (both PASS). The full ./internal/api/... run had every test PASS but the package was flagged only by the root-state guard on ~/code/.rhiza-root/usage.live.json, written by the live root daemon, unrelated to this change.

## Links

### Out

- verifies [code-check-tolerates-branch-newer-format](../capability/code-check-tolerates-branch-newer-format.md)

## History

- `01M732D2EZSS20MWDHZGG66BJ5` Fix A run whose change bumps the log entry format can never pass review/code. Run 01M72YSG33J9JB66SYG9AD0Z0F (Docs sign-off workflow, change set 01M71G46Q5378CYYVEHFF2EMP3) failed at review/code with three causes like: 'rhiza check: fail log .rhiza/log/01M71G4BKFHP1E172RR6H82WFT/01M731FHDC43E8XR3XAPNAFE1Y.json: entry: newer format: field "v": 5, this Rhiza reads 4'. The builder's own binary wrote v5 entries (pinning a change set's workflow), but run code-check runs rhiza check in-process with the installed daemon binary (internal/api/run_codecheck.go runCheck; entry.ErrNewerFormat at internal/entry/entry.go:43), which reads only v4. #417 (01M6ZT7P) skipped the reference rules for the same root cause but not this one. Expected: code-check checks a run's worktree with that worktree's own build (or otherwise accepts the format its change introduces), so format bumps can ship through the daemon. Unstick: once fixed and merged (daemon self-updates), rhiza run resume --retry 01M72YSG33J9JB66SYG9AD0Z0F.
