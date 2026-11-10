# code-check-leaves-newer-format-to-self-check

decision · `01M732H1G36Y0JY1Y2BMZ7HEXT`

## reasoning

Same root cause as #417 (code-check-leaves-reference-to-self-check): the daemon's in-process check is the installed binary, so a run whose change bumps entry.FormatVersion can never pass review/code until merged and installed. That Decision rejected building the worktree's binary in the code check (generic across projects, Go toolchain and CPU cost on the shared VM) and leaned on make self-check, which builds the branch and runs in CI before every merge; the same reasoning holds here, so newer-format entries are left to it. Limiting the tolerance to entries the branch adds keeps the check honest for the case where the installed binary is genuinely behind main. build-check already folds without the unreadable entries (store.ReadLog keeps them as problems), and a builder writes Evidence at the lowest format that spells it, so leaving them out of the fold does not hide a verifier's Evidence.

## statement

Change set 01M732D2 serves m4-runs-hold-up, as earlier run-repair change sets do. The run code check (internal/api/run_codecheck.go checkIn/runCheck) stops failing a log entry that declares a newer format than the installed binary reads (entry.ErrNewerFormat) when the run's branch adds that entry since its merge-base: such an entry is reported as a warning, not a cause. A newer-format entry already at the merge-base, and every other unreadable or invalid entry, still fails. The code check sets this through an explicit CheckRequest option, so plain rhiza check and make self-check are unchanged and still judge every entry with the branch's own build in CI before merge. The code check does not build the run worktree's binary.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [code-check-newer-format-plan](../source/code-check-newer-format-plan.md)

### In

- decided-by [code-check-tolerates-branch-newer-format](../capability/code-check-tolerates-branch-newer-format.md)

## History

- `01M732D2EZSS20MWDHZGG66BJ5` Fix A run whose change bumps the log entry format can never pass review/code. Run 01M72YSG33J9JB66SYG9AD0Z0F (Docs sign-off workflow, change set 01M71G46Q5378CYYVEHFF2EMP3) failed at review/code with three causes like: 'rhiza check: fail log .rhiza/log/01M71G4BKFHP1E172RR6H82WFT/01M731FHDC43E8XR3XAPNAFE1Y.json: entry: newer format: field "v": 5, this Rhiza reads 4'. The builder's own binary wrote v5 entries (pinning a change set's workflow), but run code-check runs rhiza check in-process with the installed daemon binary (internal/api/run_codecheck.go runCheck; entry.ErrNewerFormat at internal/entry/entry.go:43), which reads only v4. #417 (01M6ZT7P) skipped the reference rules for the same root cause but not this one. Expected: code-check checks a run's worktree with that worktree's own build (or otherwise accepts the format its change introduces), so format bumps can ship through the daemon. Unstick: once fixed and merged (daemon self-updates), rhiza run resume --retry 01M72YSG33J9JB66SYG9AD0Z0F.
