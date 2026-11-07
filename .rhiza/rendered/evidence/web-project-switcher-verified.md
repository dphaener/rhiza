# web-project-switcher-verified

evidence · `01M6SR9TA4Q9NHE35N3MTZ996B` · actor [claude-code-f888349e](../actor/claude-code-f888349e.md)

## result

PASS

## statement

At c4cee9a5 ran go test -trimpath -count=1 ./internal/web/ -run 'EveryPageKind|Switcher|Scoped|EventStream|Project' -v: all pass, incl. TestEveryPageKindIsScopedToItsProject, TestProjectPagesAreScopedWithSwitcher, TestSwitcherKeepsListKindOnDetailPages, TestEventStreamScopedToProject. go vet and gofmt clean on internal/web. The full ./internal/web/ package passes at c4cee9a5 in a clean detached worktree; in the run checkout under .git/rhiza/checkouts, TestOverviewReadsDaemonFromRepo fails, but the same test passes in a clean worktree at this commit and at 8cb4b26c and 87254fce, and this commit only adds a test file, so that failure comes from where the checkout sits, not from this leaf. Read the diff: the new test builds two projects (a, b) each with its own change set and run, fetches overview, roadmap, runs, run detail and change set pages under /p/a/ and /p/b/, asserts each shows its own titles/IDs and none of the other project's, carries a switcher link to the other project's same page kind (detail pages link the list kind) plus its own entry, and that /p/nope/<page> and the other project's run/change set under this project are 404. This closes attempt 1's finding f1.

## Links

### Out

- verifies [web-project-switcher](../capability/rhiza/web-project-switcher.md)

## History

- `01M6RN3ZBGYDW6KPBSES9A5JXS` Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher
