# run-start-new-intent

capability · `01M686V2Z913H9XPK0E1E861YG` · parent [run-start](../run-start.md)

## acceptance

A testscript with a fake agent shows change ls gaining a draft of the right kind and the run journal naming it; --dry-run predicts the ID; --new with a change set argument or with a workflow other than plan exits 2

## statement

rhiza run start --workflow plan --new <intent> [--milestone] opens a draft change set, a record with --milestone and an intention without, and starts the plan run on it

## Links

### Out

- decided-by [new-work-through-run-start](../../../../../decision/new-work-through-run-start.md)

### In

- verifies [run-start-new-intent-verified](../../../../../evidence/run-start-new-intent-verified.md)
- verifies [run-start-new-intent-qa-1](../../../../../evidence/run-start-new-intent-qa-1.md)
- verifies [pr-148-new-work-from-roadmap](../../../../../evidence/pr-148-new-work-from-roadmap.md)
- verifies [pr-301-delete-the-old-rex-apps](../../../../../evidence/pr-301-delete-the-old-rex-apps.md)

## History

- `01M686RVAD1EHGERBBDJN3S45G` The roadmap app starts a plan run for a new milestone or a new change set from a typed one-line intent
- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
