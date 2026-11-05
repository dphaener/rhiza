# roadmap-paged-drop-deleted-app-readers

decision · `01M6NQJEKEHG4QZACJ4Q774G1B`

## reasoning

The change set's intent deletes roadmapapp and modelapp and retires their capabilities, so any capability that still names them as readers is stale and the model check is right to flag it. The rest of roadmap-paged still holds: internal/api/health.go keeps healthRoadmap, and internal/api/roadmap_paged_test.go TestRoadmapOverAPageIsReadWhole already checks that a roadmap of more than 100 change sets is read whole by --all and health's roadmap check. Removing the deleted apps from roadmap-paged's statement and acceptance is the same kind of edit as commit 6d230ba2, which dropped deleted apps from other acceptances. It is within the intent and needs no code change. Changing the leaf instead would mean keeping apps the intent says to delete.

## statement

Fix the model, not the leaf: within this change set, revise roadmap-paged so it names only readers that still exist. Its acceptance clause becomes 'a roadmap of more than 100 change sets is read whole by health's roadmap check' (drop 'and the roadmap app'), and its statement ends 'and health's roadmap check reads it with all' (drop 'the roadmap TUI app and the model app's roadmap tab'). Leaf old-apps-deleted stays as planned.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [leaf-failed-01m6njpm-2](../question/leaf-failed-01m6njpm-2.md)

## History

- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
