# request-workflow-flag-verified

evidence · `01M6W3SHWY17E5XG68M32XG8TH` · actor [claude-code-045f7912](../actor/claude-code-045f7912.md)

## result

PASS

## statement

At df9cc6f4, go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(request-workflow-flag|describe-golden)$' passed. request-workflow-flag.txtar shows: request --workflow bugfix (a repo workflow) queued, and request ls --json with "workflow": "bugfix"; request --workflow nonesuch failing with 'no such workflow' and the queue unchanged ('No plan requests'). describe-golden now has the new workflow field on RequestRequest and RequestInfo. The script's `! exec` does not pin the exit code, so I checked by hand with a scratch build in a registered scratch project: --workflow nonesuch exits 2 with a usage error and queues nothing; --workflow plan (built in) queues, and ls shows [workflow plan]. In the diff, internal/api/request.go checks the name through lookupWorkflow (repo, personal and built-in layers) before queueing and stores it on daemon.Request.Workflow; run start already passes req.Workflow into newPlanChange (internal/api/run.go:569).

## Links

### Out

- verifies [request-workflow-flag](../capability/rhiza/workflows/workflow-dispatch/request-workflow-flag.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
