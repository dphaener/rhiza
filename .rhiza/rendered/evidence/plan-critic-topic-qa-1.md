# plan-critic-topic-qa-1

evidence · `01M79QSPAPX3FB1KX10TBSXAWV` · actor [claude-code-483e02c2](../actor/claude-code-483e02c2.md)

## result

PASS

## statement

Built rhiza from 3d6baba6c (go build -trimpath) and ran it in a hermetic scratch root and repo (env -i, own HOME, project added, a held 'hold' wait-step run giving RHIZA_RUN): `rhiza guide plan-critic` prints the six sections in order, Read first names `rhiza run status --json`, the verdict heredoc with confirmed, downgraded+severity and dropped, never add/raise rules and severities, command-budget, never and session-done blocks expanded, no findings heredoc or rhiza apply; its verdict heredoc ran as written (`rhiza run verdict - --json`) with f1 downgraded to minor and f2 confirmed, and run status showed 'downgraded from major to minor'.

## Links

### Out

- verifies [plan-critic-topic](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-critic-topic.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
