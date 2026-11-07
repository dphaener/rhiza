# schedule-policy-file-verified

evidence · `01M6W0ZEPGDAP28J5M854XC8YZ` · actor [claude-code-f5af258d](../actor/claude-code-f5af258d.md)

## result

PASS

## statement

At 6986e9a2 I read the diff (internal/policy, internal/check/policy.go, policy-check.txtar; repo.json untouched, toml already in go.mod) and ran go test -trimpath -count=1 ./internal/policy ./internal/check: TestRead passes its 11 cases (missing file, empty file, missing key, other table, after, before, skip, unknown value "sometimes" naming .rhiza/policy.toml and the value, empty value, wrong type, malformed TOML) and check's TestRunOrder places the policy rule before append-only. go test -trimpath ./cmd/rhiza -run TestScripts/policy-check passes: rhiza check with backlog = "sometimes" prints 'fail policy .rhiza/policy.toml: schedule.backlog "sometimes" is not one of "after", "before" or "skip"' and fails, malformed TOML fails naming the file, and backlog = "skip" (and before, and a missing key) passes with ok.

## Links

### Out

- verifies [schedule-policy-file](../capability/rhiza/orchestration/daemon/schedule-policy-file.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
