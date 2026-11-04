# plan-usage-governor

source · `01M6HHY4KK7Y87K1YY6FE8G20M`

## quotes

```json
[]
```

## raw_sha256

e0a8006b32a7e5e4f010eb4c1e21da30db8ebb2a62d3f2fc9a1dccb9a07881e7

## session

9e004580-c5cb-495c-8f99-01a1a31a84fb

## summary

Autonomous plan run for the usage governor; no owner took part. The planner read internal/api/daemon.go (its disk-low guard is the precedent for a resource hold), internal/api/health.go, internal/api/inbox.go, and internal/adapter/claudecode (TokenUse dedupes by message id, UsageLimit reads the reset time). It carved the existing usage-governor capability out of Area orchestration and linked the change set as serving anyone-builds, the umbrella's intent. The planner raised the hold rule, token weighting and hold scope; the arbiter decided them in usage-governor-hold-rules. The governor tracks a 5-hour and a 7-day window from every Claude Code transcript on the machine. It counts input, output and cache-creation tokens. Each window's limit comes from a daemon setting, else from the use seen at the last usage-limit hit; a window with no known limit never holds. The reserve defaults to 20%. The projection adds the mean of the last 10 finished runs for each live run and for the next one. A hold stops every daemon launch, like the disk-low guard. It shows as the usage-reserve skip in status, health and inbox, and clears after the reset. The planner also asked which intent to serve, how to set the budget and which windows to track; no separate answer to those came back, so the serves link follows the umbrella and the arbiter's Decision covers the rest. Nothing is left open.

## title

Plan usage governor (autonomous)

## Links

### In

- derived-from [usage-governor-hold-rules](../decision/usage-governor-hold-rules.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
