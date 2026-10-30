# inbox-app-actions

capability · `01M6734J6ESHFYKFJ09YK6TEWT` · parent [rhiza](../rhiza.md)

## acceptance

A test with fake operations and a fake exec drives each action and asserts the call it makes: o calls inbox.open with the item ID; an answer typed on a question calls session.message to its session with the typed text, and t sends the recommendation; y on a sign-off calls inbox.approve, c with text calls inbox.request-changes with it; R on a failed run calls rex window new running rhiza run resume <run>; x then y calls run.abandon, x then n calls nothing; no test runs real Rex or an agent

## statement

On the selected inbox item: enter or o focuses its Rex tab (inbox open); on a question, typing an answer and enter sends it to the asking session as a message (session message), or t takes the advisor's recommendation; on a sign-off, y approves it and c asks for changes with typed findings; on a failed, gave-up or skipped run, R retries it in a new Rex tab (rhiza run resume) and x abandons it after a y/n; each action reloads the list

## Links

### Out

- serves [m5-inbox-app](../../intent/milestone-5/m5-inbox-app.md)
- decided-by [app-runs-open-rex-tabs](../../decision/app-runs-open-rex-tabs.md)

### In

- verifies [inbox-app-actions-verified](../../evidence/inbox-app-actions-verified.md)
- verifies [inbox-app-actions-qa-1](../../evidence/inbox-app-actions-qa-1.md)
- verifies [inbox-app-actions-verified-2](../../evidence/inbox-app-actions-verified-2.md)
- verifies [inbox-app-actions-qa-2](../../evidence/inbox-app-actions-qa-2.md)
- verifies [pr-119-inbox-app](../../evidence/pr-119-inbox-app.md)

## History

- `01M6732DJSQ14B2WTZYZ9MPRA5` An inbox app lists what waits on the owner, questions, sign-offs and failed runs, and acts on each in place: answer, approve, request changes, retry, abandon or open its tab
