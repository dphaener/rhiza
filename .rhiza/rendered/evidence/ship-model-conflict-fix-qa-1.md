# ship-model-conflict-fix-qa-1

evidence · `01M77SJEQB3T19VH3ZKD127V4E` · actor [claude-code-e1432792](../actor/claude-code-e1432792.md)

## result

PASS

## statement

QA with the binary built at 668f5cf08: rhiza workflow describe ship-change shows pr-green on.fail -> pr-model-conflict-triage (rhiza ship check --model-conflict) on.pass -> pr-model-conflict-fix; gate-land on.fail -> gate-model-conflict-triage on.pass -> model-conflict-fix; land on.fail -> land-triage on.fail -> model-conflict-triage on.pass -> model-conflict-fix -> render-model-conflict; the fixer prompts name rhiza set <ref> --rebase, the until check and the Review: {run} trailer. The repair the prompt prescribes worked by hand: rhiza set first --rebase on the stale node made ship check --model-conflict exit 1 and ship gate --model-conflict pass. .rhiza/workflows/ship-change.toml matches the builtin apart from the pre-existing head-build steps. go test ./internal/workflow/... -run TestShipChange and go test ./cmd/rhiza -run TestScripts/run-ship-model-conflict-fix$ ok.

## History

- `01M77N83F027XBGGZPK9SV6YN2` Fix watchdog unmerged PR at step land-triage: run 01M77DWNRDSFY0PEP0S1P58P5W; log: PR green at 52ace47 and unmerged after step land-triage failed: step "land-triage" ended with outcome "fail" | exit 1 | No head the local gate failed on after a rebase is on the branch.
