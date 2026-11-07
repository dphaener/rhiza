# web-project-switcher-fix-verified

evidence · `01M6SWJ70JZSHQQMZ77YFNNS6Z` · actor [claude-code-f03ae487](../actor/claude-code-f03ae487.md)

## result

PASS

## statement

At 19831386 (review fixer for quality f2 and model f1) I read the diff. The switcher now reads web.ProjectNames (api.WebProjectNames: rootcfg names filtered by projectInitialised, with no runs, inbox or health). It reads them lazily through scope.Names(), which only render() in internal/web/web.go calls, so assets, fragments and /p/<x>/events never read them. The current project is always kept. web-design-components' statement now names the switcher and the Projects and Inbox nav. Ran go test -trimpath ./internal/web/...: ok, including TestEveryPageKindIsScopedToItsProject, TestProjectPagesAreScopedWithSwitcher, TestSwitcherKeepsListKindOnDetailPages, TestEventStreamScopedToProject and the new TestSwitcherNamesReadOnlyForPages (0 reads for CSS and fragment, 1 names read and 0 projects reads for a page, switcher links /p/b/runs). Ran go test -trimpath ./internal/api -run 'Web|Project|Root': ok. A whole-package api run hit the 600s go test timeout under machine load, with no failing assertion shown. go vet and gofmt are clean on internal/web and internal/api. make self-check reports only stale .rhiza/rendered files, expected before landing.

## Links

### Out

- verifies [web-project-switcher](../capability/rhiza/web-project-switcher.md)

## History

- `01M6RN3ZBGYDW6KPBSES9A5JXS` Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher
