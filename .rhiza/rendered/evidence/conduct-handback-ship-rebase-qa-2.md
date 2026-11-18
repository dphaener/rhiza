# conduct-handback-ship-rebase-qa-2

evidence · `01M7P1AFNPJHRZMNKVCQGQ1QEE` · actor [claude-code-1bba01de](../actor/claude-code-1bba01de.md)

## result

PASS

## statement

QA with the scratch build: guide conduct-handback has routes for pr-green and rebase, says the merge step never hands back and never to run rhiza ship rebase, no --no-push; guide conduct-roles gives conflict-fix for a conflict the ship steps leave in progress (CI wait, rebase step); guide workflows lists the 12 steps with pr-green and rebase. workflow check docs/examples/byo/conductor.toml valid (only the uninstalled /myteam:conduct warning); its steps from gate on are identical (comments stripped) to workflow new --from ship-change; pair workflow new --from <byo path> then workflow check valid. Failure: byo copy with handback:nope -> exit 1 naming gate, push, pr-green, rebase; guide with a typo name exit 4. Greps: no --no-rebase in the byo docs, no 'ten steps' in runs-and-daemon.md (12-step entry, failed merge goes to rebase), no 'ship rebase --no-push' in conduct-roles.md, pr-green in primitives.md.

## Links

### Out

- verifies [conduct-handback-ship-rebase](../capability/rhiza/guide/agent-guide/conduct-handback-ship-rebase.md)

## History

- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
