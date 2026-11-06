# web-runs-rebase-resolution

decision · `01M6QCDPT27EPXFR1CDGYCEM15`

## statement

Resolve the conflict with main in internal/web/web.go by keeping both branches' pages (run, runs, change, roadmap) and main's API context helpers (apiKey, withAPI(r, api), apiGet with notFoundError) from the Roadmap pages; the Runs pages drop their own copies, and run detail and the log tail tell a missing run by notFoundError. The read-only marker test lists both branches' operations.

## wording

Stated by the orchestrator acting as arbiter: the resolution above, merged as 4b34dcb2

## Links

### Out

- answers [rebase-conflict-web-ui-runs-pages](../question/rebase-conflict-web-ui-runs-pages.md)

## History

- `01M6PFH1X52KNK0C5ZCDYEEF64` Web Runs pages: runs live first then recent, and a run's detail with its steps, sessions, review summary and a live log tail
