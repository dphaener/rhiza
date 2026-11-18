# conduct-handback-ship-rebase

capability · `01M7MNGTBKNAMKGF3847C1EP3A` · parent [agent-guide](../agent-guide.md)

## acceptance

go test -trimpath ./internal/guide -run 'TestConductHandback|TestConductReadsHandback|TestConductRoles|TestConductReadsRoles|TestWorkflowsText|TestList' passes, TestConductHandback requiring the topic to name rhiza ship rebase, the conflict-fix and gate-fix roles and routes for pr-green and rebase, and no longer requiring rhiza ship rebase --no-push; (failure case) TestConductHandbackCatchesDrift fails naming the step when the topic keeps a route for merge (which no longer hands back) and when it drops the rebase route, and fails when a route tells the conductor to run rhiza ship rebase. go test -trimpath ./cmd/rhiza -run 'TestDocsPrimitives|TestDocsExampleWorkflows|TestDocsBringYourOwn|TestScripts/guide' passes, docs/examples/byo/conductor.toml passing rhiza workflow check with steps wait-ci (until rhiza ship check --settled), pr-green, merge and rebase routed as ship-change's; grep -- --no-rebase docs/examples/byo/conductor.toml docs/bring-your-own-agent.md docs/primitives.md finds nothing, and bring-your-own-agent.md no longer says a failed merge hands back to conduct. grep -n 'ten steps' docs/runs-and-daemon.md finds nothing, its ship-change entry naming the 12 steps conduct, gate, signoff, push, wait-ci, pr-green, merge, rebase, main-ci, close-out, ship, deploy-check and no longer saying every red result up to the merge hands back to conduct. grep -n 'ship rebase --no-push' internal/guide/text/conduct-roles.md finds nothing, its conflict-fix entry naming a code conflict the ship steps leave in progress. grep -n 'pr-green' internal/guide/text/workflows.md and grep -n 'rebase' internal/guide/text/workflows.md each find ship-change's step list, which names the 12 steps conduct, gate, signoff, push, wait-ci, pr-green, merge, rebase, main-ci, close-out, ship and deploy-check. grep -n 'pr-green' docs/primitives.md finds a line, its Wait and Merge sections showing the pr-green, merge and rebase steps.

## files

internal/guide/text/conduct-handback.md, internal/guide/text/conduct-roles.md, internal/guide/text/workflows.md, internal/guide/conduct_handback_test.go, internal/guide/conduct_roles_test.go, internal/guide/guide.go (the conduct-handback description), docs/primitives.md, docs/bring-your-own-agent.md, docs/examples/byo/conductor.toml, docs/runs-and-daemon.md

## must_not

Change ship-change.toml, the hand-back mechanism or handback.md's rendering; change any fixer role topic (gate-fix, ci-fix, conflict-fix); grow a topic past TopicCap; tell the conductor to fix anything itself or to run rhiza ship rebase

## statement

The conduct-handback, conduct-roles and workflows guide topics, docs/primitives.md, docs/runs-and-daemon.md, docs/bring-your-own-agent.md and docs/examples/byo/conductor.toml describe the ship steps' own rebase: conduct-handback has routes for pr-green and rebase as conduct-handback-routes-ship-rebase decides and none for merge, never tells the conductor to run rhiza ship rebase, conduct-roles gives conflict-fix for a code conflict the ship steps leave in progress, workflows lists ship-change's 12 steps, runs-and-daemon.md's ship-change entry lists the 12 steps (with pr-green and rebase) and says a failed merge goes to rebase instead of every red up to the merge handing back, primitives.md's Wait and Merge sections show wait-ci without --no-rebase and the pr-green, merge and rebase steps, and the byo conductor example copies the built-in's steps after conduct (wait-ci on rhiza ship check --settled, pr-green, merge, rebase) so bring-your-own-agent.md's "Everything after conduct is the built-in's" holds, its hand-back paragraph naming pr-green and rebase instead of a failed merge.

## Links

### Out

- depends-on [ship-steps-rebase-clean](../../workflows/built-in-workflows/ship-steps-rebase-clean.md)
- decided-by [conduct-handback-routes-ship-rebase](../../../../decision/conduct-handback-routes-ship-rebase.md)

### In

- verifies [conduct-handback-ship-rebase-verified](../../../../evidence/conduct-handback-ship-rebase-verified.md)
- verifies [conduct-handback-ship-rebase-qa-1](../../../../evidence/conduct-handback-ship-rebase-qa-1.md)
- verifies [conduct-handback-ship-rebase-qa-2](../../../../evidence/conduct-handback-ship-rebase-qa-2.md)

## History

- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
