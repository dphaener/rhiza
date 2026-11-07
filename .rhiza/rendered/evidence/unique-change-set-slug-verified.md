# unique-change-set-slug-verified

evidence · `01M6TD9A102B2YKQABSJA0DFAT` · actor [claude-code-7b3608ec](../actor/claude-code-7b3608ec.md)

## result

PASS

## statement

Verifier at 0d7a957f: read the diff of internal/session/names.go and names_test.go. ChangeSetSlug(sets, cs) takes the folded change sets, returns Slug(DisplayTitle) and appends '-' + strings.ToLower(first 12 ID chars) only when another set with a lower ID (id.Compare, state never consulted, sets without an open entry skipped) has the same title slug, so landing, abandoning or creating later sets cannot change it. Ran go test -trimpath -count=1 -run TestChangeSetSlug -v ./internal/session/ (PASS) and go test -trimpath ./internal/session/... (ok); go vet and gofmt -l clean. TestChangeSetSlug checks: the earlier (landed) set keeps the bare slug 'fix-the-usage-governor-s'; the later one gets 'fix-the-usage-governor-s-01m6t2dc1234'; a no-clash set keeps 'other-work'; a set whose ID shares 25 characters with another still gets the fixed 12-character suffix; after the earlier set is withdrawn (abandoned) the later one keeps its suffix across repeated calls; without the earlier set it is bare.

## Links

### Out

- verifies [unique-change-set-slug](../capability/rhiza/runtime/rex-sessions/unique-change-set-slug.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
