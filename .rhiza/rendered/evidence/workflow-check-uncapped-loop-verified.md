# workflow-check-uncapped-loop-verified

evidence · `01M7DZNAJZZY4ZNKXE98VG2X7H` · actor [claude-code-27256b3a](../actor/claude-code-27256b3a.md)

## result

PASS

## statement

Verified commit 4a4cef005 under a capped scope. Re-ran its Checks: go test -trimpath ./internal/workflow -run TestPlainRedLoops (exit 0; subtests: fail back to the agent step before it gives one loop, handback:build gives none, red route to a later step gives none; plus TestPlainRedLoopsNeedAnAgentStep and TestPlainRedLoopsShipChangeHasNone pass); go test -trimpath ./cmd/rhiza -run 'TestScripts/(custom-workflows-doc|run-cleanup-steps|workflow-check|workflow-check-invoke|workflow-check-risk|workflow-invoke|workflow-check-loop|workflow-new)$' (exit 0, covering the acceptance's workflow-check-loop and workflow-new); go test -trimpath ./internal/workflow ./internal/api -run 'PlainRedLoops|Workflow|Exposure|Invoke' (exit 0); rhiza ship generate (exit 0, 0 files changed); rhiza render (exit 0, only this verifier's actor record written); make self-check (exit 0). Diff read against the acceptance and Decisions uncapped-loop-is-plain-red-loop and skill-check-plan-shape: loops.go flags any non-pass outcome naming a plain step that leads back through an agent step; loopWarnings appends a warning naming step, outcome, target and suggesting handback:<agent>, never touching valid; workflow-check-loop.txtar checks the warning in text and --json with exit 0, none after switching to handback:build, none for ship-change and none for a workflow new skeleton, whose on.fail is now handback:build. The one file outside the listed touches, workflow-check-risk.txtar, narrows a '! stdout warning' to the exposure warning because the built-in plan workflow now warns about its own red loops; the commit explains it and a discovered-from draft change set covers the built-ins. No stubs or TODOs.

## Links

### Out

- verifies [workflow-check-uncapped-loop](../capability/rhiza/workflows/workflow-authoring/workflow-check-uncapped-loop.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
