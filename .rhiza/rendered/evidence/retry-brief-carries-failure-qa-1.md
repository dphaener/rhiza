# retry-brief-carries-failure-qa-1

evidence · `01M7QTPW477KJXHMG6NH1M6VYB` · actor [claude-code-b8e7291f](../actor/claude-code-b8e7291f.md)

## result

PASS

## statement

QA with a scratch build and fake rex/claude: after a retry restarting an agent step, its visit-2 brief.md held a Previous attempt section (failed step, outcome, failing command, output tail, latest round's findings); the pre-retry and never-retried briefs had none. Plan review findings show without a severity (minor finding)

## Links

### Out

- verifies [retry-brief-carries-failure](../capability/rhiza/workflows/workflow-runs/resume-failed-run/retry-brief-carries-failure.md)

## History

- `01M7NVMMPK3A9SRVA9Q63N5SR9` A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox
