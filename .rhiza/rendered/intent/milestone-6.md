# milestone-6

intent · `01M68551X2HM0EB0AY469DXYY9`

## milestone

6

## statement

Rhiza runs unattended overnight: the daemon builds, reviews, lands and merges accepted change sets, and plans its own PRs, without the orchestrator committing, rebasing, fixing ship records or closing windows by hand

## Links

### Out

- decided-by [m6-overnight-batch](../decision/m6-overnight-batch.md)

### In

- serves change set `01M684XYB96R52XXCKMV80X5T1` (When a run's PR branch is updated by someone other than the run, such as the orchestrator fixing a conflict, the run's wait-pr, rebase and rebase-merge steps follow the new head instead of stalling until wait-pr's timeout or failing)
- serves change set `01M66M3S5PZD9PXDP0X3N1ZTXH` (A run's agent-step, arbiter and advisor sessions open as tiled panes of the run's window, like its leaf, check, critic and fixer sessions, so autonomous runs with many arbiter calls add no window each)
- serves change set `01M5SWTV1D5DNN6E2B7HB1XA9Q` (genesis link checks endpoint kinds when the link is written, not only in genesis check.)
- serves change set `01M625RCFY4RB9HC0B69Y08VFV` (A plan run ends with its plan on a pushed branch and an open PR, as ship-change does, so no one renders, commits or opens the plan's PR by hand.)
- serves change set `01M625W0AJZ6AMCCWGFXQW19FN` (A plan run that ends with confirmed review findings or parked questions still open does not report success: it stops for the owner with them listed, so nothing has to be fixed by hand after a succeeded run.)
- serves change set `01M6580HBAKHN7XTR7XQV3XR0Q` (A failed review leaf fixer, or its verifier, is retried within its cycle as the change-set fixer is, so only a leaf fix failing every attempt spends the cycle)
- serves change set `01M6489HCR688XKH99YV9E2VGY` (Remove the all-apps picker, as decision rex-picker-removed says: rhiza tui launch and the rhiza_launch action written by rhiza rex install go, and the capabilities that describe the picker are retired or revised, so each app opens only from its own key in the rhiza key mode or its own action in Rex's command palette.)

## History

- `01M684XWDT3EKJD4JPXT0JVK2M` Plan milestone 6: Rhiza runs unattended overnight, with the run machinery fixes that needed the orchestrator's hand in milestones 4 and 5
