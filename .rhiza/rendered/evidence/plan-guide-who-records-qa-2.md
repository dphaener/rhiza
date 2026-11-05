# plan-guide-who-records-qa-2

evidence · `01M6MV6JZPE7QHB6B6Q4XVG4Y3` · actor [claude-code-11d93090](../actor/claude-code-11d93090.md)

## result

PASS

## statement

QA cycle 2: ran rhiza guide plan from a binary built at the branch head. It says the run's capture records the owner's AskUserQuestion picks and the planner writes no Decision for a plain owner pick; the planner writes one only when the capture's wording would be wrong (a proxy such as an orchestrator or arbiter answered, or a ruling the capture cannot word), with the proxy's wording and the asked question copied verbatim into question so the capture skips it. No line tells the planner to write Decisions for plain picks.

## Links

### Out

- verifies [plan-guide-who-records](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-guide-who-records.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
