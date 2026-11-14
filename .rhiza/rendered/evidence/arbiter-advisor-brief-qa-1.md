# arbiter-advisor-brief-qa-1

evidence · `01M7AWMPC3GK513DY3DC21TWDN` · actor [claude-code-1707fde5](../actor/claude-code-1707fde5.md)

## result

PASS

## statement

QA: advisor sessions started by real failing runs got RHIZA_BRIEF under .git/rhiza/runs/<run>/advisor/<session id>/brief.md, one directory per advisor (two advisors in one run got two), prompt unchanged; arbiter covered by go test ./internal/api -run TestArbiterBrief (ok).

## Links

### Out

- verifies [arbiter-advisor-brief](../capability/rhiza/workflows/workflow-runs/arbiter-advisor-brief.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
