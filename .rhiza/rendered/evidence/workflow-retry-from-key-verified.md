# workflow-retry-from-key-verified

evidence · `01M7QQMPXETJGX7XAR8WCNBW63` · actor [claude-code-45f8e431](../actor/claude-code-45f8e431.md)

## result

PASS

## statement

At 6f8695939 (capped scopes): go test -trimpath -p 2 ./internal/workflow/... ok (TestRetryFrom: parses retry_from; refuses itself, a missing step, done, the arbiter, and the key on arbiter/workflow steps, naming the step and retry_from, in Parse and Validate; TestIncludeRetryFrom: included retry_from comes out plan/interview, outer p/write and p resolve, skipped target follows the skip, p/nosuch and inner nosuch refused); cmd/rhiza TestScripts/workflow-describe and describe-golden pass (describe shows 'retry_from -> plan' and prefixed retry/plan, --json has retry_from); internal/api -run Workflow|Describe|Docs and internal/guide -run WorkflowKeys|Keys ok; make self-check exit 0 in a clean detached worktree at the commit (in the run worktree it fails only on the verifier's own untracked actor record). workflow-keys.md and docs/reference/workflow-file.md document the key. Diff touches no built-in workflow, no other step key, and no engine retry code.

## Links

### Out

- verifies [workflow-retry-from-key](../capability/ship-through-rhiza/workflow-step/workflow-retry-from-key.md)

## History

- `01M7NVMMPK3A9SRVA9Q63N5SR9` A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox
