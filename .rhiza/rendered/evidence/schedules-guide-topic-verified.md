# schedules-guide-topic-verified

evidence · `01M6WW7ZVA9ZG7Z3RRJHG1FXRN` · actor [claude-code-7d947278](../actor/claude-code-7d947278.md)

## result

PASS

## statement

At 05d571ac ran go test -trimpath -count=1 -v ./internal/guide/...: all tests PASS, including TestWorkflowsText, whose new schedules group requires '## Schedules', [[schedules]], CRON_TZ, catch_up, overlap, every, cron, 'main checkout', 'no change set', 'as it found it', rhiza request, rhiza health and daemon log, and the size checks TestCaps (workflows 800 tokens, cap 800), TestCapCatchesOversize and TestTokens. Read the diff from 98589c61~1: the Schedules section matches internal/dispatch (keys name, every, cron, workflow, command, catch_up once|skip, overlap skip|queue) and internal/schedule (MinEvery 1m, five-field cron with leading CRON_TZ=<zone>, else UTC). Attempt 1's finding f1 is fixed: Outside input again says 'Nothing is refused.' and lists 'narrow rhiza.daemon.permissionMode or permission_mode on the plan workflow's steps', and the test now checks both phrases. The bytes came from trimming Files/Step kinds/Targets/Dispatch wording; the points the workflow-guide-topic and dispatch-guide-topic capabilities require (resolution layers, every step kind, done/failed, includes, the workflow loop commands, both dispatch tables, defaults, ready order, workflow-unresolved skip) remain. The fix commit stages only its builder's actor record.

## Links

### Out

- verifies [schedules-guide-topic](../capability/rhiza/workflows/scheduled-triggers/schedules-guide-topic.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
