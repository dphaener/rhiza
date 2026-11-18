# retry-restarts-at-named-step-qa-1

evidence · `01M7QTPW21KMQBFH17JRZESSJR` · actor [claude-code-b8e7291f](../actor/claude-code-b8e7291f.md)

## result

PASS

## statement

QA with a scratch build: a run failing at a step with retry_from = prep was retried by rhiza run resume --retry at prep (dry run named prep), the retried journal entry recorded failed: check, prep's max_visits = 1 counted afresh, and the run succeeded

## Links

### Out

- verifies [retry-restarts-at-named-step](../capability/rhiza/workflows/workflow-runs/resume-failed-run/retry-restarts-at-named-step.md)

## History

- `01M7NVMMPK3A9SRVA9Q63N5SR9` A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox
