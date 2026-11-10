# signoff-owner-plan

source · `01M73G39FKE09Q8PVGNRNYJ42N`

## quotes

```json
[]
```

## raw_sha256

fc0b242ef24b487d15b02b752045703060c2bd715bf23d9d62fcafc56cbb3238

## session

6f082b02-8d9d-4f0f-b63d-01a1c6ff9ddf

## summary

Autonomous plan run (no owner present) for the bug where the docs workflow declares [autonomy] signoff = owner, yet run 01M73CQ4D71CW038ZBNSQB1SJF signed off under rule autonomous. Findings: the run's journal recorded the signoff stage as owner from workflow docs, but internal/run/signoff.go signs off whenever State.Autonomous is set and never reads the recorded stage autonomy; the only guard is daemonStart dropping --autonomous, which other start paths bypass. Planned one leaf: the signoff step itself waits for the owner when the run's recorded signoff stage is a workflow-declared owner, even in an autonomous run, so every start path honours it. Served intent follows the precedent of recent run fixes (m4-runs-hold-up). Setting-sourced signoff = owner keeps today's handling (the daemon starts without --autonomous); no owner words were quoted because there was no owner.

## title

Plan: workflow-declared signoff owner beats --autonomous

## Links

### In

- derived-from [signoff-step-reads-workflow-stage](../decision/signoff-step-reads-workflow-stage.md)
- derived-from [signoff-owner-beats-autonomous](../capability/rhiza/orchestration/accept-policy/signoff-owner-beats-autonomous.md)

## History

- `01M73FYE2PAG9E034MPBTTJ12D` Fix A workflow that declares 'autonomy signoff owner' is signed off autonomously anyway. The repo's docs workflow (#425) declares signoff = owner (rhiza workflow describe docs: 'autonomy signoff owner'), yet run 01M73CQ4D71CW038ZBNSQB1SJF (Concepts polish, PR #426) logged 'step signoff: signed off under rule autonomous' and went on to land; run 01M73CQ3TCG0TVWWR8ZV62VAE8 (README, PR #427) is headed the same way. Settings: rhiza.autonomy.signoff = auto (default, project); the daemon runs --autonomous. The milestone's Decisions require the owner to sign off each page, so a workflow's declared owner stage must win over the default/--autonomous (or the docs workflow is not honoured). Held by converting PRs #426 and #427 to draft so ship merge cannot merge them; the owner signs off on the PRs by hand.
