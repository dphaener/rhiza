# usage-hold-rules-live-run-remaining

decision · `01M6T3DYFGRYPHE9GPYS6S5YWX`

## reasoning

Read the question, usage-governor-hold-rules' current text, usage-live-run-remaining (01M6T2AQ, from accepted draft 01M6T1ZR) and this plan's usage-hold. usage-live-run-remaining and usage-hold both say a live run adds only the mean less its tokens so far, never below zero, plus one full mean for the next pick; the rewritten 01M6HHVZ still counts a full mean per live run. 01M6T1ZR does not touch 01M6HHVZ, so only this change set's rewrite can remove the contradiction. The rewrite already folds in 01M6T1ZR's ratchet (usage-learned-limit-ratchets), so folding in its live-run rule the same way is consistent, keeps the umbrella Decision true to the code 01M6T1ZR builds, and changes no behaviour, only the text. The alternative, leaving 01M6HHVZ stale, would keep a self-contradicting model that rhiza check reviewers flag every round.

## statement

The plan's rewrite of usage-governor-hold-rules (01M6HHVZ) changes its hold sentence to match usage-live-run-remaining and usage-hold: the daemon holds a new pick when the window's use so far, plus for each live run the mean tokens of the last 10 finished runs less that run's tokens so far (never below zero), plus one full mean for the run it would start, exceeds the limit minus the reserve. With no finished runs to average, only the use so far counts. The rest of the rewritten statement stays as planned.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6t2db-r3-f1](../question/review-01m6t2db-r3-f1.md)

## History

- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available
