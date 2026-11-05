# autonomy-app-daemon-keys-qa-1

evidence · `01M6FN7C21FFD8NXW3A1GX4J2F` · actor [claude-code-e6cae1c1](../actor/claude-code-e6cae1c1.md)

## result

PASS

## statement

QA at 78a3b478: rhiza tui app autonomy --once in scratch repo qarepo printed the settings JSON (key, value, source, default). The interactive app ran in a real Rex tab and was driven with rex send-key. It showed 'Daemon not running', every daemon and autonomy setting with value and source, and the footer 'd start/stop daemon'. d then y started the daemon in the Rex tab 'rhiza daemon', and after r it showed 'Daemon running, process Mac.localdomain:47359'. d asked 'Stop the daemon Mac.localdomain:47359? y/n', y ended that process, and the app showed 'Daemon not running'. With the cursor on rhiza.autonomy.plan, enter toggled it from owner to auto, and git config rhiza.autonomy.plan then read auto.

## Links

### Out

- verifies autonomy-app (retired)

## History

- `01M6FEMDQ1D2HWXT3GR4F3Z5H6` The daemon runs as a systemd service or in a Rex tab, as each installation chooses, and opens its runs in a fixed per-repo Rex session
