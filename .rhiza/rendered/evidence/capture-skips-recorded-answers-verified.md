# capture-skips-recorded-answers-verified

evidence · `01M6MN195YSAM7F81X26MBNAXF` · actor [claude-code-9d2d772b](../actor/claude-code-9d2d772b.md)

## result

PASS

## statement

Verified at b5f725e3. Read the diff: captureInterview (shared by session interview and session done, session_run.go:462) now skips any answer whose question equals, verbatim, the question field of a live Decision the change set's ops created, or once landed a live Decision derived-from the session's Sources, and reports it in skipped. Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/session-(interview|done)' -count=1: ok (session-interview, -landed, -source, -recorded, session-done-captures). session-interview-recorded.txtar covers the acceptance: with a Decision 'Which colour?' only answer-<sid8>-2 is proposed and answer-1 is reported skipped; a Decision 'Which store' (one character off) does not stop answer-2; with the record change set landed and a landed Decision derived-from the session's Source carrying 'Which colour?', the record capture writes only the 'Which store?' answer. Mutation check: with session_interview.go reverted to HEAD~1 in a scratch worktree, the script fails at line 34 (answer-158fe87f-1 proposed). go test ./internal/api/...: every test passes except TestReadBudget (roadmap what-if timing median over 50ms at load average ~6.5), a read-path timing budget this diff does not touch. gofmt clean.

## Links

### Out

- verifies [capture-skips-recorded-answers](../capability/rhiza/lineage/interview-capture/capture-skips-recorded-answers.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
