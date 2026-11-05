# request-app-form-verified

evidence · `01M6NACFTMBK2S85GZPQACEZ0V` · actor [claude-code-de725362](../actor/claude-code-de725362.md)

## result

PASS

## statement

Verifier at 2cde0ee9: read the diff (internal/tui/requestapp replaces requestsapp; one rexuikit Form of a kind Select and a validated TextArea, hosted by modalshell; cmd/rhiza tui wires it). Ran go test -trimpath -count=1 -v ./internal/tui/requestapp/...: TestSubmitQueuesAndCloses (fake runner gets `request --kind bug -- login fails`, shell closes), TestEmptyTextRefused (empty and blank text: Form validator message shown, no call, not closed), TestEscCallsNothing (closes, no call), TestRegistered (launcher has app request with a component and size) all PASS. Also go test ./cmd/rhiza -run 'TestScripts/(tui-request-app|rex-install)$|TestRequestsAppArgs' PASS, go vet and go build ./... clean, no remaining requestsapp references.

## Links

### Out

- verifies [request-app-form](../capability/rhiza/request-app-form.md)

## History

- `01M6MYNDJ4HAJF028CEAE7SYQN` Request app: file a feature, bug or milestone request in one rexuikit Form, then close; rebuilt on the modal shell over the request operation, replacing internal/tui/requestsapp and retiring the capabilities of the old requests app it replaces
