# plan-prompts-one-line-verified

evidence · `01M79Q9VYZM95DFW4VDM2CXN6M` · actor [claude-code-bfbb34bc](../actor/claude-code-bfbb34bc.md)

## result

PASS

## statement

At bc8cc09ca: go test -trimpath ./internal/workflow/... exit 0 (TestPlanBuiltin and TestPlanPromptsOneLine assert plan.toml and plan-milestone.toml interview/review/critic/revise prompts equal the four one-liners, contain no 'session done', and prepush-fix keeps its own text); go test -trimpath ./cmd/rhiza -run 'TestScripts/(plan-|bugfix-workflow-autonomous|arbiter-agent-tabs|status|guide|describe-golden)' exit 0 with fake agents matching /rhiza-planner, /rhiza-plan-reviewer, /rhiza-plan-critic, /rhiza-plan-reviser; describe-golden unchanged (prompts not in describe output). Diff reviewed: only the four prompts in each builtin, the tests, the scripts and the builder's actor record changed; prepush-fix untouched; session-done rules present in planner/plan-reviewer/plan-critic/plan-reviser topics; no TODOs. make self-check: every check passes except render of this verifier's own untracked actor record (claude-code-bfbb34bc, not part of the commit); rhiza render on a clean detached checkout of HEAD writes no diff.

## Links

### Out

- verifies [plan-prompts-one-line](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-prompts-one-line.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
