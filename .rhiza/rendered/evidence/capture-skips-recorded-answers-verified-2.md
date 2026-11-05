# capture-skips-recorded-answers-verified-2

evidence · `01M6MV3A2E2DMQR68C5SFMXZHF` · actor [claude-code-e7c360e6](../actor/claude-code-e7c360e6.md)

## result

PASS

## statement

At 56b21fe4 (verifier, attempt 2) I read the leaf diff (recordedQuestions in internal/api/session_interview.go; captureAfterRun now returns skipped; session done and session run report a skipped field and 'Skipped <handle>: already recorded' lines) and ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(session-interview-recorded|session-done|session-interview|session-run|describe-golden)' and go test -trimpath ./internal/api/...: both ok. session-interview-recorded.txtar, with a fake rex and a transcript holding two picked answers, checks that a change-set Decision with question 'Which colour?' verbatim makes session interview skip and report answer-<sid8>-1 while a Decision whose question 'Which store' differs by one character still lets answer-<sid8>-2 be proposed; that after landing, a landed Decision derived-from the session's Source with the first question leaves the record capture writing only the 'Which store?' answer; and that session done proposes exactly [answer-<sid8>-2] and reports skipped [answer-<sid8>-1]. go vet and gofmt are clean. The quality nit that captureInterview's doc comment does not mention the question skip is still open, but it is outside the acceptance.

## Links

### Out

- verifies [capture-skips-recorded-answers](../capability/rhiza/lineage/interview-capture/capture-skips-recorded-answers.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
