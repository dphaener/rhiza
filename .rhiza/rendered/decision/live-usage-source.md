# live-usage-source

decision · `01M6T2M4HFGZH5WZ0Z731RPDVS`

## reasoning

The intent names the status line's rate_limits (used_percentage, resets_at per window) as the source and says to derive limit and reset at every poll and fall back to settings only when no live figure exists, so live must outrank the setting (today setting outranks learned; learned stays last). Rhiza already launches every run's claude with its own --settings JSON (internal/adapter/claudecode Settings), so a statusLine command is the cheapest reliable capture and keeps the adapter boundary: daemon code never parses Claude formats directly. The 'equivalent API' would need network and OAuth credentials, which the tests forbid and the adapter's file-only UsageReader contract avoids. The governor counts tokens (internal/usage windows over transcript records), so a percentage converts to a token limit as used/percentage, using resets_at to pin the window, which also fixes window boundaries the record chain only guesses. Low percentages make that ratio unstable, hence the 1% floor. Validity tied to resets_at is deterministic and testable with RHIZA_TEST_NOW, and needs no arbitrary staleness timer since the limit is a plan property within a window. Draft 01M6T1ZR edits the same usage_hold.go and usage_learn.go functions (KindOf use in learnLimits, the live-run term in usageHoldOf); sequencing after it avoids conflicting edits and lets this change set rely on a hit teaching only its kind.

## statement

(1) Source: Claude Code's status line input. The claudecode adapter's --settings JSON (which Rhiza already writes for every run's session) gains a statusLine command, a rhiza subcommand, that reads the status line JSON on stdin and records its rate_limits (per window: used_percentage, resets_at, plus the capture time) to a live-usage file in the same state folder as the limit hits (the root's .rhiza-root when under a root, else the repo's daemon dir), latest snapshot per window kind (five_hour = session, seven_day = weekly). Reading it stays behind an optional adapter interface, like UsageReader; no network API call, no credential reads. (2) Derivation: a snapshot fixes its window exactly: Reset = resets_at, Start = resets_at minus the kind's length (replacing the inferred chain for that kind). The limit is the window's tokens from the transcript records over [Start, capture time) divided by used_percentage/100. A snapshot below 1% used, or with no tokens in its window, yields no limit (too noisy). At each poll the window's Used is still counted fresh from the records up to now and compared against that derived limit with the reserve, as today. (3) Use and fallback: a snapshot is live while its resets_at is after now; the latest live snapshot of a kind is used at every poll. Order per kind: live snapshot, then the rhiza.daemon.usage.<kind> setting, then the learned limit; with none the kind never holds. The machine-wide hold from a recorded limit hit stays as it is, independent of all three. (4) Yes: this change set lands after draft 01M6T1ZR. It fixes which kind a hit teaches and the live-run double count in the same usageHoldOf/learnLimits code; this one builds the live-snapshot path on top of that corrected code rather than racing it.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-c9536b48-1](../question/question-c9536b48-1.md)
- derived-from [live-usage-plan-interview](../source/live-usage-plan-interview.md)

### In

- relates-to [usage-governor-hold-rules](usage-governor-hold-rules.md) — Amended by live-usage-source: a live snapshot now comes first for a window's limit and fixes its start and reset; the setting and the learned limit are its fallbacks.
- decided-by [usage-live-snapshot-record](../capability/rhiza/orchestration/usage-governor/usage-live-snapshot-record.md)
- decided-by [usage-live-snapshot-read](../capability/rhiza/orchestration/usage-governor/usage-live-snapshot-read.md)
- decided-by [usage-live-limit](../capability/rhiza/orchestration/usage-governor/usage-live-limit.md)
- decided-by [usage-hold](../capability/rhiza/orchestration/usage-governor/usage-hold.md)

## History

- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available
