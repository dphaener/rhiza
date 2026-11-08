# request-exposed-workflows-qa-1

evidence · `01M6W9QDQE59RDA4KWRKK3TVW2` · actor [claude-code-afe51521](../actor/claude-code-afe51521.md)

## result

PASS

## statement

QA at bbfe87e3 with a built binary (go build -trimpath) in a scratch Rhiza root and project, real git config, no fakes: With plan auto + bypassPermissions (accept/ship auto) health and workflow check flag plan (trigger request; interview, review, critic, revise, arbiter) and ship-change (trigger ready; build/build builder and verifier, every review check, critic, fixer, verifier, ci-fix, arbiter). With plan owner + bypass, request --ref notes plan only (per decision outside-input-plan-owner-not-intake). With acceptEdits nothing is flagged anywhere.

## Links

### Out

- verifies [request-exposed-workflows](../capability/rhiza/workflows/autonomous-runs/request-exposed-workflows.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
