# rex-app-actions

capability · `01M644QW8GJBD8KTJ8Q00MMNYY` · parent [rhiza](../rhiza.md)

## acceptance

A rex-install script test: after rhiza rex install, rhiza.lua has one rex.action per registered app (inbox, request, runs, daemon and autonomy, exactly five), each named rhiza_app_<name> and running rhiza tui open <name>; internal/api TestRexInstallRegistersAnActionPerApp registers one more app and reinstalls, and finds its action in rhiza.lua

## statement

Each rhiza app is a Rex action of its own in the command palette, titled with the app's name, and running it runs that app

## Links

### Out

- serves [tui-native-actions](../../intent/rhiza-tui/tui-native-actions.md)
- decided-by [rex-app-own-action](../../decision/rex-app-own-action.md)

### In

- verifies [rex-app-actions-verified](../../evidence/rex-app-actions-verified.md)
- verifies [rex-app-actions-qa-1](../../evidence/rex-app-actions-qa-1.md)
- verifies [rex-app-actions-verified-2](../../evidence/rex-app-actions-verified-2.md)
- verifies [rex-app-actions-qa-2](../../evidence/rex-app-actions-qa-2.md)
- verifies [pr-77-native-rex-actions-for-apps](../../evidence/pr-77-native-rex-actions-for-apps.md)
- verifies [pr-301-delete-the-old-rex-apps](../../evidence/pr-301-delete-the-old-rex-apps.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
