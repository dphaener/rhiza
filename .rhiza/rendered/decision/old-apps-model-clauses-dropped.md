# old-apps-model-clauses-dropped

decision · `01M6NQG8RNZPRTDB3VHB0V9PDT`

## reasoning

The change set's intent is to delete changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps, retire their capabilities, and keep every CLI command. The branch already does this in code: internal/tui keeps only autonomyapp, daemonapp, inboxapp, requestapp and runsapp, and text-refs-display.txtar now checks rhiza show, slice and why in place of the read app. Review cycles 2 and 3 confirmed that f1 (text-ref-display) and f2 (roadmap-paged) are the same defect as one the fixer already corrected for question-answer-op and run-start-through-daemon: the capabilities still promise app behaviour the intent removes. The CLI parts of these capabilities stay valid and tested, so retiring them would lose real behaviour, and loosening old-apps-deleted would go against the intent. The fix is to revise both fields of each affected capability within this change set.

## statement

Fix the model, not the leaf: old-apps-deleted stays as written. In this change set, revise text-ref-display to drop 'and the read apps' from its statement and the 'a read app's --once JSON (rhiza tui app <read app> --once) shows @@renamed' clause from its acceptance. Also revise roadmap-paged to drop the roadmap TUI app and the model app's roadmap tab from its statement, and 'and the roadmap app' from its acceptance, so only health's roadmap check remains. Every other live capability that still names a deleted app (for example run-status-open-cycle's 'the live app counts it running') gets the same treatment: keep only its CLI behaviour, without retiring the capability or changing that behaviour.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [leaf-failed-01m6njpm-1](../question/leaf-failed-01m6njpm-1.md)

## History

- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
