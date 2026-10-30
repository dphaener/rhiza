# inbox-app-list

capability · `01M6734GPG3SZCY93EBRDS1CA8` · parent [rhiza](../rhiza.md)

## acceptance

A test in internal/tui builds the inbox app over a fake inbox result with one item of each kind and asserts the view groups them under kind headings with each row's change set, step and question or detail, shows a recommendation where present, and shows Nothing waits on you for an empty inbox; a tick message reloads it; the app appears in the launcher's registry and the Rex install script; a testscript checks --once against rhiza inbox --json

## statement

An inbox app, registered as a Rex app like the others, lists the items of rhiza inbox as a cursor list grouped by kind (questions, sign-offs, failed runs, gave-up and skipped), each row with its change set's title, step and what it asks or why it failed, and a question's advisor recommendation when it has one; it reloads every few seconds and on r, says Nothing waits on you when empty, and rhiza tui app inbox --once prints what rhiza inbox --json prints

## Links

### Out

- serves [m5-inbox-app](../../intent/milestone-5/m5-inbox-app.md)
- decided-by [m5-roadmap-inbox-first](../../decision/m5-roadmap-inbox-first.md)

### In

- verifies [inbox-app-list-verified](../../evidence/inbox-app-list-verified.md)
- verifies [inbox-app-list-qa-1](../../evidence/inbox-app-list-qa-1.md)
- verifies [pr-119-inbox-app](../../evidence/pr-119-inbox-app.md)

## History

- `01M6732DJSQ14B2WTZYZ9MPRA5` An inbox app lists what waits on the owner, questions, sign-offs and failed runs, and acts on each in place: answer, approve, request changes, retry, abandon or open its tab
