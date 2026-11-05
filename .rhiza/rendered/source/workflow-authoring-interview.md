# workflow-authoring-interview

source · `01M6HNA9EN2M1MN4JYNX0Q3GB6`

## quotes

```json
[]
```

## raw_sha256

a30bc1c655ba40a6bc130c06d838687ccc9a79c0955392283d91fd3b3d9b4559

## session

8992bae3-4669-468a-a617-01a1a3512e5a

## summary

Autonomous plan run for change set Workflow authoring; no owner took part. The planner read the change set's intent, the workflows area, the landed decisions explain-custom-workflows, workflow-writing-tools and built-in-workflows-by-name, and internal/api/run.go (resolveWorkflow, workflowResolver), which today resolve a bare name to the repo's .rhiza/workflows/<name>.toml, else the built-in. It carved workflow-authoring and workflow-overrides out of Area workflows. Four questions went to the arbiter in one batch: which intent the change set serves, whether to add a personal (user-level) override layer, what workflow describe shows, and how workflow new behaves. The arbiter's Decision workflow-authoring-shape settled them: serve power-users-first; add the personal layer with order personal ($XDG_CONFIG_HOME/rhiza/workflows, else ~/.config/rhiza/workflows), then repo, then built-in, shared by run start, workflow steps, ls, describe and check; describe shows the resolved source, the layers it shadows and the steps with kinds, targets, roles and includes, read-only, with --json; new writes .rhiza/workflows/<name>.toml from --from <workflow> or a commented skeleton, bare names only, refusing an existing file with exit 3 and validating as check does. The planner then wrote leaves: workflow ls, describe, new, check and a guide topic under workflow-authoring, and workflow-overrides as the personal-layer leaf. Nothing was left open.

## title

Plan interview: workflow authoring (autonomous)

## Links

### In

- derived-from [workflow-authoring-shape](../decision/workflow-authoring-shape.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
