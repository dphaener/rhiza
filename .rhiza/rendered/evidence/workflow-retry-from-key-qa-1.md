# workflow-retry-from-key-qa-1

evidence · `01M7QTPW01DTJGYFZ5F34726Z0` · actor [claude-code-b8e7291f](../actor/claude-code-b8e7291f.md)

## result

PASS

## statement

QA with a scratch build: rhiza workflow describe shows retry_from (prefixed plan/interview when plan is included by a workflow step; built-in plan and plan-milestone review-loop steps show retry_from -> interview); rhiza workflow check refuses a retry_from naming no step, done/failed, or itself, naming step and key

## Links

### Out

- verifies [workflow-retry-from-key](../capability/ship-through-rhiza/workflow-step/workflow-retry-from-key.md)

## History

- `01M7NVMMPK3A9SRVA9Q63N5SR9` A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox
