# roadmap-app-start-runs

capability · `01M6734CBEZPDTAJBANZVV5GXC` · parent [rhiza](../rhiza.md)

## acceptance

A test with a fake exec presses p then y on a ready change set and asserts one rex window new call whose command is rhiza run start <cs> --workflow plan with no --autonomous; s then y gives rhiza run start <cs> --autonomous --permission-mode auto; n calls nothing; o on a change set with a live run calls rex to focus that run's tab, and on one without says no live run; no test runs real Rex

## statement

On the selected change set, p starts a plan run and s starts a ship run, each in a new Rex tab running rhiza run start, the plan run interactive and the ship run autonomous with permission mode auto; the app asks y/n first, stays open, and says which tab it opened; o focuses the Rex tab of the change set's live run, or says there is none

## Links

### Out

- serves [m5-roadmap-actionable](../../intent/milestone-5/m5-roadmap-actionable.md)
- decided-by [app-runs-open-rex-tabs](../../decision/app-runs-open-rex-tabs.md)
- decided-by [app-plan-interactive-ship-autonomous](../../decision/app-plan-interactive-ship-autonomous.md)

### In

- verifies [roadmap-app-start-runs-verified](../../evidence/roadmap-app-start-runs-verified.md)
- verifies [roadmap-app-start-runs-qa-1](../../evidence/roadmap-app-start-runs-qa-1.md)
- verifies [roadmap-app-start-runs-verified-2](../../evidence/roadmap-app-start-runs-verified-2.md)
- verifies [pr-120-actionable-roadmap-app](../../evidence/pr-120-actionable-roadmap-app.md)

## History

- `01M673223RSVC166NJVVA13X75` The roadmap app is a cursor list of milestones and change sets that the owner acts on: open a change set's detail, plan it, ship it, accept it, answer its blocking questions, or jump to its run's Rex tab
