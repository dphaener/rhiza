# workflow-dispatch-shape

decision · `01M6VYH807FA05BQQ196XC8JPN`

## reasoning

Milestone: 29 is about several projects under one root and does not fit; 25 is the backlog milestone for owner-requested work the daemon builds on its own, and the sibling workflow-authoring change set 01M6HN24 served it, so this follow-up belongs there too. Rules location: internal/rootcfg's rhiza.toml is machine-wide, uncommitted and shared by several projects, so per-repo rules there would leak across projects; repo workflows already resolve from .rhiza/workflows/ (internal/api/run.go workflowsDir), and decision workflows-are-the-rules (commit 6284f5c8) says the rules are repo policy the owner defines with built-in defaults, which a committed .rhiza/ file is. A separate file rather than <name>.toml inside .rhiza/workflows/ avoids colliding with workflow names. Override: the intent says the change set records the request's workflow 'so the ship pick follows it', i.e. the override targets the ready/ship trigger; a request's plan run (internal/api/request.go starts planWorkflow) stays rule-driven so planning is not bypassed by a ship-workflow name. Manual start: one resolution path for daemon and CLI keeps 'everything runs through a workflow' consistent and matches today, where both default to ship-change.

## statement

Workflow dispatch serves @milestone-25, the backlog milestone. The dispatch rules (trigger x request kind -> workflow) live in a file committed in the repo under .rhiza/ (for example .rhiza/dispatch.toml, beside .rhiza/workflows/), not in the root's machine-wide rhiza.toml; with no repo file, built-in default rules reproduce today's behaviour (request queued -> plan, change set ready -> ship-change). `rhiza request --workflow <name>` overrides the ship pick only: the change set its plan run opens records the kind and the named workflow, and the daemon's ready trigger starts that workflow instead of the rule's; the request's plan run still follows the rules for its kind. A manual `rhiza run start <cs>` with no --workflow resolves the same way as the daemon: the change set's recorded workflow, else the rules for its kind, else the default; an explicit --workflow always wins.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-b640ca70-1](../question/question-b640ca70-1.md)
- derived-from [workflow-dispatch-plan-interview](../source/workflow-dispatch-plan-interview.md)

### In

- decided-by [workflow-dispatch](../capability/rhiza/workflows/workflow-dispatch.md)
- decided-by [dispatch-rules-file](../capability/rhiza/workflows/workflow-dispatch/dispatch-rules-file.md)
- decided-by [change-set-records-dispatch](../capability/rhiza/workflows/workflow-dispatch/change-set-records-dispatch.md)
- decided-by [request-workflow-flag](../capability/rhiza/workflows/workflow-dispatch/request-workflow-flag.md)
- decided-by [ready-trigger-dispatch](../capability/rhiza/workflows/workflow-dispatch/ready-trigger-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
