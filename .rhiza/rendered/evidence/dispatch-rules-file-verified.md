# dispatch-rules-file-verified

evidence · `01M6W13J820RMTGSXZJ37CQQB9` · actor [claude-code-83fc5c5e](../actor/claude-code-83fc5c5e.md)

## result

PASS

## statement

Verifier at 15cf9c5a ran go test -trimpath -v ./internal/dispatch/...: TestDefaultsWhenNoFile (no file gives plan for every request kind and ship-change for every ready kind), TestReadyBugOnly ([ready] bug = "bugfix" changes only that pick; other ready kinds and all request kinds keep defaults) and TestParseErrors (unknown table [other], unknown key chore in [request], non-string bug = 3, path-like "a/b", "x.toml" and empty values each fail with an error naming the table and key) all PASS. Read internal/dispatch/dispatch.go: Parse rejects tables other than [request]/[ready], keys other than feature/bug/milestone, non-strings and values with path separators, whitespace, .toml suffix, . or ..; missing file/table/key fall back to Defaults. gofmt -l and go vet clean; no new module dependency (BurntSushi/toml already in go.mod).

## Links

### Out

- verifies [dispatch-rules-file](../capability/rhiza/workflows/workflow-dispatch/dispatch-rules-file.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
