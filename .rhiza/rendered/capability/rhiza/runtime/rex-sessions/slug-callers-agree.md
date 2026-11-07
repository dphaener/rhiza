# slug-callers-agree

capability · `01M6TBRVZXT2PN5R0C1X1JC5X1` · parent [rex-sessions](../rex-sessions.md)

## acceptance

A testscript opens worktrees for two change sets with the same title and gets two distinct branches cs/<slug> and cs/<slug>-<first 12 characters of the later change set's ID, lowercased> (never its display short ID) and two distinct worktrees; internal/api tests show the ship target and abandon cleanup of the later change set name its suffixed branch; no call to session.Slug(title) remains outside the slug function

## statement

Every place that names a change set's branch or worktree (change worktree, run start and dry run, run leaves and checks sessions, ship target, abandon cleanup) uses the unique slug, never session.Slug(title) directly

## Links

### Out

- decided-by [slug-unique-by-earliest-id](../../../../decision/slug-unique-by-earliest-id.md)

### In

- verifies [slug-callers-agree-verified](../../../../evidence/slug-callers-agree-verified.md)
- verifies [slug-callers-agree-qa-1](../../../../evidence/slug-callers-agree-qa-1.md)
- verifies [slug-callers-agree-qa-2](../../../../evidence/slug-callers-agree-qa-2.md)
- verifies [pr-352-fix-two-change-sets-with-the](../../../../evidence/pr-352-fix-two-change-sets-with-the.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
