# inbox-lists-api-error-sessions

capability · `01M7A9C4G0XBK7HYT3S4EJ6254` · parent [inbox](../inbox.md)

## acceptance

A Go test in internal/session (go test -trimpath ./internal/session/...): an autonomous session idle 21 minutes with the API-error detail is not Stalled while one with a plain idle detail is. cmd/rhiza/testdata/script/inbox-api-error.txtar (go test -trimpath ./cmd/rhiza -run TestScripts/inbox-api-error): with a live run journal and a session record whose hook report is idle with the API-error detail, rhiza inbox --json has one item of kind api-error with the session, the run, the error text and the session message next step, and none of kind stalled at T+21m; the text inbox shows it; after a later working hook report the item is gone. Go tests in internal/api (go test -trimpath ./internal/api/ -run 'WebInbox|InboxNotice'): webInboxViews keeps an api-error item (webInboxKinds lists the new kind) with the change set's title and a detail naming the API error, and no session or run ID or handle in its title or detail; inboxNotice gives an api-error item its own line saying a session stopped on an API error, not the default gave-up line, with no handle or ID. A Go test in internal/tui/inboxapp (go test -trimpath ./internal/tui/inboxapp/...): an api-error item is a row, ActionFor gives it ActionShow, and Enter shows its kind and detail (the error text) and calls inbox.open on it. describe-golden.txtar is regenerated for the new kind (go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden). Touches internal/session/status.go, internal/session/status_test.go, internal/api/run.go, internal/api/inbox.go, internal/api/daemon.go, internal/api/daemon_notice_test.go, internal/api/web_inbox.go, internal/api/web_inbox_test.go, internal/tui/inboxapp/inboxapp_test.go, cmd/rhiza/testdata/script/inbox-api-error.txtar, cmd/rhiza/testdata/script/describe-golden.txtar

## statement

A live run's step session whose hook report is idle with a detail starting 'stopped on an API error: ' (session.APIErrorDetail) is listed by rhiza inbox as an item of kind api-error naming the session, its run and the error, with `rhiza session message <session> <text>` as the next step, and the daemon's notice, the web inbox and the inbox app show it; such a session is never stalled, and its next hook report clears the item

## Links

### Out

- decided-by [api-error-inbox-as-idle-detail](../../../../decision/api-error-inbox-as-idle-detail.md)
- decided-by [api-error-silence-left-to-watchdog](../../../../decision/api-error-silence-left-to-watchdog.md)

### In

- depends-on [session-run-reprompts-api-error](../../workflows/workflow-runs/run-resume/session-run-reprompts-api-error.md)
- verifies [inbox-lists-api-error-sessions-verified](../../../../evidence/inbox-lists-api-error-sessions-verified.md)
- verifies [inbox-lists-api-error-sessions-qa-1](../../../../evidence/inbox-lists-api-error-sessions-qa-1.md)
- verifies [pr-548-fix-run-sessions-that-hit-an](../../../../evidence/pr-548-fix-run-sessions-that-hit-an.md)

## History

- `01M7A949Q4VGB6D4PPY9A5GZEK` Fix Run sessions that hit an API error stop at the prompt and sit idle forever, unflagged. On 2026-11-13 around 09:47-09:52Z three agent sessions ended their turn with Claude Code's 'API Error: 403 Attention Required! | Cloudflare' (transient): verifier 4594e9fa (run 01M79TRE3W, leaf daemon-holds-steps-on-stale-guide, idle 09:47:45-10:05) and review checks quality dbb3e57f and compliance 8b78600b (run 01M7A144DD, cycle 1, idle 09:51-10:05). Their processes stayed alive, the runs showed 'running', no inbox item, no watchdog action; load fell to ~2 while three slots did nothing. Unstuck by the orchestrator: rhiza session message <session> asking each to carry on; all resumed within seconds. Wanted: the run (or watchdog) notices a session whose transcript ends in an 'API Error:' assistant message (or whose transcript has not grown for N minutes with no child process) and re-prompts it automatically, retrying a few times with backoff, then raises an inbox item. Related, different cause: 01M6N9AB3X (idle after background shell), 01M79MDD5R (idle after unknown slash command).
