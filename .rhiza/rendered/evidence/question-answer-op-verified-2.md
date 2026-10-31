# question-answer-op-verified-2

evidence · `01M69XB6WEB3NYA0AG807AHHJ4` · actor [claude-code-6f2afadc](../actor/claude-code-6f2afadc.md)

## result

PASS

## statement

At bf9829b4 (attempt 2): go test ./cmd/rhiza -run 'TestScripts/(question-answer|describe-golden)$' -count=1 passes. question-answer.txtar has questions blocking capabilities of the accepted change set; rhiza question answer <q> <words> --json drops them from ready; --dry-run reports dry_run and ready still lists the question; an unknown question returns code usage (a built binary exits 2 on it, checked by hand); two answers land as answer-store-q and answer-cache-q; the new step answers queue-q with --change <accepted cs> while change set Other is current, the Decision is in the accepted change set and not in Other, ready drops queue-q, and fold-ok passes. go test ./internal/tui/roadmapapp passes, with roadmapapp_test.go asserting the answer path calls question.answer via a fake runner (answer.go:108). describe-golden passes against the regenerated golden. Diff read: the fix commit only adds the --change scenario.

## Links

### Out

- verifies [question-answer-op](../capability/rhiza/roadmap-app-accept-answer/question-answer-op.md)

## History

- `01M686RSVTRCWEV9GPT2V16M9Y` The inbox lists open questions that block proposed and accepted change sets and the inbox app answers them in place, through one question.answer operation the roadmap app shares
