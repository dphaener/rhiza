# done-captures-answers-qa-1

evidence · `01M658X7S3YR7MYBAXSNSMV5C5` · actor [claude-code-f2616b4f](../actor/claude-code-f2616b4f.md)

## result

PASS

## statement

QA at 990e846 with a scratch build (go build ./cmd/rhiza), a scratch git repo, hermetic HOME and a fake rex: two sessions started with --permission-mode manual. A session with no transcript said done: exit 0, a 'propose decisions from the owner's answers: no transcript' warning ahead of the token-use warning, and session show reports it done. With the scenario's transcript (picked and typed answers) placed under HOME/.claude/projects, session done --dry-run proposed nothing and left no done mark; session done then printed 'Proposed decision answer-<id>-1..4', the change set held 4 node.create ops and answer-2 carried the owner's quote 'Something like DuckDB, but local'; a second done returned already:true with no decisions; session interview (the exit-capture path) proposed nothing and skipped all four, and the change set still held 4 node.create ops. The session-done-captures, session-run-tokens and describe-golden scripts and the TestSessionDone* Go tests (including TestSessionDoneCapturesBeforeMark) pass.

## Links

### Out

- verifies [done-captures-answers](../capability/rhiza/lineage/interview-capture/done-captures-answers.md)

## History

- `01M61GW89J8Y62X62RNJ8Z8HAD` A plan run captures the interview's picks and answers as Decisions when the interview step says done, not only when its session exits.
