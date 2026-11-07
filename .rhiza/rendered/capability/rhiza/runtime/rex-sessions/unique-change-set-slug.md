# unique-change-set-slug

capability · `01M6TBRVYX75117KD157BPHBGB` · parent [rex-sessions](../rex-sessions.md)

## acceptance

Unit tests for the slug function: two change sets with the same title slug give the earlier one the bare slug and the later one '<slug>-<first 12 ID characters, lowercased>'; an earlier landed or abandoned change set still counts; a change set with no clash keeps its bare slug; the result is the same on repeated calls, after the earlier change set lands, and after new change sets whose IDs share a long prefix with it are created

## statement

A change set's slug is its title slug, plus '-' and the first 12 characters of its ID in lowercase (a fixed-length prefix, not the neighbour-dependent short ID) when an earlier change set (lower ID, any state) has the same title slug, given by one function over the folded model

## Links

### Out

- decided-by [slug-unique-by-earliest-id](../../../../decision/slug-unique-by-earliest-id.md)

### In

- verifies [unique-change-set-slug-verified](../../../../evidence/unique-change-set-slug-verified.md)
- verifies [unique-change-set-slug-qa-1](../../../../evidence/unique-change-set-slug-qa-1.md)
- verifies [unique-change-set-slug-qa-2](../../../../evidence/unique-change-set-slug-qa-2.md)
- verifies [pr-352-fix-two-change-sets-with-the](../../../../evidence/pr-352-fix-two-change-sets-with-the.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
