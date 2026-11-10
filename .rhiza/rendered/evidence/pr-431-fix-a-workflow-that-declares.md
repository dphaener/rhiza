# pr-431-fix-a-workflow-that-declares

evidence · `01M73JPV60J912194NH0DAD1D5` · actor [dphaener](../actor/dphaener.md)

## head

94a9b51b2363adbfa92ffac8ab8b7655babe6230

## pr

https://github.com/dphaener/rhiza/pull/431

## result

11 checks passed: lint SUCCESS, go-install SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #431 (https://github.com/dphaener/rhiza/pull/431) passed CI at 94a9b51

## status

verified

## Links

### Out

- verifies [signoff-owner-beats-autonomous](../capability/rhiza/orchestration/accept-policy/signoff-owner-beats-autonomous.md)

## History

- `01M73FYE2PAG9E034MPBTTJ12D` Fix A workflow that declares 'autonomy signoff owner' is signed off autonomously anyway. The repo's docs workflow (#425) declares signoff = owner (rhiza workflow describe docs: 'autonomy signoff owner'), yet run 01M73CQ4D71CW038ZBNSQB1SJF (Concepts polish, PR #426) logged 'step signoff: signed off under rule autonomous' and went on to land; run 01M73CQ3TCG0TVWWR8ZV62VAE8 (README, PR #427) is headed the same way. Settings: rhiza.autonomy.signoff = auto (default, project); the daemon runs --autonomous. The milestone's Decisions require the owner to sign off each page, so a workflow's declared owner stage must win over the default/--autonomous (or the docs workflow is not honoured). Held by converting PRs #426 and #427 to draft so ship merge cannot merge them; the owner signs off on the PRs by hand.
