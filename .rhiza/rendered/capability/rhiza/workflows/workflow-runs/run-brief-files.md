# run-brief-files

capability · `01M79KPBD9GFZM83ZPXXDRY6FS` · parent [workflow-runs](../workflow-runs.md)

## acceptance

go test -trimpath ./internal/run -run 'TestBriefFiles|TestSessionDir|TestHandbackFile' shows brief.md's sections in order with the change set, its intent, every leaf with its acceptance, the decisions, the constraints, both check commands (and 'none' when neither is set) and the contract summary's last line; the same input rendering the same bytes twice; session directory names for an agent step on its first and second visit, a builder of a leaf at attempt 2, a checks step's check, two critics of one cycle (distinct directories) and one check across two cycles (distinct directories), all names different; handback.md naming the failing step and outcome with each failing gate command and its output path for a gate cause, each failed CI run and its saved log file for a CI cause, and each conflicted file for a conflict cause; failure cases: a directory that cannot be created returns an error naming its path, and a cause naming no failing step is refused. Touches: internal/run/brief.go (new), internal/run/brief_test.go (new). Must not touch internal/api or the workflow parser.

## statement

internal/run renders a session's context files into its directory <commondir>/rhiza/runs/<run id>/<step>-<visit>/<session>/, <session> named from what tells the step's sessions apart in the code: agent for an agent step; <role>-<leaf>-<attempt> for a leaves step's builder or verifier; <check>-<role>-<cycle> for a checks step's check or critic (as OpenCheckCheckout names its checkout); <role>-<leaf>-<cycle>-<attempt> for a checks step's leaf fixer or verifier; fixer-<cycle>-<attempt> for a checks step's change set fixer; so no two sessions of a run, parallel or in a later visit or cycle, share a directory, which is kept when the run ends and lies inside no worktree: brief.md, deterministic Markdown with these sections in order: the change set (title, ID, intent), the leaves (handle, statement, acceptance each), the decisions (statement each), the constraints (statement each), the check commands (rhiza.ship.prePushCheck and rhiza.ship.localGate when set, else none) and a short contract summary ending 'Full contract: rhiza guide contract'; slice.json, the bytes it is given; and handback.md, written from a hand-back cause of any failing step: the failing step and its outcome, each failing command with the path of its kept output (a gate step's), each failed CI run with its saved log file (rhiza ship check --json findings), and each conflicted file, by a writer that run-handback calls and that nothing else in this change set triggers.

## Links

### Out

- decided-by [conductor-milestone-direction](../../../../decision/conductor-milestone-direction.md)
- decided-by [brief-files-and-invoke-placement](../../../../decision/brief-files-and-invoke-placement.md)

### In

- depends-on [session-brief-wiring](session-brief-wiring.md)
- blocks [review-01m7ah6rhkmzrahgjkjfm4vtnt-r3-f1](../../../../question/review-01m7ah6rhkmzrahgjkjfm4vtnt-r3-f1.md)
- verifies [run-brief-files-verified](../../../../evidence/run-brief-files-verified.md)
- verifies [run-brief-files-qa-1](../../../../evidence/run-brief-files-qa-1.md)
- verifies [run-brief-files-verified-2](../../../../evidence/run-brief-files-verified-2.md)
- verifies [run-brief-files-qa-2](../../../../evidence/run-brief-files-qa-2.md)
- verifies [pr-565-every-agent-session-a-run](../../../../evidence/pr-565-every-agent-session-a-run.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
