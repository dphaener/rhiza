# signoff-owner-beats-autonomous

capability · `01M73G3PK7RGKSX1K0NEM0MFB2` · parent [accept-policy](../accept-policy.md)

## acceptance

go test -trimpath ./internal/run/... covers the signoff step with Autonomous set and the recorded signoff stage owner from a workflow (waits, no autonomous signoff entry), auto from a workflow and from the setting (signs off under rule autonomous), and a journal without the record; a new testscript cmd/rhiza/testdata/script/run-signoff-workflow-owner.txtar runs run start --autonomous on a workflow declaring [autonomy] signoff = owner with fake agents, sees the run stop at signoff with no rule-autonomous sign-off in rhiza run status --json, then rhiza run approve lets it finish; review-signoff-autonomous, workflow-stage-autonomy and run-autonomy-source scripts pass unchanged

## statement

A signoff step waits for the owner whenever its run's workflow declares the signoff stage owner, even when the run is autonomous: it journals no sign-off of rule autonomous, stops with the owner's sign-off wait naming the workflow's declaration, and passes only on rhiza run approve (or goes back on request-changes); a run whose workflow declares signoff auto or nothing still signs off under rule autonomous; the declaration is read from the run's recorded stage autonomy, from the top-level workflow, so a signoff step in a nested workflow obeys it too, and a journal that predates that record falls back to the started workflow definition's declaration

## Links

### Out

- decided-by [signoff-step-reads-workflow-stage](../../../../decision/signoff-step-reads-workflow-stage.md)
- derived-from [signoff-owner-plan](../../../../source/signoff-owner-plan.md)

### In

- verifies [signoff-owner-beats-autonomous-verified](../../../../evidence/signoff-owner-beats-autonomous-verified.md)
- verifies [signoff-owner-beats-autonomous-qa-1](../../../../evidence/signoff-owner-beats-autonomous-qa-1.md)
- verifies [pr-431-fix-a-workflow-that-declares](../../../../evidence/pr-431-fix-a-workflow-that-declares.md)

## History

- `01M73FYE2PAG9E034MPBTTJ12D` Fix A workflow that declares 'autonomy signoff owner' is signed off autonomously anyway. The repo's docs workflow (#425) declares signoff = owner (rhiza workflow describe docs: 'autonomy signoff owner'), yet run 01M73CQ4D71CW038ZBNSQB1SJF (Concepts polish, PR #426) logged 'step signoff: signed off under rule autonomous' and went on to land; run 01M73CQ3TCG0TVWWR8ZV62VAE8 (README, PR #427) is headed the same way. Settings: rhiza.autonomy.signoff = auto (default, project); the daemon runs --autonomous. The milestone's Decisions require the owner to sign off each page, so a workflow's declared owner stage must win over the default/--autonomous (or the docs workflow is not honoured). Held by converting PRs #426 and #427 to draft so ship merge cannot merge them; the owner signs off on the PRs by hand.
