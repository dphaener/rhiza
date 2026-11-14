# agent-step-checks-verified

evidence · `01M7C0EE3ECA1XQKXRV2YB6Q47` · actor [claude-code-c0c50494](../actor/claude-code-c0c50494.md)

## result

PASS

## statement

At 6bce0d26e re-ran the Checks block: go test -trimpath ./internal/workflow -run TestAgentStepChecks exit 0 (agent step checks parse in file order plan,code,qa; no-name and duplicate names refused); go test -trimpath ./internal/api -run TestRunFindingsAgentStepChecks exit 0 (run.findings and run.verdict --check plan record round 1 at the conduct agent step, --check scope refused naming plan, code, qa, --check at an agent step with no checks refused as not in a checks step); rhiza ship generate exit 0 with 0 files changed; rhiza render exit 0, tree clean. The behaviour itself is on main from #575 (run_review.go reviewStep/reviewCheck, workflow.go check parsing) in main's [[step.<name>.check]] name syntax per Decision agent-step-checks-already-built; this commit only renames the api test to the acceptance and adds the no-checks refusal case, plus the builder's actor record. No stubs or TODOs; nothing outside the leaf changed.

## Links

### Out

- verifies [agent-step-checks](../capability/rhiza/workflows/built-in-workflows/built-in-review/agent-step-checks.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
