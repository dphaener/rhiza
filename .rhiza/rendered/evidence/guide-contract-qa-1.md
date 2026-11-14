# guide-contract-qa-1

evidence · `01M7AWMPH74RY0SQW0EPB45EM1` · actor [claude-code-1707fde5](../actor/claude-code-1707fde5.md)

## result

PASS

## statement

QA: built rhiza guide contract prints obligations with gates, invoke placeholders, RHIZA_ variables incl. BRIEF/SLICE/HANDBACK and the three files (40 lines); rhiza guide <role> names rhiza guide contract for advisor, arbiter, builder, verifier, the four checks, ci-fix, conflict-fix, gate-fix, critic, fixer, plan-critic, plan-reviewer, plan-reviser, planner; go test ./internal/guide ok.

## Links

### Out

- verifies [guide-contract](../capability/rhiza/guide/agent-guide/guide-contract.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
