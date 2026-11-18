# conduct-handback-ship-rebase-verified

evidence · `01M7NZ62FFZZGDEAGEYPX5CVYG` · actor [claude-code-259af63c](../actor/claude-code-259af63c.md)

## result

PASS

## statement

Verified at db1ea96d1 in a detached worktree (capped scopes, RHIZA_/SHIP_/REX_ unset for cmd/rhiza). go test -p 2 -trimpath ./internal/guide/... ok; -run 'TestConductHandback|TestConductReadsHandback|TestConductRoles|TestConductReadsRoles|TestWorkflowsText|TestList' -v: all PASS incl. TestConductHandbackCatchesDrift, which now covers a merge route, a dropped rebase route and a route running rhiza ship rebase, each reported by name. go test -p 2 -trimpath ./cmd/rhiza -run 'TestDocsPrimitives|TestDocsExampleWorkflows|TestDocsBringYourOwn|TestScripts/guide': ok (TestDocsBringYourOwnConductorSteps PASS). Failure case: changing conductor.toml's merge on.fail to handback:conduct makes TestDocsBringYourOwnConductorSteps fail naming step merge (restored). A scratch build's rhiza workflow check passes on docs/examples/byo/conductor.toml (only the expected uninstalled-skill warning); wait-ci uses rhiza ship check --settled, then pr-green, merge (on.fail rebase) and rebase. Greps: --no-rebase in conductor.toml/bring-your-own-agent.md/primitives.md, 'ten steps' in runs-and-daemon.md, 'ship rebase --no-push' in conduct-roles.md and 'failed merge' in bring-your-own-agent.md all find nothing; pr-green and rebase in workflows.md find the 12-step list at line 51; pr-green in primitives.md finds the Wait/Merge sections. The diff leaves ship-change.toml, internal/workflow, handback.md and the gate-fix/ci-fix/conflict-fix topics untouched, the TopicCap tests pass, and there are no TODOs. Its only files outside the leaf's list are tests (docs_examples_test.go, workflows_test.go, guide.txtar).

## Links

### Out

- verifies [conduct-handback-ship-rebase](../capability/rhiza/guide/agent-guide/conduct-handback-ship-rebase.md)

## History

- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
