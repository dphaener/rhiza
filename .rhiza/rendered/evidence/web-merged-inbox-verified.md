# web-merged-inbox-verified

evidence · `01M6SQZSTKXHPB5DDR3GZA68FC` · actor [claude-code-baea7c60](../actor/claude-code-baea7c60.md)

## result

PASS

## statement

Verifier claude-code-baea7c60 at 6480ce30: go test -trimpath -count=1 -run MergedInbox -v ./internal/web passes TestMergedInboxOrdersTagsAndLinks (two projects alpha/beta, three items rendered oldest first across projects: Alpha early failure, Beta middle sign-off, Alpha late question; each row carries <span class="rz-project-tag"> with its project; links /p/alpha/runs/R1, /p/beta/changes/C2, /p/alpha/changes/C1) and TestMergedInboxEmpty. Read the diff: mergedInboxRows prefixes each row's href with /p/<escaped project>, tags Row.Project, stable-sorts by RFC 3339 At then project name; RootWithInbox serves GET/HEAD /inbox; WebMergedInboxSource reads each initialised registered project's inbox. go test -run 'WebRoot|WebProject|WebInbox|WebMerged' ./internal/api passes. TestOverviewReadsDaemonFromRepo fails only when run from this checkout nested in the live repo's .git; it passes at the same commit in a clean worktree under TMPDIR, so it is environmental.

## Links

### Out

- verifies [web-merged-inbox](../capability/rhiza/web-merged-inbox.md)

## History

- `01M6RN3ZBGYDW6KPBSES9A5JXS` Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher
