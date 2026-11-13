# guide-contract-qa-2

evidence · `01M7B24D09MHMJH58A2SEEJ9WC` · actor [claude-code-a1cb0102](../actor/claude-code-a1cb0102.md)

## result

PASS

## statement

QA cycle 2: scratch-built rhiza guide contract prints obligations with gates, every invoke placeholder incl. {leaf}/{attempt} on check/critic and {attempt} on fixer, the RHIZA_ variables incl. BRIEF/SLICE/HANDBACK and the three files (exit 0); every role topic (advisor, arbiter, builder, check-*, critic, fixer, ci/gate/conflict-fix, planner, plan-critic/reviewer/reviser, verifier) names rhiza guide contract; unknown topic exits 4 with a next line. go test ./internal/guide -run TestContract ok.

## Links

### Out

- verifies [guide-contract](../capability/rhiza/guide/agent-guide/guide-contract.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
