# milestone-26

intent · `01M6MZ09HMFH6H2KZF6C8EH8RH`

## milestone

26

## statement

Rhiza runs are fully managed from five small single-purpose Rex apps built on rexuikit (Inbox, Request, Runs, Daemon and Autonomy), each a small modal that opens over the current pane, does one job and closes; every other app is gone and the CLI keeps every command.

## Links

### Out

- decided-by [tui-overhaul-why](../decision/tui-overhaul-why.md)
- decided-by [tui-single-purpose-apps](../decision/tui-single-purpose-apps.md)
- decided-by [tui-small-modal](../decision/tui-small-modal.md)
- decided-by [tui-rexuikit-components](../decision/tui-rexuikit-components.md)
- decided-by [tui-milestone-first](../decision/tui-milestone-first.md)
- decided-by [tui-plan-all-now](../decision/tui-plan-all-now.md)

### In

- serves change set `01M6MYM5GB3Q8DBZST8V5AG6RQ` (Modal shell: a shared package that opens a Rhiza app as a rexuikit surface Modal sized to its content, renders rexuikit ui components (Form, List, Dialog) with the rexuikit theme, and closes on Esc or when the app's job is done)
- serves change set `01M6MYN9GRWHRJSHA31D9S838H` (Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces)
- serves change set `01M6MYNDJ4HAJF028CEAE7SYQN` (Request app: file a feature, bug or milestone request in one rexuikit Form, then close; rebuilt on the modal shell over the request operation, replacing internal/tui/requestsapp and retiring the capabilities of the old requests app it replaces)
- serves change set `01M6MYNJGY6Q735KB2ERT3V4R5` (Runs app: a read-only glance at live runs and their current steps, rebuilt on the modal shell over run.ls, replacing internal/tui/liveapp and retiring the capabilities of the old live app it replaces)
- serves change set `01M6MYNQTAMMV8YEZHDS1G5HMD` (Daemon app: pause or resume the daemon and set its run limit, rebuilt on the modal shell over daemon.settings and daemon.set (rhiza.daemon.paused, rhiza.daemon.maxRuns))
- serves change set `01M6MYPA335J1S3952RJAGZJN9` (Autonomy app: the plan, accept, ship and sign-off toggles (rhiza.autonomy.*), rebuilt on the modal shell over daemon.settings and daemon.set, replacing internal/tui/autonomyapp and retiring the capabilities of the old autonomy app it replaces)
- serves change set `01M6MYNZJSMEXAE5Y3P3S9FB6F` (Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command)
- serves change set `01M6RXK2SB0M12FXVEGXER022R` (Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.)

## History

- `01M6MYK80HN9E4DG8QREZH38DW` Milestone 26: Rhiza runs are fully managed from five small single-purpose Rex apps built on rexuikit, each a small modal that opens, does one job and closes; the Nebari import moves to milestone 27
