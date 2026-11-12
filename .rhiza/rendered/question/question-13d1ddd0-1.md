# question-13d1ddd0-1

question · `01M77XVZ8NDQW5V75QHY9SDTK2`

## session

13d1ddd0-9d8d-4ad6-8dce-01a1cfdb4993

## statement

Merge clash: change set 01M77N83 (this branch, Decision ship-land-model-conflict-repair) and 01M77N1G (landed on main via #484, Decision ship-land-model-conflict-rebase) each implement `rhiza ship check --model-conflict` and a model-conflict fixer route, and the two are incompatible. Clashing files: internal/api/ship.go, internal/model/apply.go, both ship-change.toml, workflow_test.go, cli.md and describe-golden. Main: land records one conflict (StaleBaseError → api.Error.Conflict) in state.model_conflict. --model-conflict passes only if HEAD equals the refused head and the run has made fewer than 3 commits with a Model-Conflict: trailer. It prints model_conflict with both intents. Route: land-triage → land-conflict-triage → model-conflict-fix → gate-land. This branch: land checks WhatIf before writing anything and records a model_conflicts list. --model-conflict passes while the recorded conflicts still hold in the model. ship check and `ship gate --model-conflict` also fail before push. Route: model-conflict-triage, gate-model-conflict-triage and pr-model-conflict-fix → render-model-conflict. Which design should the merge keep?

## Links

### In

- answers [merge-keeps-main-model-conflict-design](../decision/merge-keeps-main-model-conflict-design.md)

## History

- `01M77N83F027XBGGZPK9SV6YN2` Fix watchdog unmerged PR at step land-triage: run 01M77DWNRDSFY0PEP0S1P58P5W; log: PR green at 52ace47 and unmerged after step land-triage failed: step "land-triage" ended with outcome "fail" | exit 1 | No head the local gate failed on after a rebase is on the branch.
