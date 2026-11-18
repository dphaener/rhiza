# rhiza-repo-conductor-default

capability · `01M79KPE3XB2JBTEB9ZZEQPEZF` · parent [workflow-overrides](../workflow-overrides.md)

## acceptance

In this repo, rhiza workflow describe ship-change resolves to the builtin; rhiza workflow describe ship-rhiza shows the [env], the go-cleanup cleanup running sh .rhiza/workflows/go-cleanup.sh, the builtin's 12 steps and every hand-back route as handback:ship-change/conduct; rhiza workflow describe bugfix still shows the go-cleanup cleanup and its hand-back routes as handback:ship-change/conduct; rhiza workflow check passes with no new warning beyond the autonomous ones AGENTS.md expects; rhiza check accepts .rhiza/dispatch.toml naming ship-rhiza for [ready] feature and bugfix for bug; go test -trimpath ./internal/workflow -run TestShipFixPromptsOneLine passes with .rhiza/workflows/ship-change.toml gone, its repo-file loop reading ship-rhiza and docs; make self-check passes. (failure case) TestShipFixPromptsOneLine fails when .rhiza/workflows/ship-rhiza.toml is missing, and go test -trimpath ./cmd/rhiza -run TestScripts/dispatch-check shows rhiza check failing on a dispatch.toml rule naming a workflow the repo and the built-ins do not resolve.

## files

.rhiza/workflows/ship-change.toml (removed), .rhiza/workflows/ship-rhiza.toml, .rhiza/dispatch.toml, .rhiza/workflows/bugfix.toml, internal/workflow/workflow_test.go (only TestShipFixPromptsOneLine's repo-file loop)

## must_not

Change the builtin ship-change.toml or internal/api; drop bugfix.toml's [env] or its go-cleanup cleanup; change plan-bug.toml or the docs workflows; change any test in internal/workflow/workflow_test.go other than TestShipFixPromptsOneLine's repo-file loop

## statement

Rhiza's own repo runs the builtin conductor default: its .rhiza/workflows/ship-change.toml shadow is removed, its [env] (GOTMPDIR, SHIP_HEAD_RHIZA, SHIP_HEAD_BUILD) and its [[cleanup]] go-cleanup (sh .rhiza/workflows/go-cleanup.sh) move to a thin .rhiza/workflows/ship-rhiza.toml that includes the builtin ship-change, .rhiza/dispatch.toml sets [ready] feature = "ship-rhiza", and bugfix.toml keeps its own [env] and [[cleanup]] go-cleanup and includes the new builtin. Through the include, every hand-back route of the builtin names the included conduct step (handback:ship-change/conduct). TestShipFixPromptsOneLine's loop over this repo's workflow files reads ship-rhiza instead of the removed ship-change shadow.

## Links

### Out

- depends-on [builtin-conductor-default](../built-in-workflows/builtin-conductor-default.md)
- depends-on [ship-head-build-on-demand](../../orchestration/shipping/ship-head-build-on-demand.md)
- depends-on [conductor-land-render-head-build](../../orchestration/shipping/conductor-land-render-head-build.md)
- depends-on [handback-target-survives-include](../built-in-workflows/handback-target-survives-include.md)
- decided-by [conductor-promote-builtin-now](../../../../decision/conductor-promote-builtin-now.md)
- decided-by [builtin-ship-change-open-points](../../../../decision/builtin-ship-change-open-points.md)

### In

- verifies [rhiza-repo-conductor-default-verified](../../../../evidence/rhiza-repo-conductor-default-verified.md)
- verifies [rhiza-repo-conductor-default-qa-2](../../../../evidence/rhiza-repo-conductor-default-qa-2.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
