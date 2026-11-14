# skill-check-builder

capability · `01M7DW59SVX2HFDCSYGMFDV37X` · parent [workflow-authoring](../workflow-authoring.md)

## acceptance

go test -trimpath ./cmd/rhiza -run TestScripts/skill-check shows, with a fake claude and a fake rex: a builder skill that commits with the three trailers and ends with rhiza session done passes every gate and exits 0; a skill that commits without trailers fails the trailer gate by name and exits 1; a skill that never runs rhiza session done fails the session-done gate by name and exits 1; --json prints the same gates with pass or fail each; an invoke whose skill is not installed is refused naming rhiza guide install --agent claude; the current repo's git status and .rhiza/log are unchanged after every case; go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden passes after regenerating it. go test -trimpath ./internal/api -run TestSkillCheck covers the gate report. Touches internal/api/skill_check.go, internal/api/skill_check_test.go, cmd/rhiza/testdata/script/skill-check.txtar and cmd/rhiza/testdata/script/describe-golden.txtar.

## statement

rhiza skill check <invoke> --role builder [--agent <name>] builds a scratch repo outside every git repository with a one-leaf fixture change set, runs the invoke as the builder's session through a one-step workflow in the foreground, and reports each contract gate (session done, a commit carrying Change:, Actor: and Leaf: trailers, rhiza check clean) as pass or fail, exiting 1 on any fail; --json prints the same report and the current repo is never touched.

## Links

### Out

- decided-by [skill-check-scratch-one-role](../../../../decision/skill-check-scratch-one-role.md)
- decided-by [skill-check-plan-shape](../../../../decision/skill-check-plan-shape.md)

### In

- depends-on [skill-check-verifier-check](skill-check-verifier-check.md)
- verifies [skill-check-builder-verified](../../../../evidence/skill-check-builder-verified.md)
- verifies [skill-check-builder-qa-1](../../../../evidence/skill-check-builder-qa-1.md)
- verifies [skill-check-builder-qa-2](../../../../evidence/skill-check-builder-qa-2.md)
- verifies [skill-check-builder-qa-3](../../../../evidence/skill-check-builder-qa-3.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
