# ready-trigger-dispatch-verified-2

evidence · `01M6WAPD086Z0N97VF304Z8GPK` · actor [claude-code-54cca4cb](../actor/claude-code-54cca4cb.md)

## result

PASS

## statement

Verifier at 77a481ff (attempt 2, fixing quality f1 / model f2). I read the diff: readyPick's dispatch.Load error now flows out of startWorkflowOfRef as an error, so runRunStart fails instead of falling back to ship-change. daemonUnresolvedReady now skips such a candidate as workflow-unresolved, with the load error in its detail, instead of keeping it for a default start. An explicit --workflow still never reads the file. I ran go test -trimpath ./internal/api with -run TestReadyWorkflowOrder (recorded over rule over ship-change), TestStartWorkflowExplicitWins (explicit over all), TestRunStartExplicitWorkflowOverRule, TestDaemonSkipsReadyWithUnresolvedWorkflow, TestDaemonSkipsReadyWithUnloadableDispatch (new), TestDaemonSkipsRequestWithUnresolvedWorkflow, TestDaemonStartsRequestWithDispatchedWorkflow and TestRunStartWorkflowRefusals: all PASS. The package's TestMain guard then reported .rhiza-root/usage.live.json changed. That file is the live root daemon's: it changed at 00:59:10 with no test running, so the report is environmental. The whole-package run hit Go's 600s timeout at VM load ~11, and CI's make check covers it. I ran go test -trimpath ./internal/dispatch/...: ok. I ran go test -trimpath ./cmd/rhiza -run TestScripts/(run-start-ready-dispatch|dispatch-check): both PASS. run-start-ready-dispatch shows run start --dry-run on a bug change set: with no dispatch file it falls to ship-change; with [ready] bug = "bugfix" it prints Workflow: bugfix; with --workflow other it prints other; with an unloadable dispatch.toml it fails naming dispatch.toml, not ship-change.

## Links

### Out

- verifies [ready-trigger-dispatch](../capability/rhiza/workflows/workflow-dispatch/ready-trigger-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
