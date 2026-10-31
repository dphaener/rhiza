# plan-run-opens-pr-qa-1

evidence · `01M695815TN54K6CV4KTEN4D4R` · actor [claude-code-418647f6](../actor/claude-code-418647f6.md)

## result

PASS

## statement

QA cycle 1 at a3dd859: built rhiza with go build to a scratch path and ran it by hand in a scratch repo with a bare origin.git, a fake rex/claude (roles doing real rhiza writes) and a fake gh (no network). The run pushed cs/plan-reports and called gh pr create once with --title Plan reports; the body had '# Plan of run', Round 1 with 'confirmed f1 on export-csv ... — confirmed: ...', Round 2 'No findings.', Questions 'None.' and the Gate listing both change sets with states. In a second run held at wait-pr (CI pending), a new commit and RHIZA_RUN=<run> rhiza ship open-pr pushed the new head to origin, called gh pr edit 7 --body-file (no second pr create) with the refreshed body ('at step wait-pr'), and said 'Updated PR #7'. A run with red CI failed at pr-green and never called gh pr merge.

## Links

### Out

- verifies [plan-run-opens-pr](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-opens-pr.md)

## History

- `01M625RCFY4RB9HC0B69Y08VFV` A plan run ends with its plan on a pushed branch and an open PR, as ship-change does, so no one renders, commits or opens the plan's PR by hand.
