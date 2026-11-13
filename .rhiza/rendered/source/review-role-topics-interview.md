# review-role-topics-interview

source · `01M79AF2GX203GH9MK0Y5T0S8V`

## quotes

```json
[]
```

## raw_sha256

937cc0d5b76d2ba8310fed88a61c2ad4af55af41f14e08c94dab54905a97a6ec

## session

acb282c4-2fe3-4a89-a050-01a1d2a244a4

## summary

Autonomous plan run for change set 01M76XRE (Review role topics), serving m35-afk-prompts. No owner present. The planner read the change set intent, prompt-gap.md sections 1, 2.5-2.11, 5 and 6, internal/guide/text/review.md and its tests, internal/api/guide_install.go, internal/workflow/builtin/review.toml and ship-change.toml, Rhiza's .rhiza/workflows/docs-review.toml, ship-change.toml and docs.toml, the landed severity-contract leaf and the sibling Build role topics plan (01M76XR5, Decision build-role-topics-shape: shared blocks under internal/guide/text/shared/, rhiza-<topic> skills, one-line prompts, the old topic kept as an index). It raised one question bundling four choices (topic names and delivery, isolation of checks from build history, where the scope check's criteria go and whether dropping it is breaking, which ship-change fix steps use which topic). The arbiter answered with Decision review-role-topics-shape: nine topics (check-compliance, check-quality, check-qa, check-model, critic, fixer, ci-fix, gate-fix, conflict-fix) installed as rhiza-<topic> skills, a new shared severity block, review kept as an index, isolation in the topics' MUST NOT DO only, scope folded into check-compliance and dropped from review.toml and docs-review.toml without a breaking note, and the ship-change fix prompts (builtin and Rhiza's own ship-change.toml and docs.toml) shrunk to one line per topic. The planner then wrote leaves under review-guide: the severity block, the nine topics, install plus the review index, and the one-line prompts for review and for ship, and linked the change set depends-on Build role topics for the shared-block mechanism.

## title

Review role topics plan (autonomous)

## Links

### In

- derived-from [review-role-topics-shape](../decision/review-role-topics-shape.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
