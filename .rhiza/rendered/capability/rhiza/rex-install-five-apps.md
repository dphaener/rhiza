# rex-install-five-apps

capability · `01M6MZD7XSJD2ENEG7W98BKDJV` · parent [rhiza](../rhiza.md)

## acceptance

A test in internal/api counts exactly five rex.action entries, named rhiza_app_inbox, rhiza_app_request, rhiza_app_runs, rhiza_app_daemon and rhiza_app_autonomy, each with one key; the rex-install testscript, with a fake rex, greps those five and grep -count=5 for rex.action.

## statement

rhiza rex install registers exactly five Rex actions with their keys in the rhiza key mode: inbox, request, runs, daemon and autonomy.

## Links

### Out

- decided-by [tui-delete-old-apps](../../decision/tui-delete-old-apps.md)
- decided-by [tui-five-apps](../../decision/tui-five-apps.md)

### In

- verifies [rex-install-five-apps-verified](../../evidence/rex-install-five-apps-verified.md)
- verifies [rex-install-five-apps-qa-1](../../evidence/rex-install-five-apps-qa-1.md)
- verifies [rex-install-five-apps-verified-2](../../evidence/rex-install-five-apps-verified-2.md)
- verifies [rex-install-five-apps-qa-2](../../evidence/rex-install-five-apps-qa-2.md)
- verifies [pr-301-delete-the-old-rex-apps](../../evidence/pr-301-delete-the-old-rex-apps.md)

## History

- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
