# conduct-handback-ship-rebase-qa-1

evidence · `01M7P15KJ8R8FB9YPRSTQ4VJMC` · actor [claude-code-90cc7b68](../actor/claude-code-90cc7b68.md)

## result

PASS

## statement

QA: built binary's rhiza guide conduct-handback has pr-green and rebase routes, none for merge, says never run rhiza ship rebase; conduct-roles conflict-fix names a conflict the ship steps leave in progress; workflows lists the 12 steps; docs/examples/byo/conductor.toml passes rhiza workflow check in a scratch repo with wait-ci/pr-green/merge/rebase routed as ship-change's, and a broken copy fails with named steps (exit 1); greps for --no-rebase, 'ten steps', 'ship rebase --no-push' find nothing; go test ./internal/guide and TestDocsPrimitives|TestDocsExampleWorkflows|TestDocsBringYourOwn|TestScripts/guide pass

## Links

### Out

- verifies [conduct-handback-ship-rebase](../capability/rhiza/guide/agent-guide/conduct-handback-ship-rebase.md)

## History

- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
