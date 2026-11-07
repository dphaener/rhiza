# web-root-listener

capability · `01M6RN862M5W8HPN62BWDYC124` · parent [rhiza](../rhiza.md)

## acceptance

internal/web and api tests: a root daemon of two projects opens one listener, and the per-project pages answer for both projects through it.

## statement

The root daemon serves one web UI for the machine on the machine-wide web host and port, with every project's data behind it.

## Links

### Out

- decided-by [m29-one-root-daemon](../../decision/m29-one-root-daemon.md)
- decided-by [m29-web-all-projects](../../decision/m29-web-all-projects.md)

### In

- verifies [web-root-listener-verified](../../evidence/web-root-listener-verified.md)
- verifies [web-root-listener-qa-1](../../evidence/web-root-listener-qa-1.md)
- verifies [web-root-listener-verified-2](../../evidence/web-root-listener-verified-2.md)
- verifies [web-root-listener-qa-2](../../evidence/web-root-listener-qa-2.md)
- verifies [pr-341-web-ui-across-projects](../../evidence/pr-341-web-ui-across-projects.md)

## History

- `01M6RN3ZBGYDW6KPBSES9A5JXS` Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher
