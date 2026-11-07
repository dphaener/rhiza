# web-project-switcher

capability · `01M6RN8687F2WDVVXSSTW5W65A` · parent [rhiza](../rhiza.md)

## acceptance

internal/web tests fetch each page kind under /p/<a>/ and /p/<b>/ and see only that project's data, a switcher link to the other project's same page, and a 404 for an unknown project.

## statement

Every existing web page (overview, roadmap, runs, run detail, change set) is scoped to one project under /p/<project>/, with a project switcher in the header that keeps the page kind, and the live event stream is scoped to the same project.

## Links

### Out

- decided-by [m29-web-all-projects](../../decision/m29-web-all-projects.md)

### In

- verifies [web-project-switcher-verified](../../evidence/web-project-switcher-verified.md)
- verifies [web-project-switcher-qa-1](../../evidence/web-project-switcher-qa-1.md)
- verifies [web-project-switcher-fix-verified](../../evidence/web-project-switcher-fix-verified.md)
- verifies [web-project-switcher-qa-2](../../evidence/web-project-switcher-qa-2.md)
- verifies [pr-341-web-ui-across-projects](../../evidence/pr-341-web-ui-across-projects.md)

## History

- `01M6RN3ZBGYDW6KPBSES9A5JXS` Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher
