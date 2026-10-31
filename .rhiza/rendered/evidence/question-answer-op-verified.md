# question-answer-op-verified

evidence · `01M69VMJV0D9EQXXHEMDQFEAQ6` · actor [claude-code-3529bb86](../actor/claude-code-3529bb86.md)

## result

PASS

## statement

At 40d54f0d: go test ./cmd/rhiza -run 'TestScripts/(question-answer|describe-golden)$' passes (question-answer.txtar: question blocking an accepted change set drops out of ready after answer, --dry-run leaves ready unchanged, unknown question exits usage, two questions answered as answer-store-q and answer-cache-q, fold-ok); describe-golden includes question.answer and passes. go test ./internal/api/... ./internal/tui/... passes, incl. roadmapapp TestAnswerCallsQuestionAnswer (fake runner sees op question.answer, change set, args which-way|go west). By hand in a scratch repo with a scratch build and explicit --change <accepted cs> from another current change set: dry run printed dry_run true and git status unchanged; answer dropped store-q from ready; Decision answer-store-q has statement/quote 'use sqlite' and an answers link; second question got answer-other-q; a further dry-run answer to store-q proposed answer-store-q-2; unknown question exited 2. Diff read: roadmap app record() now calls m.run(question.answer) instead of hand-built apply writes.

## Links

### Out

- verifies [question-answer-op](../capability/rhiza/roadmap-app-accept-answer/question-answer-op.md)

## History

- `01M686RSVTRCWEV9GPT2V16M9Y` The inbox lists open questions that block proposed and accepted change sets and the inbox app answers them in place, through one question.answer operation the roadmap app shares
