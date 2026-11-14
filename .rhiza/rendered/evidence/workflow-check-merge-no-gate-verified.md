# workflow-check-merge-no-gate-verified

evidence · `01M7DZZGP6V0MHCDKQG4WFQK3X` · actor [claude-code-86a0bc25](../actor/claude-code-86a0bc25.md)

## result

PASS

## statement

Verified at 59967ab41. Re-ran every Checks: command, all exit 0: go test -trimpath -p 2 ./internal/workflow -run TestMergeWithoutGate (subtests merge reached only through a gate -> none, merge reachable around the gate -> merge with path build -> push -> merge, no merge -> none, plus ship-change has none); go test ./cmd/rhiza -run TestScripts/workflow-check-merge -count=1 (PASS: text warning naming step merge and path build -> push -> merge, the same in --json with valid true, exit 0; no warning once the bypass route is removed; none for ship-change); TestScripts/workflow-check; ./internal/workflow ./internal/api -run 'TestMergeWithoutGate|TestWorkflowCheck'; rhiza ship generate (0 changed); rhiza render (only this verifier's own actor). Read the diff: merge_gate.go does a BFS from the start step that never enters gate steps, detects merge steps by a regexp on rhiza ship merge / gh pr merge in command runs, and returns shortest paths in file order; workflow_check.go only appends warnings, so valid and the exit code stay the same. The code files are the four the acceptance names, plus the builder's actor log and the rendered files. No stubs or TODOs.

## Links

### Out

- verifies [workflow-check-merge-no-gate](../capability/rhiza/workflows/workflow-authoring/workflow-check-merge-no-gate.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
