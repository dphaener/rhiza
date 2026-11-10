# pr-424-fix-a-run-whose-change-bumps

evidence · `01M737EJBDPMG7NK5ZQKDDVYJ9` · actor [dphaener](../actor/dphaener.md)

## head

3d5fe2393dcc116e86f0f5ceb681179e5af895bd

## pr

https://github.com/dphaener/rhiza/pull/424

## result

11 checks passed: lint SUCCESS, go-install SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #424 (https://github.com/dphaener/rhiza/pull/424) passed CI at 3d5fe23

## status

verified

## Links

### Out

- verifies [code-check-tolerates-branch-newer-format](../capability/code-check-tolerates-branch-newer-format.md)

## History

- `01M732D2EZSS20MWDHZGG66BJ5` Fix A run whose change bumps the log entry format can never pass review/code. Run 01M72YSG33J9JB66SYG9AD0Z0F (Docs sign-off workflow, change set 01M71G46Q5378CYYVEHFF2EMP3) failed at review/code with three causes like: 'rhiza check: fail log .rhiza/log/01M71G4BKFHP1E172RR6H82WFT/01M731FHDC43E8XR3XAPNAFE1Y.json: entry: newer format: field "v": 5, this Rhiza reads 4'. The builder's own binary wrote v5 entries (pinning a change set's workflow), but run code-check runs rhiza check in-process with the installed daemon binary (internal/api/run_codecheck.go runCheck; entry.ErrNewerFormat at internal/entry/entry.go:43), which reads only v4. #417 (01M6ZT7P) skipped the reference rules for the same root cause but not this one. Expected: code-check checks a run's worktree with that worktree's own build (or otherwise accepts the format its change introduces), so format bumps can ship through the daemon. Unstick: once fixed and merged (daemon self-updates), rhiza run resume --retry 01M72YSG33J9JB66SYG9AD0Z0F.
