# plan-run-opens-pr-verified

evidence · `01M68JEEF6YN0A32A6XP8VMVBN` · actor [claude-code-4f4650ac](../actor/claude-code-4f4650ac.md)

## result

PASS

## statement

At 101d41a ran go test ./cmd/rhiza -run TestScripts/plan-workflow -count=1 (ok): the new plan-workflow-pr.txtar, with a bare origin.git and a fake gh, asserts gh.log holds exactly one pr create titled 'Plan reports' (the change set's title), pr-body.md holds the Round 1/Round 2 headings, the confirmed f1 finding line, 'No findings.' and the gate line, origin's refs/heads/cs/plan-reports equals the branch head, and the ship state records pr 7. go test ./internal/api -run TestShipOpenPR -v passed TestShipOpenPRRefusesPlanWithoutRound (a plan run with no round exits 2 'no review journaled' and calls no gh). Workflow tests pass with gate -> open-pr -> succeeded. Read the diff (plan.toml open-pr step, planSummary, open-pr branch for plan runs). make check exit 0.

## Links

### Out

- verifies [plan-run-opens-pr](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-opens-pr.md)

## History

- `01M625RCFY4RB9HC0B69Y08VFV` A plan run ends with its plan on a pushed branch and an open PR, as ship-change does, so no one renders, commits or opens the plan's PR by hand.
