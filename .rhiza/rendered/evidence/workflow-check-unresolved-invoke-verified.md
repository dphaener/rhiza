# workflow-check-unresolved-invoke-verified

evidence · `01M7DZ61D8EFTMDGER0WQ88J6W` · actor [claude-code-429fff92](../actor/claude-code-429fff92.md)

## result

PASS

## statement

At 1363302a4 re-ran the Checks: block: go test -trimpath -p 2 ./internal/api -run 'TestWorkflowCheckUnresolvedInvoke|TestSkillCheck|TestExposureWarnings|StartSkill' exit 0 (TestWorkflowCheckUnresolvedInvoke/Missing asserts one warning each for the builder /rhiza-builder with the guide install command and the agent step /my-tidy, valid true and exit 0; /Installed asserts no warning with the skills in the agent's config folder and the repo); go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/(workflow-check|run-start-skill-missing)' exit 0 (workflow-check-invoke shows the warnings in text and --json, exit 0, none once installed); go vet ./internal/api exit 0; rhiza ship generate exit 0 with 0 changed files; rhiza render changed only my own actor record. A scratch build's rhiza workflow check bugfix gives no spurious invoke warning. The diff touches only the four files the acceptance names plus rendered actor records, shares checkStartSkill's lookup as missingStartSkill, leaves valid and the exit code alone as skill-check-plan-shape decides, and has no stubs or TODOs.

## Links

### Out

- verifies [workflow-check-unresolved-invoke](../capability/rhiza/workflows/workflow-authoring/workflow-check-unresolved-invoke.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
