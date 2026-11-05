# usage-governor-hold-rules

decision · `01M6HHVZH70MTTQ6JJ2D403TY5`

## reasoning

No owner words; read the change set intent, internal/api/daemon.go and internal/adapter/claudecode. The intent asks for per-window use from the adapter's token records, a reserve setting, a skip reason in status/health/inbox, and resuming on reset. The disk-low guard (rhiza.daemon.minFree) is the existing precedent for a resource hold. It stops picks, resumes, queued work and plans, but leaves live runs alone. Copying its scope keeps the daemon consistent, and it protects the reserve, which exists for the owner's own use. Claude Code publishes no token limits, so the limit must be either configured or learned. The adapter already detects usage-limit hits (UsageLimit and PaneUsageLimit) with a reset time, so the use seen at a hit gives a calibration without any setup. No known limit means no hold, so the governor fails open instead of blocking work on a guess. Cache reads dominate the raw totals and are weighted far below fresh tokens, so counting them would swamp the signal. The account is shared, so the owner's interactive sessions eat the same window and must count. Averaging recent runs per live and next run projects whether the next pick fits. Extrapolating to the window's end would hold early in every window.

## statement

(1) When: at each poll the governor checks two windows, the 5-hour session window and the 7-day weekly window. A window starts at the first token record after the previous window ended, or at the reset time a usage-limit message gave. Each window has a limit: the daemon setting for that window (tokens) if set, else the window use at the last usage-limit hit in that kind of window. A run driver that sees a hit records its time and reset time where the daemon reads it at each poll. The hit's kind comes from its reset time: a reset 5 hours or less after the hit is the session window, a later one the weekly window. A hit with no reset time teaches nothing. A window with no known limit never holds. The reserve is a daemon setting, a percentage of the limit, default 20%. The daemon holds a new pick when the window's use so far, plus the mean tokens of the last 10 finished runs for each live run and for the one it would start, exceeds the limit minus the reserve. With no finished runs to average, only the use so far counts. (2) Tokens: input, output and cache-creation tokens count, for the session and its helpers; cache-read tokens do not. They come from every Claude Code transcript on the machine, the owner's own sessions included, since they share the account. Only records timestamped inside the window count, each message once (deduplicated by message id, as TokenUse already does). (3) Scope: a hold works like the disk-low guard. It stops everything the daemon would launch: picks, its resumes of interrupted runs, queued starts and continues, and plan requests. Each held change set shows as a skip with reason usage-reserve, giving the window, use, limit and reset time in its detail, in daemon status, health and the inbox. It never interrupts live runs, and it does not stop runs the owner starts directly outside the daemon. It clears on its own at the first poll after the window resets.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-9e004580-1](../question/question-9e004580-1.md)
- derived-from [plan-usage-governor](../source/plan-usage-governor.md)

### In

- decided-by [usage-governor](../capability/rhiza/orchestration/usage-governor.md)
- decided-by [usage-window-tokens](../capability/rhiza/orchestration/usage-governor/usage-window-tokens.md)
- decided-by [usage-windows](../capability/rhiza/orchestration/usage-governor/usage-windows.md)
- decided-by [usage-settings](../capability/rhiza/orchestration/usage-governor/usage-settings.md)
- decided-by [usage-learned-limit](../capability/rhiza/orchestration/usage-governor/usage-learned-limit.md)
- decided-by [usage-hold](../capability/rhiza/orchestration/usage-governor/usage-hold.md)
- decided-by [usage-hold-reported](../capability/rhiza/orchestration/usage-governor/usage-hold-reported.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
