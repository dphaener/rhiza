# workflow-authoring-shape

decision · `01M6HN89Y3S27AFDWTASSPBAAJ`

## reasoning

Read the question, the change set's intent (ls / describe / new --from / check + guide topic, carving workflow-authoring and workflow-overrides out of Area workflows), the umbrella's decisions and internal/api/run.go. Intent: the umbrella serves anyone-builds and power-users-first; writing custom TOML workflows is a power-user and agent feature, so power-users-first fits best. User layer: workflow-overrides says 'from the user's or the repo's own file' and its decision built-in-workflows-by-name fixes the order as 'personal, then repo, then built-in'. The change set carves workflow-overrides, so leaving the personal layer out would leave that leaf unbuilt; today resolveWorkflow and workflowResolver only know the repo file and the built-in. I followed the recorded order rather than the question's 'between repo and built-in' suggestion, because the model already decided it; the journal's recorded path keeps runs traceable. XDG_CONFIG_HOME with the ~/.config fallback is the Linux convention and the daemon runs on a Linux VM. describe: workflow-writing-tools puts describe beside --dry-run as a read-only aid, so it shows the resolved definition and its source layer, which is what an author needs to see what an override replaces. new: the change set's intent says new forks a built-in into .rhiza/workflows/<name>.toml with --from; writing into the repo keeps it reviewable in the change set, refusing to overwrite follows the never-silently-destroy norm and exit 3 is rhiza's conflict code, and validating on write ties it to check.

## statement

(1) The change set serves power-users-first: workflow authoring is a tool for agents and power users writing their own workflows. (2) Yes, it adds the personal layer: a bare workflow name resolves to $XDG_CONFIG_HOME/rhiza/workflows/<name>.toml (~/.config/rhiza/workflows when XDG_CONFIG_HOME is unset), then the repo's .rhiza/workflows/<name>.toml, then the built-in, in that order, as the landed decision built-in-workflows-by-name says (personal, then repo, then built-in). run start, workflow steps that name another workflow (workflowResolver), ls, describe and check all use that one resolution, and the run journal records the file path it resolved, so a run says which layer it came from. (3) `rhiza workflow describe <name>` resolves the name (or a path) as run start does and shows: where it resolved from (the file path or builtin:<name>) and which lower layers it shadows; its start step and each step in order with its kind, on.pass / on.fail targets, roles with their model ladders and prompts, and included workflows; --json gives the same as data. It never starts anything. (4) `rhiza workflow new <name> [--from <workflow>]` writes .rhiza/workflows/<name>.toml in the repo. With --from it copies the source of that workflow (resolved through the same layers, so --from ship-change forks the built-in) and sets its name to <name>; without --from it writes a minimal commented skeleton that passes check. <name> must be a bare name, not a path. It refuses with a conflict (exit 3) when the file already exists, never overwriting it; it validates the result as `workflow check` does and prints the path it wrote.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-8992bae3-1](../question/question-8992bae3-1.md)
- derived-from [workflow-authoring-interview](../source/workflow-authoring-interview.md)

### In

- decided-by [workflow-overrides](../capability/rhiza/workflows/workflow-overrides.md)
- decided-by [workflow-ls](../capability/rhiza/workflows/workflow-authoring/workflow-ls.md)
- decided-by [workflow-describe](../capability/rhiza/workflows/workflow-authoring/workflow-describe.md)
- decided-by [workflow-new](../capability/rhiza/workflows/workflow-authoring/workflow-new.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
