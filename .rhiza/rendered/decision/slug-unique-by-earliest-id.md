# slug-unique-by-earliest-id

decision · `01M6TBQEQ2ANYJFH8AV5EY7FB1`

## reasoning

The question requires a slug that comes out the same every time it is derived. A suffix chosen 'on collision' with whatever is in flight now is not stable: the other change set landing or being abandoned would rename the second one's worktree and branch mid-run. Ordering by ID fixes that: change set IDs are time-ordered ULIDs and change sets are never removed from the log, so the set of earlier change sets with the same slug only changes if someone retitles one, which already changes slugs today. Counting landed and abandoned change sets is what fixes the reported case (01M6T1ZR had landed when 01M6T2DC reused its slug and its merged PR #345). Always appending an ID suffix would be simpler but renames every in-flight change set's worktree and branch (13 accepted ones and live runs, this one included) at deploy; the earliest-keeps-bare-slug rule leaves every non-colliding change set's slug unchanged. The suffix must not be the display short ID: slice.Short (internal/slice/reader.go:140) is one more than the longest common prefix with the neighbouring IDs, at least 6, so its length differs between change sets and grows when a new ID lands next to it, which would rename a live worktree and branch. A fixed prefix of the ID is neighbour-independent; 12 characters cover the ULID's 10-character millisecond timestamp plus 2 random characters, so two same-titled change sets get different suffixes even when created in the same millisecond, and lowercasing fits Slug's [a-z0-9-] alphabet. Callers currently re-derive via session.Slug(t.title) in internal/api (session.go, run.go, run_leaves.go, run_checks.go, ship.go), all of which have the model in hand, so they can route through the one function.

## statement

A change set's slug is session.Slug(title), unless an earlier change set (lower ID, in any state: draft, proposed, accepted, landed or abandoned) has the same title slug; then it is that slug plus '-' and the first 12 characters of its own ID in lowercase (e.g. fix-the-usage-governor-s-01m6t2dcxxxx). The suffix is a fixed-length prefix of the ID, never the neighbour-dependent display short ID (slice.Short), so it does not change as other change sets are created. The rule is one function over the folded model (the change set's ID and title plus every other change set's ID and title), and every caller (worktree, branch, ship target, abandon cleanup, run tabs and run sessions) calls it instead of session.Slug(title) directly.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-f74b00a7-1](../question/question-f74b00a7-1.md)

### In

- decided-by [unique-change-set-slug](../capability/rhiza/runtime/rex-sessions/unique-change-set-slug.md)
- decided-by [slug-callers-agree](../capability/rhiza/runtime/rex-sessions/slug-callers-agree.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
