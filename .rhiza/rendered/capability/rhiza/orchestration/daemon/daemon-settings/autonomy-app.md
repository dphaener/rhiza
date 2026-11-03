# autonomy-app

capability · `01M6CS2HRAWMCTPDHFXBHFCX1C` · parent [daemon-settings](../daemon-settings.md)

## acceptance

tui-autonomy-app.txtar: --once prints the settings JSON; unit tests on a fake runner toggle plan from owner to auto (calling daemon set rhiza.autonomy.plan auto), set paused, edit maxRuns, and call rhiza daemon start when status says not running and rhiza daemon stop when it says running; rex-install.txtar shows rhiza_app_autonomy on key y and no existing key moved

## statement

A Rex app, autonomy, on key y, lists every daemon and autonomy setting with value and source, toggles auto/owner and paused on enter, edits numbers in place through rhiza daemon set, and on d starts the daemon with rhiza daemon start when none is running or stops a running one with rhiza daemon stop, so it runs the daemon however its installation does (daemon-start-stop)

## Links

### Out

- decided-by [autonomy-app-controls](../../../../../decision/autonomy-app-controls.md)

### In

- verifies [autonomy-app-verified](../../../../../evidence/autonomy-app-verified.md)
- verifies [autonomy-app-qa-1](../../../../../evidence/autonomy-app-qa-1.md)
- verifies [autonomy-app-verified-2](../../../../../evidence/autonomy-app-verified-2.md)
- verifies [autonomy-app-qa-2](../../../../../evidence/autonomy-app-qa-2.md)
- verifies [pr-195-autonomy-app](../../../../../evidence/pr-195-autonomy-app.md)
- verifies [autonomy-app-daemon-keys-qa-1](../../../../../evidence/autonomy-app-daemon-keys-qa-1.md)
- verifies [pr-200-the-daemon-runs-as-a-systemd](../../../../../evidence/pr-200-the-daemon-runs-as-a-systemd.md)

## History

- `01M6CRZTNDTPF1SD2WGBY4835Y` An autonomy app edits every setting, starts or stops the daemon, and the inbox nudges stalled sessions
- `01M6FEMDQ1D2HWXT3GR4F3Z5H6` The daemon runs as a systemd service or in a Rex tab, as each installation chooses, and opens its runs in a fixed per-repo Rex session
