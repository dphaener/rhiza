# plan-guide-who-records-qa-1

evidence · `01M6MN7K5J08Q6STAQ6360XCMG` · actor [claude-code-2253c893](../actor/claude-code-2253c893.md)

## result

PASS

## statement

QA at 628808db: ran rhiza guide plan from a scratch build. It has a 'Who records a pick' section saying the run's capture records the owner's AskUserQuestion picks and the planner writes no Decision for a plain owner pick; it writes one only when a proxy such as an orchestrator or arbiter answered or the answer carries a ruling the capture cannot word, with the proxy's wording and the asked question copied verbatim into its question field. The interview bullet now says 'For a plain owner pick you write no Decision', and the old 'A picked option instead' recipe is replaced by the proxy recipe. No line tells the planner to write a Decision for a plain pick. go test -trimpath ./internal/guide/...: ok.

## Links

### Out

- verifies [plan-guide-who-records](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-guide-who-records.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
