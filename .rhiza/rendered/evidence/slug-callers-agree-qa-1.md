# slug-callers-agree-qa-1

evidence · `01M6THM2Z5B88DR0ZN3EG9R29D` · actor [claude-code-1bbb24c8](../actor/claude-code-1bbb24c8.md)

## result

PASS

## statement

With the built binary and fake rex/claude: run start --dry-run on the later twin predicted worktree repo.fix-the-usage-governor-s-01m6thgetkzw; session start on both twins created distinct worktrees and branches cs/fix-the-usage-governor-s and cs/fix-the-usage-governor-s-01m6thgetkzw; ship check run in the later twin's worktree targeted cs/fix-the-usage-governor-s-01m6thgetkzw and recorded that branch in ship state

## Links

### Out

- verifies [slug-callers-agree](../capability/rhiza/runtime/rex-sessions/slug-callers-agree.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
