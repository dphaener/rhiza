# change-set-records-dispatch-verified

evidence · `01M6W3GCT9DXNDTA2P9D4BJ9Z2` · actor [claude-code-02ae7f3a](../actor/claude-code-02ae7f3a.md)

## result

PASS

## statement

At b11bed21, as verifier: go test -trimpath ./cmd/rhiza -run 'TestScripts/(change-request-kind|describe-golden)$' passes. change-request-kind.txtar shows that a change set from rhiza change new prints request_kind feature and no workflow in change show --json (and '- Request kind: feature' with no Workflow line in human output). It also shows that a queued bug request, whose queue entry in .git/rhiza/requests.json the script edits to carry workflow bugfix, opened by run start --new --request, prints request_kind bug and workflow bugfix in change show --json and human output. fold-ok passes. go test -trimpath for internal/changeset, entry, store and daemon passes (new format-4 encode/decode and Validate tests included). internal/api passes, except that TestReadBudget missed its roadmap/ready what-if 50ms budget in a run beside another agent's api.test at load ~10. Re-run alone, TestReadBudget passed at both b11bed21 and its parent with matching medians (roadmap what-if 30ms vs 27ms, ready 28ms vs 38ms), so the miss was load, not this diff. The diff adds request_kind/workflow to the open body behind entry format 4 (FormatVersion 3->4; format 3 rejects the fields), threads the request's kind and workflow through newPlanChange, and adds both to ChangeShowResponse in the fold and indexed paths. Changed Go files are gofmt-clean.

## Links

### Out

- verifies [change-set-records-dispatch](../capability/rhiza/workflows/workflow-dispatch/change-set-records-dispatch.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
