# autonomy-app

capability · `01M6CS2HRAWMCTPDHFXBHFCX1C` · parent [daemon-settings](../daemon-settings.md)

## acceptance

tui-autonomy-app.txtar: --once prints the settings JSON; unit tests on a fake runner toggle plan from owner to auto (calling daemon set rhiza.autonomy.plan auto), set paused, edit maxRuns, and start the daemon with rex window new rhiza daemon when status says not running; rex-install.txtar shows rhiza_app_autonomy on key y and no existing key moved

## statement

A Rex app, autonomy, on key y, lists every daemon and autonomy setting with value and source, toggles auto/owner and paused on enter, edits numbers in place through rhiza daemon set, and starts the daemon in a Rex tab named rhiza daemon when none is running or stops a running one on request

## Links

### Out

- decided-by [autonomy-app-controls](../../../../../decision/autonomy-app-controls.md)

### In

- verifies [autonomy-app-verified](../../../../../evidence/autonomy-app-verified.md)
- verifies [autonomy-app-qa-1](../../../../../evidence/autonomy-app-qa-1.md)
- verifies [autonomy-app-verified-2](../../../../../evidence/autonomy-app-verified-2.md)
- verifies [autonomy-app-qa-2](../../../../../evidence/autonomy-app-qa-2.md)
- verifies [pr-195-autonomy-app](../../../../../evidence/pr-195-autonomy-app.md)

## History

- `01M6CRZTNDTPF1SD2WGBY4835Y` An autonomy app edits every setting, starts or stops the daemon, and the inbox nudges stalled sessions
