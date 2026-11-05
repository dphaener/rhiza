# question-8992bae3-1

question · `01M6HN5D1F86QBKE6AA2JH5ND3`

## session

8992bae3-4669-468a-a617-01a1a3512e5a

## statement

Change set 'Workflow authoring' serves no intent yet. Its umbrella (Area workflows) serves anyone-builds and power-users-first. Which landed intent should it serve? / workflow-overrides says a built-in is replaced by name 'from the user's or the repo's own file'. The repo's .rhiza/workflows/<name>.toml already replaces a built-in. Should this change set add a user-level layer ($XDG_CONFIG_HOME/rhiza/workflows/<name>.toml, between the repo's file and the built-in)? / What does `rhiza workflow describe <name>` show? / How does `rhiza workflow new <name>` behave?

## Links

### In

- answers [workflow-authoring-shape](../decision/workflow-authoring-shape.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
