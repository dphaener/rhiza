# roadmap-new-work-requests-qa-1

evidence · `01M6D4T4S5S3MF2VKVK9KNV1Q5` · actor [claude-code-47a2bc48](../actor/claude-code-47a2bc48.md)

## result

FAIL

## statement

QA at 9127ca1a: built go build -o $TMPDIR/rhiza-qa ./cmd/rhiza, ran `rhiza-qa tui app roadmap` in a Rex pane on a fresh scratch repo (rhiza init). Pressed n, typed 'Ship the thing', enter: the app no longer asks y/n and opens no Rex tab (no new window appeared), but the notice line reads 'rhiza: unknown operation "request"'. Pressed m, typed 'Big goal', enter: same error. Nothing is queued (no .git/rhiza/requests.json). `rhiza-qa request` is also an unknown command. TestNewWorkFilesPlanRequest passes on its fake runner. So with the real binary the new change set and new milestone keys are broken: they file nothing and no longer start a plan run.

## Links

### Out

- verifies [roadmap-new-work-requests](../capability/rhiza/orchestration/daemon/plan-requests/roadmap-new-work-requests.md)

## History

- `01M6CRZQYG5PJV7X2B234DJCV4` Feature, bug and milestone requests queue in the clone and the daemon turns each into a plan run
