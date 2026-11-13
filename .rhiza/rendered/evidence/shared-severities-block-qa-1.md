# shared-severities-block-qa-1

evidence · `01M79QSQP6NPKZRNT25EKJJFD5` · actor [claude-code-483e02c2](../actor/claude-code-483e02c2.md)

## result

PASS

## statement

Built rhiza from 3d6baba6c (go build -trimpath) and ran it in a hermetic scratch root and repo (env -i, own HOME, project added, a held 'hold' wait-step run giving RHIZA_RUN): `rhiza guide plan-reviewer`, `plan-critic` and `review` each print the expanded Severities section; the plan-review mapping names missing-acceptance, contradiction, outside-intent, uncovered-intent and oversized-leaf major, matching run findings' severity output for the new kinds.

## Links

### Out

- verifies [shared-severities-block](../capability/rhiza/workflows/built-in-workflows/plan-workflow/shared-severities-block.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
