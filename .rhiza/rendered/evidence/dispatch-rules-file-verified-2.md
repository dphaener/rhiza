# dispatch-rules-file-verified-2

evidence · `01M6W9KRWCN38EJVRJEFV60VQB` · actor [claude-code-c2844eca](../actor/claude-code-c2844eca.md)

## result

PASS

## statement

At 424d4a80 (attempt 2), ran go test -trimpath -count=1 -v ./internal/dispatch/...: TestDefaultsWhenNoFile (no file gives plan for every request kind and ship-change for every ready kind), TestReadyBugOnly (a file setting only [ready] bug changes only that pick; every other pick stays default) and TestParseErrors (unknown table [other]; unknown key [request] chore; non-string [ready] bug = 3; path-like a/b, x.toml and empty values; each error names the table and key) all PASS. The new case [ready] milestone = "x" fails with '[ready] unknown key "milestone": want one of feature, bug', which closes compliance finding f1 under decision dispatch-ready-no-milestone-key. Read the diff: Parse checks keys against TableKinds(table) (ReadyKinds = feature, bug for [ready]); ForReady falls back to ship-change for a kind [ready] cannot name; check.go's dispatchProblems loops over TableKinds. go test -trimpath ./cmd/rhiza -run 'TestScripts/dispatch' passed (ok, 125s), and go vet on internal/dispatch and internal/api is clean.

## Links

### Out

- verifies [dispatch-rules-file](../capability/rhiza/workflows/workflow-dispatch/dispatch-rules-file.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
