# ship-needs-leaves

capability · `01M5Z8XSFSMXX93SBZP3ZZPZTQ` · parent [ship-through-rhiza](../ship-through-rhiza.md)

## acceptance

A testscript finds `rhiza run start ship-change` on a change set without leaves exiting non-zero with the refusal and `rhiza run start --workflow plan` as next, no run journaled and no worktree made, and the same command on a change set with a leaf starting the run.

## statement

`rhiza run start ship-change` refuses a change set with no leaves before writing anything, its next command `rhiza run start --workflow plan`.

## Links

### Out

- decided-by [model-is-plan-and-state](../../decision/model-is-plan-and-state.md)
- decided-by [answer-cb8bd934-5](../../decision/answer-cb8bd934-5.md)
- derived-from [ship-change-plan](../../source/ship-change-plan.md)

### In

- blocks [leaf-failed-01m60atz-2](../../question/leaf-failed-01m60atz-2.md)
- verifies [ship-needs-leaves-verified](../../evidence/ship-needs-leaves-verified.md)
- verifies [ship-needs-leaves-qa-1](../../evidence/ship-needs-leaves-qa-1.md)
- verifies [ship-needs-leaves-verified-3](../../evidence/ship-needs-leaves-verified-3.md)
- verifies [ship-needs-leaves-qa-3](../../evidence/ship-needs-leaves-qa-3.md)
- verifies [ship-needs-leaves-verified-4](../../evidence/ship-needs-leaves-verified-4.md)
- verifies [ship-needs-leaves-qa-4](../../evidence/ship-needs-leaves-qa-4.md)
- verifies [ship-needs-leaves-qa-5](../../evidence/ship-needs-leaves-qa-5.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
