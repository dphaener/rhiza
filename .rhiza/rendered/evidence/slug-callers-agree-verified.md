# slug-callers-agree-verified

evidence · `01M6THCGTKYPJDVY2TTGZXWA0Z` · actor [claude-code-23cb30f6](../actor/claude-code-23cb30f6.md)

## result

PASS

## statement

Verifier at a47ea982 (leaf commits cd5fe1a5, a47ea982). go test -trimpath ./cmd/rhiza -run 'TestScripts/(change-worktree-unique-slug|arbiter-agent-tabs)$' passed: two change sets titled 'Twin work' get worktrees repo.twin-work and repo.twin-work-01m5pwgjc82p and branches cs/twin-work and cs/twin-work-01m5pwgjc82p (12 lowercased ID characters, not the 9-character short ID), and run start --dry-run (existing and --new plan) predicts the suffixed worktree. go test -trimpath ./internal/api -run TestUniqueSlugShipTargetAndAbandon passed: findShipTarget names the later change set's suffixed slug/branch/worktree, and abandonCleanup names cs/twin-work-<suffix>, never cs/twin-work. Full go test -trimpath -timeout 40m ./internal/api passed (416s; an earlier run hit the default 10m timeout under load ~10) and ./internal/session passed. grep for 'Slug(' in non-test Go: session.Slug is called only inside session.ChangeSetSlug; change worktree, run start/dry run (newPlanChange), leaf/check/fix/agent prompts (t.slug), ship target and abandon cleanup all use ChangeSetSlug.

## Links

### Out

- verifies [slug-callers-agree](../capability/rhiza/runtime/rex-sessions/slug-callers-agree.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
