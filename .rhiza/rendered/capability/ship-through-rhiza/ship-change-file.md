# ship-change-file

capability · `01M5Z8XSYAD7N2R7VF3J62AXNK` · parent [ship-through-rhiza](../ship-through-rhiza.md)

## acceptance

A testscript with fake agents and a fake gh runs `rhiza run start ship-change` on a change set of two leaves, one depending on the other, and finds both built and verified, review's checks passed, the PR opened with the summary as body, the run waiting at sign-off until `rhiza run approve`, then the change set landed, the PR merged and the close-out done; workflow tests find no `shipit` in ship-change.toml; `describe-golden` is regenerated.

## statement

The built-in ship-change.toml runs a workflow step on build; a render command step that runs rhiza render and commits .rhiza/rendered and the uncommitted .rhiza/log entries on the branch; a workflow step on review skipping its sign-off; a second render command step, render-pr, that does the same right before every push, including the push after rhiza run request-changes; ship open-pr; the wait for CI; the sign-off, whose on.changes goes back to review/review; and then today's rebase, land, merge, main check, close-out and ship steps, with no shipit prompt and no `ship check --pr-open` until. Both render steps leave the pushed branch holding every actor record its rendered pages are built from, so rhiza check passes on a clone of it.

## Links

### Out

- depends-on [workflow-step](workflow-step.md)
- depends-on [ship-open-pr](ship-open-pr.md)
- depends-on [ship-signoff](ship-signoff.md)
- decided-by [four-workflows](../../decision/four-workflows.md)
- decided-by [answer-cb8bd934-1](../../decision/answer-cb8bd934-1.md)
- decided-by [answer-cb8bd934-2](../../decision/answer-cb8bd934-2.md)
- decided-by [answer-cb8bd934-5](../../decision/answer-cb8bd934-5.md)
- derived-from [ship-change-plan](../../source/ship-change-plan.md)

### In

- blocks [leaf-failed-01m60atz-1](../../question/leaf-failed-01m60atz-1.md)
- blocks [leaf-failed-01m60atz-3](../../question/leaf-failed-01m60atz-3.md)
- blocks [leaf-failed-01m60atz-4](../../question/leaf-failed-01m60atz-4.md)
- verifies [ship-change-file-verified](../../evidence/ship-change-file-verified.md)
- verifies [ship-change-file-qa-1](../../evidence/ship-change-file-qa-1.md)
- verifies [ship-change-file-verified-3](../../evidence/ship-change-file-verified-3.md)
- verifies [ship-change-file-qa-3](../../evidence/ship-change-file-qa-3.md)
- verifies [ship-change-file-verified-4](../../evidence/ship-change-file-verified-4.md)
- verifies [ship-change-file-qa-4](../../evidence/ship-change-file-qa-4.md)
- verifies [ship-change-file-verified-5](../../evidence/ship-change-file-verified-5.md)
- verifies [ship-change-file-qa-5](../../evidence/ship-change-file-qa-5.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
