# inbox-app-list-verified

evidence · `01M674CMQB8GNDK2SNZW4MQ3M6` · actor [claude-code-70d14e71](../actor/claude-code-70d14e71.md)

## result

PASS

## statement

At d29e512 I read the diff and ran go test ./internal/tui/... ./internal/api/... and the testscripts tui-inbox-app, rex-install, describe-golden and inbox-*: all pass, as does go test ./... and make fmt-check vet. internal/tui/inboxapp tests build the app over a fake runner with one item of each kind (question, signoff, failed, gave-up, skipped) and assert the Questions, Sign-offs, Failed runs, Gave up and Skipped headings in order, each row's change set title and step, the question text or failure detail, and 'recommends: sqlite'; an empty inbox shows 'Nothing waits on you'; a TickMsg and the r key each rerun the operation; the cursor keeps its item across a reload; the app is in launcher.Apps(). rex-install.txtar greps rhiza_app_inbox and the tui app inbox command in the generated Lua, and tui-inbox-app.txtar cmps rhiza tui app inbox --once with rhiza inbox --json for an empty inbox and one with a failed run. make self-check fails only on stale .rhiza/rendered actor and index files, which rhiza render rewrites after landing, not on this leaf's code.

## Links

### Out

- verifies [inbox-app-list](../capability/rhiza/inbox-app-list.md)

## History

- `01M6732DJSQ14B2WTZYZ9MPRA5` An inbox app lists what waits on the owner, questions, sign-offs and failed runs, and acts on each in place: answer, approve, request changes, retry, abandon or open its tab
