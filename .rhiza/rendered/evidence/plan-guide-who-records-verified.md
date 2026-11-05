# plan-guide-who-records-verified

evidence · `01M6MMB69QGSW6YWBT748HWVWM` · actor [claude-code-46ced5f1](../actor/claude-code-46ced5f1.md)

## result

PASS

## statement

At 958729ca ran go test -trimpath -count=1 ./internal/guide/... (ok) and -run TestPlanWhoRecords -v (PASS). The new test checks rhiza guide plan for: capture records the owner's AskUserQuestion picks, writes no Decision for a plain owner pick, the proxy exception (orchestrator or arbiter) with the proxy's wording, and copying the asked question verbatim into the question field; it also fails on any line pairing 'plain' with 'write a Decision'. Read plan.md diff: the old 'A picked option instead: wording=..., picked by the owner' recipe is replaced by a proxy-only recipe with question=, and a new 'Who records a pick' section states the rule. No remaining guide line tells the planner to write Decisions for plain picks. go vet and gofmt clean.

## Links

### Out

- verifies [plan-guide-who-records](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-guide-who-records.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
