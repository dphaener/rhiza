# code-check-tolerates-branch-newer-format-qa-1

evidence · `01M735BFFRV31PKP68730CVN4V` · actor [claude-code-7a553445](../actor/claude-code-7a553445.md)

## result

PASS

## statement

QA claude-code-7a553445 at 2aa167c9 with a scratch build (go build -trimpath ./cmd/rhiza) in a hermetic HOME and Rhiza root: rhiza init, a change set with leaf alpha, a code workflow (wait gate then command step rhiza run code-check), rhiza run start --foreground. In the run worktree a builder commit (Leaf: alpha) and verifier Evidence; baseline code-check passed. (a) Committed on the run branch a real log entry with v=5 (FormatVersion+1, id matching its file name): plain rhiza check exits 1 'fail log ...: entry: newer format: field "v": 5, this Rhiza reads 4'; rhiza run code-check exits 0 'code-check passed', and --json lists it under rule log at severity warn. (c) A branch-added corrupt (invalid JSON) entry: code-check exits 1 with cause 'rhiza check: fail log ... invalid JSON'. (b) The same v5 entry committed on main and merged into the branch, so it is at the merge-base: code-check exits 1 with cause 'rhiza check: fail log ... newer format', while the branch-added v5 entry still only warns. Then, with only the branch-added v5 entry, released the gate and rhiza run resume --foreground: step code passes and the run succeeds. Scoped tests: go test -trimpath ./internal/api -run TestCheck ok, ./internal/check/... ok, cmd/rhiza TestScripts/(review-code-check|describe-golden) ok (with GOTMPDIR outside the repo).

## Links

### Out

- verifies [code-check-tolerates-branch-newer-format](../capability/code-check-tolerates-branch-newer-format.md)

## History

- `01M732D2EZSS20MWDHZGG66BJ5` Fix A run whose change bumps the log entry format can never pass review/code. Run 01M72YSG33J9JB66SYG9AD0Z0F (Docs sign-off workflow, change set 01M71G46Q5378CYYVEHFF2EMP3) failed at review/code with three causes like: 'rhiza check: fail log .rhiza/log/01M71G4BKFHP1E172RR6H82WFT/01M731FHDC43E8XR3XAPNAFE1Y.json: entry: newer format: field "v": 5, this Rhiza reads 4'. The builder's own binary wrote v5 entries (pinning a change set's workflow), but run code-check runs rhiza check in-process with the installed daemon binary (internal/api/run_codecheck.go runCheck; entry.ErrNewerFormat at internal/entry/entry.go:43), which reads only v4. #417 (01M6ZT7P) skipped the reference rules for the same root cause but not this one. Expected: code-check checks a run's worktree with that worktree's own build (or otherwise accepts the format its change introduces), so format bumps can ship through the daemon. Unstick: once fixed and merged (daemon self-updates), rhiza run resume --retry 01M72YSG33J9JB66SYG9AD0Z0F.
