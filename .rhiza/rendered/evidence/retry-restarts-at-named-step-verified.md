# retry-restarts-at-named-step-verified

evidence · `01M7QR29H86PR5PTVRXCBTX84A` · actor [claude-code-d85cf57d](../actor/claude-code-d85cf57d.md)

## result

PASS

## statement

At d239b0421 re-ran the Checks block under the VM's capped scope, all exit 0: go test -trimpath -p 2 ./internal/run/... ./internal/workflow/...; cmd/rhiza TestScripts/run-resume-retry (also -count=1 -v: PASS) and TestScripts/(plan-review-workflow|source-raw-session|workflow-describe); ./internal/guide; ./internal/api -run 'Retry|Resume'; rhiza ship generate (0 changed) and rhiza render (only the verifier's own actor record). Read the diff against the acceptance: TestRetryRestartsAtRetryFrom fails a run at check (retry_from=prep, over max_visits 2), and after the retry the retried entry has step prep and failed check, prep is entered at visit 2 with CapVisits(prep)=1 despite max_visits 1 (counted afresh), and the run succeeds; TestRetryPlanBugCheckRestartsAtInterview parses a plan-bug-shaped workflow including the built-in plan, fails at plan/check over its cap 3 and sees the retry entered at plan/interview:2 with Retries[0]={plan/interview, failed plan/check}; RetryStep only consults retry_from when set and present in def.Steps, otherwise the existing waitFor/self logic is untouched and RetryStep's existing tests pass unchanged; TestBuiltinPlanReviewLoopRetriesFromInterview asserts plan and plan-milestone set retry_from=interview on review, has-findings, critic, fresh, check, revise, rulings, apply-rulings and on no other step. The txtar adds a CLI case (resume --retry --dry-run says 'would be retried at step prep', status shows 'prep (visit 2, retried)'). must_not respected: the toml diffs only add retry_from lines (no on.* or max_visits change), no brief, watchdog or ship-change edits; no stubs or TODOs.

## Links

### Out

- verifies [retry-restarts-at-named-step](../capability/rhiza/workflows/workflow-runs/resume-failed-run/retry-restarts-at-named-step.md)

## History

- `01M7NVMMPK3A9SRVA9Q63N5SR9` A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox
