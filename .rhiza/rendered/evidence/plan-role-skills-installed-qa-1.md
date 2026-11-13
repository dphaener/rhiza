# plan-role-skills-installed-qa-1

evidence · `01M79QSQ9FFDFQ8FANJ246NX3J` · actor [claude-code-483e02c2](../actor/claude-code-483e02c2.md)

## result

PASS

## statement

Built rhiza from 3d6baba6c (go build -trimpath) and ran it in a hermetic scratch root and repo (env -i, own HOME, project added, a held 'hold' wait-step run giving RHIZA_RUN): `rhiza guide ls` lists planner, plan-reviewer, plan-critic, plan-reviser; `rhiza guide install --agent claude --json` wrote skills/rhiza-planner, rhiza-plan-reviewer, rhiza-plan-critic, rhiza-plan-reviser SKILL.md beside rhiza-plan (rhiza-plan-critic frontmatter carries the topic's description), and a second install listed all of them unchanged, written empty; `rhiza guide plan` and `plan-review` print short indexes naming the role topics with /rhiza-<topic> and `rhiza guide <topic>` routes.

## Links

### Out

- verifies [plan-role-skills-installed](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-role-skills-installed.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
