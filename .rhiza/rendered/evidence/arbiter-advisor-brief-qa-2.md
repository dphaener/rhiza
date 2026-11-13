# arbiter-advisor-brief-qa-2

evidence · `01M7B24CV5ZBP1NZYXSXFS2ANP` · actor [claude-code-a1cb0102](../actor/claude-code-a1cb0102.md)

## result

PASS

## statement

QA cycle 2: advisor sessions started by real failing leaves and checks runs (scratch build) got RHIZA_BRIEF/RHIZA_SLICE under .git/rhiza/runs/<run>/advisor/<session id>/ with brief.md, prompt the usual advisor text; arbiter covered by go test -trimpath ./internal/api -run 'TestArbiterBrief|TestAdvisorBrief' ok.

## Links

### Out

- verifies [arbiter-advisor-brief](../capability/rhiza/workflows/workflow-runs/arbiter-advisor-brief.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
