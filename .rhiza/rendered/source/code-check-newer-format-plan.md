# code-check-newer-format-plan

source · `01M732GG4GWNEN36YE7R7025QS`

## quotes

```json
[]
```

## raw_sha256

d06054e2d1c671633f431bfbf92dd4063b9330e08fe1203c7baa49c3acfab282

## session

5799f95c-eaa5-4624-8bd3-01a1c626e0d0

## summary

Autonomous plan run (no owner) for change set 01M732D2. Run 01M72YSG (change set 01M71G46) failed review/code because its builder's binary wrote v5 log entries and the code check runs rhiza check in-process with the installed binary, which reads v4. Following the #417 decision (code-check-leaves-reference-to-self-check: no worktree build in the code check; make self-check in CI judges with the branch build), the planner settled that the change set serves m4-runs-hold-up and that the code check stops failing newer-format entries the run's branch adds since its merge-base, through an explicit check option, leaving them to make self-check. Nothing was asked of the owner.

## title

Code check and newer-format log entries, autonomous plan

## Links

### In

- derived-from [code-check-leaves-newer-format-to-self-check](../decision/code-check-leaves-newer-format-to-self-check.md)

## History

- `01M732D2EZSS20MWDHZGG66BJ5` Fix A run whose change bumps the log entry format can never pass review/code. Run 01M72YSG33J9JB66SYG9AD0Z0F (Docs sign-off workflow, change set 01M71G46Q5378CYYVEHFF2EMP3) failed at review/code with three causes like: 'rhiza check: fail log .rhiza/log/01M71G4BKFHP1E172RR6H82WFT/01M731FHDC43E8XR3XAPNAFE1Y.json: entry: newer format: field "v": 5, this Rhiza reads 4'. The builder's own binary wrote v5 entries (pinning a change set's workflow), but run code-check runs rhiza check in-process with the installed daemon binary (internal/api/run_codecheck.go runCheck; entry.ErrNewerFormat at internal/entry/entry.go:43), which reads only v4. #417 (01M6ZT7P) skipped the reference rules for the same root cause but not this one. Expected: code-check checks a run's worktree with that worktree's own build (or otherwise accepts the format its change introduces), so format bumps can ship through the daemon. Unstick: once fixed and merged (daemon self-updates), rhiza run resume --retry 01M72YSG33J9JB66SYG9AD0Z0F.
