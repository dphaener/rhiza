# question-b640ca70-2

question · `01M6VYNG696SP9M2NGXG5CEC3R`

## session

b640ca70-18c2-43e6-a7ec-01a1b7e43901

## statement

What shape does .rhiza/dispatch.toml take? There are two triggers (request queued, change set ready) and three request kinds (feature, bug, milestone). / How are bad workflow names handled: a rule in .rhiza/dispatch.toml, or a `rhiza request --workflow` value, naming a workflow no layer (personal, repo or built-in) resolves? / Which kind does a change set get when no request opened it, e.g. one opened by a milestone plan or by `rhiza change new`? The ready trigger maps kind to workflow. / Where does the owner see what dispatch picked?

## Links

### In

- answers [workflow-dispatch-rules-file](../decision/workflow-dispatch-rules-file.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
