# signoff-step-reads-workflow-stage

decision · `01M73G3PG7CCXGNJ9F2P1512ZD`

## reasoning

Run 01M73CQ4D71CW038ZBNSQB1SJF journalled signoff owner (source workflow docs) and still signed off under rule autonomous: internal/run/signoff.go checks only State.Autonomous. Guarding at the step covers every start path (daemon dispatch, queued run start, run start --autonomous, nested review workflows), which is what the milestone's owner-signs-each-page Decisions need. A signoff stage that is owner only by the rhiza.autonomy.signoff setting keeps today's behaviour, out of this bug's scope.

## statement

The signoff step itself honours the workflow's declared signoff stage: when the run recorded signoff = owner from its workflow, it waits for the owner even in an autonomous run, rather than relying on the daemon dropping --autonomous at start

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [signoff-owner-plan](../source/signoff-owner-plan.md)

### In

- decided-by [signoff-owner-beats-autonomous](../capability/rhiza/orchestration/accept-policy/signoff-owner-beats-autonomous.md)

## History

- `01M73FYE2PAG9E034MPBTTJ12D` Fix A workflow that declares 'autonomy signoff owner' is signed off autonomously anyway. The repo's docs workflow (#425) declares signoff = owner (rhiza workflow describe docs: 'autonomy signoff owner'), yet run 01M73CQ4D71CW038ZBNSQB1SJF (Concepts polish, PR #426) logged 'step signoff: signed off under rule autonomous' and went on to land; run 01M73CQ3TCG0TVWWR8ZV62VAE8 (README, PR #427) is headed the same way. Settings: rhiza.autonomy.signoff = auto (default, project); the daemon runs --autonomous. The milestone's Decisions require the owner to sign off each page, so a workflow's declared owner stage must win over the default/--autonomous (or the docs workflow is not honoured). Held by converting PRs #426 and #427 to draft so ship merge cannot merge them; the owner signs off on the PRs by hand.
