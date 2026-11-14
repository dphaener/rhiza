# inbox-lists-api-error-sessions-verified

evidence · `01M7ABGK1WVNXGQ1T889EEXVYH` · actor [claude-code-2720cd5f](../actor/claude-code-2720cd5f.md)

## result

PASS

## statement

At 4e3806b9a, re-ran every Checks: command and each acceptance command, all exit 0: go test -trimpath ./internal/session/... (TestStalled: autonomous idle api-error at 21m not stalled, plain idle at 21m stalled); go test -trimpath ./internal/tui/inboxapp/... (TestAPIErrorItemShowsDetailAndOpensTab: row, ActionShow, Enter calls inbox.open and shows kind and error text); go test -trimpath ./internal/api/ -run 'WebInbox|InboxNotice' and the wider Checks pattern (TestWebInboxViewsAPIError, api-error notice line without IDs); go test -trimpath ./cmd/rhiza -run TestScripts/inbox-api-error (api-error item with session, run, error text, rhiza session message next, no stalled at T+21m, text inbox shows it, cleared after working report) and TestScripts/describe-golden plus inbox-stalled, inbox-permission, tui-inbox-app; rhiza ship generate and rhiza render exit 0 leaving the tree clean. Read the diff: each acceptance point met and tested; changes outside the Touches list are only inboxItems signature updates in three api tests and the actor records; no TODOs or stubs; tests use fake rex/claude only.

## Links

### Out

- verifies [inbox-lists-api-error-sessions](../capability/rhiza/orchestration/inbox/inbox-lists-api-error-sessions.md)

## History

- `01M7A949Q4VGB6D4PPY9A5GZEK` Fix Run sessions that hit an API error stop at the prompt and sit idle forever, unflagged. On 2026-11-13 around 09:47-09:52Z three agent sessions ended their turn with Claude Code's 'API Error: 403 Attention Required! | Cloudflare' (transient): verifier 4594e9fa (run 01M79TRE3W, leaf daemon-holds-steps-on-stale-guide, idle 09:47:45-10:05) and review checks quality dbb3e57f and compliance 8b78600b (run 01M7A144DD, cycle 1, idle 09:51-10:05). Their processes stayed alive, the runs showed 'running', no inbox item, no watchdog action; load fell to ~2 while three slots did nothing. Unstuck by the orchestrator: rhiza session message <session> asking each to carry on; all resumed within seconds. Wanted: the run (or watchdog) notices a session whose transcript ends in an 'API Error:' assistant message (or whose transcript has not grown for N minutes with no child process) and re-prompts it automatically, retrying a few times with backoff, then raises an inbox item. Related, different cause: 01M6N9AB3X (idle after background shell), 01M79MDD5R (idle after unknown slash command).
