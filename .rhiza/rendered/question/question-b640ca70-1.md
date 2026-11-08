# question-b640ca70-1

question · `01M6VYENMSWMKSSEWJ5NEN2FKM`

## session

b640ca70-18c2-43e6-a7ec-01a1b7e43901

## statement

Which milestone should Workflow dispatch (01M6VY7G) serve? It has no serves link, and the daemon builds only change sets that serve a milestone. Milestone 29 (several projects under one root) is the open one but this work doesn't fit it. Milestone 25 is the backlog milestone ('every change set planned from the owner's requests is built and shipped by the daemon on its own, so the owner can start on new things, bugs and polish'), and the workflow-authoring change set 01M6HN24 served it. / Where do the repo's dispatch rules (trigger × request kind → workflow) live? The intent says rhiza.toml, but in the code rhiza.toml is the root's machine-wide config (internal/rootcfg: its projects and [daemon] settings). It is not committed in any repo, and one root serves several projects. / What does `rhiza request --workflow <name>` override? The intent says the change set records any workflow its request chose 'so the ship pick follows it'. / Does a manual `rhiza run start <cs>` with no --workflow follow the same dispatch rules as the daemon? Today it defaults to ship-change, the same as the daemon.

## Links

### In

- answers [workflow-dispatch-shape](../decision/workflow-dispatch-shape.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
