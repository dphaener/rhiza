# workflow-dispatch-merge-main

decision · `01M6WDA5DTCXV3R3FJDEB3Y102`

## reasoning

Rebase and force-push are not allowed on run branches, so main was merged. Main had landed request --ref (format 4 open body ref), risk warnings, scheduling as policy and the workflows-are-the-rules decision; this branch adds request kind and workflow to the same format-4 open body and the same request types, so each conflict was additive and both fields fit side by side in one unreleased format 4 (both sides bumped FormatVersion to 4 identically). daemonStartPlan keeps main's planRunArgs split (so --ref stays one word) with this branch's dispatched workflow and pick. The store test takes main's FormatVersion+1 form. Risk warnings exist to name every workflow request text reaches unseen; with dispatch that is whatever [request] and [ready] map, so the exposure list reads the rules (defaults reproduce main's behaviour, and a new test covers a bug rule). The merged guide text was 1009 tokens against an 800 cap, so it was rewritten tighter with every asserted phrase kept.

## statement

Merge origin/main into cs/workflow-dispatch instead of rebasing, keeping both sides: the open body carries ref, request_kind and workflow, all format 4; daemon.Request, RequestRequest and RequestInfo carry both Workflow and Ref; the plan run's start args take the dispatched workflow and --pick plus --ref=<ref>; check runs dispatch then policy; describe-golden is regenerated; the workflows guide topic is tightened to keep both the Dispatch and Outside input sections under its 800-token cap;, the risk-warning exposure list follows the dispatch rules for each kind instead of naming plan and ship-change alone, and plan-requests expects --pick default in the plan run args.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [rebase-conflict-workflow-dispatch](../question/rebase-conflict-workflow-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
