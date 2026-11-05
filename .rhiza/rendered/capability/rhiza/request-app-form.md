# request-app-form

capability · `01M6MZB5FB63DG4763MZBYRBPM` · parent [rhiza](../rhiza.md)

## acceptance

Go tests in the request app's package with a fake runner: submitting calls request with --kind and the typed text and the app quits; an empty text is refused by the Form's validator and calls nothing; Esc calls nothing; the app registers with the launcher as request.

## statement

The Request app is one rexuikit Form, a Select of kind (feature, bug or milestone) and a TextArea for the request's text; submitting queues it through the request operation and closes the app, and Esc closes it with nothing queued.

## Links

### Out

- decided-by [tui-five-apps](../../decision/tui-five-apps.md)
- decided-by [tui-rexuikit-components](../../decision/tui-rexuikit-components.md)

### In

- verifies [request-app-form-verified](../../evidence/request-app-form-verified.md)
- verifies [request-app-form-qa-1](../../evidence/request-app-form-qa-1.md)
- verifies [request-app-form-qa-2](../../evidence/request-app-form-qa-2.md)
- verifies [pr-294-request-app-on-the-modal-shell](../../evidence/pr-294-request-app-on-the-modal-shell.md)

## History

- `01M6MYNDJ4HAJF028CEAE7SYQN` Request app: file a feature, bug or milestone request in one rexuikit Form, then close; rebuilt on the modal shell over the request operation, replacing internal/tui/requestsapp and retiring the capabilities of the old requests app it replaces
