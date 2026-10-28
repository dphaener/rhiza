# workflow-step-qa-3

evidence · `01M60EK64DWSBZG57X6BYRG74C` · actor [claude-code-173bd7bd](../actor/claude-code-173bd7bd.md)

## result

PASS

## statement

Built binary, repo overrides run with --dry-run: a workflow step with no workflow, an unknown workflow, one including itself (w3 -> w3) and through another (w5 -> w4 -> w5), skip naming a step build lacks, and workflow on a command step were each refused with a message naming the cause. A real run of an outer workflow whose workflow step names a three-command inner workflow skipping its second ran inner/one and inner/three in order (two never ran, by its log) and the run succeeded. The built-in ship-change journaled build/build, build/check, review/code and review/review under their workflows' names.

## Links

### Out

- verifies [workflow-step](../capability/ship-through-rhiza/workflow-step.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
