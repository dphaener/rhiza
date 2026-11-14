# skill-check-verifier-check

capability · `01M7DW59XN8KSH58BTAYJHPEX6` · parent [workflow-authoring](../workflow-authoring.md)

## acceptance

go test -trimpath ./cmd/rhiza -run TestScripts/skill-check-roles shows, with a fake claude and a fake rex: a verifier skill that writes Evidence and ends with rhiza session done passes every gate; one that writes no Evidence fails the Evidence gate by name and exits 1; a check skill that records a round with rhiza run findings passes; one that records none fails the findings gate by name and exits 1; an unknown --role is refused naming builder, verifier and check; the current repo is unchanged. go test -trimpath ./internal/api -run TestSkillCheckRoles covers the fixture and the gates. Touches internal/api/skill_check.go, internal/api/skill_check_test.go and cmd/rhiza/testdata/script/skill-check-roles.txtar.

## statement

rhiza skill check --role verifier and --role check run the invoke against a fixture whose leaf already has a builder commit: the verifier's gates are session done, Evidence on the leaf written by the verifier, and rhiza check clean; the check's gates are session done, a round recorded with rhiza run findings for its check, and rhiza check clean.

## Links

### Out

- depends-on [skill-check-builder](skill-check-builder.md)
- decided-by [skill-check-scratch-one-role](../../../../decision/skill-check-scratch-one-role.md)
- decided-by [skill-check-plan-shape](../../../../decision/skill-check-plan-shape.md)

### In

- verifies [skill-check-verifier-check-verified](../../../../evidence/skill-check-verifier-check-verified.md)
- verifies [skill-check-verifier-check-qa-1](../../../../evidence/skill-check-verifier-check-qa-1.md)
- verifies [skill-check-verifier-check-qa-2](../../../../evidence/skill-check-verifier-check-qa-2.md)
- verifies [skill-check-verifier-check-qa-3](../../../../evidence/skill-check-verifier-check-qa-3.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
